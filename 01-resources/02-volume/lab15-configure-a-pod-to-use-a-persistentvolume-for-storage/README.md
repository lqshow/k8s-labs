PV 和 PVC 是 Kubernetes 中的两种资源

- PV 和 PVC 的绑定完成了实际块存储设备和存储需求的绑定。
- PV 和 PVC 是一对一的关系
- 创建没有顺序要求

## PV(Persistent Volume)
代表一块实际的后台块存储设备

## PVC（Persistent Volume Claim）
代表的是 Pod 用户对块存储的实际需求

PVC 是一种特殊的 Volume, 一个 PVC 具体是什么类型的 Volume，要在跟某个 PV 绑定之后才知道。