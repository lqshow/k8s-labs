# Overview

- Pod 是 Kubernetes 项目中最小的 API 对象（是 Kubernetes 项目的原子调度单位）
- 一组功能相关的 Container 的封装
- 共享存储和 Network Namespace
- k8s 调度和作业运行的基本单位(Scheduler 调度， Kubelet 运行)
- Pod 就是 Kubernetes 世界里的“应用”；而一个应用，可以由多个容器组成。
- Pod 这个看似复杂的 API 对象，实际上就是对容器的进一步抽象和封装而已。


Table of Contents
=================

   * [Overview](#overview)
   * [Deploy the app to Kubernetes](#deploy-the-app-to-kubernetes)
   * [Check that the Pods and Service are created](#check-that-the-pods-and-service-are-created)
   * [Pod 级别属性](#pod-级别属性)
       * [nodeName](#nodename)
       * [HostAliases](#hostaliases)
       * [NodeSelector](#nodeselector)
       * [shareProcessNamespace(TODO)](#shareprocessnamespacetodo)
       * [volumes](#volumes)
         
   * [Container 级别属性](#container-级别属性)
       * [ImagePullPolicy](#imagepullpolicy)
       * [Lifecycle](#lifecycle)
       * [健康检查](#健康检查)
            * [readinessProbe(业务探针)](#readinessprobe业务探针)
            * [livenessProbe(存活探针)](#livenessprobe存活探针)
            * [Configure Probes](#configure-probes)
   * [Reference](#reference)

# Overview

1. Pod 是 Kubernetes 里“最小”的 API 对象（原子调度单位）
2. 由于 Pod 是“最小”的对象，所以它往往都是被其他对象控制的。这种组合方式，正是 Kubernetes 进行容器编排的重要模式。
3. Pod 等价为 Kubernetes 世界里的“应用”；而一个应用，可以由多个容器组成。
4. Pod 里的所有容器，共享的是同一个 Network Namespace，并且可以声明共享同一个 Volume。Volume 的定义在 Pod 级别。
5. Pod 的生命周期只跟 Infra 容器一致，而与容器无关。

# Deploy the app to Kubernetes

```bash
kubectl apply -f nginx-pod.yml
kubectl apply -f nginx-svc.yml
```

# Check that the Pods and Service are created

```bash
kubectl get po -o wide -l app=nginx
kubectl get svc -o wide -l app=nginx
```

# Pod 级别属性
> 调度、网络、存储，以及安全相关的属性

### nodeName
调度到指定的 Node上
```yaml
apiVersion: v1
kind: Pod
...
spec:
 nodeName: k8s02
```

### HostAliases
定义了 Pod 的 hosts 文件（/etc/hosts）里的内容
```bash
kubectl apply -f hostaliases-pod.yml
```
查看 Pod 内的 hosts 内容
```bash
➜   kubectl logs hostaliases-pod
# Kubernetes-managed hosts file.
127.0.0.1       localhost
::1     localhost ip6-localhost ip6-loopback
fe00::0 ip6-localnet
fe00::0 ip6-mcastprefix
fe00::1 ip6-allnodes
fe00::2 ip6-allrouters
10.1.2.130      hostaliases-pod

# Entries added by HostAliases.
127.0.0.1       foo.local
127.0.0.1       bar.local
10.1.2.3        foo.remote
10.1.2.3        bar.remote
```
### NodeSelector
将 Pod 与 Node 进行绑定的字段

以下配置，意味着这个 Pod 永远只能运行在携带了“disktype: ssd”标签（Label）的节点上；否则，它将调度失败。
```yaml
apiVersion: v1
kind: Pod
...
spec:
 nodeSelector:
   disktype: ssd
```

### shareProcessNamespace(TODO)
将 shareProcessNamespace 设置为 true 时， 这个 Pod 里的容器要共享同一个 PID Namespace。
```bash
kubectl create -f share-process-namespace-pod.yml
```
```bash
kubectl exec -it share-process-namespace-pod -c shell -- sh
```

### volumes
一个 Volume 对应的宿主机目录对于 Pod 来说就只有一个，Pod 里的容器只要声明挂载这个 Volume，就一定可以共享这个 Volume 对应的宿主机目录。

# Container 级别属性

### ImagePullPolicy
定义了镜像拉取的策略

| value         | desc                                                         |
| ------------- | ------------------------------------------------------------ |
| 默认是 Always | 即每次创建 Pod 都重新拉取一次镜像。 |
| Never         | Pod 永远不会主动拉取这个镜像                                 |
| IfNotPresent  | 只在宿主机上不存在这个镜像时才拉取                           |

### Lifecycle
定义的是 Container Lifecycle Hooks，是在容器状态发生变化时触发一系列“钩子”

### 健康检查
默认情况下，Kubernetes 会开启进程的健康检查：如果 Kubernetes 检测到容器里面的进程退出了，那么它就会重启这个容器。
但是如果程序发生死锁了，此时进程仍然在运行，但应用其实已经不工作了。Kubernetes 提供了另外三种类型的健康检查

| key  | desc                                                         |
| ---- | ------------------------------------------------------------ |
| exec | 在容器里面执行一个命令，如果返回 0 则认为是成功的，如果返回非0值，kubelet 就会杀掉这个容器并重启它 |
| tcp  | 与容器的某个 port 建立 tcp 连接，如果连接建立成功，容器被认为是健康的 |
| http | 向容器的 http 接口发起 http 请求，如果返回的状态码（大于200小于400的返回码）是成功的，kubelet 就会认定该容器是活着的并且很健康。 |

#### readinessProbe(业务探针)

> 确定容器是否已经就绪可以接受流量，主要控制哪些 Pod 可以作为 Service 的 endpoints

探针不正常后，不会重启容器，只会拿掉 Service 后端的 endpoints

#### livenessProbe(存活探针)
> 确定何时重启容器

探针不正常后，会重启容器

#### Configure Probes 

| key                 | desc                                                         |
| ------------------- | ------------------------------------------------------------ |
| initialDelaySeconds | 容器启动后在第一次执行探测之前需要等待 X 秒钟                |
| timeoutSeconds      | 探测超时时间。为 HTTP 请求设置超时时间为 X 秒                |
| periodSeconds       | 执行探测的频率。每 X 秒执行一次健康检查                      |
| failureThreshold    | 设置失败的阈值，如果超过 X 次健康检查都失败了，容器就会被重启。 |

## 调度策略

> 为 Pod 找到一个合适的 Node

一般情况下我们部署的 Pod 是通过集群的自动调度策略来选择节点的，默认情况下调度器考虑的是资源足够,并且负载尽量平均.

## References

- [k8s 使用 Pod 来部署应用](https://github.com/lqshow/notes/issues/38)
- [配置Pod的liveness和readiness探针](https://jimmysong.io/kubernetes-handbook/guide/configure-liveness-readiness-probes.html)
