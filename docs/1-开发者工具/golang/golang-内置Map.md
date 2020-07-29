Map的使用频率很高，在使用时我们需要注意对Map的优化，比如字符串转换，以及GC扫描等。

---

[TOC]

## 预设Map大小

在声明Map时预申请一块内存，可以减少Map扩容和再Hash成本。

```go
package tests

import (
	"runtime/debug"
	"testing"
	"time"
)

func initMap(m map[int]int) {
	for i := 0; i < 10000; i++ {
		m[i] = i
	}
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m := make(map[int]int)
		b.StartTimer()

		initMap(m)
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkMap-4              1000           1261971 ns/op          687158 B/op        275 allocs/op
}

func BenchmarkCapMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		m := make(map[int]int, 10000)
		b.StartTimer()

		initMap(m)
	}
	//**********
	//goos: windows
	//goarch: amd64
	//pkg: tests
	//BenchmarkCapMap-4           2000            599823 ns/op            2709 B/op          9 allocs/op
}

```

## 释放Map空间

在释放Map时仅`delete(Map,Key)`是不够的，Map大小减为0时仍会保留内存空间以待后用。使用`Map=nil`或替换为一个新的Map对象可以释放原Map空间。

> 提示：如长期使用 map 对象（比如用作 cache 容器），偶尔换成 “新的” 或许会更好。还有，int key 要比 string key 更快。