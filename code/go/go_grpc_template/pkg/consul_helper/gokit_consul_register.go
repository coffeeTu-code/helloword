package consul_helper

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

// 注册方法
func NewConsulSdRegister(consulAddress, consulPort string, serviceName, serverPort string, tags []string, logger log.Logger) (registar sd.Registrar, err error) {

	rand.Seed(time.Now().UTC().UnixNano())

	// 服务发现域。在本例中，我们使用 Consul.
	client, err := consulClient(consulAddress, consulPort)
	if err != nil {
		return nil, err
	}
	sdClient := consulsd.NewClient(client)

	IP := localIP()
	port, _ := strconv.Atoi(serverPort)
	asr := api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v", serviceName, rand.Intn(100)), // unique service ID, 服务节点的名称
		Name:    serviceName,                                       // 服务名称
		Address: IP,                                                // 服务 IP
		Port:    port,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			//HTTP:                           fmt.Sprintf("http://%v:%v/health", IP, advertisePort),//
			GRPC:                           fmt.Sprintf("%v:%v/%v", IP, serverPort, serviceName), //需要实现 grpc_health_v1 的健康检查方法
			Interval:                       "10s",                                                // 健康检查间隔
			Timeout:                        "1s",                                                 //
			Notes:                          "Basic health checks",                                //
			DeregisterCriticalServiceAfter: (time.Duration(1) * time.Minute).String(),            // 注销时间，相当于过期时间
		},
	}

	registar = consulsd.NewRegistrar(sdClient, &asr, logger)
	return registar, nil
}

