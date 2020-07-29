> 原文链接：  
> [如何开发高性能 go 组件？https://zhuanlan.zhihu.com/p/41991119](https://zhuanlan.zhihu.com/p/41991119)  
> [在Go语言项目中使用Zap日志库 https://zhuanlan.zhihu.com/p/88856378](https://zhuanlan.zhihu.com/p/88856378)  
> [Go 每日一库之 zap https://zhuanlan.zhihu.com/p/136093026](https://zhuanlan.zhihu.com/p/136093026)

> GitHub：[https://github.com/uber-go/zap](https://github.com/uber-go/zap)  
> doc：[https://pkg.go.dev/go.uber.org/zap?tab=doc](https://pkg.go.dev/go.uber.org/zap?tab=doc)


[TOC]

---------------------------------------------------------
# zap

Zap是非常快的、结构化的，分日志级别的Go日志库。

根据Uber-go Zap的文档，它的性能比类似的结构化日志包更好——也比标准库更快。 以下是Zap发布的基准测试信息

记录一条消息和10个字段:

![zap性能基准1.png](https://upload-images.jianshu.io/upload_images/14738618-3a0c5003bd8b6287.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

记录一个静态字符串，没有任何上下文或printf风格的模板：

![zap性能基准2.png](https://upload-images.jianshu.io/upload_images/14738618-739d8c4c8600b84f.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


## why zap ?


绝大多数的代码中的写日志通常通过各式各样的日志库来实现。日志库提供了丰富的功能，对于 go 开发者来说大家常用的日志组件通常会有以下几种，下面简单的总结了常用的日志组件的特点：

- seelog: 最早的日志组件之一，功能强大但是性能不佳，不过给社区后来的日志库在设计上提供了很多的启发。
- logrus: 代码清晰简单，同时提供结构化的日志，性能较好。
- zap: uber 开源的高性能日志库，面向高性能并且也确实做到了高性能。

Zap 代码并不是很多，不到 5000 行，比 seelog 少多了（ 8000 行左右）， 但比logrus（不到 2000 行）要多很多。为什么我们会选择 zap 呢？在下文中将为大家阐述。

```
79 text files.
79 unique files.
 6 files ignored.

https://github.com/AlDanial/cloc v 1.66  T=0.95 s (77.8 files/s, 8180.4 lines/s)
-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                              70            963           1863           4821
make                             1             10              4             57
YAML                             2              0              0             51
Bourne Shell                     1              3              0             12
-------------------------------------------------------------------------------
SUM:                            74            976           1867           4941
-------------------------------------------------------------------------------
```

Zap 跟 logrus 以及目前主流的 go 语言 log 类似，提倡采用结构化的日志格式，而不是将所有消息放到消息体中，简单来讲，日志有两个概念：字段和消息。字段用来结构化输出错误相关的上下文环境，而消息简明扼要的阐述错误本身。

比如，用户不存在的错误消息可以这么打印:

```
log.Error("User does not exist", zap.Int("uid", uid))
```

上面 `User does not exist` 是消息， 而 uid 是字段。具体设计思想可以参考 logrus的文档 ，这里不再赘述。

其实我们最初的实践中并没有意识到日志框架的性能的重要性，直到开发后期进行系统的 benchmark 总是不尽人意，而且在不同的日志级别下性能差距明显。通过 go profiling 看到日志组件对于计算资源的消耗十分巨大，因此决心将其替换为一个高性能的日志框架，这也是选择用 zap 的一个重要的考量的点。

目前我们使用 zap 已有2年多的时间，zap 很好地解决了日志组件的低性能的问题。目前 zap 也从 beta 发布到了 1.8版本，对于 zap 我们不仅仅看到它的高性能，更重要的是理解它的设计与工程实践。日志属于 io 密集型的组件，这类组件如何做到高性能低成本，这也将直接影响到服务成本。


## zap, how ?


zap 具体表现如何呢？抛开 zap 的设计我们不谈，现在让我们单纯来看一个日志库究竟需要哪些元素:

1. 首先要有输入：输入的数据应该被良好的组织且易于编码，并且还要有高效的空间利用率，毕竟内存开辟回收是昂贵的。 无论是 formator 方式还是键值对 key-value 方式，本质上都是对于输入数据的组织形式。 实践中有格式的数据往往更有利于后续的分析与处理。 json 就是一种易用的日志格式
2. 其次日志能够有不同的级别：对于日志来说，基本的的日志级别: debug info warning error fatal 是必不可少的。对于某些场景，我们甚至期待类似于 assert 的 devPanic 级别。同时除了每条日志的级别，还有日志组件的级别，这可以用于屏蔽掉某些级别的日志。
3. 有了输入和不同的级别，接下来就需要组织日志的输出流：你需要一个 encoder 帮你把格式化的，经过了过滤的日志信息输出。也就是说不论你的输出是哪里，是 stdout ，还是文件，还是 NFS ，甚至是一个 tcp 连接。 Encoder 只负责高效的编码数据，其他的事情交给其他人来做。
4. 有了这些以后，我们剩下的需求就是设计一套易用的接口，来调用这些功能输出日志。 这就包含了 logger 对象和 config。

嗯，似乎我们已经知道我们要什么了，日志的组织和输出是分开的逻辑，但是这不妨碍 zapcore 将这些设计组合成 zap 最核心的接口。

![zap结构图.png](https://upload-images.jianshu.io/upload_images/14738618-d4476b922c97b95b.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

![zap包简介.png](https://upload-images.jianshu.io/upload_images/14738618-58d1d106ad09de66.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

从上文中来看一个日志库的逻辑结构已经很清晰了。


---------------------------------------------------------
# use zap


通过 zap 打印一条结构化的日志大致包含5个过程：

1. 分配日志 Entry: 创建整个结构体，此时虽然没有传参(fields)进来，但是 fields 参数其实创建了
2. 检查级别，添加core: 如果 logger 同时配置了 hook，则 hook 会在 core check 后把自己添加到 cores 中
3. 根据选项添加 caller info 和 stack 信息: 只有大于等于级别的日志才会创建checked entry
4. Encoder 对 checked entry 进行编码: 创建最终的 byte slice，将 fields 通过自己的编码方式(append)编码成目标串
5. Write 编码后的目标串，并对剩余的 core 执行操作， hook 也会在这时被调用

接下来对于我们最感兴趣的几个部分进行更加具体的分析：

- logger: zap 的接口层，包含Log 对象、Level 对象、Field 对象、config 等基本对象
- zapcore: zap 的核心逻辑，包含field 的管理、level 的判断、encode 编码日志、输出日志
- encoder: json 或者其它编码方式的实现
- utils: SubLog，Hook，SurgarLog/grpclogger/stdlogger


## logger: 对象 vs 接口

zap 对外提供的是 logger 对象和 field 和 level。 这是 zap 对外提供的基本语义: logger 对象打印 log，field 则是 log 的组织方式，level 跟打印的级别相关。 这些元素的组合是松散的但是联系确实紧密的。

有趣的是，zap 并没有定义接口。 大家可能也很容易联想到 go 自身的 log 就不是接口。 在 go-dev 很多人曾经讨论过 go 的接口，有人讨论为啥不提供接口 [Standardization around logging and related concerns](https://link.zhihu.com/?target=https%3A//groups.google.com/forum/%23%21msg/gola%3C/b%3Eng-dev/F3l9Iz1JX4g/szAb07lgFAAJ) ，甚至有人提出过草案 [Go Logging Design Proposal - Ross Light](https://link.zhihu.com/?target=https%3A//docs.google.com/document/d/1nFRxQ5SJVPpIBWTFHV-q5lBYiwGrfCMkESFGNzsrvBU/edit%23)，然而最终也难逃被 Abandon 的命运。

归根到底，reddit 上的一条评论总结最为到位:

> No one seems to be able to agree what such an interface should look like。

在 zap 的早期版本中并没有提供 zapcore 这个库。 zapcore 提供了zap 最核心的设计的逻辑封装：执行级别判断，添加 field 和 core，进行级别判断返回 checked entry。

logger 是对象不是接口，但是 zapcore 却是接口，logger 依赖 core 接口实现功能，其实是 logger 定义了接口，而 core 提供了接口的实现。 core 作为接口体现的是 zap 核心逻辑的抽象和设计理念，因此只需要约定逻辑，而实现则是多种的也是可以替换的，甚至可以基于 core 进行自定义的开发，这大大增加了灵活性。


## zap field: format vs field

对于 zap 来说，为了性能其实牺牲掉了一定的易用性。例如 log.Printf("%s", &s) format 这种方式是最自然的 log 姿势，然而对于带有反射的 go 是致命的: 反射太过耗时。

下面让我们先来看看反射和 cast 的性能对比，结果是惊人的。

![反射&cast性能对比.png](https://upload-images.jianshu.io/upload_images/14738618-605dc036f60915e9.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

通过 fmt.Sprintf() 来组合 string 是最慢的，这因为 fmt.Printf 函数族使用反射来判断类型.

```
fmt.Sprintf("%s", "hello world")
```

相比之下 string 的 + 操作基就会快很多，因为 go 中的 string 类型本质上是一个特殊的 []byte。 执行 + 会将后续的内容追加在 string 对象的后面:

```
_ = s + "hello world"
```

然而对于追求极致的 zap 而言还不够，**如果是 []byte 的 append 则还要比 + 快2倍以上**。 尽管这里其实不是准确的，因为分配 []byte 时，如果不特殊指定 capacity 是会按照 2 倍的容量预分配空间。append 追加的 slice 如果容量不足，依然会引发一次 copy, 而 我们**可以通过预分配足够大容量的 slice 来避免该问题**。zap 默认分配 1k 大小的 byte slice。

```
buf = append(buf, []byte("hello world")...)
```

表格的最下面是接口反射和直接转换的性能对比，field 通过指明类型避免了反射，zap 又针对每种类型直接提供了转换 []byte + append 的函数，这样的组合效率是极其高的。明确的调用对应类型的函数避免运行时刻的反射，可以看到 规避反射 这种类型操作是贯穿在整个 zap 的逻辑中的。

zap 的 append 家族函数封装了 strconv.AppendX 函数族，该函数用于将类型转换为 []byte 并 append 到给定的 slice 上。


## zap 高性能的秘诀

对于大部分人来说，标准库提供了覆盖最全的工具和函数。 但是标准库 为了通用有时候其实做了一些性能上的牺牲 。 而 zap 在细节上的性能调优确实下足了功夫，我们可以借鉴这些调优的思路和经验。

### 避免 GC: 对象复用

go 是提供了 gc 的语言。 gc 就像双刃剑，给你了快捷的同时又会增加系统的负担。 尽管 go 官方宣称 gc 性能很好，但是仍然无法绕开 Stop-The-World 的难题，一旦内存中的碎片较多 gc 仍然会有明显尖峰，这种尖峰对于重 io 的业务来说是致命的。 zap 每打印1条日志，至少需要2次内存分配:

- 创建 field 时分配内存。
- 将组织好的日志格式化成目标 []byte 时分配内存。

zap 通过 sync.Pool 提供的对象池，复用了大量可以复用的对象，避开了 gc 这个大麻烦。

go 的 sync.Pool 实现上提供的是 runtime 级别的绑定到 Processor 的对象池。 对象池是否高效依赖于这个池的竞争是否过多，对此我曾经做过一次对比，使用 channel 实现了一个最简单的对象池，但是 benchmark 的结果却不尽如人意，完全不如 sync.Pool 高效。 究其原因，其实也可以理解，因为使用 channel 实现的对象池在多个 Processor 之间会有强烈的并发。尽管使用 sync.Pool 涉及到一次接口的转换，性能依然是非常可观的。

zap 也是使用了 go sync.Pool 提供的标准对象池。自身的 runtime 包含了 P 的处理逻辑，每个 P 都有自己的池在调度时不会发生竞争。 这个比起代码中的软实现更加高效，是用户代码做不到的逻辑。

sync.Pool 是go提供给大家的一种优化 gc 的方式好方式尽管 go 的 gc 已经有了长足的进步但是仍然不能够绕开 gc 的 STW，因此合理的使用 pool 有助于提高代码的性能，防止过多碎片化的内存分配与回收。

之前我们对于 pool 对象的讨论中，最痛苦的一点就是是否应该包暴露 Free 函数。 最终的结论是如同 C/C++，资源的申请者应该决定何时释放。 zap 的对象池管理也深谙此道。

- buffer 实现了 io.Writer
- Put 时并不执行 Reset
- buffer 对象包含其对象池，因此可以在任何时刻将自己释放（放回对象池）

### 内建的 Encoder: 避免反射

反射是 go 提供给我们的另一个双刃剑，方便但是不够高效。 对于 zap ，规避反射贯穿在整个代码中。 对于我们来说，创建json 对象只需要简单的调用系统库即可:

```
b, err := json.Marshal(&obj)
```

对于 zap 这还不够。标准库中的 json.Marshaler 提供的是基于类型反射的拼接方式，代价是高昂的:

```
func (e *encodeState) marshal(v interface{}, opts encOpts) (err error) {
    ...
    e.reflectValue(reflect.ValueOf(v), opts) //reflect 根据 type 进行 marshal
    ...
}
```

反射的整体性能并不够高，因此通过 go 的反射可能导致额外的性能开销。 zap 选择了自己实现 json Encoder。 通过明确的类型调用，直接拼接字符串，最小化性能开销。

zap 的 json Encoder 设计的高效且较为易用，完全可以替换在代码中。 另一方面，这也是 go 长期以来缺乏泛型的一个痛点 。对于一些性能要求高的操作，如果标准库偏向于易用性。那么我们完全可以绕开标准库，通过自己的实现，规避掉额外性能开销。 同样，上文提到的 field 也是这个原因。 通过一个完整的自建类型系统，zap 提供了从组合日志到编码日志的整体逻辑，整个过程中都是可以。

> ps. 据说 go 在2.0 中就会加入泛型啦, 很期待

### 避免竞态

zap 的高效还体现在对于并发的控制上。 zap 选择了 写时复制机制 。 zap 把每条日志都抽象成了 entry。 对于 entry 还分为2种不同类型:

- Entry : 包含原始的信息 但是不包含 field
- CheckedEntry: 经过级别检查后的生成 CheckedEntry，包含了 Entry 和 Core。

CheckedEntry 的引入解决了组织日志，编码日志的竟态问题。只有经过判断的对象才会进入后续的逻辑，所有的操作 **写时触发复制** ，没有必要的操作不执行预分配。将操作与对象组织在一起，不进行资源的竞争，避免多余的竟态条件。

对于及高性能的追求者来说，预先分配的 field 尽管有 pool 加持仍然是多余的，因此 zap 提供了更高性能的接口，可以避免掉 field 的分配:

```
if ent := log.Check(zap.DebugLevel, "foo"); ent != nil {
    ent.Write(zap.String("foo", "bar"))
}
```

通过这一步判断，如果该级别日志不需要打印，那么连 field 都不必生成。 避免一切不必要的开销，zap 确实做到了而且做得很好。


## 多样的功能与简单的设计理念

### level handler:

level handler 是 zap 提供的一种 level 的处理方式，通过 http 请求动态改变日志组件级别。

对于日志组件的动态修改，seelog 最早就有提供类似功能，基于 xml 文件修改捕获新的级别。 但是 xml 文件显然不够 golang。

zap 的解决方案是 http 请求。http 是大家广泛应用的协议，zap 定义了 level handler 实现了 http.Handler 接口

go 自身的 http 服务实现起来非常的简洁:

```
http.HandleFunc("/handle/level", zapLevelHandler)
if err := http.ListenAndServe(addr, nil); err != nil {
    panic(err)
}
```

简单几行代码就能实现 http 请求控制日志级别的能力。 通过 GET 获取当前级别，PUT 设置新的级别。

### zap 的 surgar log 和易用 config 接口封装

我们的库往往希望提供事无巨细的控制能力。但是对于简单的使用者就不够友好，繁杂的配置往往容易使人第一次使用即失去耐心。同时，一个全新的 log 接口设计也容易让长期使用 format 方式打印日志的人产生疑问。在工作中发现较多的用户有这样的需求: 你的这个库怎么用?

显然只有 godoc 还不够。

zap 的 Config 非常的繁琐也非常强大，可以控制打印 log 的所有细节，因此对于我们开发者是友好的，有利于二次封装。但是对于初学者则是噩梦。因此 zap 提供了一整套的易用配置，大部分的姿势都可以通过一句代码生成需要的配置。

```
func NewDevelopmentEncoderConfig() zapcore.EncoderConfig
func NewProductionEncoderConfig() zapcore.EncoderConfig
type SamplingConfig
```

同样，对于不想付出代价学习使用 field 写格式化 log 的用户，zap 提供了 sugar log。 sugarlog 字面含义就是加糖。 给 zap 加点糖 ，sugar log 提供了 formatter 接口，可以通过 format 的方式来打印日志。sugar 的实现封装了 zap log，这样既满足了使用 printf 格式串的兼容性需求，同时也提供了更多的选择，对于不那么追求极致性能的场景提供了易用的方式。

```
sugar := log.Sugar()
sugar.Debugf("hello, world %s", "foo")
```

### zap logger 提供的 utils

zap 还在 logger 这层提供了丰富的工具包，这让整个 zap 库更加的易用:

- grpc logger：封装 zap logger 可以直接提供给 grpc 使用，对于大多数的 go 分布式程序，grpc 都是默认的 rpc 方案，grpc 提供了 SetLogger 的接口。 zap 提供了对这个接口的封装。
- hook：作为 zap。Core 的实现，zap 提供了 hook。 使用方实现 hook 然后注册到 logger，zap在合适的时机将日志进行后续的处理，例如写 kafka，统计日志错误率 等等。
- std Logger: zap 提供了将标准库提供的 logger 对象重定向到 zap logger 中的能力，也提供了封装 zap 作为标准库 logger 输出的能力。 整体上十分易用。
- sublog: 通过创建 绑定了 field 的子logger，实现了更加易用的功能。


## zap 的好帮手: RollingWriter


zap 本身提供的是设置 writer 的接口，为此我实现了一套 io.Writer，通过rolling writer 实现了 log rotate 的功能。

[rollingWriter](https://link.zhihu.com/?target=https%3A//github.com/arthurkiller/rollingWriter) 是一个 go ioWriter 用于按照需求自动滚动文件。 目的在于内置的实现 logrotate 的功能而且更加高效和易用。

具体可以见

> https://github.com/arthurkiller/rollingWrite


---------------------------------------------------------
# Polaris项目实践


```
//rawJSON

{
  "zap_config": {
    "level": "debug",
    "development": false,
    "disableCaller": false,
    "disableStacktrace": false,
    "sampling": {
      "initial": 100,
      "thereafter": 100
    },
    "encoding": "console",
    "encoderConfig": {
      "messageKey": "msg",
      "levelKey": "level",
      "timeKey": "ts",
      "nameKey": "logger",
      "callerKey": "caller",
      "stacktraceKey": "stacktrace",
      "lineEnding": "\n",
      "levelEncoder": "capitalColor",
      "timeEncoder": "rfc3339",
      "durationEncoder": "seconds",
      "callerEncoder": "short",
      "nameEncoder": ""
    },
    "outputPaths": [
      "./log/runtime.log"
    ],
    "errorOutputPaths": [
      "stderr"
    ]
  },
  "lumberjack_config": {
    "@filename": "日志文件的位置",
    "filename": "./log/runtime.log",
    "@maxsize": "在进行切割之前，日志文件的最大大小（以MB为单位）",
    "maxsize": 100,
    "@maxbackups": "保留旧文件的最大个数",
    "maxbackups": 10,
    "@maxage": "保留旧文件的最大天数",
    "maxage": 30,
    "@compress": "是否压缩/归档旧文件",
    "compress": false
  }
}
```

```

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapWithLumberjackConfig struct {
	ZapConfig        zap.Config        `json:"zap_config"`
	LumberjackConfig lumberjack.Logger `json:"lumberjack_config"`
}

func NewZapWithLumberjack(logConfigPath string) (*zap.SugaredLogger, error) {
	rawJSON, err := ioutil.ReadFile(logConfigPath)
	if err != nil {
		return nil, err
	}

	var cfg ZapWithLumberjackConfig
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.ZapConfig.OutputPaths) == 0 {
		return nil, errors.New("not set cfg.OutputPaths")
	}

	writeSyncer := zapcore.AddSync(&cfg.LumberjackConfig)
	encoder := zapcore.NewConsoleEncoder(cfg.ZapConfig.EncoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	//zap提供了几个快速创建logger的方法，zap.NewExample()、zap.NewDevelopment()、zap.NewProduction()，还有高度定制化的创建方法zap.New()。
	//创建前 3 个logger时，zap会使用一些预定义的设置，它们的使用场景也有所不同。
	//Example适合用在测试代码中，Development在开发环境中使用，Production用在生成环境。
	//logger := zap.NewExample()
	//logger, _ := zap.NewDevelopment()
	//logger, _ := zap.NewProduction()
	//logger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	//if err != nil {
	//	return nil, err
	//}

	//当然，每个字段都用方法包一层用起来比较繁琐。zap也提供了便捷的方法SugarLogger，可以使用printf格式符的方式。
	//SugaredLogger的使用比Logger简单，只是性能比Logger低 50% 左右，可以用在非热点函数中。
	sugarLogger := logger.Sugar()

	return sugarLogger, nil
}



//****************** Runtime Log **********************
//- DebugLevel Level = iota - 1
//- InfoLevel
//- WarnLevel
//- ErrorLevel
//- DPanicLevel
//- PanicLevel
//- FatalLevel

//_minLevel = DebugLevel  
//_maxLevel = FatalLevel

func (dl *PolarisLog) Debug(args ...interface{}) {
	if dl.zapRunTime != nil {
		dl.zapRunTime.Debug(args)
	}
}

func (dl *PolarisLog) Info(args ...interface{}) {
	if dl.zapRunTime != nil {
		dl.zapRunTime.Info(args)
	}
}

func (dl *PolarisLog) Infof(format string, args ...interface{}) {
	if dl.zapRunTime != nil {
		dl.zapRunTime.Infof(format, args)
	}
}

func (dl *PolarisLog) Warn(args ...interface{}) {
	if dl.zapRunTime != nil {
		dl.zapRunTime.Warn(args)
	}
}

func (dl *PolarisLog) Warnf(format string, args ...interface{}) {
	if dl.zapRunTime != nil {
		dl.zapRunTime.Warnf(format, args)
	}
}

func (dl *PolarisLog) Error(args ...interface{}) {
	if dl.zapRunTime != nil {
		dl.zapRunTime.Error(args)
	}
}


//****************** Flush **********************
func (this *PolarisLog) Flush() {

	//zap底层 API 可以设置缓存，所以一般使用defer logger.Sync()将缓存同步到文件中。
	if this.zapRunTime != nil {
		this.zapRunTime.Sync()
	}

}

```


# 总结

zap 在整体设计上有非常多精细的考量，不仅仅是在高性能上面的出色表现，更多的意义是其设计和工程实践上。此处总结下 zap 的代码之道:

- 合理的代码组织结构，结构清晰的抽象关系
- 写实复制，避免加锁
- 对象内存池，避免频繁创建销毁对象
- 避免使用 fmt json/encode 使用字符编码方式对日志信息编码，适用byte slice 的形式对日志内容进行拼接编码操作
- 其实 zap 带给我们的远不止这些，在这里建议有兴趣的朋友一定要抽时间看一下 zap 的源码，确实有很多细节需要我们细细体味。
