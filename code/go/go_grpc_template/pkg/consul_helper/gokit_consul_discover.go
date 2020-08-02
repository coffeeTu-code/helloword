package consul_helper

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
)

// 服务发现
func NewConsulSdInstance(consulAddress, consulPort string, serviceName string, tags []string, logger log.Logger) (instancer sd.Instancer, err error) {

	// 服务发现域。在本例中，我们使用 Consul.
	client, err := consulClient(consulAddress, consulPort)
	if err != nil {
		return nil, err
	}
	sdClient := consulsd.NewClient(client)

	//基于consul客户端、服务名称、服务标签等信息，
	// 创建consul的连接实例，
	// 可实时查询服务实例的状态信息
	instancer = consulsd.NewInstancer(sdClient, logger, serviceName, tags, true)
	return

}
