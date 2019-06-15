## Dashboard token 失效

Dashboard 默认的 token 失效时间为 900 秒，即 15 分钟

```bash
# 通过修改参数 token-ttl 即可调整失效时间
kubectl edit deploy kubernetes-dashboard -n kube-system
```

```yaml
spec:
  containers:
  - args:
    - --auto-generate-certificates
    - --token-ttl=86400
```

### Referneces

- [DefaultTokenTTL](https://github.com/kubernetes/dashboard/blob/master/src/app/backend/auth/api/types.go#L29)
- [Dashboard arguments](https://github.com/kubernetes/dashboard/wiki/Dashboard-arguments)

## 强制删除 PVC

有时执行了强删命令，pvc 仍然一直删不掉，一直在终结状态，比如以下删除命令
```bash
kubectl delete pv --grace-period=0 --force pvc-c162c38c-77d2-11e9-8df0-525400a03483
```
尝试的解决方案，通过  edit 命令打开了对应 pvc， 删除了 finalizers 对应的数据.
```bash
kubectl edit pvc pvc-c162c38c-77d2-11e9-8df0-525400a03483
```
```yaml
finalizers:
  -  kubernetes.io/pv-protection
```

### References

- [Kubernetes : Deleting resource like pv with — force and — grace-period=0 still keeps pvs in terminating state](https://medium.com/@miyurz/kubernetes-deleting-resource-like-pv-with-force-and-grace-period-0-still-keeps-pvs-in-3f4ad8710e51)

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
## Servie 异常
1. 访问 Service ClusterIP 失败时，可以首先确认是否有对应的 Endpoints
```bash
kubectl get endpoints <service-name>
```
2. 如果该列表为空，则有可能是该 Service 的 LabelSelector 配置错误，可以用下面的方法确认一下
```bash
# 查询 Service 的 LabelSelector
kubectl get svc <service-name> -o jsonpath='{.spec.selector}'

# 查询匹配 LabelSelector 的 Pod
kubectl get pods -l key1=value1,key2=value2
```
3. 如果 Endpoints 正常，可以进一步检查
    - Pod 的 containerPort 与 Service 的 containerPort 是否对应
    - 直接访问 podIP:containerPort 是否正常 再进一步，即使上述配置都正确无误，还有其他的原因会导致 Service 无法访问，比如
    - Pod 内的容器有可能未正常运行或者没有监听在指定的 containerPort 上
    - CNI 网络或主机路由异常也会导致类似的问题
    - kube-proxy 服务有可能未启动或者未正确配置相应的 iptables 规则，比如正常情况下名为 hostnames的服务会配置以下 iptables 规则

## References
- [Troubleshoot Kubernetes Deployments](https://docs.bitnami.com/kubernetes/how-to/troubleshoot-kubernetes-deployments/)
- [Troubleshoot Applications](https://kubernetes.io/docs/tasks/debug-application-cluster/debug-application/)
- [Kubernetes 网络排错指南](https://mp.weixin.qq.com/s/N_0cT1k_L7lRRHaULxfeZA)