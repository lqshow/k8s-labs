# Overview

TBD

1. NetworkPolicy 是一种关于 Pod 间及 Pod 与其他网络端点间所允许的通信规则的规范
2. 默认情况下集群里没有 Network Policy, 允许与 Pod 之间的所有网络连接
3. NetworkPolicy 本质上是一个允许连接的白名单
4. NetworkPolicy 是存在于特定的 namespace。在同一个 namespace 下，NetworkPolicy 的名称必须唯一，但是不限于其他 namespace
5. 因为 Pod IP 都是自动分配的，并且会频繁的变动，所以 NetworkPolicy 使用标签来选择 Pod 和 namespace

![np](https://user-images.githubusercontent.com/8086910/59970611-5df12680-959d-11e9-9485-7ab4cbf52d17.png)

## 支持 Network Policy 的网络插件

> Network Policy 只是 K8s 中的一个资源对象，具体限流是通过网络插件来实现

- Calico
- Cilium
- Weave Net
- Kube-router
- Romana

## Network Policy

> 它是通过标签来选择 Pods，然后定义具体的规则列表，该列表确定哪些类型的流量可到达所选的 Pods。
>
> podSelector 将从 NetworkPolicy 所属的 namespace 中选择 Pods，不能跨其他 namespace 选择 Pods
>
> ingress 和 egress 也遵从以上选择 pods 原则，除非同 namespace selector 搭配使用

Network Policy Spec 主要由四部分组成，除了 podSelector 是必选外，其他三个元素都是可选

1. podSelector: 一个 Pod 选择器（指示策略应用于哪个目标 Pod），它选择一组(零或多个) Pod。当一个 Pod 被一个网络策略选中时，这个网络策略被称为应用于它。
2. policyTypes: 指定该策略中包含哪些类型的策略
3. ingress: 一个入口规则列表（指示允许哪些入站流量）
4. egress: 一个出口规则列表（指示允许哪些出站流量）组成。

- 基于源 IP 的访问控制列表
    - 限制 Pod 的进/出流量

- Pod 网络隔离的一层抽像
    - label selector
    - namespace selector
    - port
    - CIDR(网段/subnet)

![1_dBQ1Zmp1fS1V9L52bYR8zQ](https://user-images.githubusercontent.com/8086910/59971465-27240c00-95af-11e9-99d8-24a51569e964.png)

注：
{} 代表允许所有
[] 代表拒绝所有

```yaml
spec:
  # 表示选择所有的 Pod
  podSelector: {}
  ingress:
  - from:
    # 表示选定所有 namespace
    - namespaceSelector: {}
  policyTypes:
  - Ingress
```

### 拒绝所有 Pod 之间的流量通信

> 这意味着只允许其他网络策略显式白名单的连接

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
spec:
  podSelector: {}
  policyTypes:
  - Ingress
```

### 拒绝所有流量进入

> 选中 app=web 这个 Pod 采用以下网络策略 

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: web-deny-all
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Ingress
```


### 限制部分流量进入

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: api-allow
spec:
  # 筛选需要隔离的 Pod，通过标签去筛选
  podSelector:
    matchLabels:
      app: bookstore
      role: api
  ingress:
  - from:
    # 指定能够进入的 Pod，通过标签去筛选
    - podSelector:
        matchLabels:
          app: bookstore
  policyTypes:
  - Ingress
```

比较常见的一个场景，比如一个数据库服务，只允许部分服务才能访问

1. 数据库服务标签： app=db
2. 显示指定能够访问数据库服务 Pod 的标签: networking/allow-db-access=true

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-db-access
spec:
  podSelector:
    matchLabels:
      app: "db"
  ingress:
  - from:
    - podSelector:
        matchLabels:
          networking/allow-db-access: "true"
  policyTypes:
  - Ingress
```

### 拒绝所有流量流出

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all-egress
spec:
  podSelector: {}
  policyTypes:
  - Egress
```

### 允许所有流量进入

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
    name: web-allow-all
spec:
  podSelector:
    matchLabels:
      app: web
  ingress:
  - {}
  policyTypes:
  - Ingress
```

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: web-allow-all
spec:
  podSelector:
    matchLabels:
      app: web
  ingress:
  - from:
    podSelector: {}
    namespaceSelector: {}
  policyTypes:
  - Ingress
```

### 允许特定 namespace 的 Pod 流量进入

> 跨 namespace 访问, 前提需要给 source namespace 打上标签

```bash
kubectl label namespace <name> networking/namespace=<name>
```

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-n1-a-to-n2-b
  namespace: N2
spec:
  podSelector:
    matchLabels:
      deployment-b-pod-label-1-key: deployment-b-pod-label-1-value
      deployment-b-pod-label-2-key: deployment-b-pod-label-2-value
  policyTypes:
  - Ingress
  ingress:
  - from:
    -  namespaceSelector:
        matchLabels:
          networking/namespace: N1
       podSelector:
        matchLabels:
          deployment-a-pod-label-1-key: deployment-a-pod-label-1-value
          deployment-a-pod-label-2-key: deployment-a-pod-label-2-value
```

### 限制流量从指定端口进入

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: web-allow-5000
spec:
  podSelector:
    matchLabels:
      app: api
  ingress:
  - ports:
    - port: 5000
      protocol: TCP
    from:
    - podSelector:
        matchLables:
          role: monitoring
  policyTypes:
  - Ingress
```

1. 隔离 name=web-allow-5000 的 pod
2. 允许带有 role=monitoring 标签的 pod 访问 pod 的 5000 TCP 端口的

### 组合策略

> 多个策略以逻辑或的方式进行

```yaml
# 在一个 from 下同时限定两组 Pod 流量进入
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default.postgres
  namespace: default
spec:
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: indexer
    - podSelector:
        matchLabels:
          app: admin
    ports:
      - port: 443
        protocol: TCP
      - port: 80
        protocol: TCP
  podSelector:
    matchLabels:
      app: postgres
  policyTypes:
  - Ingress
```

```yaml
# 在多个 from 里分别限定 Pod，作用同上
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default.postgres
  namespace: default
spec:
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: indexer
    ports:
     - port: 443
       protocol: TCP
  - from:
    - podSelector:
        matchLabels:
          app: admin
    ports:
     - port: 80
       protocol: TCP
  podSelector:
    matchLabels:
      app: postgres
  policyTypes:
  - Ingress
```

> 多个策略以逻辑与的方式进行

```yaml
# 只允许默认名称空间中的特定 pod 访问 postgres
# 即同时过滤 namespace 和 pod
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: database.postgres
  namespace: database
spec:
  podSelector:
    matchLabels:
      app: postgres
  ingress:
  - from:
    - namespaceSelector:
        matchLabels:
          namespace: default
      podSelector:
        matchLabels:
          app: admin
  policyTypes:
  - Ingress
```

## Notes

服务内的所有节点网络策略配置

### 需求

1. 集群内的其他应用(pod)不能访问服务内的所有节点(pod)
2. 服务内所有节点(pod)不能访问集群内的其他应用(pod)
3. 服务内的所有节点(pod)可以互相访问
4. 可能允许服务内的部分节点(pod)接收外部网络的流量


### 提供的解决方案

> 假设一个服务内的所有节点（Pod），在同一个 namespace 中

1. 将当前 namespace 的默认加上 default-deny-all 网络策略
    - 即当前 namespace 设置成拒绝所有 Pod 之间的流量通信
    - 同时也不能接受外部的流量

    ```yaml
    apiVersion: networking.k8s.io/v1
    kind: NetworkPolicy
    metadata:
      name: default-deny-all
      namespace: goldfinger
    spec:
      podSelector: {}
      policyTypes:
      - Ingress
    ```

2. 设置服务内的所有节点(pod)可以互相访问
    - 即允许同一个  namespace 下的所有 Pod 可以彼此通信

    ```yaml
    apiVersion: networking.k8s.io/v1
    kind: NetworkPolicy
    metadata:
      name: allow-same-namespace
      namespace: goldfinger
    spec:
      podSelector: {}
      policyTypes:
      - Ingress
      ingress:
      - from:
        - podSelector: {}
    ```

3. 允许服务内的部分节点（pod）接收外部网络的流量
    - 需要将那些接收外部流量的  Pod 打上标签: networking/allow-internet-access="true"

    ```yaml
    apiVersion: networking.k8s.io/v1
    kind: NetworkPolicy
    metadata:
      name: allow-internet-access
      namespace: goldfinger
    spec:
      podSelector:
        matchLabels:
          networking/allow-internet-access: "true"
      ingress:
      - {}
      policyTypes:
      - Ingress
    ```

## References

- [Why You Should Test Your Kubernetes Network Policies](https://www.inovex.de/blog/test-kubernetes-network-policies/)
- [Kubernetes Network Policy Recipes](https://github.com/ahmetb/kubernetes-network-policy-recipes)
- [Kubernetes Network Policies - A Detailed Security Guide](https://www.stackrox.com/post/2019/04/setting-up-kubernetes-network-policies-a-detailed-guide/)
- [An Introduction to Kubernetes Network Policies for Security People](https://medium.com/@reuvenharrison/an-introduction-to-kubernetes-network-policies-for-security-people-ba92dd4c809d)
