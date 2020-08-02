package greeter_server

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"helloword/code/go/go_grpc_template/pkg/consul_helper"
	"os"
)

// ConsulRegister 方法
func ConsulRegister(serverPort string) (registar sd.Registrar) {

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

	registar, err := consul_helper.NewConsulSdRegister(consulAddress, consulPort, serviceName, serverPort, []string{"go-kit", "greeter"}, logger)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	return registar
}
