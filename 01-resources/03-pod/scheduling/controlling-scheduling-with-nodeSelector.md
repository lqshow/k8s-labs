# Overview

nodeSelector 是通过 node label selector 实现约束 pod 运行到指定节点(精确匹配)。

该方式比较简单粗暴，使用起来不能灵活调度，控制粒度较大。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-node-selector
spec:
  nodeSelector:
    kubernetes.io/hostname: xx-bb

  containers:
    - name: pod
      image: nginx
```

## Notes

如果我们的目标节点没有可用的资源，我们的 Pod 就会一直处于 Pending 状态.

```bash
Events:
  Type     Reason            Age                  From               Message
  ----     ------            ----                 ----               -------
  Warning  FailedScheduling  29s (x1670 over 1h)  default-scheduler  0/4 nodes are available: 4 node(s) didn't match node selector.
```
