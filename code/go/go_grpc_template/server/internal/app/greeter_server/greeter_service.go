package greeter_server

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Service 描述了 greetings 这个服务
type Service interface {
	Health() bool
	Greeting(name string) string
}

// GreeterService  是 Service 接口的实现
type GreeterService struct{}

// Service 的 Health 接口实现
func (GreeterService) Health() bool {
	return true
}

// Service 的 Greeting 接口实现
func (GreeterService) Greeting(name string) (greeting string) {
	greeting = "GO-GRPC-TEMPLATE Hello " + name
	return
}

// grpc Service 健康检查实现
type Health struct{}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *Health) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (h *Health) Watch(req *grpc_health_v1.HealthCheckRequest, watch grpc_health_v1.Health_WatchServer) error {
	return nil
}
