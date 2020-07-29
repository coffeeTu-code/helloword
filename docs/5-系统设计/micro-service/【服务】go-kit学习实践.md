> Github：[https://github.com/go-kit/kit](https://github.com/go-kit/kit)
> 
> GoDoc：[https://godoc.org/github.com/go-kit/kit](https://godoc.org/github.com/go-kit/kit)
>
> 负载均衡 平滑的balancer策略：[github.com/bj-wangjia/go-kit](github.com/bj-wangjia/go-kit)
>



[TOC]



---------------------------------------------------------
# go-kit介绍


## Go-kit 包介绍


1. 本身不是一个框架，而是一套微服务工具集，是框架的底层，用它的话来说，如果你希望构建一个框架，而Go-kit 就希望成为你的框架的一部分  

2. 可以用Go-kit 做适应自己平台的框架  

3. 它自身称为toolkit，并不是framework  

4. 它主要是为了满足5大原则，单一职责原则，开放原则，封闭原则，依赖倒置原则，接口隔离原则  


## 三层模型


Go-kit最核心是提供了三层模型来解耦业务，这是我们用它的主要目的，模型由上到下分别是
transport -> endpoint -> service

1. Transport  
可以理解为是个拦截器，负责请求协议的实现和路由转发，HTTP、gRPC、thrift等相关的逻辑。

2. Endpoint  
定义Request和Response格式，并可以使用装饰器包装Service函数，以此来实现各种中间件嵌套。

3. Service  
服务功能、接口具体实现。


## 功能描述

go-kit提供以下功能：

1. Circuit breaker（熔断器）
2. Rate limiter（限流器）
3. Logging（日志）
4. Metrics（Prometheus统计）
5. Request tracing（请求跟踪）
6. Service discovery and load balancing（服务发现和负载均衡）


## 系列文章

- [go-kit微服务系列目录](https://juejin.im/post/5c861c93f265da2de7138615)

- [stringsvc教程](https://gokit.io/examples/stringsvc.html) 是一个教程，带您从基本原则开始编写服务。它可以帮助您了解Go kit设计中的决策。

- [addsvc](https://github.com/go-kit/kit/blob/master/examples/addsvc) 是原始的示例服务。它公开了所有支持的传输上的一组操作。它已完全记录，检测并使用分布式跟踪。它还演示了如何创建和使用客户端软件包。它演示了Go工具包的几乎所有功能。

- [profilesvc](https://github.com/go-kit/kit/blob/master/examples/profilesvc) 演示了如何使用Go kit编写带有REST-ish API的微服务。它使用net / http和出色的Gorilla Web工具包。

- [shipping](https://github.com/go-kit/kit/blob/master/examples/shipping) 是基于域驱动设计原则的，由多个微服务组成的完整的“真实世界”应用程序。

- [apigateway](https://github.com/go-kit/kit/blob/master/examples/apigateway) 演示了如何实现 由Consul服务发现系统支持的 API网关模式。



---------------------------------------------------------
# go-kit 微服务：http实战



## 一、Transport


Transport层负责http拦截处理。

先来看下Server定义：

```go

// Server wraps an endpoint and implements http.Handler.
type Server struct {
	e            endpoint.Endpoint
	dec          DecodeRequestFunc
	enc          EncodeResponseFunc
	before       []RequestFunc
	after        []ServerResponseFunc
	errorEncoder ErrorEncoder
	finalizer    []ServerFinalizerFunc
	errorHandler transport.ErrorHandler
}

// NewServer constructs a new server, which implements http.Handler and wraps
// the provided endpoint.
func NewServer(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	enc EncodeResponseFunc,
	options ...ServerOption,
) *Server {
	s := &Server{
		e:            e,
		dec:          dec,
		enc:          enc,
		errorEncoder: DefaultErrorEncoder,
		errorHandler: transport.NewLogErrorHandler(log.NewNopLogger()),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

```

Server处理http的流程如下：

1. finalizer
2. before
3. **dec**
4. **e**
5. after
6. **enc**

需要我们自行实现三个方法：
- dec
- e
- enc


```go

// ServeHTTP implements http.Handler.
func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if len(s.finalizer) > 0 {
		iw := &interceptingWriter{w, http.StatusOK, 0}
		defer func() {
			ctx = context.WithValue(ctx, ContextKeyResponseHeaders, iw.Header())
			ctx = context.WithValue(ctx, ContextKeyResponseSize, iw.written)
			for _, f := range s.finalizer {
				f(ctx, iw.code, r)
			}
		}()
		w = iw
	}

	for _, f := range s.before {
		ctx = f(ctx, r)
	}

	request, err := s.dec(ctx, r)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		s.errorEncoder(ctx, err, w)
		return
	}

	response, err := s.e(ctx, request)
	if err != nil {
		s.errorHandler.Handle(ctx, err)
		s.errorEncoder(ctx, err, w)
		return
	}

	for _, f := range s.after {
		ctx = f(ctx, w)
	}

	if err := s.enc(ctx, w, response); err != nil {
		s.errorHandler.Handle(ctx, err)
		s.errorEncoder(ctx, err, w)
		return
	}
}

```

使用`NewServer`注册服务的写法很简单：

- `PolarisExecute`是业务逻辑定义的接口
- `polarisServer`是实现业务逻辑接口的一个对象
- `makeGetResponseEndpoint(polarisSvc)`是将Service包装为Endpoint的方法
- `decodeHTTPRequest`是处理http request的方法
- `encodeHTTPResponse`是处理http response的方法

```
type PolarisHTTPServer struct {
	PolarisHandler http.Handler
}


func decodeHTTPRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request polaris.PolarisRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeHTTPResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

```

```
import 	httptransport "github.com/go-kit/kit/transport/http"

func PolarisServerHTTP() {
    var polarisSvc PolarisExecute
	polarisSvc = new(polarisServer)
	
    service := new(PolarisHTTPServer)
	service.PolarisHandler = httptransport.NewServer(
		makeGetResponseEndpoint(polarisSvc),
		decodeHTTPRequest,
		encodeHTTPResponse,
	)

	//启动http服务
	srv := &http.Server{Addr: "0.0.0.0:9888"}
	http.Handle("/polaris", service.PolarisHandler)

	go func() {
		<-config.StopChan
		srv.Close()
	}()

	if err := srv.ListenAndServe(); err != nil {
		// cannot panic, because this probably is an intentional close
		log.Printf("Httpserver: ListenAndServe() error: %s", err)
	}
}

```


## 二、Endpoint


将Servive包装成Endpoint后，注册到http服务中:

- `PolarisExecute`是业务逻辑的方法接口
- `polaris.PolarisRequest`是经过Transport层decodeHTTPRequest处理后的请求对象
- `srv.PolarisServerExecute`是处理业务逻辑的方法
- `resp`是业务逻辑处理结果，需Transport层encodeHTTPResponse处理后返回给http.ResponseWriter

```
func makeGetResponseEndpoint(srv PolarisExecute) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*polaris.PolarisRequest)

		resp, err := srv.PolarisServerExecute(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}
```

## 三、Service


定义业务逻辑服务方法接口。  

该接口将作为Service层级，可以被中间件服务层层包裹后暴漏给调用方。

```
type PolarisExecute interface {
	PolarisServerExecute(request *polaris.PolarisRequest) (*polaris.PolarisResponse, error)
}

```

```
type polarisServer struct {
}

func (s *polarisServer) PolarisServerExecute(request *polaris.PolarisRequest) (*polaris.PolarisResponse, error) {
    //业务处理逻辑
    ...
}
```

### 增加 限流服务

基于gokit内建的类型endpoint.Middleware，该类型实际上是一个function，使用装饰者模式实现对Endpoint的封装。定义如下：

```
# Go-kit Middleware Endpoint
type Middleware func(Endpoint) Endpoint

```

创建限流器


```
var ErrLimitExceed = errors.New("limit exceed error")

// NewTokenBucketLimitterWithBuildIn 使用x/time/rate创建限流中间件
func NewTokenBucketLimitterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
```

刚刚注册NewServer的地方，输入参数 endpoint 则可以修改为

```
NewTokenBucketLimitterWithBuildIn(
			//增加速率限制，每50毫秒补充一次，设置容量12
			rate.NewLimiter(rate.Every(time.Microsecond*50), 12))(makeGetResponseEndpoint(polarisSvc)),
```

现在代码是这样的了：

```
import 	httptransport "github.com/go-kit/kit/transport/http"

func PolarisServerHTTP() {
    var polarisSvc PolarisExecute
	polarisSvc = new(polarisServer)
	
    service := new(PolarisHTTPServer)
	service.PolarisHandler = httptransport.NewServer(
		NewTokenBucketLimitterWithBuildIn(
			//增加速率限制，每50毫秒补充一次，设置容量12
			rate.NewLimiter(rate.Every(time.Microsecond*50), 12))(makeGetResponseEndpoint(polarisSvc)),
		decodeHTTPRequest,
		encodeHTTPResponse,
	)

	//启动http服务
	srv := &http.Server{Addr: "0.0.0.0:9888"}
	http.Handle("/polaris", service.PolarisHandler)

	go func() {
		<-config.StopChan
		srv.Close()
	}()

	if err := srv.ListenAndServe(); err != nil {
		// cannot panic, because this probably is an intentional close
		log.Printf("Httpserver: ListenAndServe() error: %s", err)
	}
}
```


### 增加 日志服务


使用日志装饰器重新实现`PolarisExecute`业务方法接口，将实现业务方法接口的对象包裹进来，将包装后的`PolarisExecute`传给Endpoint。

```
type loggingMiddleware struct {
	logger seelog.LoggerInterface
	next   PolarisExecute
}

func (mw *loggingMiddleware) PolarisServerExecute(request *polaris.PolarisRequest) (response *polaris.PolarisResponse, err error) {
	defer func(begin time.Time) {
		if !request.Test {
			return
		}
		req, _ := json.Marshal(request)
		resp, _ := json.Marshal(response)
		mw.logger.Info("PolarisServerExecute", "\t", req, "\t", resp, "\t", err, "\t", time.Since(begin))
	}(time.Now())

	response, err = mw.next.PolarisServerExecute(request)
	return
}
```

现在代码是这样的了：

```
import 	httptransport "github.com/go-kit/kit/transport/http"

func PolarisServerHTTP() {
    var polarisSvc PolarisExecute
	polarisSvc = new(polarisServer)
	
	//add日志服务
	polarisSvc = &loggingMiddleware{
		logger: logger.Logger.Runtime,
		next:   polarisSvc,
	}
	
    service := new(PolarisHTTPServer)
	service.PolarisHandler = httptransport.NewServer(
		NewTokenBucketLimitterWithBuildIn(
			//增加速率限制，每50毫秒补充一次，设置容量12
			rate.NewLimiter(rate.Every(time.Microsecond*50), 12))(makeGetResponseEndpoint(polarisSvc)),
		decodeHTTPRequest,
		encodeHTTPResponse,
	)

	//启动http服务
	srv := &http.Server{Addr: "0.0.0.0:9888"}
	http.Handle("/polaris", service.PolarisHandler)

	go func() {
		<-config.StopChan
		srv.Close()
	}()

	if err := srv.ListenAndServe(); err != nil {
		// cannot panic, because this probably is an intentional close
		log.Printf("Httpserver: ListenAndServe() error: %s", err)
	}
}
```


### 增加 metrics服务

微服务中，API几乎是服务与外界的唯一交互渠道，API服务的稳定性、可靠性越来越成为不可忽略的部分。我们需要实时了解API的运行状况（请求次数、延时、失败等）。

增加 metrics 的处理方式与 增加日志服务 类似。

使用go-kit中间件机制为Service添加Prometheus监控指标采集功能。

```
type metricsMiddleware struct {
	//（1）请求数{有效请求，无效请求}；（2）响应数{有效响应，无效响应}；（3）处理时间{成功，失败}；
	requestLatency metrics.Histogram
	
	next           PolarisExecute
}

func (mw *metricsMiddleware) PolarisServerExecute(request *polaris.PolarisRequest) (response *polaris.PolarisResponse, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "PolarisServerExecute", "error", response.Status.String()}
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	response, err = mw.next.PolarisServerExecute(request)
	return
}
```

新增用于Prometheus轮循拉取监控指标的代码，开放API接口/metrics。

```
http.Handle("/metrics", promhttp.Handler())
```



现在代码是这样的了：

```
import 	httptransport "github.com/go-kit/kit/transport/http"

func PolarisServerHTTP() {
    var polarisSvc PolarisExecute
	polarisSvc = new(polarisServer)
	
	//add日志服务
	polarisSvc = &loggingMiddleware{
		logger: logger.Logger.Runtime,
		next:   polarisSvc,
	}
	
	//add监控服务
	fieldKeys := []string{"method", "error"}
	polarisSvc = &metricsMiddleware{
		requestLatency: kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "dsp",
			Subsystem: "polaris",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		next: polarisSvc,
	}
	
    service := new(PolarisHTTPServer)
	service.PolarisHandler = httptransport.NewServer(
		NewTokenBucketLimitterWithBuildIn(
			//增加速率限制，每50毫秒补充一次，设置容量12
			rate.NewLimiter(rate.Every(time.Microsecond*50), 12))(makeGetResponseEndpoint(polarisSvc)),
		decodeHTTPRequest,
		encodeHTTPResponse,
	)

	//启动http服务
	srv := &http.Server{Addr: "0.0.0.0:9888"}
	http.Handle("/polaris", service.PolarisHandler)
	//启动promhttp监控
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		<-config.StopChan
		srv.Close()
	}()

	if err := srv.ListenAndServe(); err != nil {
		// cannot panic, because this probably is an intentional close
		log.Printf("Httpserver: ListenAndServe() error: %s", err)
	}
}
```



## 附录

### 限流算法

常用的限流算法有两种：漏桶算法和令牌桶算法。


#### 漏桶算法

漏桶算法(Leaky Bucket)是网络世界中流量整形（Traffic Shaping）或速率限制（Rate Limiting）时经常使用的一种算法，它的主要目的是控制数据注入到网络的速率，平滑网络上的突发流量。漏桶算法提供了一种机制，通过它，突发流量可以被整形以便为网络提供一个稳定的流量。
漏桶算法思路很简单，水（请求）先进入到漏桶里，漏桶以一定的速度出水（接口有响应速率），当水流入速度过大会直接溢出（访问频率超过接口响应速率），然后就拒绝请求，可以看出漏桶算法能强行限制数据的传输速率。示意图如下：

![](https://user-gold-cdn.xitu.io/2019/2/21/16910727fa61ef72?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

因为漏桶的漏出速率是固定的参数,所以即使网络中不存在资源冲突（没有发生拥塞），漏桶算法也不能使流突发（burst）到端口速率。因此，漏桶算法对于存在突发特性的流量来说缺乏效率。


#### 令牌桶算法

令牌桶算法是网络流量整形（Traffic Shaping）和速率限制（Rate Limiting）中最常使用的一种算法。典型情况下，令牌桶算法用来控制发送到网络上的数据的数目，并允许突发数据的发送。
令牌桶算法的原理是系统会以一个恒定的速度往桶里放入令牌，而如果请求需要被处理，则需要先从桶里获取一个令牌，当桶里没有令牌可取时，则拒绝服务。从原理上看，令牌桶算法和漏桶算法是相反的，一个“进水”，一个是“漏水”。

![](https://user-gold-cdn.xitu.io/2019/2/21/1691072ad441f3a9?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

令牌桶的另外一个好处是可以方便的改变速度。 一旦需要提高速率，则按需提高放入桶中的令牌的速率。 一般会定时（比如100毫秒）往桶中增加一定数量的令牌，有些变种算法则实时的计算应该增加的令牌的数量。


### 监控组件 Prometheus

Prometheus （中文名称：普罗米修斯）是一套开源的系统监控报警框架。作为新一代的监控框架，Prometheus 具有以下特点：

- 提供强大的多维度数据模型，如Counter、Gauge、Histogram、Summary；
- 强大而灵活的查询语句（PromQL），可方便的实现对时间序列数据的查询、聚合操作；
- 易于管理与高效；
- 提供pull模式、push gateway方式实现时间序列数据的采集；
- 支持多种可视化图形界面：Grafana、Web UI、API clients；
- 报警规则管理、报警检测和报警推送功能。


### 可视化工具 Grafana

Grafana是一个跨平台的开源的度量分析和可视化工具，可以通过将采集的数据查询然后可视化的展示，并及时通知。它主要有以下六大特点：

- 展示方式：快速灵活的客户端图表，面板插件有许多不同方式的可视化指标和日志，官方库中具有丰富的仪表盘插件，比如热图、折线图、图表等多种展示方式；
- 数据源：Graphite，InfluxDB，OpenTSDB，Prometheus，Elasticsearch，CloudWatch和KairosDB等；
- 通知提醒：以可视方式定义最重要指标的警报规则，Grafana将不断计算并发送通知，在数据达到阈值时通过Slack、PagerDuty等获得通知；
- 混合展示：在同一图表中混合使用不同的数据源，可以基于每个查询指定数据源，甚至自定义数据源；
- 注释：使用来自不同数据源的丰富事件注释图表，将鼠标悬停在事件上会显示完整的事件元数据和标记；
- 过滤器：Ad-hoc过滤器允许动态创建新的键/值过滤器，这些过滤器会自动应用于使用该数据源的所有查询。



---------------------------------------------------------
# go-kit 微服务：grpc实战


## 一、proto定义

在 proto 文件中增加 rpc 接口定义代码
```
service RetargetService {
    rpc RetargetDsp (RetargetRequest) returns (RetargetResponse);
}
```

retarget_dsp.proto
```
syntax = "proto3";

package retarget_dsp;

//实际编译时需要写成 "base.proto"
import "base.proto";

//**************************** RetargetRequest 输入协议 ***************************
//defines the dsp Retarget request structure
message RetargetRequest {
    // 1. 竞价请求的唯一ID。
    // [Juno]
    // [Polaris]
    // [Rank]
    string request_id = 1;

    // 2.adx名称，例如doubleclick，mopub等
    // [Juno]
    // [Polaris]
    // [Rank]
    string exchange = 2;

    // 4.广告位展示信息。
    // [Juno]
    // [Polaris]
    // [Rank]
    common.Imp imp_info = 4;

    // 5.流量信息。
    // [Juno]
    // [Polaris]
    // [Rank]
    common.TrafficInfo traffic_info = 5;

    // 6.设备信息。
    // [Juno]
    // [Polaris]
    // [Rank]
    common.Device device_info = 6;

    // 7.表示使用设备的对象， 广告的受众
    // [Juno]
    // [Polaris]
    // [Rank]
    common.User user_info = 7;

    // 8.投放黑白名单限制信息。
    // [Juno]
    // [Polaris]
    // [Rank]
    common.BlackWhiteInfo black_white_info = 8;

    // 9.本次请求有效的工业，法律或政府条例
    // [Juno]
    // [Rank]
    common.Regs regs = 9;

    // 10.TC流量包名
    // [Juno] TcPkgTime，TcPackageName
    repeated string tc_package_name = 10;

    // 11.单子信息
    // [Juno]
    // [Polaris]
    // [Rank]
    repeated CampaignInfo campaign_list = 11;

    // 12.拍卖类型（胜出策略）1表示第一价格 ，2标识第二价格。
    // [Rank]
    int32 at = 12;

    // 13.超时时间设置。
    // [Juno]
    // [Polaris]
    // [Rank]
    int32 t_max = 13;

    // 14. ab test key
    // [Juno]
    // [Polaris]
    // [Rank]
    map<string, string> ab_key = 14;

    // 15.QA Debug 单子id列表，debug标志
    // [Juno]
    // [Polaris]
    // [Rank]
    repeated int64 debug_campaign_id_list = 15;

    // 16.是否抽样打印日志
    // [Juno]
    // [Polaris]
    // [Rank]
    bool open_log = 16;
}

message CampaignInfo {
    // 1.
    // [Juno]
    // [Polaris]
    // [Rank]
    int64 campaign_id = 1;

    // 2.人群包定向Id
    // [Rank]
    int64 audience_id = 2;

    // 3.重定向类型
    // [Rank]
    int32 retarget_type = 3;
}

//**************************** RetargetResponse 输入协议 ***************************
//defines the dsp Retarget response structure
message RetargetResponse {
    // 1.rt-dsp处理结果枚举值。
    int32 status = 1;

    // 2.rt-dsp处理错误消息
    string msg = 2;

    // 3.rt-dsp内部模块耗时统计
    int64 time = 3;

    // 4.算法信息
    RankInfo rank_info = 4;

    // 5.大模版信息
    common.BigTemplate big_template = 5;
}

message RankInfo {
    // 1.算法策略，算法标签
    string strategy = 1;

    // 2.二级算法策略，算法标签
    string subStrategy = 2;

    // 3.算法出价
    double price = 3;

    // 4.算法预估ecpm
    double ecpm_floor = 4;

    // 5.ivr_model
    double ivr = 5;
}

//**************************** Retarget Service ***************************
service RetargetService {
    rpc RetargetDsp (RetargetRequest) returns (RetargetResponse);
}

```


## 二、生成对应的XX.pb.go


生成对应的go语言代码文件：

```
protoc --go_out=plugins=grpc:. retarget_dsp.proto
```

生成的pb.go代码中包含 RetargetDsp 接口的客户端和服务端定义：

```
// Client API for RetargetService service

type RetargetServiceClient interface {
	RetargetDsp(ctx context.Context, in *RetargetRequest, opts ...grpc.CallOption) (*RetargetResponse, error)
}

// Server API for RetargetService service

type RetargetServiceServer interface {
	RetargetDsp(context.Context, *RetargetRequest) (*RetargetResponse, error)
}
```

源文件：

```
package retarget_dsp

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import common "gitlab.mobvista.com/mvdsp/protoc_new/pkg"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// **************************** RetargetRequest 输入协议 ***************************
// defines the dsp Retarget request structure
type RetargetRequest struct {
	// 1. 竞价请求的唯一ID。
	// [Juno]
	// [Polaris]
	// [Rank]
	RequestId string `protobuf:"bytes,1,opt,name=request_id,json=requestId" json:"request_id,omitempty"`
	// 2.adx名称，例如doubleclick，mopub等
	// [Juno]
	// [Polaris]
	// [Rank]
	Exchange string `protobuf:"bytes,2,opt,name=exchange" json:"exchange,omitempty"`
	// 4.广告位展示信息。
	// [Juno]
	// [Polaris]
	// [Rank]
	ImpInfo *common.Imp `protobuf:"bytes,4,opt,name=imp_info,json=impInfo" json:"imp_info,omitempty"`
	// 5.流量信息。
	// [Juno]
	// [Polaris]
	// [Rank]
	TrafficInfo *common.TrafficInfo `protobuf:"bytes,5,opt,name=traffic_info,json=trafficInfo" json:"traffic_info,omitempty"`
	// 6.设备信息。
	// [Juno]
	// [Polaris]
	// [Rank]
	DeviceInfo *common.Device `protobuf:"bytes,6,opt,name=device_info,json=deviceInfo" json:"device_info,omitempty"`
	// 7.表示使用设备的对象， 广告的受众
	// [Juno]
	// [Polaris]
	// [Rank]
	UserInfo *common.User `protobuf:"bytes,7,opt,name=user_info,json=userInfo" json:"user_info,omitempty"`
	// 8.投放黑白名单限制信息。
	// [Juno]
	// [Polaris]
	// [Rank]
	BlackWhiteInfo *common.BlackWhiteInfo `protobuf:"bytes,8,opt,name=black_white_info,json=blackWhiteInfo" json:"black_white_info,omitempty"`
	// 9.本次请求有效的工业，法律或政府条例
	// [Juno]
	// [Rank]
	Regs *common.Regs `protobuf:"bytes,9,opt,name=regs" json:"regs,omitempty"`
	// 10.TC流量包名
	// [Juno] TcPkgTime，TcPackageName
	TcPackageName []string `protobuf:"bytes,10,rep,name=tc_package_name,json=tcPackageName" json:"tc_package_name,omitempty"`
	// 11.单子信息
	// [Juno]
	// [Polaris]
	// [Rank]
	CampaignList []*CampaignInfo `protobuf:"bytes,11,rep,name=campaign_list,json=campaignList" json:"campaign_list,omitempty"`
	// 12.拍卖类型（胜出策略）1表示第一价格 ，2标识第二价格。
	// [Rank]
	At int32 `protobuf:"varint,12,opt,name=at" json:"at,omitempty"`
	// 13.超时时间设置。
	// [Juno]
	// [Polaris]
	// [Rank]
	TMax int32 `protobuf:"varint,13,opt,name=t_max,json=tMax" json:"t_max,omitempty"`
	// 14. ab test key
	// [Juno]
	// [Polaris]
	// [Rank]
	AbKey map[string]string `protobuf:"bytes,14,rep,name=ab_key,json=abKey" json:"ab_key,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// 15.QA Debug 单子id列表，debug标志
	// [Juno]
	// [Polaris]
	// [Rank]
	DebugCampaignIdList []int64 `protobuf:"varint,15,rep,name=debug_campaign_id_list,json=debugCampaignIdList" json:"debug_campaign_id_list,omitempty"`
	// 16.是否抽样打印日志
	// [Juno]
	// [Polaris]
	// [Rank]
	OpenLog bool `protobuf:"varint,16,opt,name=open_log,json=openLog" json:"open_log,omitempty"`
}

func (m *RetargetRequest) Reset()                    { *m = RetargetRequest{} }
func (m *RetargetRequest) String() string            { return proto.CompactTextString(m) }
func (*RetargetRequest) ProtoMessage()               {}
func (*RetargetRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RetargetRequest) GetImpInfo() *common.Imp {
	if m != nil {
		return m.ImpInfo
	}
	return nil
}

func (m *RetargetRequest) GetTrafficInfo() *common.TrafficInfo {
	if m != nil {
		return m.TrafficInfo
	}
	return nil
}

func (m *RetargetRequest) GetDeviceInfo() *common.Device {
	if m != nil {
		return m.DeviceInfo
	}
	return nil
}

func (m *RetargetRequest) GetUserInfo() *common.User {
	if m != nil {
		return m.UserInfo
	}
	return nil
}

func (m *RetargetRequest) GetBlackWhiteInfo() *common.BlackWhiteInfo {
	if m != nil {
		return m.BlackWhiteInfo
	}
	return nil
}

func (m *RetargetRequest) GetRegs() *common.Regs {
	if m != nil {
		return m.Regs
	}
	return nil
}

func (m *RetargetRequest) GetCampaignList() []*CampaignInfo {
	if m != nil {
		return m.CampaignList
	}
	return nil
}

func (m *RetargetRequest) GetAbKey() map[string]string {
	if m != nil {
		return m.AbKey
	}
	return nil
}

type CampaignInfo struct {
	// 1.
	// [Juno]
	// [Polaris]
	// [Rank]
	CampaignId int64 `protobuf:"varint,1,opt,name=campaign_id,json=campaignId" json:"campaign_id,omitempty"`
	// 2.人群包定向Id
	// [Rank]
	AudienceId int64 `protobuf:"varint,2,opt,name=audience_id,json=audienceId" json:"audience_id,omitempty"`
	// 3.重定向类型
	// [Rank]
	RetargetType int32 `protobuf:"varint,3,opt,name=retarget_type,json=retargetType" json:"retarget_type,omitempty"`
}

func (m *CampaignInfo) Reset()                    { *m = CampaignInfo{} }
func (m *CampaignInfo) String() string            { return proto.CompactTextString(m) }
func (*CampaignInfo) ProtoMessage()               {}
func (*CampaignInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

// **************************** RetargetResponse 输入协议 ***************************
// defines the dsp Retarget response structure
type RetargetResponse struct {
	// 1.rt-dsp处理结果枚举值。
	Status int32 `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	// 2.rt-dsp处理错误消息
	Msg string `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
	// 3.rt-dsp内部模块耗时统计
	Time int64 `protobuf:"varint,3,opt,name=time" json:"time,omitempty"`
	// 4.算法信息
	RankInfo *RankInfo `protobuf:"bytes,4,opt,name=rank_info,json=rankInfo" json:"rank_info,omitempty"`
	// 5.大模版信息
	BigTemplate *common.BigTemplate `protobuf:"bytes,5,opt,name=big_template,json=bigTemplate" json:"big_template,omitempty"`
}

func (m *RetargetResponse) Reset()                    { *m = RetargetResponse{} }
func (m *RetargetResponse) String() string            { return proto.CompactTextString(m) }
func (*RetargetResponse) ProtoMessage()               {}
func (*RetargetResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RetargetResponse) GetRankInfo() *RankInfo {
	if m != nil {
		return m.RankInfo
	}
	return nil
}

func (m *RetargetResponse) GetBigTemplate() *common.BigTemplate {
	if m != nil {
		return m.BigTemplate
	}
	return nil
}

type RankInfo struct {
	// 1.算法策略，算法标签
	Strategy string `protobuf:"bytes,1,opt,name=strategy" json:"strategy,omitempty"`
	// 2.二级算法策略，算法标签
	SubStrategy string `protobuf:"bytes,2,opt,name=subStrategy" json:"subStrategy,omitempty"`
	// 3.算法出价
	Price float64 `protobuf:"fixed64,3,opt,name=price" json:"price,omitempty"`
	// 4.算法预估ecpm
	EcpmFloor float64 `protobuf:"fixed64,4,opt,name=ecpm_floor,json=ecpmFloor" json:"ecpm_floor,omitempty"`
	// 5.ivr_model
	Ivr float64 `protobuf:"fixed64,5,opt,name=ivr" json:"ivr,omitempty"`
}

func (m *RankInfo) Reset()                    { *m = RankInfo{} }
func (m *RankInfo) String() string            { return proto.CompactTextString(m) }
func (*RankInfo) ProtoMessage()               {}
func (*RankInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func init() {
	proto.RegisterType((*RetargetRequest)(nil), "retarget_dsp.RetargetRequest")
	proto.RegisterType((*CampaignInfo)(nil), "retarget_dsp.CampaignInfo")
	proto.RegisterType((*RetargetResponse)(nil), "retarget_dsp.RetargetResponse")
	proto.RegisterType((*RankInfo)(nil), "retarget_dsp.RankInfo")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for RetargetService service

type RetargetServiceClient interface {
	RetargetDsp(ctx context.Context, in *RetargetRequest, opts ...grpc.CallOption) (*RetargetResponse, error)
}

type retargetServiceClient struct {
	cc *grpc.ClientConn
}

func NewRetargetServiceClient(cc *grpc.ClientConn) RetargetServiceClient {
	return &retargetServiceClient{cc}
}

func (c *retargetServiceClient) RetargetDsp(ctx context.Context, in *RetargetRequest, opts ...grpc.CallOption) (*RetargetResponse, error) {
	out := new(RetargetResponse)
	err := grpc.Invoke(ctx, "/retarget_dsp.RetargetService/RetargetDsp", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for RetargetService service

type RetargetServiceServer interface {
	RetargetDsp(context.Context, *RetargetRequest) (*RetargetResponse, error)
}

func RegisterRetargetServiceServer(s *grpc.Server, srv RetargetServiceServer) {
	s.RegisterService(&_RetargetService_serviceDesc, srv)
}

func _RetargetService_RetargetDsp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RetargetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RetargetServiceServer).RetargetDsp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/retarget_dsp.RetargetService/RetargetDsp",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RetargetServiceServer).RetargetDsp(ctx, req.(*RetargetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RetargetService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "retarget_dsp.RetargetService",
	HandlerType: (*RetargetServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RetargetDsp",
			Handler:    _RetargetService_RetargetDsp_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("retarget_dsp/retargetingDsp.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 697 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x54, 0x5f, 0x6f, 0xd3, 0x30,
	0x10, 0x57, 0xd6, 0xa6, 0x4b, 0x2f, 0x6d, 0x57, 0x79, 0xa8, 0x0a, 0x95, 0x06, 0xa3, 0x48, 0x68,
	0xbc, 0x74, 0xd2, 0x26, 0xa1, 0x89, 0x17, 0x60, 0x0c, 0xa4, 0x89, 0x81, 0x90, 0x37, 0xc4, 0x63,
	0xe4, 0x24, 0x6e, 0x66, 0xad, 0xf9, 0x83, 0xed, 0x8c, 0xf5, 0x53, 0xf0, 0x89, 0x78, 0xe4, 0x7b,
	0xe1, 0x3f, 0x49, 0xda, 0x4d, 0x82, 0xb7, 0xbb, 0xfb, 0xfd, 0xee, 0xfc, 0xf3, 0xf9, 0xce, 0xf0,
	0x8c, 0x53, 0x49, 0x78, 0x4a, 0x65, 0x98, 0x88, 0xf2, 0xb0, 0x71, 0x58, 0x9e, 0x9e, 0x89, 0x72,
	0x5e, 0xf2, 0x42, 0x16, 0x68, 0xb0, 0x49, 0x99, 0x42, 0x44, 0x04, 0xb5, 0xc8, 0xec, 0x8f, 0x0b,
	0x3b, 0xb8, 0x06, 0x31, 0xfd, 0x51, 0x51, 0x21, 0xd1, 0x1e, 0x00, 0xb7, 0x66, 0xc8, 0x92, 0xc0,
	0xd9, 0x77, 0x0e, 0xfa, 0xb8, 0x5f, 0x47, 0xce, 0x13, 0x34, 0x05, 0x8f, 0xde, 0xc5, 0xd7, 0x24,
	0x4f, 0x69, 0xb0, 0x65, 0xc0, 0xd6, 0x47, 0x2f, 0xc0, 0x63, 0x59, 0x19, 0xb2, 0x7c, 0x51, 0x04,
	0x5d, 0x85, 0xf9, 0x47, 0xfe, 0x3c, 0x2e, 0xb2, 0xac, 0xc8, 0xe7, 0xe7, 0x59, 0x89, 0xb7, 0x15,
	0x78, 0xae, 0x30, 0xf4, 0x0a, 0x06, 0x92, 0x93, 0xc5, 0x82, 0xc5, 0x96, 0xeb, 0x1a, 0xee, 0x6e,
	0xc3, 0xbd, 0xb2, 0x98, 0xa6, 0x62, 0x5f, 0xae, 0x1d, 0x74, 0x08, 0x7e, 0x42, 0x6f, 0x59, 0x4c,
	0x6d, 0x5a, 0xcf, 0xa4, 0x8d, 0x9a, 0xb4, 0x33, 0x03, 0x61, 0xb0, 0x14, 0x93, 0xf0, 0x12, 0xfa,
	0x95, 0xa0, 0xdc, 0xd2, 0xb7, 0x0d, 0x7d, 0xd0, 0xd0, 0xbf, 0x29, 0x00, 0x7b, 0x1a, 0x36, 0xd4,
	0xb7, 0x30, 0x8e, 0x96, 0x24, 0xbe, 0x09, 0x7f, 0x5e, 0x33, 0x59, 0x1f, 0xe0, 0x99, 0x8c, 0x49,
	0x93, 0x71, 0xaa, 0xf1, 0xef, 0x1a, 0x36, 0xd2, 0x46, 0xd1, 0x3d, 0x1f, 0xed, 0x43, 0x97, 0xd3,
	0x54, 0x04, 0xfd, 0xfb, 0xe7, 0x60, 0x15, 0xc3, 0x06, 0x51, 0xfd, 0xd9, 0x91, 0x71, 0x58, 0xaa,
	0x2c, 0x92, 0xd2, 0x30, 0x27, 0x19, 0x0d, 0x60, 0xbf, 0xa3, 0x5a, 0x38, 0x94, 0xf1, 0x57, 0x1b,
	0xfd, 0xa2, 0x82, 0xe8, 0x0d, 0x0c, 0x63, 0x92, 0x95, 0x84, 0xa5, 0x79, 0xb8, 0x64, 0x42, 0x06,
	0xbe, 0x62, 0xf9, 0x47, 0xd3, 0xf9, 0xe6, 0x43, 0xce, 0xdf, 0xd7, 0x14, 0x23, 0x66, 0xd0, 0x24,
	0x5c, 0x28, 0x3e, 0x1a, 0xc1, 0x16, 0x91, 0xc1, 0x40, 0x09, 0x71, 0xb1, 0xb2, 0xd0, 0x2e, 0xb8,
	0x32, 0xcc, 0xc8, 0x5d, 0x30, 0x34, 0xa1, 0xae, 0xfc, 0x4c, 0xee, 0xd4, 0x29, 0x3d, 0x12, 0x85,
	0x37, 0x74, 0x15, 0x8c, 0x4c, 0xf9, 0x83, 0xfb, 0xe5, 0x1f, 0xcc, 0xc5, 0xfc, 0x5d, 0xf4, 0x89,
	0xae, 0x3e, 0xe4, 0x92, 0xaf, 0xb0, 0x4b, 0xb4, 0x8d, 0x8e, 0x61, 0x92, 0xd0, 0xa8, 0x4a, 0xc3,
	0x56, 0x2c, 0x4b, 0xac, 0xde, 0x1d, 0x55, 0xb0, 0x83, 0x77, 0x0d, 0xda, 0xca, 0x4c, 0x8c, 0xb4,
	0xc7, 0xe0, 0x15, 0x25, 0x55, 0xf7, 0x2a, 0xd2, 0x60, 0xac, 0xd4, 0x78, 0x78, 0x5b, 0xfb, 0x17,
	0x45, 0x3a, 0x3d, 0x01, 0x58, 0x1f, 0x82, 0xc6, 0xd0, 0xd1, 0xda, 0xec, 0x00, 0x6a, 0x13, 0x3d,
	0x02, 0xf7, 0x96, 0x2c, 0xab, 0x66, 0xee, 0xac, 0xf3, 0x7a, 0xeb, 0xc4, 0x99, 0x55, 0x30, 0xd8,
	0xec, 0x06, 0x7a, 0x0a, 0xfe, 0x86, 0x26, 0x53, 0xa3, 0x83, 0x21, 0x6e, 0x95, 0x68, 0x02, 0xa9,
	0x12, 0x46, 0x73, 0x3d, 0x4b, 0x89, 0x29, 0xa8, 0x08, 0x4d, 0x48, 0x11, 0x9e, 0xc3, 0xb0, 0xed,
	0x86, 0x5c, 0x95, 0x34, 0xe8, 0x98, 0xce, 0xb5, 0xab, 0x74, 0xa5, 0x62, 0xb3, 0xdf, 0x0e, 0x8c,
	0xd7, 0x6d, 0x12, 0x65, 0x91, 0x0b, 0x8a, 0x26, 0xd0, 0x13, 0x92, 0xc8, 0x4a, 0x98, 0x63, 0x5d,
	0x5c, 0x7b, 0xfa, 0x3e, 0x99, 0x48, 0x6b, 0xed, 0xda, 0x44, 0x08, 0xba, 0x92, 0x65, 0xb6, 0x74,
	0x07, 0x1b, 0x5b, 0xf5, 0xb4, 0xcf, 0x49, 0x7e, 0xb3, 0xb9, 0x43, 0x93, 0x07, 0xef, 0xa2, 0x60,
	0xf3, 0xe4, 0x1e, 0xaf, 0x2d, 0xbd, 0x4f, 0x11, 0x4b, 0x43, 0x49, 0xb3, 0x72, 0x49, 0x24, 0x7d,
	0xb8, 0x4f, 0xa7, 0x2c, 0xbd, 0xaa, 0x21, 0xec, 0x47, 0x6b, 0x67, 0xf6, 0xcb, 0x01, 0xaf, 0x29,
	0xa7, 0x17, 0x5b, 0xa8, 0x65, 0x93, 0x34, 0x6d, 0x9a, 0xde, 0xfa, 0x6a, 0xb4, 0x7d, 0x51, 0x45,
	0x97, 0x0d, 0x6c, 0xef, 0xb0, 0x19, 0xd2, 0x6f, 0x53, 0x72, 0xb5, 0x76, 0xe6, 0x32, 0x0e, 0xb6,
	0x8e, 0xfe, 0x4b, 0x68, 0x5c, 0x66, 0xe1, 0x62, 0x59, 0x14, 0xdc, 0x5c, 0xc7, 0xc1, 0x7d, 0x1d,
	0xf9, 0xa8, 0x03, 0xba, 0x25, 0xec, 0x96, 0x1b, 0xb9, 0x0e, 0xd6, 0xe6, 0x51, 0xb8, 0xfe, 0x8f,
	0x2e, 0x29, 0xd7, 0x7b, 0x8c, 0x2e, 0xc0, 0x6f, 0x42, 0xea, 0x4b, 0x43, 0x7b, 0xff, 0x9d, 0xd2,
	0xe9, 0x93, 0x7f, 0xc1, 0xf6, 0x75, 0xa2, 0x9e, 0xf9, 0xf8, 0x8e, 0xff, 0x06, 0x00, 0x00, 0xff,
	0xff, 0xbc, 0x2c, 0xcd, 0x2b, 0x37, 0x05, 0x00, 0x00,
}
```


## 三、编写服务端代码


### 3.1 Service

RetargetServer 实现了 RetargetDsp 接口，实现了 retarget_dsp.proto 协议的业务功能，
在 go-kit 框架中，服务端接收到请求后，RetargetDsp 应将数据透传给 go-kit 的 rpc 方法 RetargetHandler.ServeGRPC，go-kit 编排处理流程，在流程中调用初始化时注册的 endpoint.Endpoint 方法。

```

type RetargetServer struct {
	listen net.Listener
	server *grpc.Server

	RetargetHandler grpc_transport.Handler
}

//NewRpcServer
//创建基于 go-kit 的 rpc 服务，MakeRetargetDspEndpoint 是主要业务逻辑实现。
//在 go-kit 推荐使用中，decodeGRPCRequest 和 encodeGRPCResponse 两部分应承担请求到 rtContext 和 rtContext 到响应的两个 pipeline。这里未实现。
func NewRpcServer() (*RetargetServer, error) {
	var err error

	//build Transport
	service := new(RetargetServer)
	service.RetargetHandler = grpc_transport.NewServer(
		//增加速率限制，每50毫秒补充一次，设置容量12
		NewTokenBucketLimiterWithBuildIn(rate.NewLimiter(
			rate.Every(time.Millisecond*time.Duration(config.RetargetCfg.ServerConfig.TokenBucketLimit)),
			config.RetargetCfg.ServerConfig.TokenBucketBurst),
		)(MakeRetargetDspEndpoint(service)),
		decodeGRPCRequest,
		encodeGRPCResponse,
	)

	service.server = grpc.NewServer()
	retarget_dsp.RegisterRetargetServiceServer(service.server, service)

	service.listen, err = net.Listen("tcp", ":"+config.RetargetCfg.ServerConfig.RpcPort)

	return service, err
}

//RetargetDsp
//服务注册方法名，实现【重定向】单子的检索和算法优选，具体实现参考内部 pipeline
//使用 go-lit 推荐将其参数透传给 ServeGRPC
func (rs *RetargetServer) RetargetDsp(ctx context.Context, request *retarget_dsp.RetargetRequest) (*retarget_dsp.RetargetResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RetargetServer RetargetDsp error:", err)
			rtmetrics.SetMetrics(rtmetrics.Panic, rtmetrics.Labels{FunctionName: "RetargetDsp"}, 1)
			debug.PrintStack()
		}
	}()

	//RetargetDsp 业务不处理空请求
	if request == nil {
		return &retarget_dsp.RetargetResponse{
			Msg:    reference.ServerStats_REQUEST_ERROR,
			Status: reference.ServerStatsCode(reference.ServerStats_REQUEST_ERROR).Enum(),
		}, errors.New("RetargetServer request is nil")
	}

	//go-kit 的 rpc 服务函数，执行 decodeGRPCRequest、MakeRetargetDspEndpoint、encodeGRPCResponse 流程
	_, response, err := rs.RetargetHandler.ServeGRPC(ctx, request)
	if err != nil {
		return &retarget_dsp.RetargetResponse{
			Msg:    reference.ServerStats_Write_Response_ERROR,
			Status: reference.ServerStatsCode(reference.ServerStats_Write_Response_ERROR).Enum(),
		}, err
	}

	return response.(*retarget_dsp.RetargetResponse), err
}

func decodeGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func encodeGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response, nil
}

// [2]中间件：限流服务
var ErrLimitExceed = errors.New("limit exceed error")

// NewTokenBucketLimiterWithBuildIn 使用x/time/rate创建限流中间件
func NewTokenBucketLimiterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				rtmetrics.SetMetrics(rtmetrics.ConcurrencyFilter, rtmetrics.Labels{FunctionName: "TokenBucketLimiter"}, 1)
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}
```


### 3.2 Endpoint

在gokit中Endpoint是可以包装到http.Handler中的特殊方法，gokit采用装饰着模式，把Service应该执行的逻辑封装到Endpoint方法中执行。Endpoint的作用是：调用Service中相应的方法处理请求对象（RetargetRequest），返回响应对象（RetargetResponse）。

```

//MakeRetargetDspEndpoint
//主要业务逻辑实现代码，触发 Rt Serve Pipeline，记录请求日志
func MakeRetargetDspEndpoint(srv retarget_dsp.RetargetServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (response interface{}, err error) {
		request := req.(*retarget_dsp.RetargetRequest)

		//pipeline
		rtContext := &pipeline.RetargetContext{
			RetargetRequest:  request,
			Condition:        new(pipeline.Condition),
			RetargetResponse: new(retarget_dsp.RetargetResponse),
		}
		polarisPipeline := rtContext.GetServePipeline(request)
		err = polarisPipeline.ProcessWithTime(rtContext)
		totalTime := polarisPipeline.TotalTime()
		rtContext.RetargetResponse.Time = totalTime

		//write log
		rtContext.WriteRequestLog(rtContext)

		//write metrics
		labels := rtmetrics.Labels{FunctionName: "RetargetDsp", Adx: rtContext.Condition.Adx, AdType: rtContext.Condition.AdType, ModelStatus: rtContext.RetargetResponse.Msg}
		rtmetrics.SetMetrics(rtmetrics.RetargetDspTime, labels, float64(totalTime))

		return rtContext.RetargetResponse, err
	}
}
```


### 3.3 Transport

Transport层用于接收用户网络请求并将其转为Endpoint可以处理的对象，然后交由Endpoint执行，最后将处理结果转为响应对象向用户响应。

为了完成这项工作，Transport需要具备两个工具方法：

- 解码器：把用户的请求内容转换为请求对象（PolarisRequest）；
- 编码器：把处理结果转换为响应对象（PolarisResponse）；
- 

```
func decodeRequest(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}

func encodeResponse(_ context.Context, req interface{}) (interface{}, error) {
	return req, nil
}
```


### 3.4 main代码

main.go

```
package retarget_service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"gitlab.mobvista.com/mvdsp/protoc_new/pkg/retarget_dsp"
	"gitlab.mobvista.com/retargetting-dsp/dsp-retarget/internal/app/retarget_service/pipeline"
	"gitlab.mobvista.com/retargetting-dsp/dsp-retarget/internal/pkg/config"
	"gitlab.mobvista.com/retargetting-dsp/dsp-retarget/internal/pkg/reference"
	"gitlab.mobvista.com/retargetting-dsp/dsp-retarget/internal/pkg/rtmetrics"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"net"
	"runtime/debug"
	"time"
)

type RetargetServer struct {
	listen net.Listener
	server *grpc.Server

	RetargetHandler grpc_transport.Handler
}

//NewRpcServer
//创建基于 go-kit 的 rpc 服务，MakeRetargetDspEndpoint 是主要业务逻辑实现。
//在 go-kit 推荐使用中，decodeGRPCRequest 和 encodeGRPCResponse 两部分应承担请求到 rtContext 和 rtContext 到响应的两个 pipeline。这里未实现。
func NewRpcServer() (*RetargetServer, error) {
	var err error

	//build Transport
	service := new(RetargetServer)
	service.RetargetHandler = grpc_transport.NewServer(
		//增加速率限制，每50毫秒补充一次，设置容量12
		NewTokenBucketLimiterWithBuildIn(rate.NewLimiter(
			rate.Every(time.Millisecond*time.Duration(config.RetargetCfg.ServerConfig.TokenBucketLimit)),
			config.RetargetCfg.ServerConfig.TokenBucketBurst),
		)(MakeRetargetDspEndpoint(service)),
		decodeGRPCRequest,
		encodeGRPCResponse,
	)

	service.server = grpc.NewServer()
	retarget_dsp.RegisterRetargetServiceServer(service.server, service)

	service.listen, err = net.Listen("tcp", ":"+config.RetargetCfg.ServerConfig.RpcPort)

	return service, err
}

//Start
//启动 rpc 服务
func (rs *RetargetServer) Start(context.Context) error {
	if err := rs.server.Serve(rs.listen); err != nil {
		return fmt.Errorf("failed to serve: %s", err.Error())
	}
	return nil
}

//RetargetDsp
//服务注册方法名，实现【重定向】单子的检索和算法优选，具体实现参考内部 pipeline
//使用 go-lit 推荐将其参数透传给 ServeGRPC
func (rs *RetargetServer) RetargetDsp(ctx context.Context, request *retarget_dsp.RetargetRequest) (*retarget_dsp.RetargetResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("RetargetServer RetargetDsp error:", err)
			rtmetrics.SetMetrics(rtmetrics.Panic, rtmetrics.Labels{FunctionName: "RetargetDsp"}, 1)
			debug.PrintStack()
		}
	}()

	//RetargetDsp 业务不处理空请求
	if request == nil {
		return &retarget_dsp.RetargetResponse{
			Msg:    reference.ServerStats_REQUEST_ERROR,
			Status: reference.ServerStatsCode(reference.ServerStats_REQUEST_ERROR).Enum(),
		}, errors.New("RetargetServer request is nil")
	}

	//go-kit 的 rpc 服务函数，执行 decodeGRPCRequest、MakeRetargetDspEndpoint、encodeGRPCResponse 流程
	_, response, err := rs.RetargetHandler.ServeGRPC(ctx, request)
	if err != nil {
		return &retarget_dsp.RetargetResponse{
			Msg:    reference.ServerStats_Write_Response_ERROR,
			Status: reference.ServerStatsCode(reference.ServerStats_Write_Response_ERROR).Enum(),
		}, err
	}

	return response.(*retarget_dsp.RetargetResponse), err
}

func decodeGRPCRequest(_ context.Context, request interface{}) (interface{}, error) {
	return request, nil
}

func encodeGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	return response, nil
}

// [2]中间件：限流服务
var ErrLimitExceed = errors.New("limit exceed error")

// NewTokenBucketLimiterWithBuildIn 使用x/time/rate创建限流中间件
func NewTokenBucketLimiterWithBuildIn(bkt *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !bkt.Allow() {
				rtmetrics.SetMetrics(rtmetrics.ConcurrencyFilter, rtmetrics.Labels{FunctionName: "TokenBucketLimiter"}, 1)
				return nil, ErrLimitExceed
			}
			return next(ctx, request)
		}
	}
}

//MakeRetargetDspEndpoint
//主要业务逻辑实现代码，触发 Rt Serve Pipeline，记录请求日志
func MakeRetargetDspEndpoint(srv retarget_dsp.RetargetServiceServer) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (response interface{}, err error) {
		request := req.(*retarget_dsp.RetargetRequest)

		//pipeline
		rtContext := &pipeline.RetargetContext{
			RetargetRequest:  request,
			Condition:        new(pipeline.Condition),
			RetargetResponse: new(retarget_dsp.RetargetResponse),
		}
		polarisPipeline := rtContext.GetServePipeline(request)
		err = polarisPipeline.ProcessWithTime(rtContext)
		totalTime := polarisPipeline.TotalTime()
		rtContext.RetargetResponse.Time = totalTime

		//write log
		rtContext.WriteRequestLog(rtContext)

		//write metrics
		labels := rtmetrics.Labels{FunctionName: "RetargetDsp", Adx: rtContext.Condition.Adx, AdType: rtContext.Condition.AdType, ModelStatus: rtContext.RetargetResponse.Msg}
		rtmetrics.SetMetrics(rtmetrics.RetargetDspTime, labels, float64(totalTime))

		return rtContext.RetargetResponse, err
	}
}

```


## 四、编写客户端代码

客户端调用 retarget_dsp.proto 服务时，通过协议生成的 newClient 方法和 RetargetDsp 方法，完成具体的远程过程调用。

```
rtClient := retarget_dsp.NewRetargetServiceClient(conn)
rtResp, err := rtClient.RetargetDsp(ctxRead, request)
```

源代码：

```
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.mobvista.com/mvdsp/protoc_new/pkg/retarget_dsp"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"path/filepath"
	"time"
)

func main() {

	var (
		addr = "localhost:9999"
		//addr = "34.197.85.30:9888"
		file = "test/example/testdata/retarget_request_sample.json"
	)

	//get request
	request := retarget_dsp.RetargetRequest{}
	path, _ := filepath.Abs(file)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("ioutil.ReadFile:", err)
		return
	}
	err = json.Unmarshal(data, &request)
	if err != nil {
		fmt.Println("json.Unmarshal:", err)
		return
	}

	for i := 0; i < 10; i++ {
		fmt.Println("----------------------- ", "request - ", i, " -----------------------")
		request.RequestId = bson.NewObjectId().Hex()
		QueryRetargetDsp(&request, addr)
	}

}

func QueryRetargetDsp(request *retarget_dsp.RetargetRequest, addr string) {
	ctxConn, celConn := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(200000))
	defer celConn()
	conn, err := grpc.DialContext(ctxConn, addr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("grpc.DialContext error:", err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	ctxRead, celRead := context.WithTimeout(context.Background(), time.Duration(1000000)*time.Millisecond)
	defer celRead()

	//run
	rtClient := retarget_dsp.NewRetargetServiceClient(conn)
	rtResp, err := rtClient.RetargetDsp(ctxRead, request)
	fmt.Println("request:", request)

	if err != nil {
		fmt.Println("rtClient.RetargetDsp error:", err)
	}

	output, _ := json.Marshal(rtResp)
	fmt.Println("resp:", string(output))
}

```

