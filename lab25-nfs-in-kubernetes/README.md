## Overview
所谓容器的 Volume，其实就是将一个宿主机上的目录，跟一个容器里的目录绑定挂载在了一起。

而所谓的“持久化 Volume”，指的就是这个宿主机上的目录，具备“持久性”。即：这个目录里面的内容，既不会因为容器的删除而被清理掉，也不会跟当前的宿主机绑定。
这样，当容器被重启或者在其他节点上重建出来之后，它仍然能够通过挂载这个 Volume，访问到这些内容。

大多数情况下，持久化 Volume 的实现，往往依赖于一个远程存储服务，比如：远程文件存储（比如，NFS、GlusterFS、Ceph）、远程块存储（比如，公有云提供的远程磁盘）等等。

下面是 NFS 的使用方法。

## Prerequisites
### CentOs

```bash
yum install nfs-utils
```

### Ubuntu

```bash
apt install nfs-common
```

## Install nfs server
```bash
kubectl apply -f ./deploy/psp.yaml

# RoleBinding 需要调整 namespace
kubectl apply -f ./deploy/rbac.yaml

# nodeName 需指定 node 节点名称
kubectl apply -f ./deploy/deployment.yaml
```

## Usage

### PV/PVC
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs-mysql-volume
spec:
  capacity:
    storage: 1Gi
  accessModes:
    - ReadWriteMany
  nfs:
  	# server 字段值为 service(nfs-provisioner) 的 ip 地址
    server: "10.102.186.53"
    path: "/mysql"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-pvc
spec:
  accessModes:
    - ReadWriteMany
  storageClassName: ""
  resources:
    requests:
      storage: 200Mi
```
### Static Volumes with the NFS
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-using-nfs
spec:
  containers:
  - name: app
    image: alpine
    volumeMounts:
      # name must match the volume name below
    - name: nfs-volume
      mountPath: "/var/nfs"
    command: ["/bin/sh"]
    args: ["-c", "while true; do date >> /var/nfs/dates.txt; sleep 5; done"]
  volumes:
  - name: nfs-volume
    nfs:
      # server 字段值为 service(nfs-provisioner) 的 ip 地址
      server: "10.102.186.53"
      path: "/abc"
      readOnly: false
```

### Dynamic Volumes with the NFS Client Provisioner
> 动态方式是基于 StorageClass 对象，而 StorageClass 对象的作用，则是充当 PV 的模板。并且，只有同属于一个 StorageClass 的 PV 和 PVC，才可以绑定在一起。
> 创建 pvc 的同时会自动创建相同 storage 大小的 pv 对象，pv 对象的名称由 pvc-uuid 组成。

#### 创建 StorageClass 对象
```bash
# 定义 StorageClass 名称，以及 provisioner 值，需同 nfs server 定义的 provisioner 值一致
kubectl apply -f ./deploy/storage-class.yaml
```

#### 创建一个 PVC 对象
```yaml
# 指定 storage-class ，accessModes 和 storage 三个值
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: dynamic-pvc
  annotations:
    volume.beta.kubernetes.io/storage-class: "basebit-nfs"
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
```
#### 在 Pod 中使用 pvc
```yaml
# 指定 claimName 的值，为以上创建的 PVC 对象名称
apiVersion: v1
kind: Pod
metadata:
  name: pod-using-dynamic-pvc
spec:
  containers:
  - name: app
    image: alpine
    volumeMounts:
      # name must match the volume name below
    - name: dynamic-pvc-volume
      mountPath: "/var/nfs"
    command: ["/bin/sh"]
    args: ["-c", "while true; do date >> /var/nfs/dates.txt; sleep 5; done"]
  volumes:
  - name: dynamic-pvc-volume
    persistentVolumeClaim:
      claimName: dynamic-pvc
```

## References
- [Mounting NFS file system in CentOS 7](https://blog.hostonnet.com/mounting-nfs-centos-7)
- [mount: wrong fs type, bad option, bad superblock](https://www.svennd.be/mount-wrong-fs-type-bad-option-bad-superblock/)
- [PV、PVC、StorageClass，这些到底在说啥？](https://time.geekbang.org/column/article/42698)
- [nfs provisioner](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs)
- [NFS in Kubernetes](https://opensource.ncsa.illinois.edu/confluence/display/~lambert8/NFS+in+Kubernetes)
