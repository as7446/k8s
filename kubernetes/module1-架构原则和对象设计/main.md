##	1. Kubernetes 架构基础

### 1.1 什么是Kubernetes（k8s）

Kubernetes 是谷歌开源的容器集群管理系统，是Google 多年大规模容器管理技术Borg的开源版本，主要功能包括：

* 基于容器的应用部署、维护和滚动升级；
* 负载均衡和服务发现；
* 跨机器和跨地区的集群调度；
* 自动伸缩；
* 无状态服务和有状态服务；
* 插件机制保证扩展性；

### 1.2 Kubernetes：声明式系统

它本身是构建于声明式系统之上的云管理平台，遵循了声明式系统以后它把Kubernetes里所有的代管对象，包括计算节点、服务、作业等等对象抽象成一个个标准API。把这些API作为统一的规范与一些大厂联合背书，现在主要的云计算厂商，阿里云、azure、腾讯都遵守同样的标准使用，那么它就成了一个事实标准，之后所有人都往这个标准上靠。通常标准只能是演进，被取代非常难，所以Kubernetes未来10年20年来看会一直存在并发展。

#### 1.2.1 kubernetes 的所有管理能力构建在对象抽象的基础上，核心对象包括：

* Node：计算节点的抽象，用来描述计算节点的资源抽象，健康状态等；
* Namespace：资源隔离的基本单位，可以简单理解为文件系统中的目录结构；
* Pod：用来描述应用实例，包括镜像地址、资源需求等。Kubernetes中最核心的对象，也是打通应用和基础架构的秘密武器；
* Service：服务如何将应用发布成服务，本质是负载均衡和域名服务的声明；



在每个Kubernetes的组件里面会有一个cAdvisor组件，它用来收集容器进程的健康状况、资源用量。

### 1.3 核心组件简述

#### 1.3.1 API server

kube-APIserver 是 Kubernetes 最重要的核心组件之一，主要提供一下功能：

* 提供集群管理的REST API 接口，包括：

  * 认证 Authentication；
  * 授权 Authorization；
  * 准入 Admission（Mutating & Valiating）。

  注意：上述三道门任何一道不过请求是不会被存储到etcd的

* 提供其他模块或者说组件之间的数据交互和通信的枢纽（其他模块通过 APIserver 查询或修改数据，只有 APIServe r才直接操作 etcd）。

* APIServer 提供etcd数据缓存以减少集群对 etcd 的访问。



APIHandler -> AuthN -> Rate LImit -> Auditing -> AuthZ -> Aggregator -> Mutating Webhook -> Schema Validation -> Validataing Webhook -> etcd 

APIserver 是Kubernetes控制面板中唯一带有用户可访问API以及用户可交互的组件。API服务器会暴露一个RESTful的Kubernetes API并使用JSON格式的清单文件。它本身是一个RESTSERVER，它注册了一些对象的Handler，当你这个对象操作任何对象的时候它实际转换成了一个REST调用，由API server接收，API server其实就是整个集群的控制面的API网关，对于任何的API网关他还有很多的附加功能。包括认证、鉴权、准入，就是说我要确保你这个请求客户端是合法的，如果是合法的，那我还要知道你有哪些操作权限，如果认证、鉴权都过了，那还要确认这个请求操作是不是合法的，比如这个请求是非法的不符合Kubernetes的规范，那就要由API server挡掉了。

#### 1.3.2 Controller Manager

* Controller Manger是集群的大脑，是确保整个集群动起来的关键；
* 作用是确保 Kubernetes 遵循声明式系统规范，确保系统的真实状态（Actual State）与用户定义的期望状态（Desired State）一致；
* Controller Manger 是多个控制器的组合，每个 Controller 事实上都是一个 control loop， 负责侦听其他管控的对象，每当对象发生变更时完成配置；
* Controller 配置失败通常会触发自动重试，整个集群会在控制器不断重试机制下确保最终一致性（Eventual Consistency）。

manager 它运行着所有处理集群日常任务的控制器，比如说deployment控制器、Replication控制器、NodeLifecycle的控制器。controller manager是让整个集群运作起来的一个核心，它是一个大脑。API server里面没有什么业务逻辑，它就接收请求，只要你有权限，这个请求有是合法的他就存储到etcd，但是这个请求被存下来以后，集体集群应该怎么动是由controller manager去做的。 

#### 1.3.3 Scheduler

特殊的 Controller， 工作原理与其他控制器无差别。

Scheduler 的特殊职责在于监控当前集群所有未调度的Pod，并且获取当前集群所有节点的健康状况和资源使用情况，为待调度 Pod 选择最佳计算节点，完成调度。

调度阶段分为：

* predict： 过滤不能满足业务需求的节点，如资源不足，端口冲突等。
* Priority：按既定要素将满足调度需求的节点评分，选择最佳节点。
* Bind：将计算节点与 Pod 绑定，完成调度。

调度器会监控新建的Pods（一组或一个容器）并将其分配给节点。

#### 1.3.4 Kubelet 

Kubernetes 的初始化系统（init system）

* 从不同源取 Pod 清单，并按需求启停 Pod 的核心组件：
  * Pod 清单可从本地文件目录，给定的 HTTPServer 或 Kube-APIServer 等源头取；
  * Kubelet 将运行时，网络和存储抽象成了 CRI、CNI、CSI。
* 负责汇报当前节点的资源信息和健康状况；
* 负责 Pod 的健康检查和状态汇报。

负责调度到对应的Pod的生命周期管理，执行任务并将 Pod 状态报告给主节点的渠道，通过容器运行时（拉取镜像、启动和停止容器等）来运行这些容器。他还会定期执行被请求的容器的健康探测程序。

#### 1.3.5 etcd

Kubernetes 使用 etcd。这是一个强大的、稳定的、高可用的键值存储，被Kubernetes用于长久存储所有的 API 对象。

#### 1.3.6 kube-proxy

他负责节点的网络，在主机上维护网络规则并执行连接转发。他还负责对正在服务的 Pods 进行负载均衡。当要去定义一个service或者说要发布一个服务的时候，要为这个服务配置负载均衡，这个是由kube-proxy去做的。

#### 1.3.7Pod

pod里面是实际就是多个容器组成，容器都是通过 Container 的 RuntimeService 去起的应用， 它都是通过标准的CNI接口，都有自己统一的cgroups的配置、namespace的配置，所有的这些容器启动都是有标准化的。

### 1.4 推荐的 Add-ones 

* kube-dns：负责为集群提供 DNS 服务；
* Ingress Controller：为服务提供外网入口；
* MetricsServer：提供资源监控；
* Dashboard： 提供GUI；
* Fluentd-Elasticsearch：提供集群日志采集、存储与查询。

### 1.5 Kubernetes 设计理念

* 高可用
  * 基于 Replicaset、statefulset 的应用高可用
  * Kubernetes 组件本身高可用
* 安全
  * 基于 TLS 提供服务
  * Serviceaccount 和 user
  * 基于 Namespace 隔离
  * secret
  * Taints，PSP，networkpolicy
* 可移植性
  * 多种 host os 选择
  * 多种基础架构的选择
  * 基于集群联邦建立混合云
* 可扩展性
  * 基于微服务部署应用
  * 横向扩容缩容
  * 自动扩容缩容

### 1.6 分层架构

* 核心层：Kubernetes 最核心的功能，对外提供 API 构建高层应用，对内提供插件式应用执行环境。
* 应用层：部署（无状态应用、有状态应用、批处理任务、集群应用）和路由（服务发现、DNS 解析等）。
* 管理层：系统度量（如基础设施、容器和网络的度量）、自动化（如自动扩展、动态 Provision等）、策略管理（RBAC、Quota、PSP、Networkpolicy等）。
* 接口层：kubectl 命令行工具、客户端 SDK 以及集群联邦。
* 生态系统：在接口层之上的庞大容器集群管理调度的生态系统，可以划分两个范畴：
  * Kubernetes 外部：日志、监控、配置管理、CI、CD、Workflow、Faas、OTS 应用、Chatops等；
  * Kubernetes 内部：CRI、CNI、CVI、镜像仓库、Cloud Provider、集群自身的配置和管理等。

### 1.7 API 设置原则

* 所有 API 都应是声明式的
  * 相对于命令式操作，声明式操作对于重复操作的效果是稳定的，这对于容易出现数据丢失或者重复的分布式环境来说是很重要的。
  * 声明式操作更易被用户使用，可以使系统向用户隐藏实现的细节，同时也保留了系统未来持续优化的可能性。
  * 此外，声明式的 API 还隐含了所有的 API 对象都是名词性质的，例如 Service、Volume 这些 API 都是名词，这些名词描述了用户所期望得到的一个目标对象。
* API 对象是彼此互补而且可以组合的
  * 这实际上鼓励 API 对象尽量实现面向对象设计是的要求，即“高内聚、低耦合”， 对业务相关的概念有一个合适的分解，提高分解出来的对象的可重用性。
* 高层 API 以操作意图为基础设计
  * 如何能够设计好 API，跟如何能用面向对象的方法设计好应用系统由相通的地方，高层设计一定是从业务出发，而不是过早的从技术实现出发。
  * 因此，针对 Kubernetes 的高层 API 设计，一定是以 Kubernetes 的业务为基础出发，也就是以系统调度管理容器的操作意图为基础设计。
* 底层的 API 根据高层 API 的控制需要设计
  * 设计实现底层 API 的目的，是为了被高层 API 使用，考虑减少冗余、提高重用性的目录，底层 API 的设计也要以需求为基础，要尽量抵抗受技术实现影响的诱惑。
* 尽量避免简单封装，不要有在外部 API 无法显式知道的内部隐藏的机制
  * 简单的封装，实际没有提供新的功能，反而增加了对所封装 API 的依赖性。
  * 例如 Statefulset 和 Replicaset ，本身就是两种 Pod 集合，那么Kubernetes 就用不同 API对象来定义他们，而不会只用同一个 Replicaset，内部在通过特殊算法再来区分这个 Replicaset 是有状态还是无状态。
* API 操作复杂度与对象数量成正比
  * API 的操作复杂度不能超过O(N)，否则系统就不具备水平伸缩性了。
* API 对象状态不能依赖网络连接状态
  * 由于众所周知，在分布式环境下，网络连接断开是经常发生的事，因此要保证 API 对象状态能应对网络的不稳定性， API 对象的状态就不能依赖于网络连接状态。
* 尽量避免让操作机制依赖于全局状态
  * 因为分布式系统中要保证全局状态同步是非常困难的。

### 1.8 架构设计原则

* 只有 APIServer可以直接访问 etcd 存储，其他服务必须通过Kubernetes API来访问集群状态；
* 单节点故障不应该影响集群的状态；
* 在没有新请求的情况下，所有组件应该在故障恢后继续执行上次最后收到的请求（比如网络分区或服务重启等）；
* 所有组件都应该在内存中保持所需的状态，APIserver 将状态写入 etcd 存储，而其他组件则通过 APIserver 更新并监听所有变化；
* 优先使用事件监听而不是轮询。 

### 1.9 引导（Bootstrapping）原则

* Self-hosting是目标。
* 减少依赖，特别是稳态运行的依赖。
* 通过分层的原则管理依赖。
* 循环依赖问题的原则：
  * 同时还接收其他地方的数据输入（比如本地文件等），这样在其他服务不可用时还可以手动配置引导服务；
  * 状态应该是可恢复或可重新发现的；
  * 支持简单的启动临时实例来创建稳态运行所需要的状态，使用分布式锁或文件锁等来协调不同状态的切换（通常称为 pivoting 技术）；
  * 自动重启异常退出的服务，比如副本或者进程管理器等。



##	2、了解kubectl

## 3、深入理解Kubernetes

## 4、核心对象概览
