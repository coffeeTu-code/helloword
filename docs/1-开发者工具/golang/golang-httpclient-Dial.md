# Go 如何使得 Web 工作

## TCP客户端

在TCP编程中，客户端和服务端的流程以及联系如下图:

![115185303067.png](https://upload-images.jianshu.io/upload_images/14738618-69206cb19c2ba80e.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

### 【相关函数】

1. ` net.ResolveTCPAddr("tcp4", "127.0.0.1:80")`

2. `net.DialTCP("tcp", nil, addr)`

3. `net.Dial("tcp", addr)`

4. ` l, err := net.ListenTCP("tcp", “:80”)`

5. `listener, err := net.Listen("tcp", ":80")`

### 【Dialer结构体定义】

```$go
type Dialer struct {
    //DeadLine和Timeout选项: 用于不成功拨号的超时设置
    Timeout time.Duration  //拨号等待连接结束的最大时间数。如果同时设置了Deadline, 可以更早失败。默认没有超时。 当使用TCP并使用多个IP地址拨号主机名，超时会在它们之间划分。使用或不使用超时，操作系统都可以强迫更早超时。例如，TCP超时一般在3分钟左右
  
    Deadline time.Time  //是拨号即将失败的绝对时间点。如果设置了Timeout, 可能会更早失败。0值表示没有截止期限， 或者依赖操作系统或使用Timeout选项

    LocalAddr Addr //真正dial时的本地地址，兼容各种类型(TCP、UDP...),如果为nil，则系统自动选择一个地址

    DualStack bool // 双协议栈，即是否同时支持ipv4和ipv6.当network值为tcp时，dial函数会向host主机的v4和v6地址都发起连接

    FallbackDelay time.Duration // 当DualStack为真，ipv6会延后于ipv4发起，此字段即为延迟时间，默认为300ms

    //KeepAlive选项: 管理连接的使用寿命(life span)
    KeepAlive time.Duration   //为活动网络连接指定保持活动的时间。如果设置为0，没有启用keep-alive。不支持keep-alive的网络协议会忽略掉这个字段

    Resolver *Resolver  //可选项，指定使用的可替代resolver

    Cancel <-chan struct{} // 可选通道，用于取消dial. 它的闭包表示拨号应该被取消。不是所有的拨号类型都支持拨号取消。 已废弃，可使用DialContext代替
}
```

底层dial函数DialContext()

```$go
func (d *Dialer) DialContext(ctx context.Context, network, address string) (Conn, error) {
    ...
    deadline := d.deadline(ctx, time.Now()) 
    //d.deadline() 比较d.deadline、ctx.deadline、now+timeout，返回其中最小.如果都为空，返回0
    ...
    subCtx, cancel := context.WithDeadline(ctx, deadline) //设置新的超时context
    defer cancel()
    ...
    // Shadow the nettrace (if any) during resolve so Connect events don't fire for DNS lookups.
    resolveCtx := ctx
    ...//给resolveCtx带上一些value

    addrs, err := d.resolver().resolveAddrList(resolveCtx, "dial", network, address, d.LocalAddr) // 解析IP地址，返回值是一个切片

    dp := &dialParam{
        Dialer:  *d,
        network: network,
        address: address,
    }

    var primaries, fallbacks addrList
    if d.DualStack && network == "tcp" { //表示同时支持ipv4和ipv6
        primaries, fallbacks = addrs.partition(isIPv4) // 将addrs分成两个切片，前者包含ipv4地址，后者包含ipv6地址
    } else {
        primaries = addrs
    }

    var c Conn
    if len(fallbacks) > 0 {//有ipv6的情况，v4和v6一起dial
        c, err = dialParallel(ctx, dp, primaries, fallbacks)
    } else {
        c, err = dialSerial(ctx, dp, primaries)
    }
    if err != nil {
        return nil, err
    }
    ...
    return c, nil
}
```

### 【Dial连接】

成功拨号之后，我们就可以如上所述的那样，将新的连接与其他的输入输出流同等对待了。我们甚至可以将连接包装进bufio.ReadWriter中，这样可以使用各种ReadWriter方法，例如ReadString(), ReadBytes, WriteString等等。

> 记住缓冲Writer在写之后需要调用Flush()方法, 这样所有的数据才会刷到底层网络连接中。

```$go
func Open(addr string) (*bufio.ReadWriter, error) {
    conn, err := net.Dial("tcp", addr)  //注意，当前函数Dial实时分配端口号进行连接，指定端口号使用DialTCP
    if err != nil {
        return nil, errors.Wrap(err, "Dialing "+addr+" failed")
    }
    // 将net.Conn对象包装到bufio.ReadWriter中
    return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
}
```

最后，每个连接对象都有一个Close()方法来终止通信。



## HTTP客户端

