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

```shell script
 go build -o ./code/go/go_grpc_template/server/cmd/go_grpc_template ./code/go/go_grpc_template/server/cmd/main.go

# 制作镜像
docker build -t go_grpc_template .
# 查看镜像
docker images
# REPOSITORY               TAG                 IMAGE ID            CREATED             SIZE
# go_grpc_template         latest              a0d3c284621e        8 minutes ago       20.7MB
# busybox                  latest              018c9d7b792b        13 days ago         1.22MB
# docker/getting-started   latest              1f32459ef038        3 weeks ago         26.8MB
# redislabs/rebloom        latest              03841e395ca0        4 weeks ago         104MB

docker run -it -p 9120:9120 go_grpc_template



```

运行客户端测试

```shell script

go run code/go/go_grpc_template/client/cmd/main.go

```
