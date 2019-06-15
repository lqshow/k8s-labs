# Overview

Kubernetes 的资源对象非常多，但是这不是重点。这里有必要提一下，要想耍好 Kubernetes，编写配置文件这个技能必须要掌握，你可以简单的认为每一个资源对象都是一份不同的配置。

换句话来说，平时使用 Kubernetes 的过程就是是基于 Yaml 在编程，当然利用 sdk api 调用的除外。

## Role of a Developer

首先作为一个开发者，肯定最关心的是如何将现有的应用程序以最小的成本迁移到 Kubernetes 集群中，下面先来讲下 Kubernetes 中与开发者关系最密切的几个资源对象。

### 配置文件相关

- [Secret](./01-secret)
- [ConfigMap](./00-configmap)
- [Volumes](./02-volume)

### 网络访问相关

- [Service](./05-service)
- [Ingress](./06-ingress)

### 任务相关

- [Job](./07-job)
- [CronJob](./07-job)

### 日志文件

不建议将日志写入指定的文件中

1. 运维人员管理这么多应用，谁还记得你的应用将日志写哪去了，排查问题还要找到应用文档，找到日志目录，很浪费时间
2. 分布式部署，需保证日志的连续性，日志应该是事件流的汇总，应该按照时间顺序去收集

具体参见：[如何在 Kubernetes 上建立 Elasticsearch、Fluentd 和 Kibana (EFK)日志栈？](../08-efk-logging-stack)

## Role of a SRE

作为运维开发工程师，需要负责维护并确保整个服务的高可用性，下面几个 Kubernetes 资源对象需要去了解。

### 控制器编排相关

- [Deployment](./04-deployment)
- [Statefulset](./08-statefulset)
- [DaemonSet](./11-daemonset)

### 账号角色相关

- [Service Accounts](./09-service-account)
- [RBAC](./10-rbac-authorization)
