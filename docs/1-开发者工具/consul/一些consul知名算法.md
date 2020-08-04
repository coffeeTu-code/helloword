# 平滑的基于权重的轮询算法


轮询算法是非常常用的一种调度/负载均衡的算法。依照百度百科上的解释:

> Round-Robin，轮询调度，通信中信道调度的一种策略，该调度策略使用户轮流使用共享资源，不会考虑瞬时信道条件。从相同数量无线资源（相同调度时间段）被分配给每条通信链路的角度讲，轮询调度可以被视为公平调度。然而，从提供相同服务质量给所有通信链路的角度而言，轮询调度是不公平的，此时，必须为带有较差信道条件的通信链路分配更多无线资源（更多时间）。此外，由于轮询调度在调度过程中不考虑瞬时信道条件，因此它将导致较低的整体系统性能，但与最大载干比调度相比，在各通信链路间具有更为均衡的服务质量。

更广泛的轮询调度应用在广度的服务调度上面，尤其在面向服务或者是面向微服务的架构中，比可以在很多知名的软件中看到它的身影，比如LVS、Nginx、Dubblo等。但是正如上面的百度百科中的介绍一样，轮询调度有一个很大的问题，那就是它认为所有的服务的性能都是一样的，每个服务器都被公平的调度，在服务器的性能有显著差别的环境中，性能比较差的服务器被调度了相同的次数，这不是我们所期望的。所以本文要介绍的是加权的轮询算法，轮询算法可以看成是加权的轮询算法的一个特例，在这种情况下，每个服务器的权重都是一样的。

本文介绍了Nginx和LVS的两种算法，比较了它们的优缺点，并提供了一个通用的 Go 语言实现的加权轮询算法库： weighted,可以用在负载均衡/调度/微服务网关等场合。

WRR(weighted round-robin) 也是周而复始地轮询分组服务资源，但不同的是WRR算法为每个服务资源分配一个权值，当轮询到某个服务的时候，将根据它所具有权值的大小决定其是否可以提供服务。由于WRR是基于轮询的，因此它只是在大于一个轮询周期的时间上才能显示是公平的。

本文介绍了Nginx和LVS基于权重的轮询算法，这两个算法都是通过巧妙的运算得到每次要提供的服务对象，但各自又有自己的特点。

## Nginx算法

Nginx基于权重的轮询算法的实现可以参考它的一次代码提交： Upstream: smooth weighted round-robin balancing。

它不但实现了基于权重的轮询算法，而且还实现了平滑的算法。所谓平滑，就是在一段时间内，不仅服务器被选择的次数的分布和它们的权重一致，而且调度算法还比较均匀的选择服务器，而不会集中一段时间之内只选择某一个权重比较高的服务器。如果使用随机算法选择或者普通的基于权重的轮询算法，就比较容易造成某个服务集中被调用压力过大。

举个例子，比如权重为{a:5, b:1, c:1)的一组服务器，Nginx的平滑的轮询算法选择的序列为{ a, a, b, a, c, a, a },这显然要比{ c, b, a, a, a, a, a }序列更平滑，更合理，不会造成对a服务器的集中访问。

算法如下：

> on each peer selection we increase current_weight of each eligible peer by its weight, select peer with greatest current_weight and reduce its current_weight by total number of weight points distributed
among peers.
```

func nextWeighted(servers []*Weighted) (best *Weighted) {
	total := 0
	for i := 0; i < len(servers); i++ {
		w := servers[i]
		if w == nil {
			continue
		}
		w.CurrentWeight += w.EffectiveWeight
		total += w.EffectiveWeight
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}
		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}
	if best == nil {
		return nil
	}
	best.CurrentWeight -= total
	return best
}
```

如果你使用weighted库，你可以通过下面的几行代码就可以使用这个算法进行调度：

```
func ExampleW_() {
	w := &W1{}
	w.Add("a", 5)
	w.Add("b", 2)
	w.Add("c", 3)
	for i := 0; i < 10; i++ {
		fmt.Printf("%s ", w.Next())
	}
	// Output: a c b a a c a b c a
}
```

## LVS算法

LVS使用的另外一种算法，它的算法的介绍可以参考它的网站的wiki。

算法用伪代码表示如下：

```
Supposing that there is a server set S = {S0, S1, …, Sn-1};
W(Si) indicates the weight of Si;
i indicates the server selected last time, and i is initialized with -1;
cw is the current weight in scheduling, and cw is initialized with zero; 
max(S) is the maximum weight of all the servers in S;
gcd(S) is the greatest common divisor of all server weights in S;


while (true) {
    i = (i + 1) mod n;
    if (i == 0) {
        cw = cw - gcd(S); 
        if (cw <= 0) {
            cw = max(S);
            if (cw == 0)
            return NULL;
        }
    } 
    if (W(Si) >= cw) 
        return Si;
}
```

可以看到它的代码逻辑比较简单，所以性能也很快，但是如果服务器的权重差别较多，就不会像Nginx那样比较平滑，可以在短时间内对权重很大的那台服务器压力过大。

使用weighted库和上面一样简单，只是把类型W换成W2即可：

```
func ExampleW_() {
	w := &W2{}
	w.Add("a", 5)
	w.Add("b", 2)
	w.Add("c", 3)
	for i := 0; i < 10; i++ {
		fmt.Printf("%s ", w.Next())
	}
	// Output: a a a c a b c a b c
}
```

## 性能比较

可以看到，上面两种方法的使用都非常的简单，只需生成一个相应的W对象，然后加入服务器和对应的权重即可，通过Next方法就可以获得下一个服务器。

如果服务器的权重差别很大，出于平滑的考虑，避免短时间内会对服务器造成冲击，你可以选择Nginx的算法，如果服务器的差别不是很大，可以考虑使用LVS的算法，因为测试可以看到它的性能要好于Nginx的算法:

```
BenchmarkW1_Next-4      20000000                50.1 ns/op             0 B/op          0 allocs/op
BenchmarkW2_Next-4      50000000                29.1 ns/op             0 B/op          0 allocs/op
```

实际上两者的性能都非常的快，十个服务器的每次调度也就是几十纳秒的级别，而且没有额外的对象分配，所以无论使用哪种算法，这个调度不应成为你整个系统的瓶颈。