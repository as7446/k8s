# kubernetes 控制平面组件

## 目录

* <a href="#P1">调度</a>

* <a href="#P2">Controller Manager</a>

* <a href="#P3">kubelet</a>

* <a href="#P4">CRI</a>

* <a href="#P5">CNI</a>

* <a href="#P6">CSI</a>

  



<a name="P1"></a>

## PART1. 调度

### kube-scheduler

kube-scheduler负责分配调度Pod到集群内的节点上，他监听kube-apiserver，查询未分配Node的pod， 然后根据调度策略为这些Pod分配节点（更新Pod的NodeName字段）。

调度器需要充分考虑多的因素：

* 公平调度；
* 资源高效利用；
* Oos
* affinity 和 anti-affinity;
* 数据本地化（data locality）；
* 内部负载干扰（initer-workload interference）；
* deadlines。



### 调度器

kube-scheduler调度分为两个阶段，predicate和priority；

* predicate：过滤不符合条件的节点；
* priority：优先级排序，选择优先级高的节点。



### Predicates 策略

![image-20220218182633909](C:\Users\pc\AppData\Roaming\Typora\typora-user-images\image-20220218182633909.png)

![image-20220218182445733](C:\Users\pc\AppData\Roaming\Typora\typora-user-images\image-20220218182445733.png)



### Predicates plugin 工作原理

![image-20220218182814482](C:\Users\pc\AppData\Roaming\Typora\typora-user-images\image-20220218182814482.png)



### Priorities 策略

* SelectorSpreadPriority：优先减少节点上属于通一个 Service 或 Replication Controller 的 Pod 数量。

* InterPodAffinityPriority：优先将 Pod调度到相同的拓扑上（如同一个节点、Rack、Zone等）。
* LeastRequestedPriority：优先调度到请求资源少的节点上。
* BalancedResourceAllocation：优先平衡各节点的资源使用。
* NodePreferAvoidPodsPriority： alpha.kubernetes.io/preferAvoidPods 字段判断，权重为 10000，避免其他优先级策略的影响。
* NodeAffinityPriority：优先调度到匹配的 NodeAffinity 的节点上。
* TainTolerAtionPriority：优先调度到匹配 TaintToleration 的节点上。
* ServiceSpreadingPriority： 尽量将同一个 service 的 Pod 分布到不同的节点上，已经被SelectorSpreadPriority 代替（默认未使用）。
* EqualPriority： 将大镜像的容器调度到已经下拉了该镜像的节点上（默认未使用）。
* MostRequestedPriority： 尽量调度到已经使用过的 Node 上，特别适用于 cluster-autoscaler （默认未使用）。

