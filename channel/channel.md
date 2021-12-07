# Go Channel

* [数据结构](#数据结构)
  + [hchan](#hchan)
  + [waitq](#waitq)
  + [sudog](#sudog)
  + [图解](#图解)
* [源码分析](#源码分析)
  + [makechan](#makechan)
    - [流程图](#makechan流程图)
  + [chansend](#chansend)
    - [流程图](#chansend流程图)
    - [核心方法---send](#核心方法---send)
    - [核心方法---sendDirect](#核心方法---senddirect)
  + [chanrecv](#chanrecv)
    - [流程图](#chanrecv流程图)
    - [核心方法---recv](#核心方法---recv)
    - [核心方法---recvDirect](#核心方法---recvdirect)
  + [closechan](#closechan)
    - [流程图](#closechan流程图)
  + [非阻塞发送和非阻塞接受](#非阻塞发送和非阻塞接受)
    - [非阻塞发送---selectnbsend](#非阻塞发送---selectnbsend)
    - [非阻塞发送---selectngrecv](#非阻塞发送---selectngrecv)
* [图解循环队列](#图解循环队列)

## 数据结构

### hchan

```go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}
```

`qcount`chan中元素数量

`datasiz`底层循环数组的长度

`buf`指向底层循环数组的指针，只针对有缓冲的chan

`elemsize`chan中元素大小

`closed`chan是否关闭的标志

`elemtype`chan中元素类型

`sendx`已发送元素在循环数组中的索引

`recvx`已接收元素在循环数组中的索引

`sendq`等待发送的goroutine队列(双向链表)

`recvq`等待接受的goroutine队列(双向链表)

`lock`锁

### waitq

```go
type waitq struct {
  first *sudog
  last  *sudog
}
```

### sudog

`sudog`表示一个在等待列表中的 Goroutine，存储着这一次阻塞的信息以及指向前后的`sudog`指针

```go
type sudog struct {
	// The following fields are protected by the hchan.lock of the
	// channel this sudog is blocking on. shrinkstack depends on
	// this for sudogs involved in channel ops.

	g *g

	next *sudog
	prev *sudog
	elem unsafe.Pointer // data element (may point to stack)

	// The following fields are never accessed concurrently.
	// For channels, waitlink is only accessed by g.
	// For semaphores, all fields (including the ones above)
	// are only accessed when holding a semaRoot lock.

	acquiretime int64
	releasetime int64
	ticket      uint32

	// isSelect indicates g is participating in a select, so
	// g.selectDone must be CAS'd to win the wake-up race.
	isSelect bool

	// success indicates whether communication over channel c
	// succeeded. It is true if the goroutine was awoken because a
	// value was delivered over channel c, and false if awoken
	// because c was closed.
	success bool

	parent   *sudog // semaRoot binary tree
	waitlink *sudog // g.waiting list or semaRoot
	waittail *sudog // semaRoot
	c        *hchan // channel
}
```

### 图解

![](https://tva1.sinaimg.cn/large/008i3skNly1gwo5pn4gkjj31a20u0af4.jpg)

## 源码分析

### makechan

#### makechan流程图

![](https://tva1.sinaimg.cn/large/008i3skNly1gwo4f9rn6cj31da0m8ace.jpg)

```go
elem := t.elem
```

`elem`是`runtime._type`类型，由`chan「type」`获取

```go
mem, overflow := math.MulUintptr(elem.size, uintptr(size))
if overflow || mem > maxAlloc-hchanSize || size < 0 {
	panic(plainError("makechan: size out of range"))
}
```

`「mem」`是`「ring buffer」`需要的大小，计算`elem.size * size`

```go
var c *hchan
switch {
case mem == 0:
	// Queue or element size is zero.
	c = (*hchan)(mallocgc(hchanSize, nil, true))
	// Race detector uses this location for synchronization.
	c.buf = c.raceaddr()
case elem.ptrdata == 0:
	// Elements do not contain pointers.
	// Allocate hchan and buf in one call.
	c = (*hchan)(mallocgc(hchanSize+mem, nil, true))
	c.buf = add(unsafe.Pointer(c), hchanSize)
default:
	// Elements contain pointers.
	c = new(hchan)
	c.buf = mallocgc(mem, elem, true)
}
```

`mem == 0`说明无`buffer`，直接分配`chan`的内存空间

`elem.ptrdata == 0`说明`type`是非指针且不包含指针的类型，分配`chan`和`buffer`连续内存空间

`default`分别分配`chan`和`buffer`的内存空间

### chansend

#### chansend流程图

![](https://tva1.sinaimg.cn/large/008i3skNly1gx46cgz1l1j317f0u0juc.jpg)

- `c`是`channel`的指针
- `ep`是待发送数据的内存地址
- `block`判断`channel`是否阻塞

```go
if !block && c.closed == 0 && full(c) {
	return false
}

func full(c *hchan) bool {
	// c.dataqsiz is immutable (never written after the channel is created)
	// so it is safe to read at any time during channel operation.
	if c.dataqsiz == 0 {
		// Assumes that a pointer read is relaxed-atomic.
		return c.recvq.first == nil
	}
	// Assumes that a uint read is relaxed-atomic.
	return c.qcount == c.dataqsiz
}
```

上面这段代码可以转化为：

```go
if !block && c.closed == 0 &&
	(c.dataqsiz == 0 && c.recvq.first == nil ||
  c.dataqsiz > 0 && c.qcount == c.dataqsiz) {
  	return false
}
```

**非阻塞性且未关闭channel**满足以下其中一个条件：

- 不存在`buffer`，同时`recvq`为空
- 存在`buffer`，并且`buffer`已满

```go
if sg := c.recvq.dequeue(); sg != nil {
	// Found a waiting receiver. We pass the value we want to send
	// directly to the receiver, bypassing the channel buffer (if any).
	send(c, sg, ep, func() { unlock(&c.lock) }, 3)
	return true
}
```

`recvq`不为空，存在`receiver`说明此时`buffer`已满，队首出队，并将`ep`的值拷贝到`receiver`

```go
if c.qcount < c.dataqsiz {
	// Space is available in the channel buffer. Enqueue the element to send.
	qp := chanbuf(c, c.sendx)
	typedmemmove(c.elemtype, qp, ep)
	c.sendx++
	if c.sendx == c.dataqsiz {
		c.sendx = 0
	}
	c.qcount++
	unlock(&c.lock)
	return true
}
```

`recvq`为空，`buffer`未满，通过`sendx`获取待发送的内存地址`qp`，将`ep`的值拷贝到`qp`，发送索引`sendx`加1，如果等于缓冲区容量，`sendx`归0，缓冲区元素加1

```go
if !block {
	unlock(&c.lock)
	return false
}
```

如果是非阻塞类型，无法进入`sendq`，直接返回

```go
// Block on the channel. Some receiver will complete our operation for us.
gp := getg()
mysg := acquireSudog()
mysg.releasetime = 0
if t0 != 0 {
	mysg.releasetime = -1
}
// No stack splits between assigning elem and enqueuing mysg
// on gp.waiting where copystack can find it.
mysg.elem = ep
mysg.waitlink = nil
mysg.g = gp
mysg.isSelect = false
mysg.c = c
gp.waiting = mysg
gp.param = nil
c.sendq.enqueue(mysg)
// Signal to anyone trying to shrink our stack that we're about
// to park on a channel. The window between when this G's status
// changes and when we set gp.activeStackChans is not safe for
// stack shrinking.
atomic.Store8(&gp.parkingOnChan, 1)
gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanSend, traceEvGoBlockSend, 2)
```

`getg()`获取当前协程指针`gp`，获取一个`sudog`，并且绑定`gp`和`channel`，再将`mysg`入队

通过`gopark`将当前协程挂起

```go
if mysg != gp.waiting {
	throw("G waiting list is corrupted")
}
gp.waiting = nil
gp.activeStackChans = false
closed := !mysg.success
gp.param = nil
if mysg.releasetime > 0 {
	blockevent(mysg.releasetime-t0, 2)
}
mysg.c = nil
releaseSudog(mysg)
if closed {
	if c.closed == 0 {
		throw("chansend: spurious wakeup")
	}
	panic(plainError("send on closed channel"))
}
```

协程被唤醒之后，释放资源

#### 核心方法---send

```go
func send(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
	if sg.elem != nil {
		sendDirect(c.elemtype, sg, ep)
		sg.elem = nil
	}
	gp := sg.g
	unlockf()
	gp.param = unsafe.Pointer(sg)
	sg.success = true
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
	goready(gp, skip+1)
}
```

`sg.elem`指向接收到的值存放的位置，`e.g. val <- ch，指向的就是&val`

`sendDirect`直接拷贝内存，从`sender -> receiver`

获取`sudog`绑定的协程`gp`

通过`goready`唤醒接受的`goroutine`

#### 核心方法---sendDirect

```go
func sendDirect(t *_type, sg *sudog, src unsafe.Pointer) {
	// src is on our stack, dst is a slot on another stack.

	// Once we read sg.elem out of sg, it will no longer
	// be updated if the destination's stack gets copied (shrunk).
	// So make sure that no preemption points can happen between read & use.
	dst := sg.elem
	typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.size)
	// No need for cgo write barrier checks because dst is always
	// Go memory.
	memmove(dst, src, t.size)
}
```

- `dst`目标`goroutine`的栈地址

- `src`当前`goroutine`的栈地址

`GC`假设对栈的写操作只能发生在`goroutine`正在运行中并且由当前`goroutine`来写，而`sendDirect`违反了这个假设，试图通过`memmove`直接修改`receiver`的栈，如果在获取了目标`goroutine`的栈地址后，发生了**栈收缩**，必然会破坏内存，所以需要增加一个写屏障`Write Barrier`

_PS: 收缩栈是在mgcmark.go中触发的，主要是在scanstack和markrootFreeGStacks函数中，也就是垃圾回收的时候会根据情况收缩栈_

### chanrecv

#### chanrecv流程图

![](https://tva1.sinaimg.cn/large/008i3skNly1gx46amomi3j317q0tmn0d.jpg)

- `c`是`channel`的指针
- `ep`是待接收数据的内存地址
- `block`判断`channel`是否阻塞

```go
if !block && empty(c) {
	if atomic.Load(&c.closed) == 0 {
		return
	}
	
  if empty(c) {
		if ep != nil {
			typedmemclr(c.elemtype, ep)
		}
		return true, false
	}
}

func empty(c *hchan) bool {
	// c.dataqsiz is immutable.
	if c.dataqsiz == 0 {
		return atomic.Loadp(unsafe.Pointer(&c.sendq.first)) == nil
	}
	return atomic.Loaduint(&c.qcount) == 0
}
```

上面这段代码同样可以转化为：

```go
if !block && (
  c.dataqsiz == 0 && atomic.Loadp(unsafe.Pointer(&c.sendq.first)) == nil || 
  c.dataqsiz > 0 && atomic.Loaduint(&c.qcount)) == 0) {
  if atomic.Load(&c.closed) == 0 {
    return false, false
  }
  
  if ep != nil {
    typedmemclr(c.elemtype, ep)
  }
  return true, false
}
```

**非阻塞性**满足以下其中一个条件：

- 不存在`buffer`，同时`sendq`为空
- 存在`buffer`，并且`buffer`为空

再判断`channel`是否关闭

- 未关 - 返回false, false
- 关闭 - 返回true, false；同时如果`ep`不为空，需要释放`ep`的内存

`chanrecv`有两个返回值：

`selected`表明在`select`中是否进入当前`channel`所在的`case`

`received`表明是否接收成功

```go
if c.closed != 0 && c.qcount == 0 {
	if ep != nil {
		typedmemclr(c.elemtype, ep)
	}
	return true, false
}
```

如果`channe`关闭，且`buffer`为空的话(说明没有待接收的数据)，返回true，false；同时如果`ep`不为空，需要释放`ep`的内存

```go
if sg := c.sendq.dequeue(); sg != nil {
	// Found a waiting sender. If buffer is size 0, receive value
	// directly from sender. Otherwise, receive from head of queue
	// and add sender's value to the tail of the queue (both map to
	// the same buffer slot because the queue is full).
	recv(c, sg, ep, func() { unlock(&c.lock) }, 3)
	return true, true
}
```

`sendq`不为空，存在`sender`说明此时`buffer`已满，队首出队，并将`ep`的值拷贝到`receiver`

```go
if c.qcount > 0 {
	// Receive directly from queue
	qp := chanbuf(c, c.recvx)
	if ep != nil {
		typedmemmove(c.elemtype, ep, qp)
	}
	typedmemclr(c.elemtype, qp)
	c.recvx++
	if c.recvx == c.dataqsiz {
		c.recvx = 0
	}
	c.qcount--
	unlock(&c.lock)
	return true, true
}
```

`sendq`为空，`buffer`未满，通过`recvx`获取待接收的内存地址`qp`，将`ep`的值拷贝到`qp`，接收索引`recv`加1，如果等于缓冲区容量，`recvx`归0，缓冲区元素减1

```go
if !block {
	unlock(&c.lock)
	return false, false
}
```

如果是非阻塞类型，无法进入`recvq`，直接返回

```go
// no sender available: block on this channel.
gp := getg()
mysg := acquireSudog()
mysg.releasetime = 0
if t0 != 0 {
	mysg.releasetime = -1
}
// No stack splits between assigning elem and enqueuing mysg
// on gp.waiting where copystack can find it.
mysg.elem = ep
mysg.waitlink = nil
gp.waiting = mysg
mysg.g = gp
mysg.isSelect = false
mysg.c = c
gp.param = nil
c.recvq.enqueue(mysg)
// Signal to anyone trying to shrink our stack that we're about
// to park on a channel. The window between when this G's status
// changes and when we set gp.activeStackChans is not safe for
// stack shrinking.
atomic.Store8(&gp.parkingOnChan, 1)
gopark(chanparkcommit, unsafe.Pointer(&c.lock), waitReasonChanReceive, traceEvGoBlockRecv, 2)
```

`getg()`获取当前协程指针`gp`，获取一个`sudog`，并且绑定`gp`和`channel`，再将`mysg`入队

通过`gopark`将当前协程挂起

```go
// someone woke us up
if mysg != gp.waiting {
	throw("G waiting list is corrupted")
}
gp.waiting = nil
gp.activeStackChans = false
if mysg.releasetime > 0 {
	blockevent(mysg.releasetime-t0, 2)
}
success := mysg.success
gp.param = nil
mysg.c = nil
releaseSudog(mysg)
return true, success
```

协程被唤醒之后，释放资源

#### 核心方法---recv

```go
func recv(c *hchan, sg *sudog, ep unsafe.Pointer, unlockf func(), skip int) {
	if c.dataqsiz == 0 {
		if ep != nil {
			// copy data from sender
			recvDirect(c.elemtype, sg, ep)
		}
	} else {
		// Queue is full. Take the item at the
		// head of the queue. Make the sender enqueue
		// its item at the tail of the queue. Since the
		// queue is full, those are both the same slot.
		qp := chanbuf(c, c.recvx)
		// copy data from queue to receiver
		if ep != nil {
			typedmemmove(c.elemtype, ep, qp)
		}
		// copy data from sender to queue
		typedmemmove(c.elemtype, qp, sg.elem)
		c.recvx++
		if c.recvx == c.dataqsiz {
			c.recvx = 0
		}
		c.sendx = c.recvx // c.sendx = (c.sendx+1) % c.dataqsiz
	}
	sg.elem = nil
	gp := sg.g
	unlockf()
	gp.param = unsafe.Pointer(sg)
	sg.success = true
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
	goready(gp, skip+1)
}
```

如果`c.dataqsiz == 0`，说明是无缓冲`channel`，如果ep不为空`recvDirect`直接拷贝内存，从`receiver -> sender`

否则`buffer`一定是满的，数据拷贝过程为：`qp -> ep`，`sg -> ep`，因为`buffer`的数据是先于`sendq`发送的，所以先将`buffer`中`recvx`指向的数据拷贝到`ep`，再将`sender`中的数据拷贝到`buffer`中`recvx`指向的内存地址，`recvx`加1，由于此时`sender`赋值到`buffer`中，并不会因为`recv`而使`qcount`减少，所以发送索引`sendx = recvx`

获取`sudog`绑定的协程`gp`

通过`goready`唤醒发送的`goroutine`

#### 核心方法---recvDirect

```go
func recvDirect(t *_type, sg *sudog, dst unsafe.Pointer) {
	// dst is on our stack or the heap, src is on another stack.
	// The channel is locked, so src will not move during this
	// operation.
	src := sg.elem
	typeBitsBulkBarrier(t, uintptr(dst), uintptr(src), t.size)
	memmove(dst, src, t.size)
}
```

与`sendDirect`类似，需要`Write Barrier`

### closechan

#### closechan流程图

![](https://tva1.sinaimg.cn/large/008i3skNly1gx48x37kc6j30z80ns40h.jpg)

```go
if c == nil {
	panic(plainError("close of nil channel"))
}

lock(&c.lock)
if c.closed != 0 {
	unlock(&c.lock)
	panic(plainError("close of closed channel"))
}
```

关闭一个`nil`的`channel`或者关闭一个`closed`的`channel`都会导致`panic`

```go
c.closed = 1
```

设置将关闭的`channel`的`closed`为1

```go
// release all readers
for {
	sg := c.recvq.dequeue()
	if sg == nil {
		break
	}
	if sg.elem != nil {
		typedmemclr(c.elemtype, sg.elem)
		sg.elem = nil
	}
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
	gp := sg.g
	gp.param = unsafe.Pointer(sg)
	sg.success = false
	if raceenabled {
		raceacquireg(gp, c.raceaddr())
	}
	glist.push(gp)
}
```

释放所有的`receiver`的资源

```go
// release all writers (they will panic)
for {
	sg := c.sendq.dequeue()
	if sg == nil {
		break
	}
	sg.elem = nil
	if sg.releasetime != 0 {
		sg.releasetime = cputicks()
	}
	gp := sg.g
	gp.param = unsafe.Pointer(sg)
	sg.success = false
	if raceenabled {
		raceacquireg(gp, c.raceaddr())
	}
	glist.push(gp)
}
```

释放所有的`sender`的资源，同时他们都会`panic`

```go
// Ready all Gs now that we've dropped the channel lock.
for !glist.empty() {
   gp := glist.pop()
   gp.schedlink = 0
   goready(gp, 3)
}
```

唤醒所有的`goroutine`

### 非阻塞发送和非阻塞接受

#### 非阻塞发送---selectnbsend

```go
func selectnbsend(c *hchan, elem unsafe.Pointer) (selected bool) {
	return chansend(c, elem, false, getcallerpc())
}
```

`select`中`channel`的发送，便是非阻塞发送

```go
select {
case c1 <- v:
  ...func1
case c2 <- v:
  ...func2
 default: 
  ...func3
}
```

等价于：

```go 
if selectnbsend(c1, &v) {
  ...func1
} else if selectnbsend(c2, &v) {
  ...func2
} else {
  ...func3
}
```

#### 非阻塞发送---selectngrecv

```go
func selectnbrecv(elem unsafe.Pointer, c *hchan) (selected bool) {
   selected, _ = chanrecv(c, elem, false)
   return
}
```

`select`中`channel`的接收，便是非阻塞接收，和`selectnbsend`同理

```go
select {
case v = <-c1:
  ...func1
case v = <-c2:
  ...func2
 default: 
  ...func3
}
```

等价于：

```go
if selectnbrecv(&v, c1) {
  ...func1
} else if selectnbrecv(&v, c2) {
  ...func2
} else {
  ...func3
}
```

## 图解循环队列

![](https://tva1.sinaimg.cn/large/008i3skNly1gx49bzivk6j30u02cwag4.jpg)


