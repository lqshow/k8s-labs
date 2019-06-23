# Overview

Taints(污点) 通常与 Tolerations(容忍) 配合使用

## Taints

Taints(污点) 是 Node 的一个属性，设置了 Taints(污点) 后，因为有了污点，所以 Kubernetes 是不会将 Pod 调度到这个 Node上 的，
Taints 可以让 Node 拒绝运行 Pod，甚至驱逐 Pod。

每个污点有一个 key 和 value 作为污点的标签，其中value可以为空，effect 描述污点的作用。

### effect

Taint effect 支持如下三个选项

| effect           | desc                                  |
| ---------------- | ------------------------------------- |
| NoSchedule       | 不允许调度，已调度的不影响                |
| PreferNoSchedule | 尽量不要调度                          |
| NoExecute        | 不仅不允许调度，还会驱逐 Node 上已有的未设置对应 Tolerate 的 Pod |

```bash
kubectl taint nodes [NODE_NAME] [KEY]=[VALUE]:[EFFECT]

# 设置污点
kubectl taint node xdp-node-37 special=gpu:NoSchedule

# 去除污点
kubectl taint node xdp-node-37 special=gpu:NoSchedule-
kubectl taint nodes xdp-node-37 special:NoSchedule-
kubectl taint node xdp-node-37 special-
```

## Tolerations

Tolerations(容忍) 是 Pod 的一个属性 ，只要 Pod 能够容忍 Node 上的污点，那么 Kubernetes 就会忽略 Node 上的污点，就能够把 Pod 调度过去。

### operator

> 如果不指定 operator，则默认为 Equal

| operator           | desc                                  |
| ---------------- | ------------------------------------- |
| Equal       | key 与 value 之间的关系是 equal                        |
| Exists | value 属性可省略                          |

> tolerations 属性下各值必须使用引号

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-taints
spec:
  tolerations:
  - key: "special"
    operator: "Equal"
    value: "gpu"
    effect: "NoSchedule"
    tolerationSeconds: 180 

  containers:
    - name: pod-tains
      image: nginx
```

### 关于 master

使用 kubeadm 搭建的集群默认就给 master 节点添加了一个污点标记， 所以我们看到我们平时的 pod 都没有被调度到 master 上去

```bash
Taints:             node-role.kubernetes.io/master:NoSchedule
                    node-role.kubernetes.io/master:PreferNoSchedule
```
让 master 节点调度其他 Pod 的方法

1. 去掉 master 上的 污点(一般不建议这么做)

    ```bash
    # 去除污点
    kubectl taint node k8s-master node-role.kubernetes.io/master-

    # 恢复 master 污点
    kubectl taint node k8s-master node-role.kubernetes.io/master=""
    ```

2. 如果需要某些特殊 Pod 调度到 master 节点上，需要加上相应 tolerations 属性

    TODO(test)

    ```yaml
    tolerations:
    - key: "node-role.kubernetes.io/master"
      operator: "Exists"
      effect: "NoSchedule"
    - key: "node-role.kubernetes.io/master"
      operator: "Exists"
      effect: "PreferNoSchedule"
    ```

## Notes

1. tolerations 中的 key、value、effect 与 Node的 Taint 设置需保持一致
2. 空的 key 如果再配合 Exists 就能匹配所有的 key 与 value ，也就是能容忍所有 node 的所有 Taints。
3. 空的 effect 匹配所有的 effect。
4. 一个 node 上可以有多个污点, 相对应的 Pod 必须同时匹配多个污点

## References

- [Controlling scheduling with node taints](https://cloud.google.com/kubernetes-engine/docs/how-to/node-taints)
- [Taints and Tolerations](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/)
- [Kubernetes之Taints与Tolerations 污点和容忍](https://blog.51cto.com/newfly/2067531)
