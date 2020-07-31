//******************************************************
//> 原文链接：
//> [如何开发高性能 go 组件？https://zhuanlan.zhihu.com/p/41991119](https://zhuanlan.zhihu.com/p/41991119)
//> [在Go语言项目中使用Zap日志库 https://zhuanlan.zhihu.com/p/88856378](https://zhuanlan.zhihu.com/p/88856378)
//> [Go 每日一库之 zap https://zhuanlan.zhihu.com/p/136093026](https://zhuanlan.zhihu.com/p/136093026)
//
//> GitHub：[https://github.com/uber-go/zap](https://github.com/uber-go/zap)
//> doc：[https://pkg.go.dev/go.uber.org/zap?tab=doc](https://pkg.go.dev/go.uber.org/zap?tab=doc)
//
//******************************************************

package log_helper

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//**********************************************

func NewZapWithLumberjack(outputPath string) *zap.SugaredLogger {
	if len(outputPath) == 0 {
		outputPath = "stdout"
	}

	writeSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   outputPath, //日志文件的位置
		MaxSize:    100,        //在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 10,         //保留旧文件的最大个数
		MaxAge:     3,          //保留旧文件的最大天数
		Compress:   false,      //是否压缩/归档旧文件
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	//zap提供了便捷的方法SugarLogger，可以使用printf格式符的方式。
	//SugaredLogger的使用比Logger简单，只是性能比Logger低 50% 左右，可以用在非热点函数中。
	sugarLogger := logger.Sugar()

	return sugarLogger
}

//zap提供了几个快速创建logger的方法，zap.NewExample()、zap.NewDevelopment()、zap.NewProduction()，还有高度定制化的创建方法zap.New()。
//创建前 3 个logger时，zap会使用一些预定义的设置，它们的使用场景也有所不同。
//Example适合用在测试代码中，Development在开发环境中使用，Production用在生成环境。

func NewZapExample() *zap.Logger {
	return zap.NewExample()
}

func NewZapDevelopment() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func NewZapProduction() (*zap.Logger, error) {
	return zap.NewProduction()
}
