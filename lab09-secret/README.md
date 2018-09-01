Secret 对象类型用来保存敏感信息，比如密码，OAuth 令牌和 ssh key。

Pod 中有3种方式使用 secret
1. 环境变量
2. 作为 volume 中的文件挂载到 Pod 中
3. Pod 中拉取私有镜像时使用

### Secret 有三种类型

| type                                | desc                                    |
| ----------------------------------- | --------------------------------------- |
| Opaque                              | 用来存储密钥等，value 使用 base64 编码格式    |
| kubernetes.io/dockerconfigjson      | 用来存储私有 docker registry 的认证信息 |
| kubernetes.io/service-account-token | 创建 sa 时会默认创建对应的 secret       |

### 将 Secret 作为环境变量引入

##### 创建 Pod 以及 secret
```bash
kubectl create -f using-secret-as-env.yml
```

##### 查看创建的 secret
```bash
➜  kubectl get secret consume-secret-env-secret -o yaml
apiVersion: v1
data:
  password: MTIzNDU2Cg==
  username: TFEK
kind: Secret
metadata:
  creationTimestamp: 2018-09-01T07:43:38Z
  name: consume-secret-env-secret
  namespace: default
  resourceVersion: "515073"
  selfLink: /api/v1/namespaces/default/secrets/consume-secret-env-secret
  uid: bd8523a6-adba-11e8-8166-025000000001
type: Opaque
```

##### 从环境变量中消费 secret
```bash
➜  kubectl exec $(kubectl get pods consume-secret-env-pod  -o=name|cut -d "/" -f2) env |grep SECRET
SECRET_USERNAME=LQ
SECRET_PASSWORD=123456
```

### 将 Secret 挂载到 Volume

##### 创建 Pod 以及 secret
```bash
kubectl create -f using-secret-in-volume.yml
```

##### 查看创建的 secret
```bash
➜  kubectl get secret consume-secret-in-volume-secret -o yaml
apiVersion: v1
data:
  password: MTIzNDU2Cg==
  username: TFEK
kind: Secret
metadata:
  creationTimestamp: 2018-09-01T06:52:01Z
  name: consume-secret-in-volume-secret
  namespace: default
  resourceVersion: "511689"
  selfLink: /api/v1/namespaces/default/secrets/consume-secret-in-volume-secret
  uid: 870ed0a5-adb3-11e8-8166-025000000001
type: Opaque
```

##### 从 volume 消费 secret
```bash
➜  kubectl exec $(kubectl get pods consume-secret-in-volume-pod -o=name|cut -d "/" -f2) ls /etc/foobar
password
username
```

```bash
➜  kubectl exec $(kubectl get pods consume-secret-in-volume-pod -o=name|cut -d "/" -f2) cat /etc/foobar/username
LQ
➜  kubectl exec $(kubectl get pods consume-secret-in-volume-pod -o=name|cut -d "/" -f2) cat /etc/foobar/password
123456
```
### kubernetes.io/dockerconfigjson 类型

kubernetes.io/dockerconfigjson 类型的 Secret 是将包含 Docker Registry 凭证传递给 Kubelet 的一种方式，可以用来为 Pod 拉取私有镜像

##### 创建 imagePullSecret
```bash
# 以下变量需用实际环境的值来替换
kubectl create \
    secret docker-registry reg-secret \
    --docker-server=DOCKER_REGISTRY_SERVER \
    --docker-username=DOCKER_USER \
    --docker-password=DOCKER_PASSWORD \
    --docker-email=DOCKER_EMAIL
```
##### 查看创建的 secret
```bash
➜  kubectl get secrets reg-secret -o yaml
apiVersion: v1
data:
  .dockerconfigjson: eyJhdXRocyI6eyJET0NLRVJfUkVHSVNUUllfU0VSVkVSIjp7InVzZXJuYW1lIjoiRE9DS0VSX1VTRVIiLCJwYXNzd29yZCI6IkRPQ0tFUl9QQVNTV09SRCIsImVtYWlsIjoiRE9DS0VSX0VNQUlMIiwiYXV0aCI6IlJFOURTMFZTWDFWVFJWSTZSRTlEUzBWU1gxQkJVMU5YVDFKRSJ9fX0=
kind: Secret
metadata:
  creationTimestamp: 2018-09-01T08:19:01Z
  name: reg-secret
  namespace: default
  resourceVersion: "517396"
  selfLink: /api/v1/namespaces/default/secrets/reg-secret
  uid: ae93aeb1-adbf-11e8-8166-025000000001
type: kubernetes.io/dockerconfigjson
```
我们可以看下.dockerconfigjson实际的值
```bash
➜   echo -n "eyJhdXRocyI6eyJET0NLRVJfUkVHSVNUUllfU0VSVkVSIjp7InVzZXJuYW1lIjoiRE9DS0VSX1VTRVIiLCJwYXNzd29yZCI6IkRPQ0tFUl9QQVNTV09SRCIsImVtYWlsIjoiRE9DS0VSX0VNQUlMIiwiYXV0aCI6IlJFOURTMFZTWDFWVFJWSTZSRTlEUzBWU1gxQkJVMU5YVDFKRSJ9fX0=" |base64 -D |jq
{
  "auths": {
    "DOCKER_REGISTRY_SERVER": {
      "username": "DOCKER_USER",
      "password": "DOCKER_PASSWORD",
      "email": "DOCKER_EMAIL",
      "auth": "RE9DS0VSX1VTRVI6RE9DS0VSX1BBU1NXT1JE"
    }
  }
}
```
##### kubernetes.io/dockerconfigjson 类型 Secret 使用
> 在 Pod 定义的 imagePullSecrets 中使用这个 Secret
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: foo
spec:
  containers:
    - name: foo
      image: janedoe/awesomeapp:v1
  imagePullSecrets:
    - name: reg-secret
```

### 总结
1. 挂载到 Volume 中的 secret 被更新时，被映射的 key 也将被更新。Kubelet 在周期性同步时会检查被挂载的 secret 是不是最新的


### 参考
- [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/#decoding-a-secret)
- [Images](https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod)