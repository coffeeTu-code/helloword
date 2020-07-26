> 原文链接 [Git常用命令](https://juejin.im/post/5db25db7518825644402dc0e)
> 
> 原文链接 [服务器 ssh key 以及 git 的配置](https://shanyue.tech/op/ssh-setting)



# Git常用命令

![](https://user-gold-cdn.xitu.io/2019/10/25/16e00bdfba6174f2?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)



> 资源列表  
> [Git Book](https://git-scm.com/book/zh/v2)  
> [深入浅出Git教程（转载）](https://www.cnblogs.com/syp172654682/p/7689328.html)  
> [Git使用详细教程](http://www.admin10000.com/document/5374.html)  



[TOC]



## 名词介绍

![](https://user-gold-cdn.xitu.io/2019/10/25/16e00be00de88281?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

> Workspace：工作区
> 
> Index/Stage：暂存区，也叫索引
> 
> Repository：仓库区（或本地仓库），也存储库
> 
> Remote：远程仓库


## git常用命令

```
$ git --help

```

### 创建SSH Key

```
ssh-keygen -t rsa -C "youremail@example.com"

```

### 仓库

```
# 在当前目录新建一个Git代码库
$ git init

# 新建一个目录，将其初始化为Git代码库
$ git init [project-name]

# 下载一个项目和它的整个代码历史。git clone 实际上是一个封装了其他几个命令的命令。
# 它创建了一个新目录，切换到新的目录，
# 然后 git init 来初始化一个空的 Git 仓库， 
# 然后为你指定的 URL 添加一个（默认名称为 origin 的）远程仓库（git remote add），
# 再针对远程仓库执行 git fetch，
# 最后通过 git checkout 将远程仓库的最新提交检出到本地的工作目录。
$ git clone [url]

# 添加远程地址的时候带上用户名及密码
$ git clone http://yourname:password@git.coding.net/name/project.git
$ git clone http://yourname@git.coding.net/name/project.git

```

### 增加/删除文件

```
# 添加指定文件到暂存区
$ git add [file1] [file2] ...

# 添加指定目录到暂存区，包括子目录
$ git add [dir]

# 添加当前目录的所有文件到暂存区
$ git add .

# 添加每个变化前，都会要求确认
# 对于同一个文件的多处变化，可以实现分次提交
$ git add -p

# 删除工作区文件，并且将这次删除放入暂存区
$ git rm [file1] [file2] ...

# 停止追踪指定文件，但该文件会保留在工作区
$ git rm --cached [file]

# 改名文件，并且将这个改名放入暂存区
$ git mv [file-original] [file-renamed]

```

### git stash

```
# 执行存储时，添加备注，方便查找，只有git stash 也要可以的，但查找时不方便识别。
$ git stash save "save message"

# 查看stash了哪些存储
$ git stash list 

# 显示做了哪些改动，默认show第一个存储,如果要显示其他存贮，后面加stash@{$num}，比如第二个 git stash show stash@{1}
$ git stash show

# 显示第一个存储的改动，如果想显示其他存存储，命令：git stash show  stash@{$num}  -p ，比如第二个：git stash show  stash@{1}  -p
$ git stash show -p

# 应用某个存储,但不会把存储从存储列表中删除，默认使用第一个存储,即stash@{0}，如果要使用其他个，git stash apply stash@{$num} ， 比如第二个：git stash apply stash@{1} 
$ git stash apply

# 命令恢复之前缓存的工作目录，将缓存堆栈中的对应stash删除，并将对应修改应用到当前的工作目录下,默认为第一个stash,即stash@{0}，如果要应用并删除其他stash，命令：git stash pop stash@{$num} ，比如应用并删除第二个：git stash pop stash@{1}
$ git stash pop

# 丢弃stash@{$num}存储，从列表中删除这个存储
$ git stash drop stash@{$num}

# 删除所有缓存的stash
$ git stash clear

```

### 代码提交

```
# 提交暂存区到仓库区
$ git commit -m [message]

# 提交暂存区的指定文件到仓库区
$ git commit [file1] [file2] ... -m [message]

# 提交工作区自上次commit之后的变化，直接到仓库区
$ git commit -a

# 提交时显示所有diff信息
$ git commit -v

# 使用一次新的commit，替代上一次提交
# 如果代码没有任何新变化，则用来改写上一次commit的提交信息
$ git commit --amend -m [message]

# 重做上一次commit，并包括指定文件的新变化
$ git commit --amend [file1] [file2] ...

```

### 分支

```
# 列出所有本地分支
$ git branch

# 列出所有远程分支
$ git branch -r

# 列出所有本地分支和远程分支
$ git branch -a

# 新建一个分支，但依然停留在当前分支
$ git branch [branch-name]

# 新建一个分支，并切换到该分支
$ git checkout -b [branch]

# 新建一个分支，指向指定commit
$ git branch [branch] [commit]

# 新建一个分支，与指定的远程分支建立追踪关系
$ git branch --track [branch] [remote-branch]

# 创建远程分支newBranch
$ git push origin [newBranch] :[newBranch]

# 切换到指定分支，并更新工作区
$ git checkout [branch-name]

# 切换到上一个分支
$ git checkout -

# 建立追踪关系，在现有分支与指定的远程分支之间
$ git branch --set-upstream [branch] [remote-branch]

# 建立联系
$ git push --set-upstream origin [newBranch] 

# 合并指定分支到当前分支
$ git merge [branch]

# 选择一个commit，合并进当前分支
$ git cherry-pick [commit]

# 删除分支
$ git branch -d [branch-name]

# 删除远程分支
$ git push origin --delete [branch-name]
$ git branch -dr [remote/branch]

```

### 标签

```
# 列出所有tag
$ git tag

# 新建一个tag在当前commit
$ git tag [tag]

# 新建一个tag在指定commit
$ git tag [tag] [commit]

# 删除本地tag
$ git tag -d [tag]

# 删除远程tag
$ git push origin :refs/tags/[tagName]

# 查看tag信息
$ git show [tag]

# 提交指定tag
$ git push [remote] [tag]

# 提交所有tag
$ git push [remote] --tags

# 新建一个分支，指向某个tag
$ git checkout -b [branch] [tag]

```

### 查看信息

```
# 显示有变更的文件
$ git status

# 显示当前分支的版本历史
$ git log

# 显示commit历史，以及每次commit发生变更的文件
$ git log --stat

# 搜索提交历史，根据关键词
$ git log -S [keyword]

# 显示某个commit之后的所有变动，每个commit占据一行
$ git log [tag] HEAD --pretty=format:%s

# 显示某个commit之后的所有变动，其"提交说明"必须符合搜索条件
$ git log [tag] HEAD --grep feature

# 显示某个文件的版本历史，包括文件改名
$ git log --follow [file]
$ git whatchanged [file]

# 显示指定文件相关的每一次diff
$ git log -p [file]

# 显示过去5次提交
$ git log -5 --pretty --oneline

# 显示所有提交过的用户，按提交次数排序
$ git shortlog -sn

# 显示指定文件是什么人在什么时间修改过
$ git blame [file]

# 显示暂存区和工作区的差异
$ git diff

# 显示暂存区和上一个commit的差异
$ git diff --cached [file]

# 显示工作区与当前分支最新commit之间的差异
$ git diff HEAD

# 显示两次提交之间的差异
$ git diff [first-branch]...[second-branch]

# 显示今天你写了多少行代码
$ git diff --shortstat "@{0 day ago}"

# 显示某次提交的元数据和内容变化
$ git show [commit]

# 显示某次提交发生变化的文件
$ git show --name-only [commit]

# 显示某次提交时，某个文件的内容
$ git show [commit]:[filename]

# 显示当前分支的最近几次提交
$ git reflog

```

### 远程同步

```
# 下载远程仓库的所有变动
$ git fetch [remote]

# 显示所有远程仓库
$ git remote -v

# 显示某个远程仓库的信息
$ git remote show [remote]

# 增加一个新的远程仓库，并命名
$ git remote add [shortname] [url]

# 取回远程仓库的变化，并与本地分支合并
$ git pull [remote] [branch]

# 上传本地指定分支到远程仓库
$ git push [remote] [branch]

# 强行推送当前分支到远程仓库，即使有冲突
$ git push [remote] --force

# 推送所有分支到远程仓库
$ git push [remote] --all

```

### 撤销

```
# 恢复暂存区的指定文件到工作区
$ git checkout [file]

# 恢复某个commit的指定文件到暂存区和工作区
$ git checkout [commit] [file]

# 恢复暂存区的所有文件到工作区
$ git checkout .

# 重置暂存区的指定文件，与上一次commit保持一致，但工作区不变
$ git reset [file]

# 重置暂存区与工作区，与上一次commit保持一致
$ git reset --hard

# 重置当前分支的指针为指定commit，同时重置暂存区，但工作区不变
$ git reset [commit]

# 重置当前分支的HEAD为指定commit，同时重置暂存区和工作区，与指定commit一致
$ git reset --hard [commit]

# 重置当前HEAD为指定commit，但保持暂存区和工作区不变
$ git reset --keep [commit]

# 新建一个commit，用来撤销指定commit
# 后者的所有变化都将被前者抵消，并且应用到当前分支
$ git revert [commit]

# 暂时将未提交的变化移除，稍后再移入
$ git stash
$ git stash pop

```



## 服务器 ssh key 与 git 配置

虽然 git 可以工作在 ssh 与 https 两种协议上，但为了安全性，更多时候会选择 ssh。

> 如果采用 https，则每次 git push 都需要验证身份


### Permission denied (publickey).

如果没有设置 public key 直接 git clone 的话，会有权限问题

可以使用 `ssh -T` 测试连通性

```
$ git clone git@github.com:vim/vim.git
Cloning into 'vim'...
Warning: Permanently added the RSA host key for IP address '13.229.188.59' to the list of known hosts.
Permission denied (publickey).
fatal: Could not read from remote repository.

Please make sure you have the correct access rights
and the repository exists.

# 不过有一个更直接的命令去查看是否有权限
$ ssh -T git@github.com
Permission denied (publickey).


```

### 生成一个新的 ssh key

使用 ssh-keygen 可以生成配对的 `id_rsa` 与 `id_rsa.pub` 文件

```
# 生成一个 ssh-key
# -t: 可选择 dsa | ecdsa | ed25519 | rsa | rsa1，代表加密方式
# -C: 注释，一般写自己的邮箱
$ ssh-keygen -t rsa -C "shanyue"

# 生成 id_rsa/id_rsa.pub: 配对的私钥与公钥
$ ls ~/.ssh
authorized_keys  config  id_rsa  id_rsa.pub  known_hosts


```

### 在 github 设置里新添一个 ssh key

在云服务器中复制 ~/.ssh/id_rsa.pub 中文件内容，并粘贴到 github 的配置 中

在 github 的 ssh keys 设置中：github.com/settings/ke… 点击 New SSH key 添加刚才的key。

```
$ cat ~/.ssh/id_rsa.pub
ssh-rsa AAAAB3SSSSSSSSSSSSSSSSSSSSSBAQDcM4aOo9qlrHOnh0+HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHuM9cYmdKq5ZMfO0dQ5PB53nqZQ1YAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAc1w7bC0PD02M706ZdQm5M9Q9VFzLY0TK1nz19fsh2I2yuKwHJJeRxsFAUJKgrtNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNNN7nm6B/9erp5n4FDKJFxdnFWuhqqUwMzRa9rUfhOX1qJ1SYAWUryQ90rpxOwXt9Pfq0Y13VsWk3QQ8nyaEJzytEXG7OR9pf9zDQph4r4rpJbXCwNjXn/ThL shanyue


```

### 设置成功

使用 ssh -T 测试成功， 此时可以成功的面向 github 编程了

```
$ ssh -T git@github.com
Hi shfshanyue! You've successfully authenticated, but GitHub does not provide shell access.

$ git clone git@github.com:shfshanyue/vim-config.git
Cloning into 'vim-config'...
remote: Enumerating objects: 183, done.
remote: Total 183 (delta 0), reused 0 (delta 0), pack-reused 183
Receiving objects: 100% (183/183), 411.13 KiB | 55.00 KiB/s, done.
Resolving deltas: 100% (100/100), done.


```