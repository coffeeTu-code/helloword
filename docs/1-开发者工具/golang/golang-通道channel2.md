作为内置类型，通道（channel）从运行时得到很多支持，其自身设计也算得上精巧。但不管怎么说，它本质上依旧是一种队列，当多个 goroutine 并发操作时，免不了要使用锁。某些时候，这种竞争机制，会导致性能问题。

---

下面是一个简单利用 channel 收发数据的示例，为便于 “准确” 测量收发操作性能，我们将 make channel 操作放到外部，尽可能避免额外消耗。

在研究 go runtime 源码实现过程中，会看到大量利用 “批操作” 来提升性能的样例。在此，我们可借鉴一下，看看效果对比。

从测试结果看，性能提升很高，可见批操作是一种有效方案。

```go
package test

import "testing"

const (
	max     = 500000
	bufsize = 100
)

func channel_test(data chan int, done chan struct{}) int {
	count := 0

	//接收，统计。
	go func() {
		for x := range data {
			count += x
		}

		close(done)
	}()

	//发送
	for i := 0; i < max; i++ {
		data <- i
	}

	close(data)

	<-done
	return count
}

const block = 500

func testBlock(data chan [block]int, done chan struct{}) int {
	count := 0

	//接收，统计。
	go func() {
		for a := range data {
			for _, x := range a {
				count += x
			}
		}

		close(done)
	}()

	//发送
	for i := 0; i < max; i += block {
		//按块打包。
		var b [block]int
		for n := 0; n < block; n++ {
			b[n] = i + n
			if i+n == max-1 {
				break
			}
		}
		data <- b
	}

	close(data)

	<-done
	return count
}

func BenchmarkTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make(chan int, bufsize)
		done := make(chan struct{})
		b.StartTimer()

		_ = channel_test(data, done)
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkTest-4               20          81382200 ns/op             102 B/op          1 allocs/op
}

func BenchmarkBlockTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make(chan [block]int, bufsize)
		done := make(chan struct{})
		b.StartTimer()

		_ = testBlock(data, done)
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkBlockTest-4         500           3019946 ns/op              36 B/op          1 allocs/op
}

```