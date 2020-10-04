#第一行必须指令基于的基础镜像
FROM scratch

#维护者信息
MAINTAINER xuefeng.han  xuefeng.han@mobvista.com

#镜像的操作指令
ADD ./code/go/go_grpc_template/server/cmd/go_grpc_template /go_grpc_template

RUN ["/bin/bash", "-c", "echo hello"]
RUN ["/go_grpc_template"]

#容器启动时执行指令
CMD "echo" "-c" "Hello docker!"
