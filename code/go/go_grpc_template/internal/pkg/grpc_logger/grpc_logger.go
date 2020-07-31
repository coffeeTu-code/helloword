package grpc_logger

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"helloword/code/go/go_grpc_template/pkg/log_helper"
)

var Grpc_logger RPCLog

type RPCLog struct {
	Request *logrus.Logger

	RunTime *zap.SugaredLogger
}

func (log *RPCLog) Init() (err error) {

	log.Request = log_helper.NewLogrusWithFileRotatelogs("./log/request.log")

	log.RunTime = log_helper.NewZapWithLumberjack("./log/runtime.log")

	return err
}

//****************** Request Log **********************
func (log *RPCLog) Req(fields logrus.Fields) {
	if log.Request != nil {
		log.Request.WithFields(fields).Info("")
	}
}

//****************** Runtime Log **********************
func (log *RPCLog) Debug(args ...interface{}) {
	if log.RunTime != nil {
		log.RunTime.Debug(args)
	}
}

func (log *RPCLog) Info(args ...interface{}) {
	if log.RunTime != nil {
		log.RunTime.Info(args)
	}
}

func (log *RPCLog) Warn(args ...interface{}) {
	if log.RunTime != nil {
		log.RunTime.Warn(args)
	}
}

func (log *RPCLog) Infof(format string, args ...interface{}) {
	if log.RunTime != nil {
		log.RunTime.Infof(format, args)
	}
}

func (log *RPCLog) Warnf(format string, args ...interface{}) {
	if log.RunTime != nil {
		log.RunTime.Warnf(format, args)
	}
}

func (log *RPCLog) Error(args ...interface{}) {
	if log.RunTime != nil {
		log.RunTime.Error(args)
	}
}

//****************** Flush **********************
func (log *RPCLog) Flush() {

	//zap底层 API 可以设置缓存，所以一般使用defer logger.Sync()将缓存同步到文件中。
	if log.RunTime != nil {
		log.RunTime.Sync()
	}

}
