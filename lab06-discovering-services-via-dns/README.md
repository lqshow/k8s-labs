# 服务发现 - DNS

### 集群内部 DNS 域名格式有两种

1. <自定义的访问方式名称>.<工作负载所在命名空间> （例如：redis-svc.default）
2. <自定义的访问方式名称>.<工作负载所在命名空间>.svc.cluster.local（例如：redis-svc.default.svc.cluster.local）

### 从现有的 Golang 源代码构建 Docker 镜像，并将其推送到Docker Hub

```bash
cd `PWD`/src
docker build -t lqshow/discovering-services-via-dns:0.0.1 .
docker push lqshow/discovering-services-via-dns:0.0.1
```

### 将工作负载部署到 Kubernetes
```bash
kubectl create -f `PWD`/kubernetes/
```

### 验证
#### 通过域名1验证
```bash
➜  kubectl exec -it net-test-dff6845bb-wb889 curl discovering-services-via-dns-svc.default:3000
Welcome to the HomePage!

Hostname is discovering-services-via-dns-deploy-55659c6749-hkmlg
```
#### 通过域名2验证
```bash
➜  kubectl exec -it net-test-dff6845bb-wb889 curl discovering-services-via-dns-svc.default.svc.cluster.local:3000
Welcome to the HomePage!

Hostname is discovering-services-via-dns-deploy-55659c6749-hkmlg
```

### 删除 Service 和 Deplodment
```bash
kubectl delete -f `PWD`/kubernetes/
```