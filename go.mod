module helloword

go 1.14

// 超过v1.26.0 版本的发布拉取代码有限制
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect

	// go-kit 微服务组件
	github.com/go-kit/kit v0.10.0
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.7.3

	// leetcode 题库
	github.com/halfrost/LeetCode-Go v0.0.0-20200716140546-1b39219e2b81
	github.com/hashicorp/consul/api v1.3.0
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jonboulle/clockwork v0.2.0 // indirect

	// 更快的 json 组件，支持原生
	github.com/json-iterator/go v1.1.8
	github.com/kr/text v0.2.0 // indirect

	// 文件切割器，通常搭配 logrus 一起使用
	github.com/lestrrat-go/file-rotatelogs v2.3.0+incompatible
	github.com/lestrrat-go/strftime v1.0.3 // indirect
	github.com/miekg/dns v1.1.27 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect

    // 文件切割器，通常搭配 zap 一起使用
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/oklog/oklog v0.3.2
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shirou/gopsutil v2.19.10+incompatible

	// logrus-logger 日志组件
	github.com/sirupsen/logrus v1.6.0

	// go-test 框架
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/tebeka/strftime v0.1.4 // indirect

	// zap-logger 日志组件
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f // indirect
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2
	golang.org/x/sys v0.0.0-20200615200032-f1bc736245b1 // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e
	golang.org/x/tools v0.0.0-20200117065230-39095c1d176c // indirect
	google.golang.org/grpc v1.30.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)
