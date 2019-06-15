# Overview

使用 kubectl 与集群交互

## Basic Commands

| command | desc                                              | example                                                      |
| ------- | ------------------------------------------------- | ------------------------------------------------------------ |
| create  | 从文件或stdin创建资源                             | kubectl create -f using-projected-volume.yml                 |
| delete  | 删除指定资源，支持文件名、资源名、label selector  | kubectl delete po -l foo=bar<br />kubectl delete -f using-projected-volume.yml |
| edit    | 使用系统编辑器编辑资源                            | kubectl edit deploy/foo                                      |
| apply   | 从文件或stdin创建/更新资源，建议直接用 apply 创建 | kubectl apply -f using-projected-volume.yml                  |
| get     | 最基本的查询命令                                  | kubectl get rs<br />kubectl get deploy<br />kubectl get svc<br />kubectl get rs/foo |
| explain | 查看资源定义                                      | kubectl explain po                                           |

## Troubleshooting and Debugging Commands

| command      | desc                      | example                                                      |
| ------------ | ------------------------- | ------------------------------------------------------------ |
| describe     | 查看资源详情              | kubectl describe pod {POD-NAME}                              |
| logs         | 查看pod内容器的日志       | kubectl logs -f {POD-NAME}<br />kubectl logs -f -p {POD-NAME}<br />kubectl logs -f {POD-NAME} -c {CONTAINER-NAME}<br />kubectl logs -f --tail 500 {POD-NAME} -c {CONTAINER-NAME} |
| exec         | 在指定容器内执行命令      | kubectl exec -it {POD-NAME} /bin/bash -c {CONTAINER-NAME}    |
| port-forward | 为pod创建本地端口映射     | kubectl port-forward nginx-po 3000:80<br />将 localhost:3000的请求转发到 nginx-pod Pod 的80端口 |
| port-forward | 为service创建本地端口映射 | kubectl port-forward svc/svc-name 3000                       |

## Detailed Commands

### Cluster

```bash
# 查看客户端和Server端K8S版本
kubectl version

# 查看集群整体状态
kubectl cluster-info

# 检查集群常规健康状况
kubectl get componentstatus

# 查看集群中节点
kubectl get node -o wide

# 查看节点信息
kubectl describe node <NODE-NAME>

# 查看资源占用率
kubectl top pod <POD-NAME>

# 查看支持的API版本
kubectl api-versions
```

### Pod

```bash
# 获取 Pod 运行状态
kubectl get pods
kubectl get pod
kubectl get po

# 通过 yaml 配置文件创建 Pod
kubectl create -f pod.yaml
kubectl apply -f pod.yaml

# 列出所有 Pod，并包含节点信息
kubectl get po -o wide

# 列出所有 Pod，且包含 label
kubectl get po --show-labels -o wide

# 通过标签过滤 Pod
kubectl get pods -l <LABEL-KEY>=<LABEL-VALUE>

# 查看 Pod 中所有容器的运行状态
kubectl get pod <POD-NAME> -o yaml 

# 查看 Pod 的详细信息(包括容器的信息以及 Pod 相关的事件)
kubectl describe pod <POD-NAME>

# 查看 Pod 日志
kubectl logs <POD-NAME>

# 查看 Pod 中某个容器日志
kubectl logs <POD-NAME> -c <CONTAINER-NAME>

# 进入 Pod 中的容器调试，默认情况进入第一个容器
kubectl exec -it <POD-NAME> /bin/bash

# 指定进入具体的容器
kubectl exec -it <POD-NAME> /bin/bash -c <CONTAINER-NAME>

# 通过指定 Pod 名称来删除 Pod
# 如果 Pod 是通过 Deployment 创建的，直接删除 Pod，则 Deployment 将会重新创建该 Pod。
# 不能直接删除 Pod，需使用 kubectl 删除拥有该 Pod 的 Deployment。
# 由 RC 或者 RS 创建的 Pod 也是同样的情况
kubectl delete pods <POD-NAME>

# 通过 yaml 配置文件删除 Pod
kubectl delete -f pod.yaml
```

### Deployment

```bash
# 获取所有 Deployment
kubectl get deployments
kubectl get deployment
kubectl get deploy

# 通过 yaml 配置文件创建一个 Deployment
# record标识会记录当前的命令，有利于查看部署历史
kubectl create -f ./deployment.yaml --record

# 更改 Deployment image 版本配置，更新 Deployment
kubectl apply -f ./deployment-update.yaml

# 扩容(通过更改 deployment.yaml 中的 replicas 值来扩容)
# replicas: 4
kubectl apply -f deployment.yaml --record

# 通过kubectl scale指令来扩容/缩容
kubectl scale deploy/<DEPLOY-NAME> --replicas=4

# template部分调整
# 更新镜像版本
kubectl set image deployment/<DEPLOY-NAME> nginx=nginx:1.9.1

# 也可手动调整配置文件再执行 edit
kubectl edit deployment/<DEPLOY-NAME>

# 强制更新(先删除后替换)
kubectl replace -f deployment.yaml --force

# 通过指定 Deployment 名称来删除 Deployment(同时删除 Pod)
kubectl delete deployment <DEPLOY-NAME>

# 通过配置文件删除
kubectl delete -f ./deployment.yaml

# 列出 Deployment 的部署历史
kubectl rollout history deployment/<DEPLOY-NAME>

# 查看具体某个部署版本
kubectl rollout history deployment/nginx-deployment --revision=2

# 回滚到上一个版本
kubectl rollout undo deployment/<DEPLOY-NAME>

# 回滚到指定版本
kubectl rollout undo deployment/<DEPLOY-NAME> --to-revision=2

# 查看升级状态
kubectl rollout status deployment/<DEPLOY-NAME>

# 暂停部署
kubectl rollout pause deployment/<DEPLOY-NAME>

# 恢复部署
kubectl rollout resume deployment/<DEPLOY-NAME>
```

### Services

```bash
# 获取所有 Service
kubectl get services
kubectl get service
kubectl get svc

# 通过标签过滤删除服务
kubectl delete service -l <LABEL-KEY>=<LABEL-VALUE>

# 通过指定名称删除服务
kubectl delete service <SERVER-NAME>

# 查看服务的 cluster-ip
kubectl get svc <SERVER-NAME> -o go-template='{{.spec.clusterIP}}'

# 查看服务的端口
kubectl get svc <SERVER-NAME> -o go-template='{{(index .spec.ports 0).port}}'

# 获取 service 的 endpoints(endpoints分别对应服务的 pod)
kubectl get ep
```

### Service Account

```bash
# 获取所有服务账号
kubectl get serviceAccounts

# 查看某一个服务账号
kubectl get serviceaccounts/lqshow -o yaml
 
# 创建一个新的serviceaccount
kubectl create serviceaccount lqshow

# 获取 secret 信息
kubectl get secrets | grep ^lqshow | cut -f1 -d ' '

# 获取 token
kubectl describe secret $(kubectl get secrets | grep ^lqshow | cut -f1 -d ' ') | grep -E '^token' | cut -f2 -d':' | tr -d " "

kubectl get secret $(kubectl get secrets | grep ^lqshow | cut -f1 -d ' ') -o json | jq -r '.data["token"]' | base64 -D

# 获取 ca 证书
kubectl get secret $(kubectl get secrets | grep ^lqshow | cut -f1 -d ' ') -o json |  jq -r '.data["ca.crt"]' | base64 -D > ca.crt

# 通过指定名称删除服务账号
kubectl delete serviceaccount/lqshow
```