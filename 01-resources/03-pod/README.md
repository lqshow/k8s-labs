# Overview

- Pod 是 Kubernetes 项目中最小的 API 对象（是 Kubernetes 项目的原子调度单位）
- 一组功能相关的 Container 的封装
- 共享存储和 Network Namespace
- k8s 调度和作业运行的基本单位(Scheduler 调度， Kubelet 运行)
- Pod 就是 Kubernetes 世界里的“应用”；而一个应用，可以由多个容器组成。
