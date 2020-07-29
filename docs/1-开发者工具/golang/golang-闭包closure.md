闭包（closure）也是很常见的编码模式，因它隐式携带上下文环境变量，因此可让算法代码变得更加简洁。

但任何 “便利” 和 “优雅” 的背后，往往都是更复杂的实现机制，无非是语法糖或编译器隐藏了相关细节。最终，这些都会变成额外成本在运行期由 CPU、runtime 负担。甚至因不合理使用，造成性能问题。


---

```go
package tests

import "testing"

func test(x int) int {
	return x * 2
}

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = test(i)
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkNormal-4       2000000000               0.85 ns/op            0 B/op          0 allocs/op
}

func BenchmarkClosure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func() int {
			return i * 2
		}
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkClosure-4      1000000000               2.45 ns/op            0 B/op          0 allocs/op
}

func BenchmarkAnonymous(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = func(x int) int {
			return i * 2
		}(i)
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkAnonymous-4    1000000000               2.54 ns/op            0 B/op          0 allocs/op
}

```

首先，闭包引用原环境变量，导致 y 逃逸到堆上，这必然增加了 GC 扫描和回收对象的数量。

接下来，同样是因为闭包引用原对象，造成数据竞争（data race）。

可见，闭包未必总能将事情 “简单化”。