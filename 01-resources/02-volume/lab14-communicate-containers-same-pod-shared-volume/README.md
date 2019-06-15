# Pod 内的两个容器使用一个卷（Volume）进行通信

## emptyDir 类型
等同于 Docker 的隐式 Volume 参数，即：不显式声明宿主机目录的 Volume。
所以，Kubernetes 也会在宿主机上创建一个临时目录，这个目录将来就会被绑定挂载到容器所声明的 Volume 目录上。

Pod 中的容器，使用的是 volumeMounts 字段来声明自己要挂载哪个 Volume，并通过 mountPath 字段来定义容器内的 Volume 目录。

## hostPath 类型
提供了显式的 Volume 定义

## References
- [Multi-Container Pods in Kubernetes](https://linchpiner.github.io/k8s-multi-container-pods.html)
- [同 Pod 内的容器使用共享卷通信](https://kubernetes.io/cn/docs/tasks/access-application-cluster/communicate-containers-same-pod-shared-volume/)
- [Kubernetes: How to Share Disk Storage Between Containers in a Pod](https://www.stratoscale.com/blog/kubernetes/kubernetes-how-to-share-disk-storage-between-containers-in-a-pod/)
- [Kubernetes volumes by example](http://kubernetesbyexample.com/volumes/)