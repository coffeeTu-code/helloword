> 原文链接：  
> [Golang logrus 日志包及日志切割的实现 https://www.jb51.net/article/180448.htm](https://www.jb51.net/article/180448.htm)  
> [Logrus基本用法 https://www.jianshu.com/p/2d90b32acade](https://www.jianshu.com/p/2d90b32acade)  
> [bilibili代码基础日志组件 https://github.com/bilibili/sniper/blob/master/util/log/log.go#L61](https://github.com/bilibili/sniper/blob/master/util/log/log.go#L61)

> GitHub：[https://github.com/sirupsen/logrus](https://github.com/sirupsen/logrus)  
> doc：[https://godoc.org/github.com/sirupsen/logrus](https://godoc.org/github.com/sirupsen/logrus)

# Retarget Dsp 项目实践

自定义日志格式实现请求日志记录。

```

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func NewLogrusWithFileRotatelogs(logOutputPath string) (*logrus.Logger, error) {

	absPath, _ := filepath.Abs(logOutputPath)
	err := mkLogDir(absPath)
	if err != nil {
		return nil, err
	}

	writer, err := rotatelogs.New(
		absPath+".%Y-%m-%d-%H",
		rotatelogs.WithLinkName(absPath),       // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(4*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetFormatter(&RequestLogFormatter{
		//time               //0
		FieldsOrder: []string{
			"deviceIp",      //1
			"offerType",     //2
			"exchanges",     //3
			"elapsed",       //4
			"reqData",       //5
			"ab_flag",       //6
			"requestId",     //7
			"bid",           //8
			"price",         //9
			"describe",      //10
			"algInfo",       //11
			"subAlgInfo",    //12
			"mvappId",       //13
			"ivr",           //14
			"devIds",        //15
			"auctionType",   //16
			"unitId",        //17
			"placementId",   //18
			"publisherId",   //19
			"appId",         //20
			"appName",       //21
			"extProcess",    //22
			"category",      //23
			"imageSize",     //24
			"reqPackages",   //25
			"reqOffers",     //26
			"make",          //27
			"model",         //28
			"os",            //29
			"osv",           //30
			"deviceType",    //31
			"cncType",       //32
			"countryCode",   //33
			"googleAdId",    //34
			"extRequest",    //35
			"city",          //36
			"reqType",       //37
			"keywords",      //38
			"yob",           //39
			"gender",        //40
			"userAgent",     //41
			"trafficType",   //42
			"carrier",       //43
			"bidFloor",      //44
			"extOffer",      //45
			"campaignId",    //46
			"cInstallPrice", //47
			"cAppName",      //48
			"cPackageName",  //49
			"creativeInfo",  //50
		},
		HideKeys: true,
	})
	logger.SetOutput(writer)
	logger.SetLevel(logrus.InfoLevel)

	return logger, nil
}

func mkLogDir(logPath string) error {
	dir, _ := filepath.Split(logPath)
	if len(dir) > 0 {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

//****************** Request Log **********************
func (this *RetargetLog) Req(fields logrus.Fields) {
	if this.GetRequestLogger() != nil {
		this.GetRequestLogger().WithFields(fields).Info("")
	}
}

//****************** 自定义的日志格式 **********************
type RequestLogFormatter struct {
	// FieldsOrder - default: fields sorted alphabetically
	FieldsOrder []string

	// HideKeys - show [fieldValue] instead of [fieldKey:fieldValue]
	HideKeys bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool
}

func (f *RequestLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// output buffer
	b := &bytes.Buffer{}
	// write time
	b.WriteString(entry.Time.Format("[2006-01-02 15:04:05]"))

	// write fields
	f.writeOrderedFields(b, entry)

	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *RequestLogFormatter) writeOrderedFields(b *bytes.Buffer, entry *logrus.Entry) {
	length := len(entry.Data)
	foundFieldsMap := map[string]bool{}
	for _, field := range f.FieldsOrder {
		if _, ok := entry.Data[field]; ok {
			foundFieldsMap[field] = true
			length--
			f.writeField(b, entry, field)
		} else {
			if f.HideKeys {
				b.WriteString("\t-")
			} else {
				b.WriteString("\t")
				b.WriteString(field)
				b.WriteString(":-")
			}

		}
	}

	if length > 0 {
		notFoundFields := make([]string, 0, length)
		for field := range entry.Data {
			if foundFieldsMap[field] == false {
				notFoundFields = append(notFoundFields, field)
			}
		}

		sort.Strings(notFoundFields)

		for _, field := range notFoundFields {
			f.writeField(b, entry, field)
		}
	}
}
func (f *RequestLogFormatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	if f.HideKeys {
		//fmt.Fprintf(b, "\t%v", entry.Data[field])
		if v, ok := rtutil.GetStringFromInterface(entry.Data[field]); ok != nil {
			fmt.Fprintf(b, "\t%v", entry.Data[field])
		} else {
			b.WriteString("\t")
			b.WriteString(v)
		}

	} else {

		if v, ok := rtutil.GetStringFromInterface(entry.Data[field]); ok != nil {
			fmt.Fprintf(b, "\t%s:%v", field, entry.Data[field])
		} else {
			b.WriteString("\t")
			b.WriteString(field)
			b.WriteString(":")
			b.WriteString(v)
		}
	}
}

```
