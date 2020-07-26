package cmap

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

// -- Put -- //

func BenchmarkCmapPutAbsent(b *testing.B) {
	var number = 20
	var testCases = genNoRepetitiveTestingPairs(number)
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	b.ResetTimer()
	for _, tc := range testCases {
		key := tc.Key()
		element := tc.Element()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cm.Put(key, element)
			}
		})
	}
	//BenchmarkCmapPutAbsent
	//BenchmarkCmapPutAbsent/9a4bcc50
	//BenchmarkCmapPutAbsent/9a4bcc50-12         	 3670723	       318 ns/op
	//BenchmarkCmapPutAbsent/7d69ff59
	//BenchmarkCmapPutAbsent/7d69ff59-12         	 3746776	       326 ns/op
	//BenchmarkCmapPutAbsent/56a64179
	//BenchmarkCmapPutAbsent/56a64179-12         	 3773296	       313 ns/op
	//BenchmarkCmapPutAbsent/2d29d94a
	//BenchmarkCmapPutAbsent/2d29d94a-12         	 3754586	       313 ns/op
	//BenchmarkCmapPutAbsent/c495ab3d
	//BenchmarkCmapPutAbsent/c495ab3d-12         	 3760123	       315 ns/op
	//BenchmarkCmapPutAbsent/960d0e50
	//BenchmarkCmapPutAbsent/960d0e50-12         	 3777882	       317 ns/op
	//BenchmarkCmapPutAbsent/ca26a91f
	//BenchmarkCmapPutAbsent/ca26a91f-12         	 3780439	       313 ns/op
	//BenchmarkCmapPutAbsent/70b07144
	//BenchmarkCmapPutAbsent/70b07144-12         	 3757180	       316 ns/op
	//BenchmarkCmapPutAbsent/dfbbbe76
	//BenchmarkCmapPutAbsent/dfbbbe76-12         	 3778299	       320 ns/op
	//BenchmarkCmapPutAbsent/89ff0d30
	//BenchmarkCmapPutAbsent/89ff0d30-12         	 3739047	       315 ns/op
	//BenchmarkCmapPutAbsent/f21fcd76
	//BenchmarkCmapPutAbsent/f21fcd76-12         	 3769461	       317 ns/op
	//BenchmarkCmapPutAbsent/14ac9407
	//BenchmarkCmapPutAbsent/14ac9407-12         	 3722270	       316 ns/op
	//BenchmarkCmapPutAbsent/ed00b241
	//BenchmarkCmapPutAbsent/ed00b241-12         	 3778554	       314 ns/op
	//BenchmarkCmapPutAbsent/cf8e8e19
	//BenchmarkCmapPutAbsent/cf8e8e19-12         	 3778910	       314 ns/op
	//BenchmarkCmapPutAbsent/e2204e7d
	//BenchmarkCmapPutAbsent/e2204e7d-12         	 3787832	       315 ns/op
	//BenchmarkCmapPutAbsent/3f559e04
	//BenchmarkCmapPutAbsent/3f559e04-12         	 3748010	       316 ns/op
	//BenchmarkCmapPutAbsent/08dcb255
	//BenchmarkCmapPutAbsent/08dcb255-12         	 3757221	       316 ns/op
	//BenchmarkCmapPutAbsent/9ba20943
	//BenchmarkCmapPutAbsent/9ba20943-12         	 3776740	       318 ns/op
	//BenchmarkCmapPutAbsent/4a1c022b
	//BenchmarkCmapPutAbsent/4a1c022b-12         	 3744590	       314 ns/op
	//BenchmarkCmapPutAbsent/e51d9c08
	//BenchmarkCmapPutAbsent/e51d9c08-12         	 3764926	       317 ns/op
}

func BenchmarkCmapPutPresent(b *testing.B) {
	var number = 20
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	key := "invariable key"
	b.ResetTimer()
	for i := 0; i < number; i++ {
		element := strconv.Itoa(i)
		b.Run(key, func(b *testing.B) {
			for j := 0; j < b.N; j++ {
				cm.Put(key, element)
			}
		})
	}
	//BenchmarkCmapPutPresent
	//BenchmarkCmapPutPresent/invariable_key
	//BenchmarkCmapPutPresent/invariable_key-12  	 2965962	       374 ns/op
	//BenchmarkCmapPutPresent/invariable_key#01
	//BenchmarkCmapPutPresent/invariable_key#01-12         	 3198027	       372 ns/op
	//BenchmarkCmapPutPresent/invariable_key#02
	//BenchmarkCmapPutPresent/invariable_key#02-12         	 3208993	       369 ns/op
	//BenchmarkCmapPutPresent/invariable_key#03
	//BenchmarkCmapPutPresent/invariable_key#03-12         	 3191559	       369 ns/op
	//BenchmarkCmapPutPresent/invariable_key#04
	//BenchmarkCmapPutPresent/invariable_key#04-12         	 3181662	       369 ns/op
	//BenchmarkCmapPutPresent/invariable_key#05
	//BenchmarkCmapPutPresent/invariable_key#05-12         	 3185995	       376 ns/op
	//BenchmarkCmapPutPresent/invariable_key#06
	//BenchmarkCmapPutPresent/invariable_key#06-12         	 3213631	       378 ns/op
	//BenchmarkCmapPutPresent/invariable_key#07
	//BenchmarkCmapPutPresent/invariable_key#07-12         	 3197047	       382 ns/op
	//BenchmarkCmapPutPresent/invariable_key#08
	//BenchmarkCmapPutPresent/invariable_key#08-12         	 3081370	       389 ns/op
	//BenchmarkCmapPutPresent/invariable_key#09
	//BenchmarkCmapPutPresent/invariable_key#09-12         	 3118737	       380 ns/op
	//BenchmarkCmapPutPresent/invariable_key#10
	//BenchmarkCmapPutPresent/invariable_key#10-12         	 3156633	       374 ns/op
	//BenchmarkCmapPutPresent/invariable_key#11
	//BenchmarkCmapPutPresent/invariable_key#11-12         	 3161607	       375 ns/op
	//BenchmarkCmapPutPresent/invariable_key#12
	//BenchmarkCmapPutPresent/invariable_key#12-12         	 3180688	       379 ns/op
	//BenchmarkCmapPutPresent/invariable_key#13
	//BenchmarkCmapPutPresent/invariable_key#13-12         	 3124522	       393 ns/op
	//BenchmarkCmapPutPresent/invariable_key#14
	//BenchmarkCmapPutPresent/invariable_key#14-12         	 3055022	       395 ns/op
	//BenchmarkCmapPutPresent/invariable_key#15
	//BenchmarkCmapPutPresent/invariable_key#15-12         	 3125641	       393 ns/op
	//BenchmarkCmapPutPresent/invariable_key#16
	//BenchmarkCmapPutPresent/invariable_key#16-12         	 3090996	       386 ns/op
	//BenchmarkCmapPutPresent/invariable_key#17
	//BenchmarkCmapPutPresent/invariable_key#17-12         	 3080146	       394 ns/op
	//BenchmarkCmapPutPresent/invariable_key#18
	//BenchmarkCmapPutPresent/invariable_key#18-12         	 3111006	       387 ns/op
	//BenchmarkCmapPutPresent/invariable_key#19
	//BenchmarkCmapPutPresent/invariable_key#19-12         	 3080784	       386 ns/op
}

func BenchmarkMapPut(b *testing.B) {
	var number = 10
	var testCases = genNoRepetitiveTestingPairs(number)
	m := make(map[string]interface{})
	b.ResetTimer()
	for _, tc := range testCases {
		key := tc.Key()
		element := tc.Element()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m[key] = element
			}
		})
	}
	//BenchmarkMapPut
	//BenchmarkMapPut/76e73219
	//BenchmarkMapPut/76e73219-12                          	120761281	         9.95 ns/op
	//BenchmarkMapPut/d8577e2d
	//BenchmarkMapPut/d8577e2d-12                          	100000000	        11.1 ns/op
	//BenchmarkMapPut/c660d864
	//BenchmarkMapPut/c660d864-12                          	95452438	        12.2 ns/op
	//BenchmarkMapPut/c37fa218
	//BenchmarkMapPut/c37fa218-12                          	93326871	        12.9 ns/op
	//BenchmarkMapPut/92750d7e
	//BenchmarkMapPut/92750d7e-12                          	72524024	        14.1 ns/op
	//BenchmarkMapPut/7dbcaa3b
	//BenchmarkMapPut/7dbcaa3b-12                          	70548055	        17.0 ns/op
	//BenchmarkMapPut/6065284f
	//BenchmarkMapPut/6065284f-12                          	62770119	        16.8 ns/op
	//BenchmarkMapPut/3a0f9b20
	//BenchmarkMapPut/3a0f9b20-12                          	63150764	        18.8 ns/op
	//BenchmarkMapPut/c86ce65e
	//BenchmarkMapPut/c86ce65e-12                          	72556610	        14.1 ns/op
	//BenchmarkMapPut/e5bb110c
	//BenchmarkMapPut/e5bb110c-12                          	79057086	        14.2 ns/op
}

// -- Get -- //

func BenchmarkCmapGet(b *testing.B) {
	var number = 100000
	var testCases = genNoRepetitiveTestingPairs(number)
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	for _, p := range testCases {
		cm.Put(p.Key(), p.Element())
	}
	b.ResetTimer()
	for i := 0; i < 10; i++ {
		key := testCases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cm.Get(key)
			}
		})
	}
	//BenchmarkCmapGet
	//BenchmarkCmapGet/ea7d2424
	//BenchmarkCmapGet/ea7d2424-12                         	 6633852	       179 ns/op
	//BenchmarkCmapGet/e950e965
	//BenchmarkCmapGet/e950e965-12                         	 5881538	       202 ns/op
	//BenchmarkCmapGet/cf74b435
	//BenchmarkCmapGet/cf74b435-12                         	 6733189	       177 ns/op
	//BenchmarkCmapGet/b0148968
	//BenchmarkCmapGet/b0148968-12                         	 6693403	       187 ns/op
	//BenchmarkCmapGet/93ea334c
	//BenchmarkCmapGet/93ea334c-12                         	 6828248	       181 ns/op
	//BenchmarkCmapGet/3278c518
	//BenchmarkCmapGet/3278c518-12                         	 6551282	       180 ns/op
	//BenchmarkCmapGet/b3ea9c48
	//BenchmarkCmapGet/b3ea9c48-12                         	 6603882	       179 ns/op
	//BenchmarkCmapGet/50082204
	//BenchmarkCmapGet/50082204-12                         	 6564352	       184 ns/op
	//BenchmarkCmapGet/d22c2571
	//BenchmarkCmapGet/d22c2571-12                         	 6730406	       176 ns/op
	//BenchmarkCmapGet/7a2a2268
	//BenchmarkCmapGet/7a2a2268-12                         	 6181437	       180 ns/op
}

func BenchmarkMapGet(b *testing.B) {
	var number = 100000
	var testCases = genNoRepetitiveTestingPairs(number)
	m := make(map[string]interface{})
	for _, p := range testCases {
		m[p.Key()] = p.Element()
	}
	b.ResetTimer()
	for i := 0; i < 10; i++ {
		key := testCases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = m[key]
			}
		})
	}
	//BenchmarkMapGet
	//BenchmarkMapGet/0d9ae405
	//BenchmarkMapGet/0d9ae405-12                          	77593770	        13.6 ns/op
	//BenchmarkMapGet/4cde7340
	//BenchmarkMapGet/4cde7340-12                          	139293402	         8.51 ns/op
	//BenchmarkMapGet/cd996f7d
	//BenchmarkMapGet/cd996f7d-12                          	117144183	        10.2 ns/op
	//BenchmarkMapGet/49b7f17b
	//BenchmarkMapGet/49b7f17b-12                          	104900542	        11.5 ns/op
	//BenchmarkMapGet/58e10a1a
	//BenchmarkMapGet/58e10a1a-12                          	75679015	        13.8 ns/op
	//BenchmarkMapGet/8860ac74
	//BenchmarkMapGet/8860ac74-12                          	141961966	         9.60 ns/op
	//BenchmarkMapGet/d51e5d13
	//BenchmarkMapGet/d51e5d13-12                          	100000000	        10.2 ns/op
	//BenchmarkMapGet/5e706625
	//BenchmarkMapGet/5e706625-12                          	100000000	        10.0 ns/op
	//BenchmarkMapGet/83d21773
	//BenchmarkMapGet/83d21773-12                          	144613375	         8.05 ns/op
	//BenchmarkMapGet/6f640b73
	//BenchmarkMapGet/6f640b73-12                          	125002148	         9.62 ns/op
}

// -- Delete -- //

func BenchmarkCmapDelete(b *testing.B) {
	var number = 100000
	var testCases = genNoRepetitiveTestingPairs(number)
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	for _, p := range testCases {
		cm.Put(p.Key(), p.Element())
	}
	b.ResetTimer()
	for i := 0; i < 20; i++ {
		key := testCases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cm.Delete(key)
			}
		})
	}
	//BenchmarkCmapDelete
	//BenchmarkCmapDelete/03fe2e60
	//BenchmarkCmapDelete/03fe2e60-12                      	 5312269	       223 ns/op
	//BenchmarkCmapDelete/74099377
	//BenchmarkCmapDelete/74099377-12                      	 5355393	       220 ns/op
	//BenchmarkCmapDelete/63bc2f0d
	//BenchmarkCmapDelete/63bc2f0d-12                      	 5266490	       226 ns/op
	//BenchmarkCmapDelete/2e9e1023
	//BenchmarkCmapDelete/2e9e1023-12                      	 5277662	       223 ns/op
	//BenchmarkCmapDelete/6439664c
	//BenchmarkCmapDelete/6439664c-12                      	 5246576	       222 ns/op
	//BenchmarkCmapDelete/c5446562
	//BenchmarkCmapDelete/c5446562-12                      	 5411587	       225 ns/op
	//BenchmarkCmapDelete/cfe1e03f
	//BenchmarkCmapDelete/cfe1e03f-12                      	 4030590	       295 ns/op
	//BenchmarkCmapDelete/8930c37e
	//BenchmarkCmapDelete/8930c37e-12                      	 4087266	       298 ns/op
	//BenchmarkCmapDelete/42ccbf09
	//BenchmarkCmapDelete/42ccbf09-12                      	 5226162	       224 ns/op
	//BenchmarkCmapDelete/87a66056
	//BenchmarkCmapDelete/87a66056-12                      	 5139304	       256 ns/op
	//BenchmarkCmapDelete/1e31f552
	//BenchmarkCmapDelete/1e31f552-12                      	 4131214	       270 ns/op
	//BenchmarkCmapDelete/5a32ea70
	//BenchmarkCmapDelete/5a32ea70-12                      	 5273348	       227 ns/op
	//BenchmarkCmapDelete/0d5deb42
	//BenchmarkCmapDelete/0d5deb42-12                      	 4114441	       297 ns/op
	//BenchmarkCmapDelete/723ead05
	//BenchmarkCmapDelete/723ead05-12                      	 5592582	       214 ns/op
	//BenchmarkCmapDelete/9fdc7d51
	//BenchmarkCmapDelete/9fdc7d51-12                      	 5544966	       223 ns/op
	//BenchmarkCmapDelete/f8fc8b51
	//BenchmarkCmapDelete/f8fc8b51-12                      	 5325492	       225 ns/op
	//BenchmarkCmapDelete/583fea0f
	//BenchmarkCmapDelete/583fea0f-12                      	 5275300	       220 ns/op
	//BenchmarkCmapDelete/736cc66b
	//BenchmarkCmapDelete/736cc66b-12                      	 4251702	       290 ns/op
	//BenchmarkCmapDelete/24e59e0e
	//BenchmarkCmapDelete/24e59e0e-12                      	 4200740	       287 ns/op
	//BenchmarkCmapDelete/a3eeed78
	//BenchmarkCmapDelete/a3eeed78-12                      	 4076577	       286 ns/op
}

func BenchmarkMapDelete(b *testing.B) {
	var number = 100000
	var testCases = genNoRepetitiveTestingPairs(number)
	m := make(map[string]interface{})
	for _, p := range testCases {
		m[p.Key()] = p.Element()
	}
	b.ResetTimer()
	for i := 0; i < 20; i++ {
		key := testCases[rand.Intn(number)].Key()
		b.Run(key, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				delete(m, key)
			}
		})
	}
	//BenchmarkMapDelete
	//BenchmarkMapDelete/03b5e037
	//BenchmarkMapDelete/03b5e037-12                       	81133131	        13.7 ns/op
	//BenchmarkMapDelete/af958311
	//BenchmarkMapDelete/af958311-12                       	56303046	        20.1 ns/op
	//BenchmarkMapDelete/eada1452
	//BenchmarkMapDelete/eada1452-12                       	60101767	        18.6 ns/op
	//BenchmarkMapDelete/afe39810
	//BenchmarkMapDelete/afe39810-12                       	59188128	        18.7 ns/op
	//BenchmarkMapDelete/8228836e
	//BenchmarkMapDelete/8228836e-12                       	88888980	        13.8 ns/op
	//BenchmarkMapDelete/e30d6d52
	//BenchmarkMapDelete/e30d6d52-12                       	97752171	        12.8 ns/op
	//BenchmarkMapDelete/3d5d0916
	//BenchmarkMapDelete/3d5d0916-12                       	81254740	        12.7 ns/op
	//BenchmarkMapDelete/26038959
	//BenchmarkMapDelete/26038959-12                       	85685335	        12.8 ns/op
	//BenchmarkMapDelete/8c84c05d
	//BenchmarkMapDelete/8c84c05d-12                       	75768514	        13.6 ns/op
	//BenchmarkMapDelete/84ab8024
	//BenchmarkMapDelete/84ab8024-12                       	88027096	        12.8 ns/op
	//BenchmarkMapDelete/bef49257
	//BenchmarkMapDelete/bef49257-12                       	91708706	        12.8 ns/op
	//BenchmarkMapDelete/a95abc5e
	//BenchmarkMapDelete/a95abc5e-12                       	87654511	        12.9 ns/op
	//BenchmarkMapDelete/78684c3a
	//BenchmarkMapDelete/78684c3a-12                       	79324100	        12.9 ns/op
	//BenchmarkMapDelete/3c009101
	//BenchmarkMapDelete/3c009101-12                       	84596385	        12.6 ns/op
	//BenchmarkMapDelete/87188011
	//BenchmarkMapDelete/87188011-12                       	86870695	        12.1 ns/op
	//BenchmarkMapDelete/1c3de75f
	//BenchmarkMapDelete/1c3de75f-12                       	87865599	        12.3 ns/op
	//BenchmarkMapDelete/c6f73361
	//BenchmarkMapDelete/c6f73361-12                       	55307238	        18.8 ns/op
	//BenchmarkMapDelete/c819ce39
	//BenchmarkMapDelete/c819ce39-12                       	83579887	        12.6 ns/op
	//BenchmarkMapDelete/95a3202b
	//BenchmarkMapDelete/95a3202b-12                       	91030519	        12.4 ns/op
	//BenchmarkMapDelete/d6701500
	//BenchmarkMapDelete/d6701500-12                       	75591631	        14.5 ns/op
}

// -- Len -- //

func BenchmarkCmapLen(b *testing.B) {
	var number = 100000
	var testCases = genNoRepetitiveTestingPairs(number)
	concurrency := number / 4
	cm, _ := NewConcurrentMap(concurrency, nil)
	for _, p := range testCases {
		cm.Put(p.Key(), p.Element())
	}
	b.ResetTimer()
	for i := 0; i < 5; i++ {
		b.Run(fmt.Sprintf("Len%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				cm.Len()
			}
		})
	}
	//BenchmarkCmapLen
	//BenchmarkCmapLen/Len0
	//BenchmarkCmapLen/Len0-12                             	191855707	         6.24 ns/op
	//BenchmarkCmapLen/Len1
	//BenchmarkCmapLen/Len1-12                             	197390983	         6.19 ns/op
	//BenchmarkCmapLen/Len2
	//BenchmarkCmapLen/Len2-12                             	189480876	         6.22 ns/op
	//BenchmarkCmapLen/Len3
	//BenchmarkCmapLen/Len3-12                             	192469284	         6.26 ns/op
	//BenchmarkCmapLen/Len4
	//BenchmarkCmapLen/Len4-12                             	184335219	         6.25 ns/op
}

func BenchmarkMapLen(b *testing.B) {
	var number = 100000
	var testCases = genNoRepetitiveTestingPairs(number)
	m := make(map[string]interface{})
	for _, p := range testCases {
		m[p.Key()] = p.Element()
	}
	b.ResetTimer()
	for i := 0; i < 5; i++ {
		b.Run(fmt.Sprintf("Len%d", i), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = len(m)
			}
		})
	}
	//BenchmarkMapLen
	//BenchmarkMapLen/Len0
	//BenchmarkMapLen/Len0-12                              	1000000000	         0.246 ns/op
	//BenchmarkMapLen/Len1
	//BenchmarkMapLen/Len1-12                              	1000000000	         0.247 ns/op
	//BenchmarkMapLen/Len2
	//BenchmarkMapLen/Len2-12                              	1000000000	         0.246 ns/op
	//BenchmarkMapLen/Len3
	//BenchmarkMapLen/Len3-12                              	1000000000	         0.245 ns/op
	//BenchmarkMapLen/Len4
	//BenchmarkMapLen/Len4-12                              	1000000000	         0.245 ns/op
}
