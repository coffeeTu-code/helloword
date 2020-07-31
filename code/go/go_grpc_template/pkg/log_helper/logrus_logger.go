//******************************************************
//> 原文链接：
//> [Golang logrus 日志包及日志切割的实现 https://www.jb51.net/article/180448.htm](https://www.jb51.net/article/180448.htm)
//> [Logrus基本用法 https://www.jianshu.com/p/2d90b32acade](https://www.jianshu.com/p/2d90b32acade)
//> [bilibili代码基础日志组件 https://github.com/bilibili/sniper/blob/master/util/log/log.go#L61](https://github.com/bilibili/sniper/blob/master/util/log/log.go#L61)
//
//> GitHub：[https://github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)
//> doc：[https://godoc.org/github.com/sirupsen/logrus](https://godoc.org/github.com/sirupsen/logrus)
//
//******************************************************

package log_helper

import (
	"bytes"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"sort"
	"time"
)

func NewLogrusWithFileRotatelogs(outputPath string) *logrus.Logger {
	if len(outputPath) == 0 {
		outputPath = "stdout"
	}

	writer, err := rotatelogs.New(
		outputPath+".%Y-%m-%d-%H",
		rotatelogs.WithLinkName(outputPath),                        // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Duration(3)*time.Hour),          // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Duration(60)*time.Minute), // 日志切割时间间隔
		//rotatelogs.WithRotationCount()                            // 保存日志个数，默认 7，不能与 MaxAge 同时设置
	)
	if err != nil {
		fmt.Println("rotatelogs.New ", outputPath, "error=", err)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(writer)
	logger.SetLevel(logrus.InfoLevel)

	return logger
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
