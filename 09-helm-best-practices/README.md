# Overview

Helm 是  Kubernetes 的一个包管理工具，就类似 Ubuntu上 的 APT，和 CentOs 上的 yum 命令。
具有如下功能：

- 创建新的 chart
- chart 打包成 tgz 格式
- 上传 chart 到 chart 仓库或从仓库中下载 chart
- 在 `Kubernetes`集群中安装或卸载 chart
- 管理用 `Helm`安装的 chart 的发布周期

## Basic concept

### Chart

一个 Helm 包，包含了运行 Kubernetes 一个应用实例所需要的镜像、依赖和资源定义等，是描述一组相关 Kubernetes 资源的文件集合。

#### Chart 的基本结构

```bash
.
├── Chart.yaml				# 用于描述这个 chart， 包括名字，描述信息以及版本号
├── README.md
├── templates				# Kubernetes 模板
│   ├── NOTES.txt			# 用于介绍 chart 部署后的信息。例如介绍如何使用该 chart，缺省设置等
│   ├── _helpers.tpl
│   ├── deployment.yaml
│   ├── pvc.yaml
│   ├── secrets.yaml
│   └── svc.yaml
└── values.yaml				# 用于存储 templates 目录中模板文件中用到的变量
```

### Config

包含了应用发布配置信息

#### version and appVersion

> helm 安装是指定 version 安装的，并不关注 appVersion，起作用的是 chart version

| version    | Desc                                         |
| ---------- | -------------------------------------------- |
| version    | Chart 的版本号                               |
| appVersion | 指应用程序的版本，通常与镜像的版本号保持一致 |

以上两个版本号需遵循以下原则

1. appVersion 变更的同时需要变更 version
2. version 变更不需要改变 appVersion
3. 两个版本号约定 major 和 minor，需要保持一致，patch 可以任意调整。这样看起来比较直观些(helm 原则上没有做强制约定)

### Release

在 Kubernetes 集群上运行的 Chart 及其配置的一个实例。

在同一个集群上，一个 Chart 可以安装很多次。每次安装都会创建一个新的 release。

### Repository

用于发布和存储 Chart 的仓库。

## Components

> client 管理 charts， server 管理发布 release

### Helm Client

用户命令行工具，用来创建，拉取，搜索和验证 Charts，初始化 Tiller 服务。

1. 本地 chart 开发
2. 仓库管理
3. 与 Tiller sever 交互
4. 发送预安装的 chart
5. 查询 release 信息
6. 要求升级或卸载已存在的 release

### Tiller Server

部署在  Kubernetes 集群内部的 server，与 Helm client、Kubernetes API server 进行交互。

1. 监听来自 Helm client 的请求

2. 通过 chart 及其配置构建一次发布

3. 安装 chart 到 `Kubernetes`集群，并跟踪随后的发布

4. 通过与 `Kubernetes`交互升级或卸载 chart

## Installing Helm

### Installing Helm Client

```bash
# FROM SCRIPT
curl https://raw.githubusercontent.com/helm/helm/master/scripts/get > get_helm.sh
chmod 700 get_helm.sh
./get_helm.sh
```

### Installing Tiller Server

```bash
# 安装 Tiller
helm init

# 覆盖 Tiller image
helm init --upgrade --tiller-image gcr.io/kubernetes-helm/tiller:v2.10.0

# 升级 Tiller
helm init --upgrade

# 确认 Tiller 是否在集群安装好
kubectl get pods --namespace kube-system -l app=helm
```

### ClusterRoleBinding

```yaml
# 生成以下资源
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tiller
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: tiller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: tiller
    namespace: kube-system
```

```bash
# 调整 service account
kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'
```

## Helm Chart Repository

Helm Chart Repository 可以简单理解为一个用来托管 index.yaml 文件和 chart package 的 web 服务器。

### Installing chartmuseum

#### Via Helm Chart

定义指定的 yaml 文件

```yaml
# vaules-production.yaml

fullnameOverride: "chartmuseum"
env:
  open:
    DISABLE_API: false
persistence:
  enabled: true
  accessMode: ReadWriteMany
  storageClass: basebit-nfs
ingress:
  enabled: true
  hosts:
  - name: chartmuseum.domain.com
    path: /
```

安装命令

```bash
# 安装
helm upgrade -i chartmuseum -f vaules-production.yaml stable/chartmuseum

# 验证
1. 访问 http://chartmuseum.domain.com/
2. 访问 http://chartmuseum.domain.com/health
```

#### Add *ChartMuseum* installation to the local repository list

```bash
helm repo add chartmuseum http://chartmuseum.domain.com
```

#### Installing helm push plugin

```bash
helm plugin install https://github.com/chartmuseum/helm-push
```

### Usage

```bash
# 方式一：推送目录(会自动打包上传)
helm push mychart/ chartmuseum

# 方式二：推送已经打好的包
helm push mychart-1.0.0.tgz chartmuseum

# 方式三：推送到指定的远程 URL
helm push mychart-1.0.0.tgz http://chartmuseum.domain.com

# 方式四：强制推送
helm push --force mychart/ chartmuseum

# 方法五：通过 curl 来上传（如果没装 push plugin 情况下）
curl --data-binary xxx.tar http://chartmuseum.domain.com/api/charts;

# 更新 repo
# 需要明确的是本地推送到远端 repo 后，同时需要在本地做更新，这样才能和远端保持一致
helm repo update

# 查看 chartmuseum repo 中的所有 chart, 确认 charts 是否上传成功
helm search chartmuseum/

# 安装
helm install chartmuseum/mychart

# 安装指定 version chart
helm install --version 2.3.0 chartmuseum/mychart
```

## Notes

1. helm release name 是全局唯一的，不区分 namespace，集群内 release name 需根据使用场景加上不同前缀来区分。
2. ingress.enabled 值默认需设置成 false , 这样可以避免用户尝试安装，导致 ingress host 冲突。
3. 严禁使用 helm delete -purge 删除 chartmuseum，该操作会同时将  pvc 全部删掉，所有 chart 都会白传，当然后续你重新传一遍所有 charts 问题也不大。

## References

- [探索托管 Helm Charts 的正确方式](https://mp.weixin.qq.com/s/DFqS4sYtTjLoUSzLu2khtA)
- [Helm 安装使用](https://mp.weixin.qq.com/s/Ub6xsTw2HNXCs7WcEeF03w)
- [HELM Best practices](https://codefresh.io/docs/docs/new-helm/helm-best-practices)
