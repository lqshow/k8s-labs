## Overview

现在用 kubernetes 的用户越来越多，每个公司一般都存在多个 kubernetes 集群，这种情况下一般会遇到以下几个问题:
1. 怎么才能让用户轻松地在多个集群之间自由切换呢?
2. 我只想给某个用户提供查看 Pod 日志的权限，其他权限一概不提供，怎么操作呢？

用 kubeconfig 文件即可解决这个问题，它主要是利用 Service Account + RBAC 创建不同的 kubeconfig 文件为不同角色的用户访问集群。

## Create a Role
创建一个 Role 配置，比如我要生成的 kubeconfig 文件，只提供用户查看 Pod 以及 Pod 日志的权限。
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: read-pod
  namespace: dev
rules:
- apiGroups:
  - '*'
  resources:
  - pods
  - pods/log
  verbs:
  - get
  - watch
  - list
```

## Creating a kubeconfig file for Kubernetes cluster
我将涉及的命令整合到了 Makefile 文件中，只要执行以下三个命令，即可快速生成一个 kubeconfig 文件。

当然，环境变量的 value 值需要自行修改调整。

```bash
# create service account
make create-serviceaccount -e SERVICE_ACCOUNT=lqshow \
                           -e NAMESPACE=dev

# set kube config
make setup-kubeconfig -e SERVICE_ACCOUNT=lqshow \
					  -e KUBE_APISERVER=https://localhost:6443 \
					  -e NAMESPACE=dev \
					  -e CLUSTER_NAME=test-staging \
					  -e CONTEXT_NAME=test-context

# create a RoleBinding
make create-rolebinding -e ROLEBINDING_NAME=read-pod-rolebinding \
						-e ROLE_NAME=read-pod \
						-e NAMESPACE=dev \
						-e SERVICE_ACCOUNT=lqshow
```

## The KUBECONFIG environment variable
环境变量 KUBECONFIG 保存一个 kubeconfig 文件列表, kubectl 默认的 kubeconfig 文件位$HOME/.kube/config。
将以上生成的 `k8s-lqshow-config` 文件追加到该环境变量上。

```bash
export KUBECONFIG=$HOME/.kube/config:/Users/linqiong/k8s-lqshow-config
```
## Tools
1. [kubectx](https://github.com/ahmetb/kubectx)：用于切换 kubernetes context
2. [kubens](https://github.com/ahmetb/kubectx)：用于切换 kubernetes namespace
3. [kube-ps1](https://github.com/jonmosco/kube-ps1)：为命令行终端增加 $PROMPT 字段
4. [kube-shell](https://github.com/cloudnativelabs/kube-shell)：交互式带命令提示的kubectl终端


## References
- [k8s 利用 Service Account + RBAC 访问资源](https://github.com/lqshow/notes/issues/45)
- [打造高效的Kubernetes命令行终端](https://zhuanlan.zhihu.com/p/34357028)
