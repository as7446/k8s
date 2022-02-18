# API server

kube-apiserver是kubernetes最重要的核心组件之一，主要提供一下功能：
* 提供集群管理的REST API接口，包括认证、授权、数据校验以及集群状态变更等；
* 提供其他模块或者说组件之间的数据交互和通信枢纽（其他模块通过API Server查询或修改数据，只有API Server才能直接操作集群）。
