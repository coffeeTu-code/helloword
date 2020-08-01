package greeter_server

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
