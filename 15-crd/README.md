

# Overview

1. CRD 是 Kubernetes 的一种资源类型，是 Custom Resource Definition 的缩写
2. CRD 是自定义资源的定义，用来描述自定义资源的具体的 spec
3. Kubernetes 允许用户自定义资源 CRD，向 Kubernetes 集群注册一种新资源，用于扩展 Kubernetes 集群能力
4. 当我们想要将自己的对象引入 Kubernetes 集群以完全满足我们的需求时，就需要使用 CRD。一旦我们在 Kubernetes 中创建了CRD，我们就可以像使用其他 Kubernetes 内置资源类型一样使用 kubectl 操作。


## Create a CRD

以下是一个 CRD 样例

关于自定义验证对象
1. 在 apiextensions.k8s.io/v1 中是必填字段，在 apiextensions.k8s.io/v1beta1 中可选的
2. 在 apiextensions.k8s.io/v1 中可指定版本做不同定义

```yaml
# crd.yaml

# Deprecated in v1.16 in favor of apiextensions.k8s.io/v1
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  # 用于定义 CRD 的名字, 名称必须与下面的 spec 字段匹配：<spec.names.plural>.<spec.group>
  name: dagnoderunners.enigma.basebit.me
spec:
  group: enigma.basebit.me
  names:
    # CamelCased 格式的单数类型
    kind: DagNodeRunner
    listKind: DagNodeRunnerList
    plural: dagnoderunners
    singular: dagnoderunner
    shortNames:
      - dnr
  version: v1
  versions:
  - name: v1
    served: true
    storage: true
  # CRD 可以是命名空间的，也可以是集群范围的，通过 scope 来指定
  scope: Namespaced
  # openAPIV3Schema is the schema for validating custom objects
  validation:
    openAPIV3Schema:
      description: DagNodeRunner is the Schema for the dagnoderunners API
      type: object
      properties:
        required: ["spec"]
        apiVersion:
          description: 'APIVersion defines the versioned schema of this representation
            of an object. Servers should convert recognized schemas to the latest
            internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources'
          type: string
        kind:
          description: 'Kind is a string value representing the REST resource this
            object represents. Servers may infer this from the endpoint the client
            submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds'
          type: string
        metadata:
          type: object
        spec:
          description: DagNodeRunnerSpec defines the desired state of DagNodeRunner
          type: object
          required: [foo]
          properties:
            foo: {type: string, minimum: 1}
            barz: {type: boolean}
        status:
          description: DagNodeRunnerStatus defines the observed state of DagNodeRunner
          type: object
```

创建成功后，可访问 API endpoint: /apis/enigma.basebit.me/v1

```bash
➜   curl -k  https://xx.xx.xx.xx:6443/apis/enigma.basebit.me/v1 |jq
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   473  100   473    0     0   2620      0 --:--:-- --:--:-- --:--:--  2627
{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "enigma.basebit.me/v1",
  "resources": [
    {
      "name": "dagnoderunners",
      "singularName": "dagnoderunner",
      "namespaced": true,
      "kind": "DagNodeRunner",
      "verbs": [
        "delete",
        "deletecollection",
        "get",
        "list",
        "patch",
        "create",
        "update",
        "watch"
      ],
      "shortNames": [
        "dnr"
      ]
    }
  ]
}
```

## Create a CR

使用 CRD 中定义的类型，创建一个自定义 DagNodeRunner 实例（CR）

```yaml
# DagNodeRunner-kind.yaml
apiVersion: enigma.basebit.me/v1
kind: DagNodeRunner
metadata:
  name: dagnoderunner-sample
spec:
  foo: bar
```

## View CRD

```bash
# 查看集群内的 CRD 资源

➜  kubectl get crd
NAME                                   CREATED AT
dagnoderunners.enigma.basebit.me       2019-10-03T08:22:05Z
meshpolicies.authentication.istio.io   2019-02-19T16:37:28Z
policies.authentication.istio.io       2019-02-19T16:37:28Z
```

```bash
# 查看创建的 DagNodeRunner 资源

➜   kubectl get dnr
NAME                   AGE
dagnoderunner-sample   2m
```

## Notes

1. 一个 API 对象在 Etcd 里的完整资源路径，是由：Group（API 组）、Version（API 版本）和 Resource（API 资源类型）三个部分组成。
2. CRD 仅仅是自定义资源的定义，如果只是创建了 CRD 并没有多大用处，必须配合 Custom Controller 一起使用，Custom Controller 可以去监听 CRD 的 CRUD 事件来添加自定义业务逻辑。

## References

- [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [Future of CRDs: Structural Schemas](https://kubernetes.io/blog/2019/06/20/crd-structural-schema/)
- [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/)
- [Extending Kubernetes APIs with Custom Resource Definitions (CRDs)](https://medium.com/velotio-perspectives/extending-kubernetes-apis-with-custom-resource-definitions-crds-139c99ed3477)
- [使用CRD扩展Kubernetes API](https://jimmysong.io/kubernetes-handbook/concepts/crd.html)
