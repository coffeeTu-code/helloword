# grpc template

一个基于 GRPC 构建微服务的模板, 实现 GRPC 请求返回问候语的功能。

# 目录

```
├── api/                          XX.proto、XX.pb.go，服务接口文件
├── pkg/                          GRPC 服务的依赖，logger、consul、grpc配置 等
├── client/                       GRPC 客户端的代码实现  
├── server/                       GRPC 服务端的代码实现
```


# 启动

本地启动 consul

```shell script

consul agent -dev

```

启动服务端

```shell script

go run code/go/go_grpc_template/server/cmd/main.go

```

运行客户端测试

```shell script

go run code/go/go_grpc_template/client/cmd/main.go

```
