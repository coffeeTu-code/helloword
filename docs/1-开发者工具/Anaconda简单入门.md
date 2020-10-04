# Anaconda简单入门

## Anaconda简介

Anaconda是一个方便的python包管理和环境管理软件，一般用来配置不同的项目环境。

我们常常会遇到这样的情况，正在做的项目A和项目B分别基于python2和python3，而第电脑只能安装一个环境，这个时候Anaconda就派上了用场，它可以创建多个互不干扰的环境，分别运行不同版本的软件包，以达到兼容的目的。

Anaconda通过管理工具包、开发环境、Python版本，大大简化了你的工作流程。不仅可以方便地安装、更新、卸载工具包，而且安装时能自动安装相应的依赖包，同时还能使用不同的虚拟环境隔离不同要求的项目。

anconda使用conda作为包管理工具，也就是anaconda把一些常用的python包统一管理，可以方便的安装、更新和卸载包。
conda常用命令

```shell script

conda --version 查看conda版本
conda -V
conda --help   获取帮助
conda -h
conda update --help 
conda remove --help
--help 都可以换成-h
```

## Anaconda环境管理

接下来我们先研究单一环境下的Anaconda环境管理。

conda env --help

刚刚系统默认创建了名叫base的默认环境，我们可以使用conda命令查看当前有多少环境

conda env list 查看所有环境

或者
conda info --envs

我们可以看到只有一个base，也就是刚刚安装的时候勾选的带有python3.5的环境，还有一些py3的包
如果现在我需要一个python2.7的和tensorflow1.0的环境该怎么办呢

```
输入

conda create --name your_env_name  
或者  
conda create -n your_env_name
```

your_env_name 就是你新创建的环境名，你可以在里面安装其他包但不会与现有环境冲突

如果你要在创建环境时指定包内容， 可以用

```shell script

conda create -n your_env_name python=3.5
```

如果你要指定多个包 可以用

```
conda create -n your_env_name python=3.5 numpy pandas
```

要指定特殊版本号加上=版本号就行，默认是最新的
对了，安装前为了保障你查询到最新包情况，最好使用
`conda update --all 更新包信息。

例如我现在要创建一个名叫 learningpy的基于py3的环境

```shell script

conda update --all
conda create -n learningpy python=3.7
```

系统会询问你是否创建，输入y回车后，系统将列出必要安装的包

conda有一点好处是，如果你需要安装一个包，系统将自动检查这个包需要的前置包并且安装，比如你要安装TensorFlow，而TensorFlow会用到很多像前置包像pandas、matiplot等，如果你在单纯的python下没有安装pandas等包就直接安装TensorFlow，那么和有可能无法使用，而使用conda安装TensorFlow将会询问你并自动帮你把缺少的前置包安装好

那么现在我们有多个环境了，如何切换环境呢？

```shell script

windows
activate 环境名

退出时记得退出命令哦
deactivate
```

```shell script

linux和mac用户的命令不一样
conda activate 环境名(cr)

```

```shell script

创建一个新环境想克隆一部分旧的环境
conda create -n your_env_name --clone oldname

删除某个环境
conda remove -n your_env_name --all

导出环境配置（非常有用，比如你想帮朋友安装和你一模一样的环境，你可以直接导出一个配置文件给他，就能免除很多人力安装调试)
conda env export > environment.yml

将会在当前目录生成一个environment.yml,你把它交给小伙伴或拷到另一台机器，小伙伴只需要对这个文件执行命令  
conda env create -f environment.yml

就可以生成和你原来一模一样的环境啦
```

## anaconda包管理

上文我们提到了创建环境时的包管理，那么我们创建好环境后如何进行包的安装、更新和卸载呢？

```
conda list 列举当前环境下的所有包
conda list -n packagename 列举某个特定名称包
conda install packagename 为当前环境安装某包
conda install -n envname packagename 为某环境安装某包
conda search packagename 搜索某包
conda updata packagename 更新当前环境某包
conda update -n envname packagename 更新某特定环境某包
conda remove packagename 删除当前环境某包
conda remove -n envname packagename 删除某环境环境某包
conda本身和anaconda、python本身也算包
conda update conda
conda update anaconda
conda update python
conda默认源可能速度比较慢
```

可以添加其他源，常用的有清华TUNA

```

conda config --add channels https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/free/
conda config --add channels https://mirrors.tuna.tsinghua.edu.cn/anaconda/pkgs/main/
conda config --set show_channel_urls yes 在包后面显示来源
第三条执行安装包时会显示来自哪个源，一目了然
source.png
教育网用户可以添加ipv6源，速度很快

conda config --add channels https://mirrors6.tuna.tsinghua.edu.cn/anaconda/pkgs/free/
conda config --add channels https://mirrors6.tuna.tsinghua.edu.cn/anaconda/pkgs/main/
conda config --set show_channel_urls yes 在包后面显示来源
anaconda实现原理解析
anaconda在目录下的envs文件夹保存了环境配置，也就是把所有的安装在这个环境下的包放在同一个文件夹中
当创建一个新环境时，anaconda将在envs中创建一个新的文件夹，这个文件夹包括了你安装在这个环境中的所有包
anaconda通过巧妙的包管理解决的一个大难题，确实方便了很多。
下一期会讲如何在第三方软件中使用anaconda的不同环境配置。
```
