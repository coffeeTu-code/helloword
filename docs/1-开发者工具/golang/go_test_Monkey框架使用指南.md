# Monkey框架使用指南

[TOC]

## 序言

在软件开发中，产品代码的正确性通过测试代码来保证，而测试代码的正确性谁来保证？答案是毫无争议的，肯定是程序员自己。这就要求测试代码必须足够简单且表达力强，让错误无处藏身。我们要有一个好鼻子，能够嗅出测试的坏味道，及时的进行测试重构，从而让测试代码易于维护。

要写出好的测试代码，必须精通相关的测试框架。对于Golang的程序员来说，至少需要掌握下面四个测试框架：

- GoConvey
- GoStub
- GoMock
- Monkey

本文将主要介绍**Monkey框架**的基本使用方法，从而指导读者更好的进行测试实践，最终写出简单优雅的测试代码。

## Monkey简介

通过前面四篇文章，我们已经掌握了框架GoConvey + GoStub + GoMock组合使用的正确姿势，同时已经知道：

1. 全局变量可通过GoStub框架打桩
2. 过程可通过GoStub框架打桩
3. 函数可通过GoStub框架打桩
4. interface可通过GoMock框架打桩

但还有两个问题比较棘手：

1. 方法（成员函数）无法通过GoStub框架打桩，当产品代码的OO设计比较多时，打桩点可能离被测函数比较远，导致UT用例写起来比较痛
2. 过程或函数通过GoStub框架打桩时，对产品代码有侵入性

下面我们举两个例子，阐述GoStub框架对产品代码的侵入性

例一：函数定义侵入

```
func Exec(cmd string, args ...string) (string, error) {
    ...
}
```

上面的函数Exec的定义为常规方式，但这时不能通过GoStub框架对函数Exec进行打桩，除非将函数Exec定义为非常规方式（侵入性）：

```
var Exec = func(cmd string, args ...string) (string, error) {
    ...
}
```

例二：适配层侵入

产品代码中很多函数都会调用Golang的库函数或第三方的库函数，这些库函数的定义显然是常规方式，要想通过GoStub框架对这些函数打桩，一般会在适配层定义相关的变量（侵入性）：

```
package adapter

var Stat = os.Stat
var Marshal = json.Marshal
var UnMarshal = json.Unmarshal
...
```

本文将介绍第四个框架Monkey的使用方法，目的是解决这两个棘手的问题，同时考虑将GoStub的优点集成到Monkey。

> Monkey不是线程安全的，不要将Monkey用于并发的测试中。

## 安装

在命令行运行下面的命令：

```
go get github.com/bouk/monkey
```

运行完后你会发现，在$GOPATH/src/github.com目录下，新增了bouk/monkey子目录，这就是本文的主角。

## 基本使用方法

Monkey框架的使用场景很多，依次为：

1. 基本场景：为一个函数打桩
2. 基本场景：为一个过程打桩
3. 基本场景：为一个方法打桩
4. 复合场景：由任意相同或不同的基本场景组合而成
5. 特殊场景：桩中桩的一个案例

### 为一个函数打桩

Exec是infra层的一个操作函数，实现很简单，代码如下所示：

```go
// infra/os-encap/exec.go

func Exec(cmd string, args ...string) (string, error) {
    cmdpath, err := exec.LookPath(cmd)
    if err != nil {
        fmt.Errorf("exec.LookPath err: %v, cmd: %s", err, cmd)
        return "", infra.ErrExecLookPathFailed
    }

    var output []byte
    output, err = exec.Command(cmdpath, args...).CombinedOutput()
    if err != nil {
        fmt.Errorf("exec.Command.CombinedOutput err: %v, cmd: %s", err, cmd)
        return "", infra.ErrExecCombinedOutputFailed
    }
    fmt.Println("CMD[", cmdpath, "]ARGS[", args, "]OUT[", string(output), "]")
    return string(output), nil
}
```

Exec函数的实现中调用了库函数exec.LoopPath和exec.Command，因此Exec函数的返回值和运行时的底层环境密切相关。在UT中，如果被测函数调用了Exec函数，则应根据用例的场景对Exec函数打桩。

Monkey的API非常简单和直接，我们直接看打桩代码：

```
import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
    . "github.com/bouk/monkey"
    "infra/osencap"
)

const any = "any"

func TestExec(t *testing.T) {
    Convey("test has digit", t, func() {
        Convey("for succ", func() {
            outputExpect := "xxx-vethName100-yyy"
            guard := Patch(osencap.Exec, func(_ string, _ ...string) (string, error) {
                return outputExpect, nil
            })
            defer guard.Unpatch()
            output, err := osencap.Exec(any, any)
            So(output, ShouldEqual, outputExpect)
            So(err, ShouldBeNil)
        })
    })
}
```

Patch是Monkey提供给用户用于函数打桩的API：

- 第一个参数是目标函数的函数名
- 第二个参数是桩函数的函数名，习惯用法是匿名函数或闭包
- 返回值是一个PatchGuard对象指针，主要用于在测试结束时删除当前的补丁

### 为一个过程打桩

当一个函数没有返回值时，该函数我们一般称为过程。很多时候，我们将资源清理类函数定义为过程。

我们对过程DestroyResource的打桩代码为：

```go
guard := Patch(DestroyResource, func(_ string) {

})
defer guard.Unpatch()
```

### 为一个方法打桩

当微服务有多个实例时，先通过Etcd选举一个Master实例，然后Master实例为所有实例较均匀的分配任务，并将任务分配结果Set到Etcd，最后Master和Node实例Watch到任务列表，并过滤出自身需要处理的任务列表。

我们用类Etcd的方法Get来模拟获取任务列表的功能，入参为instanceId：

```go
type Etcd struct {

}

func (e *Etcd) Get(instanceId string) []string {
    taskList := make([]string, 0)
    ...
    return taskList
}
```

我们对Get方法的打桩代码如下：

```
var e *Etcd
guard := PatchInstanceMethod(reflect.TypeOf(e), "Get", func(_ *Etcd, _ string) []string {
    return []string{"task1", "task5", "task8"}
})
defer guard.Unpatch()
```

PatchInstanceMethod API是Monkey提供给用户用于方法打桩的API：

- 在使用前，先要定义一个目标类的指针变量x
- 第一个参数是reflect.TypeOf(x)
- 第二个参数是字符串形式的函数名
- 返回值是一个PatchGuard对象指针，主要用于在测试结束时删除当前的补丁

### 任意相同或不同的基本场景组合

假设Px为用于函数、过程或方法打桩的API调用，则任意相同或不同基本场景组合的打桩过程形式化表达为：

```
Px1
defer UnpatchAll()
Px2
...
Pxn
```

该测试执行完后，函数UnpatchAll将删除所有的补丁。

#### 桩中桩的一个案例

在某些特殊场景下（比如反序列化），函数或方法既有返回值，又有出参。出参一般为指针类型，包括具体的指针类型（比如*int)和抽象的指针类型（一般为interface{}）。我们常用的库函数json.Unmarshal就属于这种情况。

笔者在实践中遇到的出参类型大多是具体的指针类型，其指针变量指向的内存不管在传入前确定还是在传入后确定，都将影响后面的代码逻辑。

下面呈现桩中桩的一个案例，以便大家灵活使用Monkey框架。

> 何谓桩中桩？  
>  
> interface中声明了一个方法，既有返回值，又有出参。在测试中，先通过GoMock框架打桩多态到mock方法，然后又通过Monkey框架跳转到补丁方法，最终修改出参并返回。在这个过程中，mock方法可以看作一个桩，补丁方法又可以看作mock方法的一个桩，即补丁方法是一个桩中桩。

定义一个具体类型Movie：

```
type Movie struct {
    Name string
    Type string
    Score int
}
```

定义一个interface类型Repository：

```
type Repository interface {
    Retrieve(key string, movie *Movie) error
    ...
}
```

桩中桩的一个测试用例：

```
func TestDemo(t *testing.T) {
    Convey("test demo", t, func() {
        Convey("retrieve movie", func() {
            ctrl := NewController(t)
            defer ctrl.Finish()
            mockRepo := mock_db.NewMockRepository(ctrl)
            mockRepo.EXPECT().Retrieve(Any(), Any()).Return(nil)
            Patch(redisrepo.GetInstance, func() Repository {
                return mockRepo
            })
            defer UnpatchAll()
            PatchInstanceMethod(reflect.TypeOf(mockRepo), "Retrieve", func(_ *mock_db.MockRepository, name string, movie *Movie) error {
                movie = &Movie{Name: name, Type: "Love", Score: 95}
                return nil
            })
            repo := redisrepo.GetInstance()
            var movie *Movie
            err := repo.Retrieve("Titanic", movie)
            So(err, ShouldBeNil)
            So(movie.Name, ShouldEqual, "Titanic")
            So(movie.Type, ShouldEqual, "Love")
            So(movie.Score, ShouldEqual, 95)
        })
        ...
    })
}
```

我们先通过Monkey框架的Patch API将mock对象注入，然后通过Monkey框架的PatchInstanceMethod API将mock方法跳转到补丁方法，间接完成对指针变量movie的内存分配及赋值，并返回nil。


## Monkey的缺陷及解决方案

### inline函数

Golang中虽然没有inline关键字，但仍存在inline函数，一个函数是否是inline函数由编译器决定。inline函数的特点是简单短小，在源代码的层次看有函数的结构，而在编译后却不具备函数的性质。inline函数不是在调用时发生控制转移，而是在编译时将函数体嵌入到每一个调用处，所以inline函数在调用时没有地址。

**inline函数没有地址的特性导致了Monkey框架的第一个缺陷：对inline函数打桩无效。**

模拟一个简单的inline函数：

```
func IsEqual(a, b string) bool {
    return a == b
}
```

对HasDigit函数进行打桩测试：

```
func TestIsEqual(t *testing.T) {
    Convey("test is equal", t, func() {
        Convey("for patch true", func() {
            guard := Patch(IsEqual, func(_, _ string) bool {
                return true
            })
            defer guard.Unpatch()
            ok := IsEqual("hello", "world")
            So(ok, ShouldBeTrue)
        })
    })
}
```

在命令行运行这个测试，结果不符合期望：

```
$ go test -v func_test.go -test.run TestIsEqual
=== RUN   TestIsEqual

  test is equal 
    for patch true ✘


Failures:

  * /Users/zhangxiaolong/Desktop/D/go-workspace/src/test/monkey/func_test.go 
  Line 67:
  Expected: true
  Actual:   false


1 total assertion

--- FAIL: TestIsEqual (0.00s)
FAIL
exit status 1
FAIL    command-line-arguments  0.006s
```

**解决方案：通过命令行参数-gcflags=-l禁止inline**

在命令行增加参数-gcflags=-l重新运行测试，结果符合期望：

```
go test -gcflags=-l -v func_test.go -test.run TestIsEqual
=== RUN   TestIsEqual

  test is equal 
    for patch true ✔


1 total assertion

--- PASS: TestIsEqual (0.00s)
PASS
ok      command-line-arguments  0.007s
```

### 方法名首字母小写

这一年多，Golang的版本在快速演进，上个月已经发布了go1.9版本。然而，一些团队可能一直还在用go1.6版本，并有计划在近期升级到go1.7或以上版本。

Monkey框架的实现中大量使用了反射机制，尤其是方法的补丁实现函数PatchInstanceMethod。但是，go1.6版本和更高版本（比如go1.7）的反射机制有些差异：在go1.6版本中反射机制会导出所有方法（不论首字母是大写还是小写），而在更高版本中反射机制仅会导出首字母大写的方法。

**反射机制的这种差异导致了Monkey框架的第二个缺陷：在go1.6版本中可以成功打桩的首字母小写的方法，当go版本升级后Monkey框架会显式触发panic，表示unknown method:**

```
m, ok := target.MethodByName(methodName)
if !ok {
    panic(fmt.Sprintf("unknown method %s", methodName))
}
```

说明：反射机制的差异并不波及Patch函数的实现，所以go版本升级前后首字母小写的函数名的打桩不受影响。

正交设计四原则告诉我们，要向稳定的方向依赖。首字母小写的方法或函数不是public的，仅在包内可见，不是一个稳定的依赖方向。如果在UT测试中对首字母小写的方法或函数打桩的话，会导致重构的成本比较大。

解决方案：不管现在团队使用的go版本是哪一个，都不要对首字母小写的方法或函数打桩，不但可以确保测试用例在go版本升级前后的稳定性，而且能有效降低重构的成本。

## API

在讨论Monkey的API之前，我们先回顾一下GoStub框架的API。

GoStub框架的API既包括函数API，也包括方法API。由于Monkey框架的API只涉及函数API，所以在这里我们只回顾GoStub框架的函数API。

我们先看GoStub框架的第一个函数API：

```
func Stub(varToStub interface{}, stubVal interface{}) *Stubs
```

这个API我们一般用于对全局变量打桩：

```
stubs := Stub(&num, 150)
defer stubs.Reset()
```

然而，这个API也可以用于函数打桩：

```
stubs := Stub(&osencap.Exec, func(_ string, _ ...string) (string, error) {
            return "xxx-vethName100-yyy", nil
})
defer stubs.Reset()
```

GoStub框架的Stub API对函数的打桩方法是不是和Monkey框架的API的使用方法很像？这是毋庸置疑的，这样的API才是原生的API，StubFunc API是专门针对函数或过程打桩的改进版：

```
func StubFunc(funcVarToStub interface{}, stubVal ...interface{}) *Stubs
```

StubFunc替代Stub对函数的打桩示例：

```
stubs := StubFunc(&osencap.Exec,"xxx-vethName100-yyy", nil)
defer stubs.Reset()
```

是不是简洁优雅了很多？

> 说明：一般情况下，Golang的桩函数都关注的是返回值，所以这种封装很适用。但在特殊场景下，即桩函数在关注返回值的同时也关注出参，这时就要用原生的API。

## 小结

本文主要介绍了Monkey框架的使用方法，基本上解决了序言中提到的那两个棘手的问题，同时针对Monkey框架的三个缺陷，分别提供了解决方案。

至此，我们已经知道：

1. 全局变量可通过GoStub框架打桩
2. 过程可通过Monkey框架打桩
3. 函数可通过Monkey框架打桩
4. 方法可通过Monkey框架打桩
5. interface可通过GoMock框架打桩

我们在测试实践中要举一反三，深度掌握GoConvey + GoStub + GoMock + Monkey框架组合使用的正确姿势，写出高质量的测试代码。

我们在产品代码中，尽量不要使用全局变量

