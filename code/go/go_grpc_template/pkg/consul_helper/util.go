package consul_helper

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
)

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknow"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}

func consulClient(consulAddress, consulPort string) (*api.Client, error) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = fmt.Sprintf("%v:%v", consulAddress, consulPort)
	client, err := api.NewClient(consulConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}
