package pipeline

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"helloword/code/go/go_grpc_template/server/internal/pkg/grpc_logger"
)

//快速的 json 库，替代原生 json 编解码
var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ServerContext struct {
}

//GetServePipeline
func (ctx *ServerContext) GetServePipeline() *WallTimePipeline {
	return &WallTimePipeline{
		Name: "Greeter Server",
		Filters: []Filter{

		},
	}
}

//WriteRequestLog
//记录请求日志，`\t`分割，字段顺序参考 rtlogger.go logger.SetFormatter(&RequestLogFormatter{}) 部分处理逻辑
//函数中填充的 fields 的 key 需要与注册 logger formatter 时保持一致，否则会打乱日志字段顺序
func (ctx *ServerContext) WriteRequestLog() {
	if ctx == nil {
		return
	}

	var (
		fields = logrus.Fields{}
	)

	//add fields with ServerContext

	grpc_logger.Grpc_logger.RequestInfo(fields)
}
