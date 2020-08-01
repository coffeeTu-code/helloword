package greeter_server

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	greeter "helloword/code/go/go_grpc_template/api/greeter/go"
)

type (
	// Endpoints 收集构成 Greeting Service 的所有端点。
	// 它应该被用作助手 struct ，将所有端点收集到一个参数中。
	Endpoints struct {
		// Consul 用这个端点做健康检查(Health Check)
		HealthEndpoint   endpoint.Endpoint
		GreetingEndpoint endpoint.Endpoint
	}

	HealthRequest struct {
	}

	HealthResponse struct {
		Healthy bool  `json:"health,omitempty"`
		Err     error `json:"err,omitempty"`
	}

	// Failer是一个应该由响应类型实现的接口。
	// 响应编码器可以检查响应是否是一个 Failer。
	// 如果响应是一个 Failer, 那么就说明响应已经失败了，根据错误使用单独的写入路径对响应进行编码。
	Failer interface {
		Failed() error
	}
)

// MakeServiceEndpoints 返回服务端点, 将所有已提供的中间件连接起来
func MakeServerEndpoints(s Service) Endpoints {
	var healthEndpoint endpoint.Endpoint
	{
		healthEndpoint = MakeHealthEndpoint(s)
	}

	var greetingEndpoint endpoint.Endpoint
	{
		greetingEndpoint = MakeGreetingEndpoint(s)
	}

	return Endpoints{
		HealthEndpoint:   healthEndpoint,
		GreetingEndpoint: greetingEndpoint,
	}
}

// MakeHealthEndpoints 构造一个 Health 端点，将服务包装为一个端点
func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		service := new(GreeterService)
		return HealthResponse{Healthy: service.Health()}, nil
	}
}

// MakeGreetingEndpoints 构造一个 Greeter 端点，将 Greeting 服务包装为一个端点
func MakeGreetingEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*greeter.GreetingRequest)
		resp := new(greeter.GreetingResponse)
		service := new(GreeterService)
		resp.Greeting = service.Greeting(req.GetName())
		return resp, nil
	}
}
