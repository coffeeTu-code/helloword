使用 slice 代替 array，企图避免数据拷贝，提升性能。这一做法在某些时候怕会适得其反，倒造成不必要的性能损失。
 
以下测试用例指出：  
array 非但拥有更好的性能，还避免了堆内存分配，也就是说减轻了 GC 压力。为什么会这样？

整个 array 函数完全在栈上完成，而 slice 函数则需执行 makeslice，继而在堆上分配内存，这就是问题所在。

---

```go
package tests

import "testing"

const capacity = 1024

func array() [capacity]int {
	var d [capacity]int
	for i := 0; i < len(d); i++ {
		d[i] = 1
	}
	return d
}

func slice() []int {
	d := make([]int, capacity)
	for i := 0; i < len(d); i++ {
		d[i] = 1
	}
	return d
}

func BenchmarkArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = array()
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkArray-4         1000000              2026 ns/op               0 B/op          0 allocs/op
}

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = slice()
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkSlice-4          500000              3355 ns/op            8192 B/op          1 allocs/op
}
```