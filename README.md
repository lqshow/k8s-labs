# Overview

Kubernetes 是 Google 开源的容器集群管理系统，是 Google 多年大规模容器管理技术 Borg 的开源版本，也是 CNCF 最重要的项目之一。

## Why use Kubernetes

首选我们来看下 Kubernetes 能够解决的问题是什么？

- 编排
- 调度
- 容器云
- 集群管理

Kubernetes 的主要功能有以下几块

- 基于容器的应用部署、维护和滚动升级
- 内置了负载均衡和服务发现
- 跨机器和跨地区的集群调度
- 实现容器应用的自动伸缩
- 无状态服务和有状态服务
- 广泛的 Volume 支持
- 插件机制保证扩展性

## Step by step play with Kubernetes

- [00-concepts](./00-concepts)：介绍 Kubernetes 的基本概念
- [01-resources](./01-resources)：介绍 Kubernetes 的常用资源对象
- [02-commands](./02-commands)：介绍 Kubernetes 常用命令
- [03-networking](./03-networking)：介绍 Kubernetes 网络访问情况
- [04-sharing-nfs-persistent-volume-in-kubernetes](./04-sharing-nfs-persistent-volume-in-kubernetes)：Kubernetes 容器持久化存储举例(nfs)
- [05-container-design-patterns-for-kubernetes](./05-container-design-patterns-for-kubernetes)：介绍 Kubernetes 常用容器设计模式
- [06-creating-kubeconfig-for-cluster](./06-creating-kubeconfig-for-cluster)：如何切换并操作不同的 Kubernetes 集群？
- [07-migrating-app-to-kubernetes-best-practices](./07-migrating-app-to-kubernetes-best-practices)：如何将传统应用迁移到 Kubernetes 集群中？
- [08-efk-logging-stack](./08-efk-logging-stack)：如何在 Kubernetes 上建立 Elasticsearch、Fluentd 和 Kibana (EFK)日志栈？
- [09-helm-best-practices](./09-helm-best-practices)： 介绍 Kubernetes 包管理器入门实践
- [10-kustomize](./10-kustomize)：TBD
- [11-dev-tools](./11-dev-tools)：TBD
- [12-traefik](./12-traefik)：TBD
- [13-traefik](./13-monitoring)：TBD
- [14-troubleshooting](./14-troubleshooting)：使用 Kubernetes 中碰到的一些常见问题以及解决方案。
