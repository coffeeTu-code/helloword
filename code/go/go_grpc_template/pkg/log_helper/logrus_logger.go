package log_helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"sort"
	"time"
)

//******************************************************

var defaultLogrusLogConfig = `
{
  "logrus_config": {
    "outputPaths": "./log/retarget_request.log"
  },
  "lumberjack_config": {
    "@filename": "日志文件的位置",
    "filename": "./log/retarget_request.log",
    "@rotation_time": "日志切割时间间隔，单位：m",
    "rotation_time": 60,
    "@maxage": "保留旧文件的最大时常，单位：h",
    "maxage": 3
  },
  "option": {
    "formatter": "text",
    "fields_order": [
      
    ]
  }
}
`

type LogrusWithFileRotatelogsConfig struct {
	LogrusConfig         LogrusConfig         `json:"logrus_config"`
	FileRotatelogsConfig FileRotatelogsConfig `json:"file_rotatelogs_config"`
	Option               struct {
		Formatter   string   `json:"formatter"` // json | test | requestFormatter
		FieldsOrder []string `json:"fields_order"`
	}
}

type LogrusConfig struct {
	OutputPaths string `json:"outputPaths"`
}

type FileRotatelogsConfig struct {
	Filename     string `json:"filename"`
	Maxage       int    `json:"maxage"`
	RotationTime int    `json:"rotation_time"`
}

func NewDefaultLogrusWithFileRotatelogs() (*logrus.Logger, error) {
	var cfg LogrusWithFileRotatelogsConfig
	if err := json.Unmarshal([]byte(defaultLogrusLogConfig), &cfg); err != nil {
		return nil, err
	}
	return NewLogrusWithFileRotatelogs(cfg)
}

func NewLogrusWithFileRotatelogs(cfg LogrusWithFileRotatelogsConfig) (*logrus.Logger, error) {
	if len(cfg.LogrusConfig.OutputPaths) == 0 {
		return nil, errors.New("not set cfg.OutputPaths")
	}

	writer, err := rotatelogs.New(
		cfg.LogrusConfig.OutputPaths+".%Y-%m-%d-%H",
		rotatelogs.WithLinkName(cfg.FileRotatelogsConfig.Filename),                                    // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(cfg.FileRotatelogsConfig.Maxage)*time.Hour),               // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(cfg.FileRotatelogsConfig.RotationTime)*time.Minute), // 日志切割时间间隔
		//rotatelogs.WithRotationCount()                                                               // 保存日志个数，默认 7，不能与 MaxAge 同时设置
	)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()

	switch cfg.Option.Formatter {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	case "request":
		logger.SetFormatter(&RequestLogFormatter{
			FieldsOrder: cfg.Option.FieldsOrder,
		})
	}

	logger.SetOutput(writer)
	logger.SetLevel(logrus.InfoLevel)

	return logger, nil
}

// -------- RequestLogFormatter ----------
type RequestLogFormatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool
}

func (f *RequestLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	// write
	b.WriteString(entry.Time.Format("[2006-01-02 15:04:05]"))
	f.writeOrderedFields(b, entry)
	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *RequestLogFormatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	noOrderField := len(entry.Data)
	foundFieldsMap := make(map[string]bool, len(entry.Data))

	for _, field := range f.FieldsOrder {
		if value, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			noOrderField--
			f.appendValue(b, value)
		} else {
			f.appendValue(b, "-")
		}
	}

	if noOrderField > 0 {
		notFoundFields := make([]string, 0, noOrderField)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.appendValue(b, entry.Data[field])
		}
	}
}

func (f *RequestLogFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	b.WriteString("\t")

	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func (f *RequestLogFormatter) needsQuoting(text string) bool {
	if len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}
