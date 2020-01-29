# 概念

## Kubernetes 架构

一个 Kubernetes 集群由 Master Node（主节点） 和 Worker Node （工作节点）两种节点组成，而这两种角色分别对应着控制节点和计算节点。

1. 在一个高可用的集群当中，Master 节点和 Worker 节点一般都是由多个节点构成，这些节点可以是物理机，也可以是虚拟机。

2. Worker 节点提供的资源单位为 Pod，简单理解 Pod 是 kubernetes 集群提供的虚拟机，是最小的调度单位。

3. Pod 里一般只包含一个应用容器，但是其他场景下（典型的如 Service Mesh） Pod 可以包含一个主容器或多个辅助容器。

下面是一张 Kubernetes 架构概览示意图：

![kubernetes-cluster](https://user-images.githubusercontent.com/8086910/73325011-d896c880-4287-11ea-8ce0-0090a504a14a.png)

| 节点   | 作用                                          |
| ------ | --------------------------------------------- |
| Master | 管理和调度集群资源                            |
| Worker | 提供集群资源， 是 kubernetes 集群资源的提供者 |


其中，控制节点，即 Master 节点，由三个紧密协作的独立组件组合而成，它们分别是负责 API 服务的 kube-apiserver、负责调度的 kube-scheduler，以及负责容器编排的 kube-controller-manager。整个集群的持久化数据，则由 kube-apiserver 处理后保存在 Etcd 中。而计算节点上最核心的部分，则是一个叫作 kubelet 的组件。

## Master 节点核心组件

### Kubernetes API server

所有的 K8s 操作都是通过 API Server， API 通过标准的 HTTP Web Service 实现，包含了 Rest 与 WebSocket 等等 API 设计

### Kubernetes Controller Manager

Kubernetes Controller Manager 在 Kubernetes Master 中，负责所有的控制功能。

### kube-scheduler

主要负责资源(Pod)调度，每个 Pod 最终被调度到哪台服务器上是由 Scheduler 决定

### ectd

Kubernetes 集群所有的资源对象的数据以及状态保存在 etcd 中， 它可以理解成是 Kubernetes 集群的数据库

## Node 节点组件

### Kubelet

Kubelet 安装在每一个 Node 上，负责与 API Server 沟通。

同时主要负责同容器运行时（比如 Docker 项目）打交道。而这个交互所依赖的，是一个称作 CRI（Container Runtime Interface）的远程调用接口，这个接口定义了容器运行时的各项核心操作，比如：启动一个容器需要的所有参数。

### kube-proxy

实现 Kubernetes 上 Service 的通信及负载均衡
