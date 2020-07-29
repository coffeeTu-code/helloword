Go 使用 channel 实现 CSP 模型。处理双方仅关注通道和数据本身，无需理会对方身份和数量，以此实现结构性解耦。

> channel 多数时候适用于结构层面，而非单个区域的数据处理。原话中 “communicate” 本就表明一种 “message-passing”，而非 “lock-free”。所以，它并非用来取代 mutex，各自有不同的使用场景。

---

```go
package tests

import (
	"sync"
	"sync/atomic"
	"testing"
)

func chanCounter() chan int {
	c := make(chan int)

	go func() {
		for x := 1; ; x++ {
			c <- x
		}
	}()

	return c
}

func mutexCounter() func() int {
	var (
		m sync.Mutex
		x int
	)

	return func() (n int) {
		m.Lock()
		x++
		n = x
		m.Unlock()
		return
	}
}

func atomicCounter() func() int {
	var x int64

	return func() int {
		return int(atomic.AddInt64(&x, 1))
	}
}

func BenchmarkChanCounter(b *testing.B) {
	c := chanCounter()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = <-c
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkChanCounter-4           3000000               626 ns/op               0 B/op          0 allocs/op
}

func BenchmarkMutexCounter(b *testing.B) {
	f := mutexCounter()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = f()
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkMutexCounter-4         30000000                38.6 ns/op             0 B/op          0 allocs/op
}

func BenchmarkAtomicCounter(b *testing.B) {
	f := mutexCounter()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = f()
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkAtomicCounter-4        50000000                28.1 ns/op             0 B/op          0 allocs/op
}

```

`channel` 适用于结构层面解耦，`mutex` 则适合保护语句级别的数据安全，至于 `atomic`处理起来要复杂得多（比如 ABA 等问题），也未必就比 `mutex` 快很多。