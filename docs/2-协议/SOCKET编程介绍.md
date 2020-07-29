# go Socket编程

本地的进程间通信（IPC）有很多种方式，但可以总结为下面4类：
1. 消息传递（管道、FIFO、消息队列）
2. 同步（互斥量、条件变量、读写锁、文件和写记录锁、信号量）
3. 共享内存（匿名的和具名的）
4. 远程过程调用（Solaris门和Sun RPC）

网络中进程之间如何通信？首要解决的问题是如何唯一标识一个进程，否则通信无从谈起！在本地可以通过进程PID来唯一标识一个进程，但是在网络中这是行不通的。其实TCP/IP协议族已经帮我们解决了这个问题，网络层的“ip地址”可以唯一标识网络中的主机，而传输层的“协议+端口”可以唯一标识主机中的应用程序（进程）。这样利用三元组（ip地址，协议，端口）就可以标识网络的进程了，网络中的进程通信就可以利用这个标志与其它进程进行交互。

Socket起源于Unix，而Unix基本哲学之一就是“一切皆文件”，都可以用“打开open –> 读写write/read –> 关闭close”模式来操作。Socket就是该模式的一个实现，网络的Socket数据传输是一种特殊的I/O，Socket也是一种文件描述符。Socket也具有一个类似于打开文件的函数调用：Socket()，该函数返回一个整型的Socket描述符，随后的连接建立、数据传输等操作都是通过该Socket实现的。

常用的Socket类型有两种：流式Socket（SOCK_STREAM）和数据报式Socket（SOCK_DGRAM）。流式是一种面向连接的Socket，针对于面向连接的TCP服务应用；数据报式Socket是一种无连接的Socket，对应于无连接的UDP服务应用。

### socket基础

在Go的et包中有很多函数来操作IP，但是其中比较有用的也就几个，其中`ParseIP(s string) IP`函数会把一个IPv4或者IPv6的地址转化成IP类型。

IPv4地址 `127.0.0.1 `  
IPv6地址 `2002:c0e8:82e7:0:0:0:c0e8:82e7`

`func ResolveTCPAddr(net, addr string) (*TCPAddr, os.Error)`函数可以获取一个TCPAddr。

### 客户端代码

我们的客户端代码如下所示：

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)
	// result, err := ioutil.ReadAll(conn)
	result := make([]byte, 256)
	_, err = conn.Read(result)
	checkError(err)
	fmt.Println(string(result))
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```

### 服务端代码

```go
package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"strings"
)

func main() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	request := make([]byte, 128) // set maxium request length to 128B to prevent flood attack
	defer conn.Close()  // close connection before exit
	for {
		read_len, err := conn.Read(request)

		if err != nil {
			fmt.Println(err)
			break
		}

    		if read_len == 0 {
    			break // connection already closed by client
    		} else if strings.TrimSpace(string(request[:read_len])) == "timestamp" {
    			daytime := strconv.FormatInt(time.Now().Unix(), 10)
    			conn.Write([]byte(daytime))
    		} else {
    			daytime := time.Now().String()
    			conn.Write([]byte(daytime))
    		}

    		request = make([]byte, 128) // clear last read content
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
```

### 控制TCP连接

TCP有很多连接控制函数，我们平常用到比较多的有如下几个函数：

`func DialTimeout(net, addr string, timeout time.Duration) (Conn, error)`


设置建立连接的超时时间，客户端和服务器端都适用，当超过设置时间时，连接自动关闭。

`func (c *TCPConn) SetReadDeadline(t time.Time) error`

`func (c *TCPConn) SetWriteDeadline(t time.Time) error`


用来设置写入/读取一个连接的超时时间。当超过设置时间时，连接自动关闭。

`func (c *TCPConn) SetKeepAlive(keepalive bool) os.Error`


设置keepAlive属性，是操作系统层在tcp上没有数据和ACK的时候，会间隔性的发送keepalive包，操作系统可以通过该包来判断一个tcp连接是否已经断开，在windows上默认2个小时没有收到数据和keepalive包的时候人为tcp连接已经断开，这个功能和我们通常在应用层加的心跳包的功能类似。

###  其他项

> 主机字节序：就是我们平常说的大端和小端模式：不同的CPU有不同的字节序类型，这些字节序是指整数在内存中保存的顺序，这个叫做主机序。引用标准的Big-Endian和Little-Endian的定义如下：
> + Little-Endian就是低位字节排放在内存的低地址端，高位字节排放在内存的高地址端。
> + Big-Endian就是高位字节排放在内存的低地址端，低位字节排放在内存的高地址端。

> 网络字节序：4个字节的32 bit值以下面的次序传输：首先是0～7bit，其次8～15bit，然后16～23bit，最后是24~31bit。这种传输次序称作大端字节序。由于TCP/IP首部中所有的二进制整数在网络中传输时都要求以这种次序，因此它又称作网络字节序。字节序，顾名思义字节的顺序，就是大于一个字节类型的数据在内存中的存放顺序，一个字节的数据没有顺序的问题了。所以：在将一个地址绑定到socket的时候，请先将主机字节序转换成为网络字节序，而不要假定主机字节序跟网络字节序一样使用的是Big-Endian。由于这个问题曾引发过血案！公司项目代码中由于存在这个问题，导致了很多莫名其妙的问题，所以请谨记对主机字节序不要做任何假定，务必将其转化为网络字节序再赋给socket。
