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

### 4.1 核心技术概念和 API 对象

API 对象是 Kubernetes集群中的管理操作单元。

Kubernetes 集群系统每支持一项新功能，引入一项新技术，一定会新引入对象的 API 对象，支持对核心功能的管理操作。

每个 API 对象都有四大类属性：

* TypeMeta
* MetaData 
* Spec
* Status

##### 4.1.1 TypeMeta

Kubernetes 对象的最基本定义，它通过引入 GKV （Group、Kind、Version）模型定义了一个对象的模型。

1、Group

Kubernetes 定义了非常多的对象，如何将这些对象进行归类是一门学问，将对象依据其功能范围归入不同的分组，比如把支撑最基本功能的对象归入 core 组，把与应用部署有关的对象归入 apss 组，会使这些对象的可维护性和可理解性更高。

2、Kind

定义一个对象的基本类型，比如 Node、Pod、Deployment等。

3、Version

社区每个季度会推出一个 Kubernetes 版本，随着 Kubernetes 版本的演进，对象从创建之初到能够完全生产化就绪的版本是不断变化的。与软件版本类似，通常社区提出一个模型定义以后，随着该对象不断成熟，其版本可能会从v1alpha1到v1alpha2，或者到v1beta1，最终变成生产就绪版本 v1。

##### 4.1.2 Metata

Metadata 中有两个最重要的属性：Namespace 和 Name ，分别定义了对象的 Namespace 归属名字，这两个属性唯一定义了某个对象实例。

1、Label

顾名思义就是给对象打标签，一个对象可以有任意对标签，其存在形式是键值对。Label 定义了对象的可识别属性， Kubernetes API 支持以 Label 作为过滤条件查询对象。

* Label 是识别 Kubernetes 对象的标签，以 key/value 的方式附加在对象上。
* key 最长不能超过 63 字节，value 可以为空，也可以是不超过 253 字节的字符串。
* Label 不提供唯一性，并且实际上经常是很多对象（如 Pods）都使用相同的 label 来标志具体应用。
* Label 定义好后其他对象可以使用 Label Seletor 来选择一组相同 label 的对象。
* Label Seletor 支持一下几种方式：
  * 等式，如 app=nginx 和 env!=production;
  * 集合，如 env in (production, qa);
  * 多个 label（它们之间是 AND 关系），如 app=nginx，env=test。

2、Annotation

Annotation 与 Label 一样用键值对来定义，但 Annotation 是作为属性扩展，更多面向对于系统管理员和开发人员，因此需要像其他属性一样做合理归类。

* Annotations 是 key/value 形式附加于对象的注解。
* 不用于 Labels 用于标志和选择对象，Annotations 则是用来记录一些附加信息，用辅助应用部署、安全策略以及调度策略等。
* 比如 deployment 使用 annotations 来记录rolling update的状态。

3、Finalizer

Finalizer 本质上是一个资源锁， Kubernetes在接收某个对象删除请求时，会检查 Finalizer 是否为空，如果不为空则只对其做逻辑删除，即只会更新对象中的 metadata.deletionTimestamp 字段。

4、ResourceVersion

ResourceVersion 可以被看作一种乐观锁，每个对象在任意时刻都有其 ResourceVersion， 当 Kubernetes 对象被客户端读取以后，ResourceVersion 信息也被一并读取。此机制确保了分布式系统中任意多线程能够无锁并发访问对象，极大提升了系统的整体效率。

##### 4.1.3 Spec 和 Status

* Spec 和 Status 才是对象的核心。
* Spec 是用户的期望状态，由创建对象的用户端来定义。
* Status 是对象的实际状态，由对应的控制器收集实际状态并更新。
* 与 TypeMeta 和 Metadata 等通用属性不同，Spec 和 Status 是每个对象独有的。

### 4.2 Node

* Node 是 Pod 真正运行的主机，可以物理机，也可以是虚拟机。
* 为了管理 Pod，每个 Node 节点上至少要运行 container runtime（比如 Docker 或者 Rkt）、Kubelet 和 Kube-proxy 服务

### 4.3 Namespace

Namespace 是对一组资源和对象的抽象集合，比如可以用来将系统内部的对象划分为不同的项目和用户组。

常见的 Pods，servers、replication controllers 和 deployments 等都是属于某一个 Namespace 的（默认是default），而 Node，persistentVolumes 等则不属于任何 Namespace。

### 4.4 什么是 Pod

* Pod 是一组紧密关联的容器集合，它们共享 PID、IPC、Network 和 UTS namespace， 是 Kubernetes 调度的基本单位。
* Pod 的设计理念是支持多个容器在一个 Pod 共享网络和文件系统，可以通过进程间通信和文件共享这种简单高效的方式组合完成服务。
* 同一个 Pod 中的不同容器可以共享资源：
  * 共享网络 Namespace；
  * 可通过挂在存储卷共享；
  * 共享 Security Context。

通过 Pod 对象定义支撑应用运行

* 环境变量：
  * 直接设置值；
  * 读取 Pod Spec 的某些属性；
  * 从 ConfigMap 读取某个值；
  * 从Secret 读取某个值；
* 存储卷
  * 通过存储卷可以将外挂存储挂载到 Pod 内部使用。
  * 存储卷定义包括两个部分： Volume 和 VolumeMouts。
    * Volume：定义 Pod 可以使用的存储卷来源；
    * VolumeMouts：定义存储卷如何挂载到容器内部。

* Pod 网络

​		Pod 的多个容器是共享网络 Namespace 的，这意味着：

​		*	同一个 Pod 中的不同容器可以彼此通过 Lookback 地址访问：

​			*	在第一个容器中起了一个服务http://127.0.0.1。

​			*	在第二个容器内，是可以通过 httpGet http://127.0.0.1 访问到该地址。

​		这种方法常用于不同容器的相互协作。

* 资源限制

​	Kubernetes通过 Cgroups 提供容器资源管理的功能，可以限制每个容器的 CPU 和内存使用，比如对于刚才创建的 deployment， 可以通过下面的命令限制 nginx 容器最多只用 50% 的 CPU 和 128MB 的内存：

$ kubectl set resources deployment nginx-app -c=nginx --limits=cpu=50m,memory=128Mi

* 健康检查

  Kubernetes 作为一个面向应用的集群管理工具，需要确保容器在部署后确实处于正常的运行状态。

  1、探针类型：

  * LivenessProbe
    * 探测应用是否处于健康状态，如果不健康则删除并重新创建容器。
  * ReadinessProbe
    * 探测应用是否就绪并且处于正常服务状态，如果不正常则不会接收来自 Kubernetes Service 的流量。
  * startupProbe 
    * 探测应用是否启动完成，如果在 failureThreshold*periodSeconds 周期内未就绪，则应用进程会被重启。

  2、探活方式：

  * Exec
  * TCP socket
  * HTTP

### 4.5 configMap

* configMap 用来将非机密性的数据保存到键值对中。
* 使用时， Pod 可以将其用作环境变量、命令行参数或者存储卷中的配置文件。
* ConfigMap 将环境配置信息和容器镜像解耦，便于应用配置的修改。

### 4.6 密钥对象（secret）

* Secret 是用来保存和传递密码、密钥、认证凭证这些敏感信息的对象。
* 使用 Secret 的好处是可以避免把敏感信息明文写在配置文件里。
* Kubernetes 集群中配置和使用服务不可避免的要用到各种敏感信息实现登录、认证登录功能，例如访问 AWS 存储的用户名密码。
* 为了避免将类似的敏感信息写在所有需要使用的配置文件中，可以将这些信息存储一个 Secret 对象，而在配置文件中通过 Secret 对象引用这些敏感信息。
* 这种方式的好处包括：意图明确，避免重复，减少暴露机会。

###  4.7 用户（User Account）&服务账户（Service Account）

* 顾名思义，用户账户为人提供账户标识，而服务账户计算进程和 Kubernetes 集群中运行的 Pod 提供账户标识。
* 用户账户和服务账户的一个区别是作用范围：
  * 用户账户对应的是人的身份，人的身份与服务的 Namespace 无关，所以用户账户是跨 Namespace 的；
  * 而服务账户对应的是一个运行中程序的身份，与特定的 Namespace 是相关的。

###  4.8 Service

Service 是应用服务的抽象，通过 Labels 为应用提供负载均衡和服务发现。

匹配 Labels 的 Pod IP 和端口列表组成 ednpoints，由 Kube-proxy 负责将服务 IP 负载均衡到这些 endpoints 上。

每个 Service 都会自动分配一个 CLuser IP （仅在集群内部可以访问的虚拟地址）和 DNS 名，其他容器可以通过该地址或 DNS 来访问服务，而不需要了解后端容器的运行。

### 4.9 副本集（Replica Set）

* Pod 只是单个应用实例的抽象，要构建高可用应用，通常要构建多个同样的副本，提供同一个服务。
* Kubernetes 为此抽象出副本集 Replicaset， 其允许用户定义 Pod 的副本数量，每个 Pod都会被当做一个无状态的成员进行管理，Kubernetes 保证总是有用户期望的数量的 Pod 正常运行。
* 当某个副本宕机以后，控制器将会创建一个新的副本。
* 当因业务负载发生变更而需要调整扩缩容是，可以方便的调整副本数量。

### 4.10 部署（Deployment）

* 部署标识用户对 Kubernetes 集群的一次更新操作。
* 部署是一个比 RS 应用模式更广的 API 对象，可以是创建一个新的服务，更新一个新的服务，也可以是滚动升级一个服务。
* 滚动升级一个服务，实际是创建一个新的 RS，然后逐渐将新 RS 中副本数量调整到理想状态，将就 Rs 中的副本数减少到0的复合操作。
* 这样一个复合操作用一个 RS 是不太好描述的，所以用一个更通用的 Deployment来描述。
* 以 Kubernetes 的发展方向，未来对所有长期伺服型的业务的管理，都会通过 Deployment 管理。

### 4.11 有状态服务集（StatefulSet）

* 对于 SatefulSet 中的 Pod，每个 Pod 挂载自己独立的存储，如果一个 Pod 出现故障，从其他节点启动一个同样名字的 Pod，要挂载上原来 Pod 的存储继续以它的状态提供服务。
* 适合于 StatefulSet 的业务包括数据库服务 MySQL 和 PostgreSql，集群化管理服务 Zookeeper、etcd 等有状态服务。
* 使用 StatefulSet，Pod 仍然可以通过漂移到不同节点提供高可用，而存储也可以通过外挂的存储来提高可靠性，StatefulSet 做的只是将确定的 Pod 与确定的存储关联起来保证状态的连续性。

### 4.12 Statefulset 与 Deployment的差异

* 身份标识
  * StatefulSet Controller为每个 Pod 编号，序号从 0 开始
* 数据存储
  * Statefulset 允许用户定义 volumeClaimTemplate， Pod 被创建的同时，Kubernetes 会以volumeClaimTemplate 中定义的模板创建存储卷，并挂载 Pod。
* StatefulSet 的升级策略不同
  * onDelete
  * 滚动升级
  * 分片升级

### 4.13 任务（Job）

* Job 是 Kubernetes用来控制批处理型任务的 APi 对象。
* Job 管理的 Pod 根据用户的设置把任务成功完成后就自动退出。
* 成功完成的标志根据不同的 spec.complatetions 策略而不同：
  * 单 Pod 型任务有一个 Pod 成功就标志成功；
  * 定数成功型任务保证有 N 个任务全部成；
  * 工作队列型任务根据应用确认的全局成功而标志成功。

### 4.14 后台支撑服务集（DaemonSet）

* 长期伺服型和批处理型服务的核心在业务应用，可能有些节点运行多个同类业务的 Pod，有些节点上又没有这类 Pod 运行；
* 而后台支撑型服务的核心关注点在 Kubernetes 集群中的节点（物理机或虚拟机），要保证每个节点上都有一个此同类 Pod 运行。
* 节点可能是所有集群节点也可能是通过 nodeSelector 选定的一些特定节点。
* 典型的后台支撑型服务包括存储、日志和监控等在每个节点上支撑 Kubernetes 集群运行的服务。

### 4.15 存储 PV 和 PVC

* PersistentVolume（PV）是集群中的一块存储卷，可以由管理员手动设置，或当用户创建 PersistenVolumeClaim（PVC）时根据 StorageClass 动态设置。
* PV 和 PVC 与 Pod 生命周期无关。也就是说，当 Pod 中的容器重启、重新调度或者删除时，PV 和 PVC 不会受到影响， Pod 存储于 PV 里的数据得以保留。
* 对于不同的使用场景，用户通常需要不同属性（例如性能、访问模式等）的 PV。



### 4.16 CustomResourceDefinition 

* CRD 就像数据库的开放式表结构，允许用户自定义 Schema。
* 有了这种开放式设计，用户可以基于 CRD 定义一些需要的模型，满足不同业务的需求。
* 社区鼓励基于 CRD 的业务抽象，众多主流的扩展应用都是基于 CRD 构建的，比如 Istio、Knaive。
* 甚至基于 CRD 推出了 Operator Mode 和 Operator SDK，可以以极低的开发成本定义新对象，并构建新对象的控制器。
