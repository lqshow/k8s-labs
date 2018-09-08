# Ingress resources

### 从现有的 Golang 源代码构建 Docker 镜像，并将其推送到Docker Hub
> 该镜像的作用： 展示当前启动的服务端口号
```bash
cd `PWD`/src
docker build -t lqshow/web-server:0.0.1 .
docker push lqshow/web-server:0.0.1
```

### 部署服务
#### 后端服务
> 启动两个服务，一个 port 为3000， 另外一个 port 为3001
```bash
kubectl create -f backend-svc-1-deploy.yml
kubectl create -f backend-svc-2-deploy.yml
```

#### 默认后端服务
> 对于未知请求全部负载到这个默认后端上，这个后端啥也不干，就是返回 404
```bash
kubectl create -f default-backend.yml
```

获取创建后的 Pod 信息
```bash
➜  kubectl get po -o wide
NAME                                    READY     STATUS    RESTARTS   AGE       IP           NODE
backend-svc-1-deploy-b7b568545-nbb8r    1/1       Running   0          1h        10.1.1.131   docker-for-desktop
backend-svc-2-deploy-58b874cbf5-t7lcc   1/1       Running   0          1h        10.1.1.132   docker-for-desktop
default-http-backend-5c6d95c48-8vc97    1/1       Running   0          1h        10.1.1.134   docker-for-desktop
default-http-backend-5c6d95c48-v6ftl    1/1       Running   0          1h        10.1.1.135   docker-for-desktop
net-test-dff6845bb-jc42c                1/1       Running   0          1h        10.1.1.133   docker-for-desktop
```

获取创建后的 Service 信息
```bash
➜  kubectl get svc -o wide
NAME                   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE       SELECTOR
backend-svc-1          ClusterIP   10.98.82.238     <none>        80/TCP    1h        app=backend-svc-1
backend-svc-2          ClusterIP   10.103.160.153   <none>        80/TCP    1h        app=backend-svc-2
default-http-backend   ClusterIP   10.110.199.163   <none>        80/TCP    1h        app=default-http-backend
```

通过 DNS 名称验证服务访问是否正常
```bash
➜  kubectl exec -it net-test-dff6845bb-jc42c curl backend-svc-1
This service is listening on port 3000
```
```bash
➜  kubectl exec -it net-test-dff6845bb-jc42c curl backend-svc-2
This service is listening on port 3001
```
```bash
➜  kubectl exec -it net-test-dff6845bb-jc42c curl default-http-backend
default backend - 404
```

### 部署 Ingress Controller
获取创建后的 Pod 信息
```bash
➜  kubectl get po -l app=nginx-ingress-lb -o wide
NAME                                       READY     STATUS    RESTARTS   AGE       IP           NODE
nginx-ingress-controller-8c64f6f87-kxl6w   1/1       Running   0          16m       10.1.1.136   docker-for-desktop
```

直接访问路由到默认的后端服务
```bash
➜  kubectl exec -it net-test-dff6845bb-jc42c curl http://10.1.1.136
default backend - 404
```

### 部署 Ingress
```bash
➜  kubectl create -f nginx-ingress.yml
ingress.extensions "nginx-ingress" created
```
查看创建的 Ingress
```bash
➜  kubectl get ing -o wide
NAME             HOSTS                                         ADDRESS   PORTS     AGE
drinks-ingress   *                                                       80        2d
nginx-ingress    svc1.k8s.local,svc2.k8s.local,www.k8s.local             80        29s
```
验证结果如下
```bash
➜  kubectl exec -it net-test-dff6845bb-jc42c -- curl -H Host:www.k8s.local http://10.1.1.136/healthz
ok

➜  kubectl exec -it net-test-dff6845bb-jc42c -- curl -H Host:svc1.k8s.local http://10.1.1.136
This service is listening on port 3000

➜  kubectl exec -it net-test-dff6845bb-jc42c -- curl -H Host:svc2.k8s.local http://10.1.1.136
This service is listening on port 3001

➜  kubectl exec -it net-test-dff6845bb-jc42c -- curl -H Host:www.k8s.local http://10.1.1.136
default backend - 404
```


### 参考
- [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
- [Kubernetes Nginx Ingress 教程](https://mritd.me/2017/03/04/how-to-use-nginx-ingress/)
- [Kubernetes Networking 101 – Ingress resources](http://www.dasblinkenlichten.com/kubernetes-networking-101-ingress-resources/)
- [nginx-ingress-controller](https://console.cloud.google.com/gcr/images/google-containers/GLOBAL/nginx-ingress-controller)
- [defaultbackend](https://console.cloud.google.com/gcr/images/google-containers/GLOBAL/defaultbackend?gcrImageListsize=50)