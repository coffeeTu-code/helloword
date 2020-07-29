延迟调用（defer）确实是一种 “优雅” 机制。可简化代码，并确保即便发生 panic 依然会被执行。如将 panic/recover 比作 try/except，那么 defer 似乎可看做 finally。

如同异常保护被滥用一样，defer 被无限制使用的例子比比皆是。

---

```go
package tests

import (
	"sync"
	"testing"
)

var m sync.Mutex

func call() {
	m.Lock()
	m.Unlock()
}

func defercall() {
	m.Lock()
	defer m.Unlock()
}

func BenchmarkCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		call()
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkCall-4         50000000                29.1 ns/op             0 B/op          0 allocs/op
}

func BenchmarkDeferCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		defercall()
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkDeferCall-4    20000000               101 ns/op               0 B/op          0 allocs/op
}

```

只需稍稍了解 defer 实现机制，就不难理解会有这样的性能差异。

编译器通过 runtime.deferproc “注册” 延迟调用，除目标函数地址外，还会复制相关参数（包括 receiver）。在函数返回前，执行 runtime.deferreturn 提取相关信息执行延迟调用。这其中的代价自然不是普通函数调用一条 CALL 指令所能比拟的。

当多个 goroutine 执行该函数时，原本的并发设计，因为错误的 defer 调用变成 “串行”。

> 单个函数里过多的 defer 调用可尝试合并。最起码，在并发竞争激烈时，mutex.Unlock 不应该使用 defer，而应尽快执行，仅保护最短的代码片段。