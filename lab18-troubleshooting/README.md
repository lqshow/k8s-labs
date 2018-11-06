## 查看 Pod 运行状态
1. 运行 kubectl get pods，查看 Pod 是否正常运行，如若正常跳过以下两步。
2. 运行 kubectl describe pod `<POD-NAME>`，主要观察相关的 Events，查看哪个步骤除了问题，是调度除了问题还是资源不足等等。
3. 运行 kubectl kubectl logs `<POD-NAME>`，查看 Pod 内容器的日志，如果 app 做的好的话，基本从日志层面就能看出来哪里出了问题，是依赖的服务没起，或者缺了某个环境变量啥的。

### Pod 异常状态

| status                          | desc               |
| ------------------------------- | ------------------ |
| Pending                         | CPU/memory 等资源不足     |
| ImagePullBackOff / ErrImagePull | 拉不到镜像或者错误的容器镜像     |
| CrashLoopBackOff                | 应用启动之后又挂掉 |

## Pod 调试

Pod YAML 定义如下，容器的端口为 80
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-po
  labels:
    app: nginx
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 80
```
### 方法一: 通过 port-forward 为 Pod 创建本地端口映射测试
待 Pod 运行成功后，我们可以通过 Kubernetes 提供的端口转发功能来验证 Pod 内的容器是否正常运行。

```bash
# 使用 kubectl 端口转发
➜  kubectl port-forward nginx-po 3000:80
Forwarding from 127.0.0.1:3000 -> 80
Forwarding from [::1]:3000 -> 80
```

```bash
# kubectl port-forward 将 localhost:3000的请求转发到 nginx-pod Pod 的80端口。验证如下
➜  curl localhost:3000
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
### 方法二：通过 Pod 与 Pod 间的通信是否成功来验证
以下是创建一个 buysbox pod 来验证容期间的通信
```bash
# 创建 busybox pod
kubectl apply -f busybox-pod.yml
```
```bash
➜  kubectl exec -it busybox-pod -- curl $(kubectl get pod nginx-po -o go-template='{{.status.podIP}}')
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
### 方法三：进入指定容器内执行命令调试
如果 Pod 运行正常，但是通信异常，一直没有返回响应，最有可能的是容器中的 app 进程没起。
最有效的方式就是进入到指定的容器，调试代码后再做通信测试。
```bash
# 进入 Pod 中的容器调试，默认情况进入第一个容器
kubectl exec -it <POD-NAME> -- /bin/bash

# 指定进入具体的容器
kubectl exec -it <POD-NAME> -- /bin/bash -c <CONTAINER-NAME>
```

## Serive 调试

## References
- [Troubleshoot Kubernetes Deployments](https://docs.bitnami.com/kubernetes/how-to/troubleshoot-kubernetes-deployments/)
- [Troubleshoot Applications](https://kubernetes.io/docs/tasks/debug-application-cluster/debug-application/)