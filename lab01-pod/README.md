# Deploy the app to Kubernetes

```bash
kubectl create -f nginx-pod.yml
kubectl create -f nginx-svc.yml
```

# Check that the Pods and Service are created

```bash
kubectl get po -o wide -l app=nginx
kubectl get svc -o wide -l app=nginx
```

## Pod 级别属性
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
kubectl create -f hostaliases-pod.yml
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

## Container 级别属性

### ImagePullPolicy
定义了镜像拉取的策略

| value         | desc                                                         |
| ------------- | ------------------------------------------------------------ |
| 默认是 Always | 即每次创建 Pod 都重新拉取一次镜像。 |
| Never         | Pod 永远不会主动拉取这个镜像                                 |
| IfNotPresent  | 只在宿主机上不存在这个镜像时才拉取                           |

### Lifecycle
定义的是 Container Lifecycle Hooks，是在容器状态发生变化时触发一系列“钩子”

## 健康检查
### readinessProbe(业务探针)
探针不正常后，不会重启容器，只会拿掉服务后端的 endpoints

### livenessProbe(存活探针)
探针不正常后，会重启容器

# Reference
- [k8s 使用 Pod 来部署应用](https://github.com/lqshow/notes/issues/38)

