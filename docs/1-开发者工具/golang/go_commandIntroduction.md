# Go 命令介绍

Go语言自带有一套完整的命令操作工具，你可以通过在命令行中执行go来查看它们：
![go-command.png](https://upload-images.jianshu.io/upload_images/14738618-58d93d211c9dbb15.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

[TOC]  

### `go build`

这个命令主要用于编译代码。在包的编译过程中，若有必要，会同时编译与之相关联的包。
> go build会忽略目录下以“_”或“.”开头的go文件。
> 如果你的源代码针对不同的操作系统需要不同的处理，那么你可以根据不同的操作系统后缀来命名文件。
```$go
array_linux.go 
array_darwin.go 
array_windows.go 
array_freebsd.go
```

+ -o 指定输出的文件名，可以带上路径，例如 go build -o a/b/c
+ -i 安装相应的包，编译+go install
+ -a 更新全部已经是最新的包的，但是对标准包不适用
+ -n 把需要执行的编译命令打印出来，但是不执行，这样就可以很容易的知道底层是如何运行的
+ -p n 指定可以并行可运行的编译数目，默认是CPU数目
+ -race 开启编译的时候自动检测数据竞争的情况，目前只支持64位的机器
+ -v 打印出来我们正在编译的包名
+ -work 打印出来编译时候的临时文件夹名称，并且如果已经存在的话就不要删除
+ -x 打印出来执行的命令，其实就是和-n的结果类似，只是这个会执行
+ -ccflags 'arg list' 传递参数给5c, 6c, 8c 调用
+ -compiler name 指定相应的编译器，gccgo还是gc
+ -gccgoflags 'arg list' 传递参数给gccgo编译连接调用
+ -gcflags 'arg list' 传递参数给5g, 6g, 8g 调用
+ -installsuffix suffix 为了和默认的安装包区别开来，采用这个前缀来重新安装那些依赖的包，-race的时候默认已经是-installsuffix race,大家可以通过-n命令来验证
+ -ldflags 'flag list' 传递参数给5l, 6l, 8l 调用
+ -tags 'tag list' 设置在编译的时候可以适配的那些tag，详细的tag限制参考里面的 Build Constraints

### `go clean`

这个命令是用来移除当前源码包和关联源码包里面编译生成的文件。这些文件包括
```$go
_obj/            旧的object目录，由Makefiles遗留
_test/           旧的test目录，由Makefiles遗留
_testmain.go     旧的gotest文件，由Makefiles遗留
test.out         旧的test记录，由Makefiles遗留
build.out        旧的test记录，由Makefiles遗留
*.[568ao]        object文件，由Makefiles遗留

DIR(.exe)        由go build产生
DIR.test(.exe)   由go test -c产生
MAINFILE(.exe)   由go build MAINFILE.go产生
*.so             由 SWIG 产生
```

+ -i 清除关联的安装的包和可运行文件，也就是通过go install安装的文件
+ -n 把需要执行的清除命令打印出来，但是不执行，这样就可以很容易的知道底层是如何运行的
+ -r 循环的清除在import中引入的包
+ -x 打印出来执行的详细命令，其实就是-n打印的执行版本

### `go fmt`

使用go fmt命令，其实是调用了gofmt，而且需要参数-w，否则格式化结果不会写入文件。`gofmt -w -l src`，可以格式化整个项目。所以go fmt是gofmt的上层一个包装的命令。

+ -l 显示那些需要格式化的文件
+ -w 把改写后的内容直接写入到文件中，而不是作为结果打印到标准输出。
+ -r 添加形如“a[b:len(a)] -> a[b:]”的重写规则，方便我们做批量替换
+ -s 简化文件中的代码
+ -d 显示格式化前后的diff而不是写入文件，默认是false
+ -e 打印所有的语法错误到标准输出。如果不使用此标记，则只会打印不同行的前10个错误。
+ -cpuprofile 支持调试模式，写入相应的cpufile到指定的文件

### `go get`

这个命令是用来动态获取远程代码包的，目前支持的有BitBucket、GitHub、Google Code和Launchpad。这个命令在内部实际上分成了两步操作：第一步是下载源码包，第二步是执行`go install`。

+ -d 只下载不安装
+ -f 只有在你包含了-u参数的时候才有效，不让-u去验证import中的每一个都已经获取了，这对于本地fork的包特别有用
+ -fix 在获取源码之后先运行fix，然后再去做其他的事情
+ -t 同时也下载需要为运行测试所需要的包
+ -u 强制使用网络去更新包和它的依赖包
+ -v 显示执行的命令

### `go test`

执行这个命令，会自动读取源码目录下面名为*_test.go的文件，生成并运行测试用的可执行文件。

+ -bench regexp 执行相应的benchmarks，例如 -bench=.
+ -cover 开启测试覆盖率
+ -run regexp 只运行regexp匹配的函数，例如 -run=Array 那么就执行包含有Array开头的函数
+ -v 显示测试的详细命令

### `go tool`

+ `go tool fix` . 用来修复以前老版本的代码到新版本，例如go1之前老版本的代码转化到go1,例如API的变化
+ `go tool vet directory|files` 用来分析当前目录的代码是否都是正确的代码,例如是不是调用fmt.Printf里面的参数不正确，例如函数里面提前return了然后出现了无用代码之类的。

### `go generate`

go generate和go build是完全不一样的命令，通过分析源码中特殊的注释，然后执行相应的命令。这些命令都是很明确的，没有任何的依赖在里面。

### `godoc`

在命令行执行 godoc -http=:端口号 比如`godoc -http=:8080`。然后在浏览器中打开`127.0.0.1:8080`，你将会看到一个golang.org的本地copy版本，通过它你可以查询pkg文档等其它内容。

### `go version` 

查看go当前的版本

### `go env`

查看当前go的环境变量

### `go list`
 
列出当前全部安装的package

### `go run`
 
 编译并运行Go程序