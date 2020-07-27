package log_helper

import (
	"encoding/json"
	"errors"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//**********************************************

var defaultZapLogConfig = `
{
  "zap_config": {
    "level": "debug",
    "development": false,
    "disableCaller": false,
    "disableStacktrace": false,
    "sampling": {
      "initial": 100,
      "thereafter": 100
    },
    "encoding": "console",
    "encoderConfig": {
      "messageKey": "msg",
      "levelKey": "level",
      "timeKey": "ts",
      "nameKey": "logger",
      "callerKey": "caller",
      "stacktraceKey": "stacktrace",
      "lineEnding": "\n",
      "levelEncoder": "capitalColor",
      "timeEncoder": "rfc3339",
      "durationEncoder": "seconds",
      "callerEncoder": "short",
      "nameEncoder": ""
    },
    "outputPaths": [
      "./log/retarget_runtime.log"
    ],
    "errorOutputPaths": [
      "stderr"
    ]
  },
  "lumberjack_config": {
    "@filename": "日志文件的位置",
    "filename": "./log/retarget_runtime.log",
    "@maxsize": "在进行切割之前，日志文件的最大大小（以MB为单位）",
    "maxsize": 100,
    "@maxbackups": "保留旧文件的最大个数",
    "maxbackups": 10,
    "@maxage": "保留旧文件的最大天数",
    "maxage": 3,
    "@compress": "是否压缩/归档旧文件",
    "compress": false
  },
  "option": {
    "caller_skip": 1
  }
}
`

type ZapWithLumberjackConfig struct {
	ZapConfig        zap.Config        `json:"zap_config"`
	LumberjackConfig lumberjack.Logger `json:"lumberjack_config"`
	Option           struct {
		CallerSkip int `json:"caller_skip"`
	} `json:"option"`
}

func NewDefaultZapWithLumberjack() (*zap.SugaredLogger, error) {
	var cfg ZapWithLumberjackConfig
	if err := json.Unmarshal([]byte(defaultZapLogConfig), &cfg); err != nil {
		return nil, err
	}
	return NewZapWithLumberjack(cfg)
}

func NewZapWithLumberjack(cfg ZapWithLumberjackConfig) (*zap.SugaredLogger, error) {
	if len(cfg.ZapConfig.OutputPaths) == 0 {
		return nil, errors.New("not set cfg.OutputPaths")
	}

	writeSyncer := zapcore.AddSync(&cfg.LumberjackConfig)
	encoder := zapcore.NewConsoleEncoder(cfg.ZapConfig.EncoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(cfg.Option.CallerSkip))
	sugarLogger := logger.Sugar()

	return sugarLogger, nil

	//zap提供了几个快速创建logger的方法，zap.NewExample()、zap.NewDevelopment()、zap.NewProduction()，还有高度定制化的创建方法zap.New()。
	//创建前 3 个logger时，zap会使用一些预定义的设置，它们的使用场景也有所不同。
	//Example适合用在测试代码中，Development在开发环境中使用，Production用在生成环境。
	//logger := zap.NewExample()
	//logger, _ := zap.NewDevelopment()
	//logger, _ := zap.NewProduction()
	//logger, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	//if err != nil {
	//	return nil, err
	//}

	//当然，每个字段都用方法包一层用起来比较繁琐。zap也提供了便捷的方法SugarLogger，可以使用printf格式符的方式。
	//SugaredLogger的使用比Logger简单，只是性能比Logger低 50% 左右，可以用在非热点函数中。
}
