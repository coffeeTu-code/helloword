# sync.mutex 源代码分析



[TOC]



针对 Golang 1.10.3 的 sync.Mutex 进行分析，代码位置：`sync/mutex.go`

![sync_mutex.jpeg](https://upload-images.jianshu.io/upload_images/14738618-5075e24d24364a2f.jpeg?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)


## 结构体定义


```go

type Mutex struct {
	state int32  // 指代mutex锁当前的状态
	sema  uint32 // 信号量，用于唤醒goroutine
}

```

Mutex 中的 state 用于指代锁当前的状态，如下所示

```

1111 1111 ...... 1111 1111
\_________29__________/|||
 存储等待 goroutine 数量 
 
[0]第0个bit 表示当前 mutex 是否已被某个goroutine所拥有，1=加锁
[1]第1个bit 表示当前 mutex 是否被唤醒，也就是有某个唤醒的goroutine要尝试获取锁
[2]第2个bit 表示 mutex 当前是否处于饥饿状态

```


## 常量定义


```go

const (
	mutexLocked = 1 << iota
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota
	starvationThresholdNs = 1e6
)

```

- mutexLocked 值为1，根据 mutex.state & mutexLocked 得到 mutex 的加锁状态，结果为1表示已加锁，0表示未加锁

- mutexWoken 值为2（二进制：10），根据 mutex.state & mutexWoken 得到 mutex 的唤醒状态，结果为1表示已唤醒，0表示未唤醒

- mutexStarving 值为4（二进制：100），根据 mutex.state & mutexStarving 得到 mutex 的饥饿状态，结果为1表示处于饥饿状态，0表示处于正常状态

- mutexWaiterShift 值为3，根据 mutex.state >> mutexWaiterShift 得到当前等待的 goroutine 数目

- starvationThresholdNs 值为1e6纳秒，也就是1毫秒，当等待队列中队首 goroutine 等待时间超过 starvationThresholdNs，mutex 进入饥饿模式


## 饥饿模式与正常模式


根据Mutex的注释，当前的Mutex有如下的性质。这些注释将极大的帮助我们理解Mutex的实现。

> Mutex 有两种工作模式：正常模式和饥饿模式
>
>
> **正常模式**
>
> 在正常模式中，所有等待锁的 goroutine 按照 FIFO 的顺序排队获取锁，但是一个被唤醒的等待者有时候并不能获取 mutex，它还需要和新到来的 goroutine 们竞争 mutex 的使用权。新到来的 goroutine 具有优势，它们已经在 CPU 上运行且它们数量很多，因此一个被唤醒的等待者有很大的概率获取不到锁。在这种情况下，这个被唤醒的 goroutine 会加入到等待队列的前面。如果一个等待的 goroutine 超过 1ms 没有获取锁，它就会将 mutex 切换到饥饿模式
>
> **饥饿模式**
>
> 在饥饿模式中，锁的所有权将从 unlock 的 gorutine 直接交给交给等待队列中的第一个。新来的 goroutine 将不会尝试去获得锁，即使锁看起来是 unlock 状态, 也不会去尝试自旋操作，而是放在等待队列的尾部。
>
>
> 如果一个等待的goroutine获取了锁，并且满足一以下其中的任何一个条件，它会将锁的状态转换为正常状态。
> 1) 它是等待队列中的最后一个;
> 2) 它等待的时间少于1ms。


## 函数


在分析源代码之前， 我们要从多线程(goroutine)的并发场景去理解为什么实现中有很多的分支。

当一个goroutine获取这个锁的时候， 有可能这个锁根本没有竞争者， 那么这个goroutine轻轻松松获取了这个锁。

而如果这个锁已经被别的goroutine拥有， 就需要考虑怎么处理当前的期望获取锁的goroutine。

同时， 当并发goroutine很多的时候，有可能会有多个竞争者， 而且还会有通过信号量唤醒的等待者。

以下代码已经去掉了与核心代码无关的 race 代码。

### Lock

Lock 方法申请对 mutex 加锁，Lock 执行的时候，分三种情况：

1. **无冲突** 通过 CAS 操作把当前状态设置为加锁状态  
2. **有冲突 开始自旋**，并等待锁释放，如果其他 goroutine 在这段时间内释放了该锁，直接获得该锁；如果没有释放，进入3  
3. **有冲突，且已经过了自旋阶段** 通过调用 semacquire 函数来让当前 goroutine 进入等待状态  

```go

func (m *Mutex) Lock() {
	// 如果 mutext 的 state 没有被锁，也没有等待/唤醒的 goroutine , 锁处于正常状态，那么获得锁，返回.
	// 比如锁第一次被 goroutine 请求时，就是这种状态。或者锁处于空闲的时候，也是这种状态。
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		return
	}

	var waitStartTime int64 // 当前 goroutine 开始等待的时间
	starving := false       // 当前 goroutine 是否已经处于饥饿状态
	awoke := false          // 当前 goroutine 是否被唤醒
	iter := 0               // 自旋迭代的次数
	old := m.state          // old 保存当前 mutex 的状态
	for {
		// 第一个条件是 state 已被锁，但是不是饥饿状态。如果是饥饿状态，自旋是没有用的，锁的拥有权直接交给了等待队列的第一个。
		// 第二个条件是还可以自旋，多核、压力不大并且在一定次数内可以自旋， 
		// 具体的条件可以参考 `sync_runtime_canSpin` 的实现（汇编实现，内部持续调用 PAUSE 指令，消耗 CPU 时间）。
		// 如果满足这两个条件，不断自旋来等待锁被释放、或者进入饥饿状态、或者不能再自旋。
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// 自旋的过程中如果发现 state 还没有设置 woken 标识，则设置它的 woken 标识， 并标记自己为被唤醒。
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				awoke = true
			}
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		// 到了这一步， state的状态可能是：
		// 1. 锁还没有被释放，锁处于正常状态
		// 2. 锁还没有被释放， 锁处于饥饿状态
		// 3. 锁已经被释放， 锁处于正常状态
		// 4. 锁已经被释放， 锁处于饥饿状态
		// 并且本gorutine的 awoke可能是true, 也可能是false (其它goutine已经设置了state的woken标识)

		// new 复制 state的当前状态， 用来设置新的状态。old 是锁当前的状态
		new := old
		// 如果 old state 状态不是饥饿状态, new state 设置锁， 尝试通过CAS获取锁。
		// 如果 old state 状态是饥饿状态, 则不设置 new state 的锁，因为饥饿状态下锁直接转给等待队列的第一个.
		if old&mutexStarving == 0 {
			new |= mutexLocked
		}
		// 当 mutex 处于加锁状态或饥饿状态的时候，新到来的 goroutine 进入等待队列
		if old&(mutexLocked|mutexStarving) != 0 {
			new += 1 << mutexWaiterShift
		}
		// 当前 goroutine 将 mutex 切换为饥饿状态，但如果当前 mutex 未加锁，则不需要切换
		// Unlock 操作希望饥饿模式存在等待者
		if starving && old&mutexLocked != 0 {
			new |= mutexStarving
		}
		// 如果本 goroutine 已经设置为唤醒状态, 需要清除 new state 的唤醒标记,
		// 因为本 goroutine 要么获得了锁，要么进入休眠，总之state的新状态不再是woken状态.
		if awoke {
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			new &^= mutexWoken
		}
		// 调用 CAS 更新 state 状态
		// 注意 new 的锁标记不一定是 true, 也可能只是标记一下锁的 state 是饥饿状态.
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			// 如果 old state 的状态是未被锁状态，并且锁不处于饥饿状态,
			// 那么当前 goroutine 已经获取了锁的拥有权，返回
			if old&(mutexLocked|mutexStarving) == 0 {
				break
			}
			// 设置/计算本 goroutine 的等待时间
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
			// 既然未能获取到锁， 那么就使用 sleep 原语阻塞本 goroutine
			// 如果是新来的 goroutine, queueLifo=false, 加入到等待队列的尾部，耐心等待
			// 如果是唤醒的 goroutine, queueLifo=true, 加入到等待队列的头部
			runtime_SemacquireMutex(&m.sema, queueLifo)
			// 如果当前 goroutine 等待时间超过 starvationThresholdNs，mutex 进入饥饿模式
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			old = m.state
			// 如果当前的 state 已经是饥饿状态
			// 那么锁应该处于 Unlock 状态，锁被直接交给了本 goroutine
			if old&mutexStarving != 0 {
				// 如果当前的 state 已被锁，或者已标记为唤醒， 或者等待的队列中不为空,
				// 那么 state 是一个非法状态
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				// 等待状态的 goroutine - 1
				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				// 如果不是饥饿模式了或者当前等待着只剩下一个，退出饥饿模式
				if !starving || old>>mutexWaiterShift == 1 {
					delta -= mutexStarving
				}
				// 更新状态
				// 因为已经获得了锁，退出、返回
				atomic.AddInt32(&m.state, delta)
				break
			}
			// 如果当前的锁是正常模式，本 goroutine 被唤醒，自旋次数清零，从 for 循环开始处重新开始
			awoke = true
			iter = 0
		} else {
			// 如果CAS不成功，重新获取锁的 state, 从 for 循环开始处重新开始
			old = m.state
		}
	}
}


```


### Unlock

Unlock方法释放所申请的锁


```go

func (m *Mutex) Unlock() {
	// 如果 state 不是处于锁的状态, 那么就是 Unlock 根本没有加锁的 mutex, panic
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if (new+mutexLocked)&mutexLocked == 0 {
		throw("sync: unlock of unlocked mutex")
	}

	// 释放锁，并通知其它等待者
	// 锁如果处于饥饿状态，直接交给等待队列的第一个, 唤醒它，让它去获取锁
	// mutex 正常模式
	if new&mutexStarving == 0 {
		old := new
		for {
			// 如果没有等待者，或者已经存在一个 goroutine 被唤醒或得到锁、或处于饥饿模式
			// 直接返回.
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				return
			}
			// 将等待的 goroutine-1，并设置 woken 标识
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			// 设置新的 state, 这里通过信号量会唤醒一个阻塞的 goroutine 去获取锁.
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				runtime_Semrelease(&m.sema, false)
				return
			}
			old = m.state
		}
	} else {
		// mutex 饥饿模式，直接将 mutex 拥有权移交给等待队列最前端的 goroutine
		// 注意此时 state 的 mutex 还没有加锁，唤醒的 goroutine 会设置它。
		// 在此期间，如果有新的 goroutine 来请求锁， 因为 mutex 处于饥饿状态， mutex 还是被认为处于锁状态，
		// 新来的 goroutine 不会把锁抢过去.
		runtime_Semrelease(&m.sema, true)
	}
}

```


# sync.RWMutex 源码分析

RWMutex 是读写互斥锁，锁可以由任意数量的读取器或单个写入器来保持

RWMutex 的零值是一个解锁的互斥锁

RWMutex 是抢占式的读写锁，写锁之后来的读锁是加不上的

代码位置：`sync/rwmutex.go`


## 结构体定义

```go

type RWMutex struct {
    w           Mutex  // 互斥锁
    writerSem   uint32 // 写锁信号量
    readerSem   uint32 // 读锁信号量
    readerCount int32  // 读锁计数器
    readerWait  int32  // 获取写锁时需要等待的读锁释放数量
}

```


## 常量定义

```

const rwmutexMaxReaders = 1 << 30   // 支持最多2^30个读锁

```


## 函数


以下是 sync.RWMutex 提供的4个方法

### Lock

提供写锁加锁操作

```

func (rw *RWMutex) Lock() {
	// 使用 Mutex 锁
	rw.w.Lock()
	// 将当前的 readerCount 置为负数，告诉 RUnLock 当前存在写锁等待
	r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
	if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
		// 等待读锁释放
		runtime_Semacquire(&rw.writerSem)
	}
}

```

### Unlock

提供写锁释放操作

```

func (rw *RWMutex) Unlock() {
	// 加上 Lock 的时候减去的 rwmutexMaxReaders
	r := atomic.AddInt32(&rw.readerCount, rwmutexMaxReaders)
	// 没执行Lock调用Unlock，抛出异常
	if r >= rwmutexMaxReaders {
		race.Enable()
		throw("sync: Unlock of unlocked RWMutex")
	}
	// 通知当前等待的读锁
	for i := 0; i < int(r); i++ {
		runtime_Semrelease(&rw.readerSem, false)
	}
	// 释放 Mutex 锁
	rw.w.Unlock()
}

```

### RLock

提供读锁操作

```

func (rw *RWMutex) RLock() {
	// 每次 goroutine 获取读锁时，readerCount+1
    // 如果写锁已经被获取，那么 readerCount 在 -rwmutexMaxReaders 与 0 之间，这时挂起获取读锁的 goroutine
    // 如果写锁没有被获取，那么 readerCount > 0，获取读锁, 不阻塞
    // 通过 readerCount 判断读锁与写锁互斥, 如果有写锁存在就挂起goroutine, 多个读锁可以并行
	if atomic.AddInt32(&rw.readerCount, 1) < 0 {
		// 将 goroutine 排到G队列的后面,挂起 goroutine
		runtime_Semacquire(&rw.readerSem)
	}
}

```

### RUnLock

对读锁进行解锁

```

func (rw *RWMutex) RUnlock() {
	// 写锁等待状态，检查当前是否可以进行获取
	if r := atomic.AddInt32(&rw.readerCount, -1); r < 0 {
		// r + 1 == 0表示直接执行RUnlock()
		// r + 1 == -rwmutexMaxReaders表示执行Lock()再执行RUnlock()
		// 两总情况均抛出异常
		if r+1 == 0 || r+1 == -rwmutexMaxReaders {
			race.Enable()
			throw("sync: RUnlock of unlocked RWMutex")
		}
		// 当读锁释放完毕后，通知写锁
		if atomic.AddInt32(&rw.readerWait, -1) == 0 {
			// The last reader unblocks the writer.
			runtime_Semrelease(&rw.writerSem, false)
		}
	}
}

```

# 总结

**sync.Mutex**

- 同一个时刻只有一个线程能够拿到锁

> 注意：
> 1. 不要重复锁定互斥锁
> 2. 不要忘记解锁互斥锁
> 3. 不要在多个函数之间直接传递互斥锁


**sync.RWMutex**

- 如果设置了一个写锁，那么其它读的线程以及写的线程都拿不到锁，这个时候，与互斥锁的功能相同
- 如果设置了一个读锁，那么其它写的线程是拿不到锁的，但是其它读的线程是可以拿到锁


> 读写互斥锁的实现比较有技巧性一些，需要几点
> - 读锁不能阻塞读锁，引入readerCount实现
> - 读锁需要阻塞写锁，直到所有读锁都释放，引入readerSem实现
> - 写锁需要阻塞读锁，直到所有写锁都释放，引入wirterSem实现
> - 写锁需要阻塞写锁，引入Metux实现

