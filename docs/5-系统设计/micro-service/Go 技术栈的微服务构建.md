[TOC]



# 基于 Go 技术栈的微服务构建

- 基于 Go 技术栈的微服务构建[https://www.infoq.cn/article/BRTyU40J1qxENh53mHSH?utm_source=related_read_bottom&utm_medium=article](https://www.infoq.cn/article/BRTyU40J1qxENh53mHSH?utm_source=related_read_bottom&utm_medium=article)  
本文的素材来源于我们在开发中的一些最佳实践案例，从开发、监控、日志这三个角度介绍了一些我们基于 Go 技术栈的微服务构建经验。  
grpc+TLS+trace, Prometheus, logrus


- 罗辑思维首席架构师：Go 微服务改造实践[https://www.infoq.cn/article/2018/05/luojisiwei?utm_source=related_read_bottom&utm_medium=article](https://www.infoq.cn/article/2018/05/luojisiwei?utm_source=related_read_bottom&utm_medium=article)  
对于系统改造来说，首先需要知道，系统需要改成什么样子。因此我们需要一个架构的蓝图。上面就是我们的架构蓝图。首先需要的是一个统一对外的 API GATEWAY，向下是对外的业务服务+基础资源服务，最下层是公用服务的一些基础设施，搭配一些通用的框架和中间件。