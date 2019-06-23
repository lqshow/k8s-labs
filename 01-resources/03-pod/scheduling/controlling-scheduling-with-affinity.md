# Overview

亲和性有分成节点亲和性( nodeAffinity)和 Pod 亲和性( podAffinity)

亲和性调度可以分成软策略和硬策略两种方式:

1. 软策略就是如果你没有满足调度要求的节点的话，pod 就会忽略这条规则，继续完成调度过程，说白了就是满足条件最好了，没有的话也无所谓了的策略
2. 硬策略就比较强硬了，如果没有满足条件的节点的话，就不断重试直到满足条件为止，简单说就是你必须满足我的要求，不然我就不干的策略。

### preferredDuringSchedulingIgnoredDuringExecution

软策略

### requiredDuringSchedulingIgnoredDuringExecution

硬策略

### operator

| operator     | desc                     |
| ------------ | ------------------------ |
| In           | label 的值在某个列表中   |
| NotIn        | label 的值不在某个列表中 |
| Gt           | label 的值大于某个值     |
| Lt           | label 的值小于某个值     |
| Exists       | 某个 label 存在          |
| DoesNotExist | 某个 label 不存在        |

## NodeAffinity

NodeAffinity 节点亲和性，是 Pod上 定义的一种属性，使 Pod 能够按我们的要求调度到某个 Node 上

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: node-affinity
spec:
  containers:
    - name: pod
      image: nginx
  affinity:
    nodeAffinity:
      # 硬性要求不能运行在 k8s01 这个节点上
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: kubernetes.io/hostname
            operator: NotIn
            values:
            - k8s01
      # 如果存在 k8s02 这个节点的话，优先调度到这个节点上
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 1
        preference:
          matchExpressions:
          - key: kubernetes.io/hostname
            operator: In
            values:
            - k8s02
```

## podAffinity

Pod 亲和性主要解决 pod 可以和哪些 pod 部署在同一个拓扑域中的问题

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-affinity
spec:
  containers:
    - name: pod
      image: nginx
  affinity:
    podAffinity:
      # 硬性要求: 将当前 Pod 调度到运行了 app=nginx 的节点上
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - nginx
        topologyKey: kubernetes.io/hostname
```

### Notes

如果集群内不存在 app=nginx 的 pod，那我们的 Pod 就会一直处于 Pending 状态，因为采用的是硬策略。

```bash
Events:
  Type     Reason            Age               From               Message
  ----     ------            ----              ----               -------
  Warning  FailedScheduling  1m (x25 over 2m)  default-scheduler  0/4 nodes are available: 2 node(s) didn't match pod affinity rules, 2 node(s) didn't match pod affinity/anti-affinity, 2 node(s) had taints that the pod didn't tolerate.
```

## podAntiAffinity

Pod 反亲和性主要是解决 pod 不能和哪些 pod 部署在同一个拓扑域中的问题

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-affinity
spec:
  containers:
    - name: pod
      image: nginx
  affinity:
    podAntiAffinity:
      # 硬性要求: 如果一个节点上存在一个 app=nginx 的 Pod，当前 Pod 不会调度到这个节点上
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: app
            operator: In
            values:
            - nginx
        topologyKey: kubernetes.io/hostname
```

## References

- [终于明白了 K8S 亲和性调度](https://mp.weixin.qq.com/s/HBxyO9k615x9--BVawOnSw)
