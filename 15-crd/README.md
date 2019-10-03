

# Overview

1. CRD 是 Kubernetes 的一种资源类型，是 Custom Resource Definition 的缩写
2. CRD 是自定义资源的定义
3. Kubernetes 允许用户自定义资源 CRD，向 Kubernetes 集群注册一种新资源，用于扩展 Kubernetes 集群能力
4. 当我们想要将自己的对象引入 Kubernetes 集群以完全满足我们的需求时，就需要使用 CRD。一旦我们在 Kubernetes 中创建了CRD，我们就可以像使用其他 Kubernetes 内置资源类型一样使用它


## Create a CRD

以下是一个 CRD 样例

```yaml
# crd.yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: dagnoderunners.enigma.basebit.me
spec:
  group: enigma.basebit.me
  names:
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
  # made sure that this CRD is a namespaced and not cluster wide
  scope: Namespaced
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
        status:
          description: DagNodeRunnerStatus defines the observed state of DagNodeRunner
          type: object
```

使用 CRD 中定义的类型，创建一个 DagNodeRunner 对象

```yaml
# DagNodeRunner-kind.yaml
apiVersion: enigma.basebit.me/v1
kind: DagNodeRunner
metadata:
  name: dagnoderunner-sample
spec:
  foo: bar
```

## 查看 CRD

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

## References

- [Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [Extend the Kubernetes API with CustomResourceDefinitions](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/)
- [Extending Kubernetes APIs with Custom Resource Definitions (CRDs)](https://medium.com/velotio-perspectives/extending-kubernetes-apis-with-custom-resource-definitions-crds-139c99ed3477)
