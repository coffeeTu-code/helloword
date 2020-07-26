package cmap

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

// 插入不存在的 key (粗糙的锁)
func BenchmarkSingleInsertAbsentBuiltInMap(b *testing.B) {
	myMap = &MyMap{
		m: make(map[string]interface{}, 32),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myMap.BuiltinMapStore(strconv.Itoa(i), "value")
	}
	//BenchmarkSingleInsertAbsentBuiltInMap
	//BenchmarkSingleInsertAbsentBuiltInMap-12      	 3238304	       352 ns/op
}

// 插入不存在的 key (分段锁)
func BenchmarkSingleInsertAbsent(b *testing.B) {
	m := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i), "value")
	}
	//BenchmarkSingleInsertAbsent
	//BenchmarkSingleInsertAbsent-12                	 3008502	       401 ns/op
}

// 插入不存在的 key (syncMap)
func BenchmarkSingleInsertAbsentSyncMap(b *testing.B) {
	syncMap := &sync.Map{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		syncMap.Store(strconv.Itoa(i), "value")
	}
	//BenchmarkSingleInsertAbsentSyncMap
	//BenchmarkSingleInsertAbsentSyncMap-12         	 1591418	       632 ns/op
}

// 插入存在 key (粗糙锁)
func BenchmarkSingleInsertPresentBuiltInMap(b *testing.B) {
	myMap = &MyMap{
		m: make(map[string]interface{}, 32),
	}
	myMap.BuiltinMapStore("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myMap.BuiltinMapStore("key", "value")
	}
	//BenchmarkSingleInsertPresentBuiltInMap
	//BenchmarkSingleInsertPresentBuiltInMap-12     	44342637	        26.8 ns/op
}

// 插入存在 key (分段锁)
func BenchmarkSingleInsertPresent(b *testing.B) {
	m := New()
	m.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Set("key", "value")
	}
	//BenchmarkSingleInsertPresent
	//BenchmarkSingleInsertPresent-12               	16114506	        74.9 ns/op
}

// 插入存在 key (syncMap)
func BenchmarkSingleInsertPresentSyncMap(b *testing.B) {
	syncMap := &sync.Map{}
	syncMap.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		syncMap.Store("key", "value")
	}
	//BenchmarkSingleInsertPresentSyncMap
	//BenchmarkSingleInsertPresentSyncMap-12        	14371869	        84.0 ns/op
}

// 读取存在 key (粗糙锁)
func BenchmarkSingleGetPresentBuiltInMap(b *testing.B) {
	myMap = &MyMap{
		m: make(map[string]interface{}, 32),
	}
	myMap.BuiltinMapStore("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myMap.BuiltinMapLookup("key")
	}
	//BenchmarkSingleGetPresentBuiltInMap
	//BenchmarkSingleGetPresentBuiltInMap-12        	37317054	        31.0 ns/op
}

// 读取存在 key (分段锁)
func BenchmarkSingleGetPresent(b *testing.B) {
	m := New()
	m.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get("key")
	}
	//BenchmarkSingleGetPresent
	//BenchmarkSingleGetPresent-12                  	19969479	        59.7 ns/op
}

// 读取存在 key (syncMap)
func BenchmarkSingleGetPresentSyncMap(b *testing.B) {
	syncMap := &sync.Map{}
	syncMap.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		syncMap.Load("key")
	}
	//BenchmarkSingleGetPresentSyncMap
	//BenchmarkSingleGetPresentSyncMap-12           	45729268	        26.0 ns/op
}

// 删除存在 key (粗糙锁)
func BenchmarkDeleteBuiltInMap(b *testing.B) {
	myMap = &MyMap{
		m: make(map[string]interface{}, 32),
	}
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for pb.Next() {
			// The loop body is executed b.N times total across all goroutines.
			k := r.Intn(100000000)
			myMap.BuiltinMapDelete(strconv.Itoa(k))
		}
	})
	//BenchmarkDeleteBuiltInMap
	//BenchmarkDeleteBuiltInMap-12                  	 9289084	       127 ns/op
}

// 删除存在 key (分段锁)
func BenchmarkDelete(b *testing.B) {
	m := New()
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for pb.Next() {
			// The loop body is executed b.N times total across all goroutines.
			k := r.Intn(100000000)
			m.Remove(strconv.Itoa(k))
		}
	})
	//BenchmarkDelete
	//BenchmarkDelete-12                            	 7641771	       157 ns/op
}

// 删除存在 key (syncMap)
func BenchmarkDeleteSyncMap(b *testing.B) {
	syncMap := &sync.Map{}
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for pb.Next() {
			// The loop body is executed b.N times total across all goroutines.
			k := r.Intn(100000000)
			syncMap.Delete(strconv.Itoa(k))
		}
	})
	//BenchmarkDeleteSyncMap
	//BenchmarkDeleteSyncMap-12                     	100000000	        11.8 ns/op
}

// 并发的插入不存在的 key-value (粗糙锁)
func BenchmarkMultiInsertDifferentBuiltInMap(b *testing.B) {
	myMap = &MyMap{
		m: make(map[string]interface{}, 32),
	}
	finished := make(chan struct{}, b.N)

	set := func(key, value string) {
		for i := 0; i < 10; i++ {
			myMap.BuiltinMapStore(key, value)
		}
		finished <- struct{}{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
	//BenchmarkMultiInsertDifferentBuiltInMap
	//BenchmarkMultiInsertDifferentBuiltInMap-12    	 1236994	       968 ns/op
}

// 并发的插入不存在的 key-value (分段锁)
func benchmarkMultiInsertDifferent(b *testing.B) {
	m := New()
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertDifferent_1_Shard(b *testing.B) {
	runWithShards(benchmarkMultiInsertDifferent, b, 1)
	//BenchmarkMultiInsertDifferent_1_Shard
	//BenchmarkMultiInsertDifferent_1_Shard-12      	  797124	      1458 ns/op
}
func BenchmarkMultiInsertDifferent_16_Shard(b *testing.B) {
	runWithShards(benchmarkMultiInsertDifferent, b, 16)
	//BenchmarkMultiInsertDifferent_16_Shard
	//BenchmarkMultiInsertDifferent_16_Shard-12     	  790224	      1464 ns/op
}
func BenchmarkMultiInsertDifferent_32_Shard(b *testing.B) {
	runWithShards(benchmarkMultiInsertDifferent, b, 32)
	//BenchmarkMultiInsertDifferent_32_Shard
	//BenchmarkMultiInsertDifferent_32_Shard-12     	  773660	      1453 ns/op
}
func BenchmarkMultiInsertDifferent_256_Shard(b *testing.B) {
	runWithShards(benchmarkMultiInsertDifferent, b, 256)
	//BenchmarkMultiInsertDifferent_256_Shard
	//BenchmarkMultiInsertDifferent_256_Shard-12    	  788391	      1456 ns/op
}

// 并发的插入不存在的 key-value (syncMap)
func BenchmarkMultiInsertDifferentSyncMap(b *testing.B) {
	syncMap := &sync.Map{}
	finished := make(chan struct{}, b.N)

	set := func(key, value string) {
		for i := 0; i < 10; i++ {
			syncMap.Store(key, value)
		}
		finished <- struct{}{}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
	//BenchmarkMultiInsertDifferentSyncMap
	//BenchmarkMultiInsertDifferentSyncMap-12       	  541359	      2376 ns/op
}

// 并发的插入相同的 key-value (粗糙锁)
func BenchmarkMultiInsertSameBuiltInMap(b *testing.B) {
	myMap = &MyMap{
		m: make(map[string]interface{}, 32),
	}
	finished := make(chan struct{}, b.N)

	set := func(key, value string) {
		for i := 0; i < 10; i++ {
			myMap.BuiltinMapStore(key, value)
		}
		finished <- struct{}{}
	}
	myMap.BuiltinMapStore("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
	//BenchmarkMultiInsertSameBuiltInMap
	//BenchmarkMultiInsertSameBuiltInMap-12         	 2248771	       525 ns/op
}

// 并发的插入相同的 key-value (分段锁)
func BenchmarkMultiInsertSame(b *testing.B) {
	m := New()
	finished := make(chan struct{}, b.N)
	_, set := GetSet(m, finished)
	m.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
	//BenchmarkMultiInsertSame
	//BenchmarkMultiInsertSame-12                   	 1200621	      1015 ns/op
}

// 并发的插入相同的 key-value (syncMap)
func BenchmarkMultiInsertSameSyncMap(b *testing.B) {
	syncMap := &sync.Map{}
	finished := make(chan struct{}, b.N)

	set := func(key, value string) {
		for i := 0; i < 10; i++ {
			syncMap.Store(key, value)
		}
		finished <- struct{}{}
	}
	syncMap.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
	//BenchmarkMultiInsertSameSyncMap
	//BenchmarkMultiInsertSameSyncMap-12            	  855370	      1385 ns/op
}

// 并发的 get (粗糙锁)
func BenchmarkMultiGetSameBuiltInMap(b *testing.B) {
	myMap = &MyMap{
		m: make(map[string]interface{}, 32),
	}
	finished := make(chan struct{}, b.N)
	get := func(key, value string) {
		for i := 0; i < 10; i++ {
			myMap.BuiltinMapLookup(key)
		}
		finished <- struct{}{}
	}
	myMap.BuiltinMapStore("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		get("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
	//BenchmarkMultiGetSameBuiltInMap
	//BenchmarkMultiGetSameBuiltInMap-12            	 3513026	       336 ns/op
}

// 并发的 get (分段锁)
func BenchmarkMultiGetSame(b *testing.B) {
	m := New()
	finished := make(chan struct{}, b.N)
	get, _ := GetSet(m, finished)
	m.Set("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		get("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
	//BenchmarkMultiGetSame
	//BenchmarkMultiGetSame-12                      	 1926654	       619 ns/op
}

// 并发的 get (syncMap)
func BenchmarkMultiGetSameSyncMap(b *testing.B) {
	syncMap := &sync.Map{}
	finished := make(chan struct{}, b.N)
	get := func(key, value string) {
		for i := 0; i < 10; i++ {
			syncMap.Load(key)
		}
		finished <- struct{}{}
	}
	syncMap.Store("key", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		get("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
	//BenchmarkMultiGetSameSyncMap
	//BenchmarkMultiGetSameSyncMap-12               	 3947632	       297 ns/op
}

// 并发的 get 和 set (粗糙锁)
func BenchmarkMultiGetSetDifferentBuiltInMap(b *testing.B) {
	myMap = &MyMap{
		m: make(map[string]interface{}, 32),
	}
	finished := make(chan struct{}, 2*b.N)
	get := func(key, value string) {
		for i := 0; i < 10; i++ {
			myMap.BuiltinMapLookup(key)
		}
		finished <- struct{}{}
	}
	set := func(key, value string) {
		for i := 0; i < 10; i++ {
			myMap.BuiltinMapStore(key, value)
		}
		finished <- struct{}{}
	}
	myMap.BuiltinMapStore("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i-1), "value")
		get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
	//BenchmarkMultiGetSetDifferentBuiltInMap
	//BenchmarkMultiGetSetDifferentBuiltInMap-12    	 1000000	      1386 ns/op
}

// 并发的 get 和 set（分段锁）
func benchmarkMultiGetSetDifferent(b *testing.B) {
	m := New()
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	m.Set("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i-1), "value")
		get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent_1_Shard(b *testing.B) {
	runWithShards(benchmarkMultiGetSetDifferent, b, 1)
	//BenchmarkMultiGetSetDifferent_1_Shard
	//BenchmarkMultiGetSetDifferent_1_Shard-12      	  510382	      2314 ns/op
}
func BenchmarkMultiGetSetDifferent_16_Shard(b *testing.B) {
	runWithShards(benchmarkMultiGetSetDifferent, b, 16)
	//BenchmarkMultiGetSetDifferent_16_Shard
	//BenchmarkMultiGetSetDifferent_16_Shard-12     	  506209	      2314 ns/op
}
func BenchmarkMultiGetSetDifferent_32_Shard(b *testing.B) {
	runWithShards(benchmarkMultiGetSetDifferent, b, 32)
	//BenchmarkMultiGetSetDifferent_32_Shard
	//BenchmarkMultiGetSetDifferent_32_Shard-12     	  464994	      2342 ns/op
}
func BenchmarkMultiGetSetDifferent_256_Shard(b *testing.B) {
	runWithShards(benchmarkMultiGetSetDifferent, b, 256)
	//BenchmarkMultiGetSetDifferent_256_Shard
	//BenchmarkMultiGetSetDifferent_256_Shard-12    	  504132	      2306 ns/op
}

// 并发的 get 和 set (syncMap)
func BenchmarkMultiGetSetDifferentSyncMap(b *testing.B) {
	syncMap := &sync.Map{}
	finished := make(chan struct{}, 2*b.N)
	get := func(key, value string) {
		for i := 0; i < 10; i++ {
			syncMap.Load(key)
		}
		finished <- struct{}{}
	}
	set := func(key, value string) {
		for i := 0; i < 10; i++ {
			syncMap.Store(key, value)
		}
		finished <- struct{}{}
	}
	syncMap.Store("-1", "value")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i-1), "value")
		get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
	//BenchmarkMultiGetSetDifferentSyncMap
	//BenchmarkMultiGetSetDifferentSyncMap-12       	  290384	      6152 ns/op
}

// get set 已经存在的一些 key (粗糙锁)
func BenchmarkMultiGetSetBlockBuiltInMap(b *testing.B) {
	myMap = &MyMap{
		m: make(map[string]interface{}, 32),
	}
	finished := make(chan struct{}, 2*b.N)
	get := func(key, value string) {
		for i := 0; i < 10; i++ {
			myMap.BuiltinMapLookup(key)
		}
		finished <- struct{}{}
	}
	set := func(key, value string) {
		for i := 0; i < 10; i++ {
			myMap.BuiltinMapStore(key, value)
		}
		finished <- struct{}{}
	}
	for i := 0; i < b.N; i++ {
		myMap.BuiltinMapStore(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i%100), "value")
		get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
	//BenchmarkMultiGetSetBlockBuiltInMap
	//BenchmarkMultiGetSetBlockBuiltInMap-12        	 1341573	       888 ns/op
}

// get set 已经存在的一些 key（分段锁）
func benchmarkMultiGetSetBlock(b *testing.B) {
	m := New()
	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(m, finished)
	for i := 0; i < b.N; i++ {
		m.Set(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i%100), "value")
		get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock_1_Shard(b *testing.B) {
	runWithShards(benchmarkMultiGetSetBlock, b, 1)
	//BenchmarkMultiGetSetBlock_1_Shard
	//BenchmarkMultiGetSetBlock_1_Shard-12          	  792862	      1546 ns/op
}
func BenchmarkMultiGetSetBlock_16_Shard(b *testing.B) {
	runWithShards(benchmarkMultiGetSetBlock, b, 16)
	//BenchmarkMultiGetSetBlock_16_Shard
	//BenchmarkMultiGetSetBlock_16_Shard-12         	  775627	      1531 ns/op
}
func BenchmarkMultiGetSetBlock_32_Shard(b *testing.B) {
	runWithShards(benchmarkMultiGetSetBlock, b, 32)
	//BenchmarkMultiGetSetBlock_32_Shard
	//BenchmarkMultiGetSetBlock_32_Shard-12         	  799513	      1490 ns/op
}
func BenchmarkMultiGetSetBlock_256_Shard(b *testing.B) {
	runWithShards(benchmarkMultiGetSetBlock, b, 256)
	//BenchmarkMultiGetSetBlock_256_Shard
	//BenchmarkMultiGetSetBlock_256_Shard-12        	  802521	      1460 ns/op
}

// get set 已经存在的一些 key (syncMap)
func BenchmarkMultiGetSetBlockSyncMap(b *testing.B) {
	syncMap := &sync.Map{}
	finished := make(chan struct{}, 2*b.N)
	get := func(key, value string) {
		for i := 0; i < 10; i++ {
			syncMap.Load(key)
		}
		finished <- struct{}{}
	}
	set := func(key, value string) {
		for i := 0; i < 10; i++ {
			syncMap.Store(key, value)
		}
		finished <- struct{}{}
	}
	for i := 0; i < b.N; i++ {
		syncMap.Store(strconv.Itoa(i%100), "value")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i%100), "value")
		get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
	//BenchmarkMultiGetSetBlockSyncMap
	//BenchmarkMultiGetSetBlockSyncMap-12           	  763969	      1546 ns/op
}

func GetSet(m ConcurrentMap, finished chan struct{}) (set func(key, value string), get func(key, value string)) {
	return func(key, value string) {
			for i := 0; i < 10; i++ {
				m.Get(key)
			}
			finished <- struct{}{}
		}, func(key, value string) {
			for i := 0; i < 10; i++ {
				m.Set(key, value)
			}
			finished <- struct{}{}
		}
}

func runWithShards(bench func(b *testing.B), b *testing.B, shardsCount int) {
	oldShardsCount := SHARD_COUNT
	SHARD_COUNT = shardsCount
	bench(b)
	SHARD_COUNT = oldShardsCount
}
