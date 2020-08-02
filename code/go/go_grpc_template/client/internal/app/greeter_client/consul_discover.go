package greeter_client

import (
	"github.com/go-kit/kit/log"
	"helloword/code/go/go_grpc_template/pkg/consul_helper"
	"os"
)

func ConsulDiscover() *consul_helper.ConsulResolver {

	// 日志相关
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var (
		consulAddress = ""
		consulPort    = "8500"
		serviceName   = "go-kit-srv-greeter"
	)
	resolver, err := consul_helper.NewConsulResolver(consulAddress, consulPort, serviceName, "go-kit-client", consul_helper.SetInterval(10))
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	return resolver
}
