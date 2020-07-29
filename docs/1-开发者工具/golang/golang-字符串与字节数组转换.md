字符串（string）作为一种不可变类型，在与字节数组（slice, [ ]byte）转换时需付出 “沉重” 代价，根本原因是对底层字节数组的复制。这种代价会在以万为单位的高并发压力下迅速放大，所以对它的优化常变成 “必须” 行为。

---

```go
package tests

import (
	"strings"
	"testing"
	"unsafe"
)

func str2bytes(s string) (b []byte) {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	b = *(*[]byte)(unsafe.Pointer(&h))
	return
}

func bytes2str(b []byte) (s string) {
	s = *(*string)(unsafe.Pointer(&b))
	return
}

var s = strings.Repeat("a", 1024)

func BenchmarkTest_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := []byte(s)
		_ = string(b)
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkTest_1-4        2000000               646 ns/op            2048 B/op          2 allocs/op
}

func BenchmarkTest_2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := str2bytes(s)
		_ = bytes2str(b)
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkTest_2-4       2000000000               0.99 ns/op            0 B/op          0 allocs/op
}



```