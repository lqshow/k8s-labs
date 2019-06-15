# Overview
Deployment，是一个定义多副本应用（即多个副本 Pod）的对象。
同时 Deployment 还负责在 Pod 定义发生变化时，对每个副本进行滚动更新（Rolling Update）。

像这样使用一种 API 对象（Deployment）管理另一种 API 对象（Pod）的方法，在 Kubernetes 中，叫作“控制器”模式（controller pattern）。Deployment 扮演的正是 Pod 的控制器的角色。通过 Label 来识别被管理的对象。


# Deploy the app to Kubernetes

```bash
kubectl create -f webapp-deploy.yml --record
kubectl create -f webapp-svc.yml
```

# Check that the Deploy, Replica Set, Pod and Service are created

```bash
kubectl get deploy -o wide
kubectl get rs
kubectl get po -o wide -l app=webapp --show-labels
kubectl get svc -o wide -l app=webapp
```
# Test communication
```bash
➜  kubectl run busybox-curl --rm -ti --image=lqshow/busybox-curl curl $(kubectl get pod nginx-po -o go-template='{{.status.podIP}}')
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Welcome to nginx!</h1>
<p>If you see this page, the nginx web server is successfully installed and
working. Further configuration is required.</p>

<p>For online documentation and support please refer to
<a href="http://nginx.org/">nginx.org</a>.<br/>
Commercial support is available at
<a href="http://nginx.com/">nginx.com</a>.</p>

<p><em>Thank you for using nginx.</em></p>
</body>
</html>
```

# Reference
- [k8s 使用 Deployment 部署应用](https://github.com/lqshow/notes/issues/39)

