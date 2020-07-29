# Standard Go项目布局

[原文指引](https://github.com/golang-standards/project-layout)  

这是Go应用程序项目的基本布局。它不是核心Go开发团队定义的官方标准; 然而，它是Go生态系统中一组常见的历史和新兴项目布局模式。

## Go目录

### `/ cmd`

该项目的主要应用。

每个应用程序的目录名称应与您想要的可执行文件的名称相匹配（例如，`/ cmd / myapp`）。

不要在应用程序目录中放入大量代码。  
如果您认为代码可以导入并在其他项目中使用，那么它应该存在于`/ pkg`目录中。如果代码不可重用或者您不希望其他人重用它，请将该代码放在`/ internal`目录中。  
通常有一个小的`main`函数可以导入和调用`/ internal`和`/ pkg`目录中的代码，而不是别的。

### `/ internal`

私有应用程序和库代码。这是您不希望其他人在其应用程序或库中导入的代码。

将您的实际应用程序代码放在`/ internal / app`目录中（例如，`/ internal / app / myapp`）以及这些应用程序在`/ internal / pkg`目录中共享的代码（例如，`/ internal / pkg / myprivlib`）。

### `/ pkg` 

可以由外部应用程序使用的库代码（例如，`/ pkg / mypubliclib`）。其他项目将导入这些库，期望它们可以工作，所以在你把东西放在这里之前要三思而后行:-)

### `/ vendor` 

应用程序依赖项（手动管理或由您最喜欢的依赖管理工具管理，如[`dep`]（https://github.com/golang/dep））。

### `/ api` 

OpenAPI / Swagger规范，JSON模式文件，协议定义文件。

## Web应用程序目录

### `/ web` 

Web应用程序特定组件：静态Web资产，服务器端模板和SPA。

## Common Application Directories 

### `/ configs` 

配置文件模板或默认配置。

### `/ init` 

系统init（systemd，upstart，sysv）和进程管理器/主管（runit，supervisord）配置。

### `/ scripts` 

脚本执行各种构建，安装，分析等操作。

### `/ build` 

打包和持续集成。

将您的云（AMI），容器（Docker），OS（deb，rpm，pkg）包配置和脚本放在`/ build / package`目录中。

将CI（travis，circle，drone）配置和脚本放在`/ build / ci`目录中。

### `/ deployments`

IaaS，PaaS，系统和容器编排部署配置和模板（docker-compose，kubernetes / helm，mesos，terraform，bosh）。

### `/ test` 

其他外部测试应用和测试数据。您可以随意构建`/ test`目录。对于更大的项目，有一个数据子目录是有意义的。例如，如果需要Go忽略该目录中的内容，可以使用`/ test / data`或`/ test / testdata`。请注意，Go也会忽略以“。”开头的目录或文件。或“_”，因此您在命名测试数据目录方面具有更大的灵活性。

### `/ tools` 

为该项目提供支持工具。请注意，这些工具可以从`/ pkg`和`/ internal`目录导入代码。

### `/ examples` 

您的应用程序和/或公共库的示例。

### `/ third_party` 

外部帮助工具，分叉代码和其他第三方实用程序（例如，Swagger UI）。

### `/ githooks` 

Git钩子。

### `/ assets` 

与您的存储库一起使用的其他资产（图像，徽标等）。

### `/ website`

如果您不使用Github页面，这是放置项目的网站数据的地方。

