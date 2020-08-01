package greeter_server

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

// ConsulRegister 方法
func ConsulRegister(consulAddress, consulPort string, advertisePort string, serverPort string) (registar sd.Registrar) {

	// 日志相关
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// 服务发现域。在本例中，我们使用 Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		consulConfig.Address = consulAddress + ":" + consulPort
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	port, _ := strconv.Atoi(serverPort)
	num := rand.Intn(100) // to make service ID unique
	asr := api.AgentServiceRegistration{
		ID:      "go-kit-srv-greeter-" + strconv.Itoa(num), //unique service ID
		Name:    "go-kit-srv-greeter",
		Address: "",
		Port:    port,
		Tags:    []string{"go-kit", "greeter"},
		Check: &api.AgentServiceCheck{
			HTTP:     "http://" + "" + ":" + advertisePort + "/health",
			Interval: "10s",
			Timeout:  "1s",
			Notes:    "Basic health checks",
		},
	}

	registar = consulsd.NewRegistrar(client, &asr, logger)
	return
}
