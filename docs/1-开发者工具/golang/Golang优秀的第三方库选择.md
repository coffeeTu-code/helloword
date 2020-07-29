> 原文链接：[go语言几个最快最好运用最广的web框架比较](https://www.cnblogs.com/desmond123/p/9821687.html)
>
> 原文链接：[golang比较优秀的第三方库收集](https://www.jianshu.com/p/94ccafe2a982)



# golang比较优秀的第三方库收集



[TOC]



##  Go Web 框架选择


如果你为自己设计一个小应用程序，你可能不需要一个Web框架，但如果你正在进行生产，那么你肯定需要一个，一个好的应用程序。

虽然您认为自己拥有必要的知识和经验，但您是否愿意自行编写所有这些功能的代码？
您是否有时间找到生产级外部包来完成这项工作？ 您确定这将与您应用的其余部分保持一致吗？

这些是推动我们使用框架的原因，如果其他人已经做了这些繁琐且艰苦的工作，我们不想自己编写所有必要的功能。


| 开源web框架 | - | Github地址 | 文档地址 | 学习曲线 |
| :- | :- | :- | :- | :- |
| Beego | Go编程语言的开源，高性能Web框架。 |  https://github.com/astaxie/beego  | https://beego.me |  https://beego.me/docs 示例49|
| Buffalo | 快速Web开发w/Go。 | https://github.com/gobuffalo/buffalo  |https://gobuffalo.io | https://gobuffalo.io/docs/installation 示例6 |
| Echo | 高性能，极简主义的Go Web框架。 | https://github.com/labstack/echo   | https://echo.labstack.com   | https://echo.labstack.com/cookbook/hello-world 示例20 |
| Gin | 用Go（Golang）编写的HTTP Web框架。它具有类似Martini的API，具有更好的性能。 | https://github.com/gin-gonic/gin   | https://gin-gonic.github.io/gin   | https://github.com/gin-gonic/gin/tree/master/examples 示例15 |
| Iris | Go in the Universe中最快的Web框架。MVC功能齐全。 | https://github.com/kataras/iris   | https://iris-go.com   | https://github.com/kataras/iris/tree/master/_examples 示例92 |
| Revel | Go语言的高生产力，全栈Web框架。 | https://github.com/revel/revel   | https://revel.github.io   | http://revel.github.io/examples/index.html 示例6 |


## NOSQL 数据库操作


NoSQL(Not Only SQL)，指的是非关系型的数据库。目前流行的NOSQL主要有redis、mongoDB、Cassandra和Membase等。这些数据库都有高性能、高并发读写等特点，目前已经广泛应用于各种应用中。

https://github.com/garyburd/redigo (推荐)
https://github.com/go-redis/redis
https://github.com/hoisie/redis
https://github.com/alphazero/Go-Redis
https://github.com/simonz05/godis

```
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
)

var (
	Pool *redis.Pool
)

func init() {
	redisHost := ":6379"
	Pool = newPool(redisHost)
	close()
}

func newPool(server string) *redis.Pool {

	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		}
	}
}

func close() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		Pool.Close()
		os.Exit(0)
	}()
}

func Get(key string) ([]byte, error) {

	conn := Pool.Get()
	defer conn.Close()

	var data []byte
	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return data, fmt.Errorf("error get key %s: %v", key, err)
	}
	return data, err
}

func main() {
	test, err := Get("test")
	fmt.Println(test, err)
}
```

MongoDB是一个高性能，开源，无模式的文档型数据库，是一个介于关系数据库和非关系数据库之间的产品，是非关系数据库当中功能最丰富，最像关系数据库的。他支持的数据结构非常松散，采用的是类似json的bjson格式来存储数据，因此可以存储比较复杂的数据类型。Mongo最大的特点是他支持的查询语言非常强大，其语法有点类似于面向对象的查询语言，几乎可以实现类似关系数据库单表查询的绝大部分功能，而且还支持对数据建立索引。

```
package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

func main() {
	session, err := mgo.Dial("server1.example.com,server2.example.com")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("people")
	err = c.Insert(&Person{"Ale", "+55 53 8116 9639"},
		&Person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := Person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}
```


## log 


| - | Github地址 |
| :- | :- |
| logrus | https://github.com/sirupsen/logrus |
| file-rotatelogs | https://github.com/lestrrat-go/file-rotatelogs |


## JSON解析


| - | Github地址 | - |
| :- | :- | - |
| json-iterator | https://github.com/json-iterator/go | 号称最快的go json解析器。跟官方的写法兼容，我目前基本都使用这个解析器。 |
| tidwall/gjson | https://github.com/tidwall/gjson | 直接像其他语言一样根据Key来获取数据，方便很多。 |


## 随机数


| - | Github地址 | - |
| :- | :- | - |
| NebulousLabs/fastrand | https://github.com/NebulousLabs/fastrand | 性能非常好的随机数库，实测确实比官方随机数快很多。|


## 数据库


golang数据库自带连接池，操作数据库时需要事先了解一下这个概念。

| - | Github地址 | - |
| :- | :- | - |
| go-sql-driver/mysql | https://github.com/go-sql-driver/mysql | 我想大多数人都使用这个mysql连接驱动包。|
| gomodule/redigo | https://github.com/gomodule/redigo | 不错的redis客户端。 |


## protobuf


| - | Github地址 | - |
| :- | :- | - |
| golang/protobuf | https://github.com/golang/protobuf | 官方的protobuf库。 |
| gogo/protobuf | https://github.com/gogo/protobuf | 据说比楼上官方的快很多，没有实际使用过。|


## 配置库


| - | Github地址 | - |
| :- | :- | - |
| BurntSushi/toml | https://github.com/BurntSushi/toml | 大多数go程序都使用toml做为配置文件。|
| viper | https://github.com/spf13/viper | 支持从多种配置格式的文件、环境变量、命令行等读取配置。|


## 爬虫


| - | Github地址 | - |
| :- | :- | - |
| PuerkitoBio/goquery | https://github.com/PuerkitoBio/goquery | |



## 静态打包工具


| - | Github地址 | - |
| :- | :- | - |
| packr | https://github.com/gobuffalo/packr | 将静态资源打包成go文件，直接编译进最终程序。 |



## 图表chart


| - | Github地址 | - |
| :- | :- | - |
| go-echarts | https://github.com/go-echarts/go-echarts | 非常好的chart图表，是百度开源js echarts的go封装，发布程序时需要用到上面的packr把资源打包成go文件，然后再编译发布。|


## Hashids


| - | Github地址 | - |
| :- | :- | - |
| go-hashids | https://github.com/speps/go-hashids |  |


## SnowFlake 算法


| - | Github地址 | - |
| :- | :- | - |
| goSnowFlake | https://github.com/zheng-ji/goSnowFlake | 分布式id生成器 |



## mail发送客户端


| - | Github地址 | - |
| :- | :- | - |
| gomail | https://github.com/go-gomail/gomail |  |



## 验证码服务


| - | Github地址 | - |
| :- | :- | - |
| captcha | https://github.com/dchest/captcha ||
| gocaptcha | https://github.com/hanguofeng/gocaptcha | http://ju.outofmemory.cn/entry/48563 |

