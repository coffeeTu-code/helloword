# Go 并发

goroutine是通过Go的runtime管理的一个线程管理器。

goroutine是Go并行设计的核心。goroutine说到底其实就是协程，但是它比线程更小，十几个goroutine可能体现在底层就是五六个线程，Go语言内部帮你实现了这些goroutine之间的内存共享。执行goroutine只需极少的栈内存(大概是4~5KB)，当然会根据相应的数据伸缩。也正因为如此，可同时运行成千上万个并发任务。goroutine比thread更易用、更高效、更轻便。


## Go goroutine 控制

### channels

goroutine运行在相同的地址空间，因此访问共享内存必须做好同步。goroutine之间通过channel进行数据通信。

channel可以与Unix shell 中的双向管道做类比：可以通过它发送或者接收值。这些值只能是特定的类型：channel类型。定义一个channel时，也需要定义发送到channel的值的类型。注意，必须使用make 创建channel。

```$go
ci := make(chan int)
cs := make(chan string)
cf := make(chan interface{})
```

channel通过操作符<-来接收和发送数据

```$go
ch <- v    // 发送v到channel ch.
v := <-ch  // 从ch中接收数据，并赋值给v
```
> 记住应该在生产者的地方关闭channel，而不是消费的地方去关闭它，这样容易引起panic
> 
> 另外记住一点的就是channel不像文件之类的，不需要经常去关闭，只有当你确实没有任何发送数据了，或者你想显式的结束range循环之类的

### `Select`

Go里面通过select可以监听channel上的数据流动。

select默认是阻塞的，只有当监听的channel中有发送或接收可以进行时才会运行，当多个channel都准备好的时候，select是随机的选择一个执行的。

```$go
func main() {
	c := make(chan int)
	o := make(chan bool)
	go func() {
		for {
			select {
				case v := <- c:
					println(v)
				case <- time.After(5 * time.Second):
					println("timeout")
					o <- true
					break
			}
		}
	}()
	<- o
}
```

### `Goexit`

退出当前执行的goroutine，但是defer函数还会继续调用。

### `Gosched`

让出当前goroutine的执行权限，调度器安排其他等待的任务运行，并在下次某个时候从该位置恢复执行。

### `NumCPU`

返回 CPU 核数量。

### `NumGoroutine`

返回正在执行和排队的任务总数。

### `GOMAXPROCS`

用来设置可以并行计算的CPU核数的最大值，并返回之前的值。