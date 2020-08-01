# Go 微服务框架对比：Go Micro, Go Kit, Gizmo, Kite

>
> Go 微服务框架对比：Go Micro, Go Kit, Gizmo, Kite
 [https://learnku.com/go/t/36973](https://learnku.com/go/t/36973)
>
>
>
>
>
>
>
>
>

##  Go Micro

什么是 Go Micro？它是一个可插入的 RPC 框架，用于在 Go 中编写微服务。开箱即用，您将收到：

- 服务发现 - 应用程序自动注册到服务发现系统。
- 负载平衡 - 客户端负载平衡，用于平衡服务实例之间的请求。
- 同步通信 - 提供请求 / 响应传输层。
- 异步通信 - 内置发布 / 订阅功能。
- 消息编码 - 基于消息的内容类型头的编码 / 解码。
- RPC 客户机 / 服务器包 - 利用上述功能并公开接口来构建微服务。


Go 微体系结构可以描述为三层堆栈。

![https://segmentfault.com/img/bVblToP?w=753&h=164](https://segmentfault.com/img/bVblToP?w=753&h=164)

顶层由客户端 - 服务器模型和服务抽象组成。服务器是用于编写服务的构建块。客户端提供了向服务请求的接口。

底层由以下类型的插件组成：

- 代理 - 为异步发布 / 订阅通信提供消息代理的接口。
- 编解码器 - 用于编码 / 解码消息。支持的格式包括 json，bson，protobuf，msgpack 等。
- 注册表 - 提供服务发现机制（默认为 Consul）。
- 选择器 - 建立在注册表上的负载平衡抽象。它允许使用诸如随机，轮循，最小康等算法来 “选择” 服务。
- 传输 - 服务之间同步请求 / 响应通信的接口。
- Go Micro 还提供了 Sidecar 等功能。这使您可以使用以 Go 以外的语言编写的服务。 Sidecar 提供服务注册，gRPC 编码 / 解码和 HTTP 处理程序。它支持多种语言。

## Go Kit

Go Kit 是一个用于在 Go 中构建微服务的编程工具包。与 Go Micro 不同，它被设计为一个用于导入二进制包的库。

Go Kit 遵循简单的规则，例如:

- 没有全局状态
- 声明式组合
- 显式依赖关系
- 接口即约定
- 领域驱动设计

在 Go Kit 中，您可以找到以下的包:

- 认证 - Basic 认证和 JWT 认证
- 传输 - HTTP、Nats、gRPC 等等。
- 日志记录 - 用于结构化服务日志记录的通用接口。
- 指标 - CloudWatch、Statsd、Graphite 等。
- 追踪 - Zipkin 和 Opentracing。
- 服务发现 - Consul、Etcd、Eureka 等等。
- 断路器 - Hystrix 的 Go 实现。

