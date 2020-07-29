# 谈谈 Websocket



[TOC]



## 一.WebSocket的使用场景

从下面我列举的这些场景来看，一个共同点就是，**高实时性**！

1. 社交聊天

最著名的就是微信，QQ，这一类社交聊天的app。这一类聊天app的特点是低延迟，高即时。即时是这里面要求最高的，如果有一个紧急的事情，通过IM软件通知你，假设网络环境良好的情况下，这条message还无法立即送达到你的客户端上，紧急的事情都结束了，你才收到消息，那么这个软件肯定是失败的。

2. 弹幕

说到这里，大家一定里面想到了A站和B站了。确实，他们的弹幕一直是一种特色。而且弹幕对于一个视频来说，很可能弹幕才是精华。发弹幕需要实时显示，也需要和聊天一样，需要即时。

3. 多玩家游戏

剑灵、DNF...

4. 协同编辑

现在很多开源项目都是分散在世界各地的开发者一起协同开发，此时就会用到版本控制系统，比如Git，SVN去合并冲突。但是如果有一份文档，支持多人实时在线协同编辑，那么此时就会用到比如WebSocket了，它可以保证各个编辑者都在编辑同一个文档，此时不需要用到Git，SVN这些版本控制，因为在协同编辑界面就会实时看到对方编辑了什么，谁在修改哪些段落和文字。

5. 股票基金实时报价

金融界瞬息万变——几乎是每毫秒都在变化。如果采用的网络架构无法满足实时性，那么就会给客户带来巨大的损失。几毫秒钱股票开始大跌，几秒以后才刷新数据，一秒钟的时间内，很可能用户就已经损失巨大财产了。

6. 体育实况更新

全世界的球迷，体育爱好者特别多，当然大家在关心自己喜欢的体育活动的时候，比赛实时的赛况是他们最最关心的事情。这类新闻中最好的体验就是利用Websocket达到实时的更新！

7. 视频会议/聊天

视频会议并不能代替和真人相见，但是他能让分布在全球天涯海角的人聚在电脑前一起开会。既能节省大家聚在一起路上花费的时间，讨论聚会地点的纠结，还能随时随地，只要有网络就可以开会。

8. 基于位置的应用

越来越多的开发者借用移动设备的GPS功能来实现他们基于位置的网络应用。如果你一直记录用户的位置(比如运行应用来记录运动轨迹)，你可以收集到更加细致化的数据。

9. 在线教育

在线教育近几年也发展迅速。优点很多，免去了场地的限制，能让名师的资源合理的分配给全国各地想要学习知识的同学手上，Websocket是个不错的选择，可以视频聊天、即时聊天以及其与别人合作一起在网上讨论问题...

10. 智能家居

这也是我一毕业加入的一个伟大的物联网智能家居的公司。考虑到家里的智能设备的状态必须需要实时的展现在手机app客户端上，毫无疑问选择了Websocket。



## 二.WebSocket诞生由来

### 1.最开始的轮询Polling阶段
![](https://img.halfrost.com/Blog/ArticleImage/8_2.png)

这种方式下，是不适合获取实时信息的，客户端和服务器之间会一直进行连接，每隔一段时间就询问一次。客户端会轮询，有没有新消息。这种方式连接数会很多，一个接受，一个发送。而且每次发送请求都会有Http的Header，会很耗流量，也会消耗CPU的利用率。

### 2.改进版的长轮询Long polling阶段

![](https://img.halfrost.com/Blog/ArticleImage/8_3.png)

长轮询是对轮询的改进版，客户端发送HTTP给服务器之后，有没有新消息，如果没有新消息，就一直等待。当有新消息的时候，才会返回给客户端。在某种程度上减小了网络带宽和CPU利用率等问题。但是这种方式还是有一种弊端：例如假设服务器端的数据更新速度很快，服务器在传送一个数据包给客户端后必须等待客户端的下一个Get请求到来，才能传递第二个更新的数据包给客户端，那么这样的话，客户端显示实时数据最快的时间为2×RTT（往返时间），而且如果在网络拥塞的情况下，这个时间用户是不能接受的，比如在股市的的报价上。另外，由于http数据包的头部数据量往往很大（通常有400多个字节），但是真正被服务器需要的数据却很少（有时只有10个字节左右），这样的数据包在网络上周期性的传输，难免对网络带宽是一种浪费。

### 3.WebSocket诞生

现在急需的需求是能支持客户端和服务器端的双向通信，而且协议的头部又没有HTTP的Header那么大，于是，Websocket就诞生了！

![](https://img.halfrost.com/Blog/ArticleImage/8_4.png)

上图就是Websocket和Polling的区别，从图中可以看到Polling里面客户端发送了好多Request，相对的，WebSocket只有一个Upgrade，非常简洁高效。至于消耗方面的比较就要看下图了

！[](https://img.halfrost.com/Blog/ArticleImage/8_5.png)

上图中，我们先看蓝色的柱状图，是Polling轮询消耗的流量，这次测试，HTTP请求和响应头信息开销总共包括871字节。当然每次测试不同的请求，头的开销不同。这次测试都以871字节的请求来测试。

```
**Use case A:**1,000 clients polling every second: Network throughput is (871 x 1,000) = 871,000 bytes = 6,968,000 bits per second (6.6 Mbps)

Use case B: 10,000 clients polling every second: Network throughput is (871 x 10,000) = 8,710,000 bytes = 69,680,000 bits per second (66 Mbps)

**Use case C:**100,000 clients polling every 1 second: Network throughput is (871 x 100,000) = 87,100,000 bytes = 696,800,000 bits per second (665 Mbps)
```

而Websocket的Frame是 just two bytes of overhead instead of 871，仅仅用2个字节就代替了轮询的871字节！

```
**Use case A:**1,000 clients receive 1 message per second: Network throughput is (2 x 1,000) = 2,000 bytes = 16,000 bits per second (0.015 Mbps)

**Use case B:**10,000 clients receive 1 message per second: Network throughput is (2 x 10,000) = 20,000 bytes = 160,000 bits per second (0.153 Mbps)

**Use case C:**100,000 clients receive 1 message per second: Network throughput is (2 x 100,000) = 200,000 bytes = 1,600,000 bits per second (1.526 Mbps)
```

相同的每秒客户端轮询的次数，当次数高达10W/s的高频率次数的时候，Polling轮询需要消耗665Mbps，而Websocket仅仅只花费了1.526Mbps，将近435倍！！



## 三.谈谈WebSocket协议原理

Websocket是应用层第七层上的一个应用层协议，它必须依赖 HTTP 协议进行一次握手 ，握手成功后，数据就直接从 TCP 通道传输，与 HTTP 无关了。

Websocket的数据传输是frame形式传输的，比如会将一条消息分为几个frame，按照先后顺序传输出去。这样做会有几个好处：

1. 大数据的传输可以分片传输，不用考虑到数据大小导致的长度标志位不足够的情况。
2. 和http的chunk一样，可以边生成数据边传递消息，即提高传输效率。
 
```
0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-------+-+-------------+-------------------------------+
 |F|R|R|R| opcode|M| Payload len |    Extended payload length    |
 |I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
 |N|V|V|V|       |S|             |   (if payload len==126/127)   |
 | |1|2|3|       |K|             |                               |
 +-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
 |     Extended payload length continued, if payload len == 127  |
 + - - - - - - - - - - - - - - - +-------------------------------+
 |                               |Masking-key, if MASK set to 1  |
 +-------------------------------+-------------------------------+
 | Masking-key (continued)       |          Payload Data         |
 +-------------------------------- - - - - - - - - - - - - - - - +
 :                     Payload Data continued ...                :
 + - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
 |                     Payload Data continued ...                |
 +---------------------------------------------------------------+



    FIN      1bit 表示信息的最后一帧，flag，也就是标记符
    RSV 1-3  1bit each 以后备用的 默认都为 0
    Opcode   4bit 帧类型，稍后细说
    Mask     1bit 掩码，是否加密数据，默认必须置为1 （这里很蛋疼）
    Payload  7bit 数据的长度
    Masking-key      1 or 4 bit 掩码
    Payload data     (x + y) bytes 数据
    Extension data   x bytes  扩展数据
    Application data y bytes  程序数据
```

具体的规范，还请看官网的[RFC 6455文档](https://tools.ietf.org/html/rfc6455)给出的详细定义。这里还有一个[翻译版本](https://www.gitbook.com/book/chenjianlong/rfc-6455-websocket-protocol-in-chinese/details)。



## 四.WebSocket 和 Socket的区别与联系

首先，Socket 其实并不是一个协议。它工作在 OSI 模型会话层（第5层），是为了方便大家直接使用更底层协议（一般是 TCP 或 UDP ）而存在的一个抽象层。Socket是对TCP/IP协议的封装，Socket本身并不是协议，而是一个调用接口(API)。

![](https://img.halfrost.com/Blog/ArticleImage/8_6.png)

Socket通常也称作”套接字”，用于描述IP地址和端口，是一个通信链的句柄。网络上的两个程序通过一个双向的通讯连接实现数据的交换，这个双向链路的一端称为一个Socket，一个Socket由一个IP地址和一个端口号唯一确定。应用程序通常通过”套接字”向网络发出请求或者应答网络请求。

Socket在通讯过程中，服务端监听某个端口是否有连接请求，客户端向服务端发送连接请求，服务端收到连接请求向客户端发出接收消息，这样一个连接就建立起来了。客户端和服务端也都可以相互发送消息与对方进行通讯，直到双方连接断开。