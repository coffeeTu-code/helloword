> Consul 简介和快速入门 [https://book-consul-guide.vnzmi.com/](https://book-consul-guide.vnzmi.com/)


# 安装Consul

## OS X

如果你使用homebrew作为包管理器,你可以使用命令

```shell script
brew install consul

```

来进行安装.

## 验证安装

完成安装后,通过打开一个新终端窗口检查consul安装是否成功.通过执行 consul --version你应该看到类似下面的输出

```shell script

Consul v1.8.1
Revision 12f574c+CHANGES
Protocol 2 spoken by default, understands 2 to 3 (agent will automatically use protocol >2 when speaking to compatible agents)

```

# 运行Agent

完成Consul的安装后,必须运行agent. agent可以运行为server或client模式.每个数据中心至少必须拥有一台server . 建议在一个集群中有3或者5个server.部署单一的server,在出现失败时会不可避免的造成数据丢失.

其他的agent运行为client模式.一个client是一个非常轻量级的进程.用于注册服务,运行健康检查和转发对server的查询.agent必须在集群中的每个主机上运行.

## 启动 Agent

为了更简单,现在我们将启动Consul agent的开发模式.这个模式快速和简单的启动一个单节点的Consul.这个模式不能用于生产环境,因为他不持久化任何状态.

```shell script
[root@hdp2 ~]# consul agent -dev
==> Starting Consul agent...
==> Starting Consul agent RPC...
==> Consul agent running!
         Node name: 'hdp2'
        Datacenter: 'dc1'
            Server: true (bootstrap: false)
       Client Addr: 127.0.0.1 (HTTP: 8500, HTTPS: -1, DNS: 8600, RPC: 8400)
      Cluster Addr: 10.0.0.52 (LAN: 8301, WAN: 8302)
    Gossip encrypt: false, RPC-TLS: false, TLS-Incoming: false
             Atlas: <disabled>

==> Log data will now stream in as it occurs:

    2016/08/17 15:20:41 [INFO] serf: EventMemberJoin: hdp2 10.0.0.52
    2016/08/17 15:20:41 [INFO] serf: EventMemberJoin: hdp2.dc1 10.0.0.52
    2016/08/17 15:20:41 [INFO] raft: Node at 10.0.0.52:8300 [Follower] entering Follower state
    2016/08/17 15:20:41 [INFO] consul: adding LAN server hdp2 (Addr: 10.0.0.52:8300) (DC: dc1)
    2016/08/17 15:20:41 [INFO] consul: adding WAN server hdp2.dc1 (Addr: 10.0.0.52:8300) (DC: dc1)
    2016/08/17 15:20:41 [ERR] agent: failed to sync remote state: No cluster leader
    2016/08/17 15:20:42 [WARN] raft: Heartbeat timeout reached, starting election
    2016/08/17 15:20:42 [INFO] raft: Node at 10.0.0.52:8300 [Candidate] entering Candidate state
    2016/08/17 15:20:42 [DEBUG] raft: Votes needed: 1
    2016/08/17 15:20:42 [DEBUG] raft: Vote granted from 10.0.0.52:8300. Tally: 1
    2016/08/17 15:20:42 [INFO] raft: Election won. Tally: 1
    2016/08/17 15:20:42 [INFO] raft: Node at 10.0.0.52:8300 [Leader] entering Leader state
    2016/08/17 15:20:42 [INFO] raft: Disabling EnableSingleNode (bootstrap)
    2016/08/17 15:20:42 [DEBUG] raft: Node 10.0.0.52:8300 updated peer set (2): [10.0.0.52:8300]
    2016/08/17 15:20:42 [INFO] consul: cluster leadership acquired
    2016/08/17 15:20:42 [DEBUG] consul: reset tombstone GC to index 2
    2016/08/17 15:20:42 [INFO] consul: member 'hdp2' joined, marking health alive
    2016/08/17 15:20:42 [INFO] consul: New leader elected: hdp2
    2016/08/17 15:20:43 [INFO] agent: Synced service 'consul'

```

如你所见,Consul Agent 启动并输出了一些日志数据.从这些日志中你可以看到,我们的agent运行在server模式并且声明作为一个集群的领袖.额外的本地镀锌被标记为一个健康的成员.

> OS X用户注意: Consul 使用你的主机hostname作为默认的节点名字.如果你的主机名包含时间,到这个节点的DNS查询将不会工作.为了避免这个情况,使用-node参数来明确的设置node名.

## 集群成员

新开一个终端窗口运行consul members, 你可以看到Consul集群的成员.下一节我们将讲到加入集群.现在你应该只能看到一个成员,就是你自己:

```shell script

[root@hdp2 ~]# consul members
Node  Address         Status  Type    Build  Protocol  DC
hdp2  10.0.0.52:8301  alive   server  0.6.4  2         dc1

```

这个输出显示我们自己的节点.运行的地址,健康状态,自己在集群中的角色,版本信息.添加-detailed选项可以查看到额外的信息.

members命令的输出是基于gossip协议是最终一致的.意味着,在任何时候,通过你本地agent看到的结果可能不是准确匹配server的状态.为了查看到一致的信息,使用HTTP API(将自动转发)到Consul Server上去进行查询:

```shell script

[root@hdp2 ~]#  curl localhost:8500/v1/catalog/nodes
[{"Node":"hdp2","Address":"10.0.0.52","TaggedAddresses":{"wan":"10.0.0.52"},"CreateIndex":3,"ModifyIndex":4}]

```

除了HTTP API ,DNS 接口也可以用来查询节点.注意,你必须确定将你的DNS查询指向Consul agent的DNS服务器,这个默认运行在 8600端口.DNS条目的格式(例如:"Armons-MacBook-Air.node.consul")将在后面讲到.

```shell script
$ dig @127.0.0.1 -p 8600 Armons-MacBook-Air.node.consul
...

;; QUESTION SECTION:
;Armons-MacBook-Air.node.consul.    IN  A

;; ANSWER SECTION:
Armons-MacBook-Air.node.consul. 0 IN    A   172.20.20.11

```

# 注册服务

## 定义一个服务

可以通过提供服务定义或者调用HTTP API来注册一个服务.服务定义文件是注册服务的最通用的方式.所以我们将在这一步使用这种方式.我们将会建立在前一步我们覆盖的代理配置。

首先,为Consul配置创建一个目录.Consul会载入配置文件夹里的所有配置文件.在Unix系统中通常类似 /etc/consul.d (.d 后缀意思是这个路径包含了一组配置文件).

```shell script
sudo mkdir /etc/consul.d

```

然后,我们将编写服务定义配置文件.假设我们有一个名叫web的服务运行在 80端口.另外,我们将给他设置一个标签.这样我们可以使用他作为额外的查询方式:

```shell script
echo '{"service": {"name": "web", "tags": ["rails"], "port": 80}}' \
    >/etc/consul.d/web.json

```

现在重启agent , 设置配置目录:

```shell script
$ consul agent -dev -config-dir /etc/consul.d
==> Starting Consul agent...
...
    [INFO] agent: Synced service 'web'
...

```

你可能注意到了输出了 "synced" 了 web这个服务.意思是这个agent从配置文件中载入了服务定义,并且成功注册到服务目录.

如果你想注册多个服务,你应该在Consul配置目录创建多个服务定义文件.

## 查询服务

一旦agent启动并且服务同步了.我们可以通过DNS或者HTTP的API来查询服务.

### DNS API

让我们首先使用DNS API来查询.在DNS API中,服务的DNS名字是 NAME.service.consul. 虽然是可配置的,但默认的所有DNS名字会都在consul命名空间下.这个子域告诉Consul,我们在查询服务,NAME则是服务的名称.

对于我们上面注册的Web服务.它的域名是 web.service.consul :

```shell script
[root@hdp2 consul.d]# dig @127.0.0.1 -p 8600 web.service.consul

; <<>> DiG 9.8.2rc1-RedHat-9.8.2-0.47.rc1.el6 <<>> @127.0.0.1 -p 8600 web.service.consul
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 46501
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;web.service.consul.           IN         A

;; ANSWER SECTION:
web.service.consul.        0          IN         A          10.0.0.52

;; Query time: 0 msec
;; SERVER: 127.0.0.1#8600(127.0.0.1)
;; WHEN: Wed Aug 17 19:07:05 2016
;; MSG SIZE  rcvd: 70

```

如你所见,一个A记录返回了一个可用的服务所在的节点的IP地址.`A记录只能设置为IP地址. 有也可用使用 DNS API 来接收包含 地址和端口的 SRV记录:

```shell script

[root@hdp2 ~]# dig @127.0.0.1 -p 8600 web.service.consul SRV

; <<>> DiG 9.8.2rc1-RedHat-9.8.2-0.47.rc1.el6 <<>> @127.0.0.1 -p 8600 web.service.consul SRV
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 33415
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;web.service.consul.           IN         SRV

;; ANSWER SECTION:
web.service.consul.        0          IN         SRV        1 1 80 hdp2.node.dc1.consul.

;; ADDITIONAL SECTION:
hdp2.node.dc1.consul.      0          IN         A          10.0.0.52

;; Query time: 1 msec
;; SERVER: 127.0.0.1#8600(127.0.0.1)
;; WHEN: Thu Aug 18 10:40:48 2016
;; MSG SIZE  rcvd: 130

```

SRV记录告诉我们 web 这个服务运行于节点hdp2.node.dc1.consul 的80端口. DNS额外返回了节点的A记录.

最后,我们也可以用 DNS API 通过标签来过滤服务.基于标签的服务查询格式为TAG.NAME.service.consul. 在下面的例子中,我们请求Consul返回有 rails标签的 web服务.我们成功获取了我们注册为这个标签的服务:

```shell script
[root@hdp2 ~]# dig @127.0.0.1 -p 8600 rails.web.service.consul SRV

; <<>> DiG 9.8.2rc1-RedHat-9.8.2-0.47.rc1.el6 <<>> @127.0.0.1 -p 8600 rails.web.service.consul SRV
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 3517
;; flags: qr aa rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;rails.web.service.consul.         IN         SRV

;; ANSWER SECTION:
rails.web.service.consul. 0        IN         SRV        1 1 80 hdp2.node.dc1.consul.

;; ADDITIONAL SECTION:
hdp2.node.dc1.consul.      0          IN         A          10.0.0.52

;; Query time: 1 msec
;; SERVER: 127.0.0.1#8600(127.0.0.1)
;; WHEN: Thu Aug 18 11:26:17 2016
;; MSG SIZE  rcvd: 142

```

### HTTP API

除了DNS API之外,HTTP API也可以用来进行服务查询:

```shell script

[root@hdp2 ~]# curl http://localhost:8500/v1/catalog/service/web
[{"Node":"hdp2","Address":"10.0.0.52","ServiceID":"web","ServiceName":"web","ServiceTags":["rails"],"ServiceAddress":"","ServicePort":80,"ServiceEnableTagOverride":false,"CreateIndex":4,"ModifyIndex":254}]

```

目录API给出所有节点提供的服务.稍后我们会像通常的那样带上健康检查进行查询.就像DNS内部处理的那样.这是只查看健康的实例的查询方法:

```shell script
[root@hdp2 ~]# curl http://localhost:8500/v1/catalog/service/web?passing
[{"Node":"hdp2","Address":"10.0.0.52","ServiceID":"web","ServiceName":"web","ServiceTags":["rails"],"ServiceAddress":"","ServicePort":80,"ServiceEnableTagOverride":false,"CreateIndex":4,"ModifyIndex":254}]

```

## 更新服务

服务定义可以通过配置文件并发送SIGHUP给agent来进行更新.这样你可以让你在不关闭服务或者保持服务请求可用的情况下进行更新.

另外 HTTP API可以用来动态的添加,移除和修改服务.


