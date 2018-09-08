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

