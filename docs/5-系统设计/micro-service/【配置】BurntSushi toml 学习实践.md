> 原文链接：  
> [toml格式配置文件详解 https://mojotv.cn/2018/12/26/what-is-toml](https://mojotv.cn/2018/12/26/what-is-toml)




[TOC]


# BurntSushi toml


## toml文件格式
```
#
#   Retarget Dsp 系统配置
#

#dsp_retarget.toml 所在路径 pwd
#ConfPath = "/data/recommend/advanced_search/conf/polaris_conf/"
ConfPath = "/Users/coffee/go/src/polaris/config/"
#ConfPath = "/home/mobdev/xuefeng.han/polaris/config/"

#
#   ServerConfig
#   系统启动web服务配置（支持http、grpc两种方式），consul注册配置
#
[ServerConfig]

ServerType = "http"
Port = "9999"

#
#   LogConfig
#   系统日志XML格式配置[runtime、request]，debug模式开关
#
[LogConfig]

RequestLogOutputPath = "./log/polaris_request.log"

RunTimeLogPath = "run_log_config.json"

#{true=debug模式，false=非debug模式。置成true等于请求test为true}
RunTimeDebug = false

#
#   JunoConfig
#   juno rpc
#
[JunoConfig]
Addr = ""

UseConsul = true
ConsulAddr = "vg-consul-internl-ali.mobvista.com:8500"
ServiceRemote = "juno_server"
MyService = "retarget_dsp_server"
Interval = 1
ServiceRatio = 3.0
CpuThreshold = 60.0

ConnectionTimeout = 100
ReadTimeout = 100

#
#   JunoConfig
#   juno rpc
#
[PolarisConfig]
LibConfigPath = "../../configs/service/polaris_dsp.toml"

ReadTimeout = 80
```


## go struct

```
//系统配置
type RetargetConfig struct {
	ConfPath      string       //全局配置文件绝对路径
	ServerConfig  ServerConfig //系统启动web服务配置（支持http、grpc两种方式），consul注册配置
	LogConfig     LogConfig    //系统日志XML格式配置[runtime、request]，debug模式开关
	JunoConfig    ConsulConfig
	PolarisConfig LibConfig
}

type ServerConfig struct {
	ServerType string //web服务类型，[http，grpc]
	Port       string //web服务端口
}

type LogConfig struct {
	RequestLogOutputPath string
	RunTimeLogPath       string
	RunTimeDebug         bool //收集debug信息，系统调试。
}

type ConsulConfig struct {
	//线上部署，consul模式
	UseConsul     bool
	ConsulAddr    string
	ServiceRemote string
	MyService     string
	Interval      int
	ServiceRatio  float64
	CpuThreshold  float64

	//单机测试用，非consul模式
	Addr string

	ConnectionTimeout int //连接超时时间，单位ms
	ReadTimeout       int //读超时，单位ms
}

type LibConfig struct {
	LibConfigPath string
	ReadTimeout   int //读超时，单位ms
}
```


## 解析代码
```

func (rtCfg *RetargetConfig) LoadConfig(confPath string) (err error) {
	//读取配置文件
	_, err = toml.DecodeFile(confPath, rtCfg)
	if err != nil {
		return err
	}

	return nil
}

```