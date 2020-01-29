# Overview

1. Custom controllers 允许用户自定义控制器的逻辑，基于已有的资源定义更高阶的控制器，实现 Kubernetes 集群原生不支持的功能
2. Custom controllers 的逻辑其实是很简单的：watch CRD 实例（以及关联的资源）的 CRUD 事件，然后开始执行相应的业务逻辑


## Controller model

> 一个无限循环不断地对比实际状态和期望状态，如果有出入则进行调谐逻辑将实际状态调整为期望状态，最终达到与申明一致

实际状态: 来自于 Kubernetes 集群本身
期望状态: 来自于 用户提交的 YAML 文件


```golang
for {
  实际状态 := 获取集群中对象 X 的实际状态（Actual State）
  期望状态 := 获取集群中对象 X 的期望状态（Desired State）
  if 实际状态 == 期望状态{
    // do nothing
  } else {
    // reconcile loop (执行编排动作，将实际状态调整为期望状态)
  }
}
```

## flow

1. 编写 CRD
2. 编写 Custom Controller


## Kubernetes code-generator

Kubernetes 使用 CRD+Controller 来扩展集群功能，官方提供了 CRD 代码的自动生成器 code-generator


### code-generator

> 其中 informer-gen 和 lister-gen 是构建 controller 的基础

| Code-gen     | Desc                                                         |
| ------------ | ------------------------------------------------------------ |
| deepcopy-gen | 生成`func` `(t *T)` `DeepCopy()` `*T` 和`func` `(t *T)` `DeepCopyInto(*T)`方法 |
| client-gen   | 创建类型化客户端集合(typed client sets)                      |
| informer-gen | 为CR创建一个informer , 当CR有变化的时候, 这个informer可以基于事件接口获取到信息变更 |
| lister-gen   | 为CR创建一个listers , 就是为`GET` and `LIST`请求提供read-only caching layer |


code-generator 脚本

- [generate-groups.sh](https://github.com/kubernetes/code-generator/blob/master/generate-groups.sh)
- [generate-internal-groups.sh](https://github.com/kubernetes/code-generator/blob/master/generate-internal-groups.sh)


```bash
# 初始化项目
mkdir code-generator-sample && cd code-generator-sample
go mod init github.com/lqshow/code-generator-sample

# 初始化crd资源类型
mkdir -p pkg/api/foobar/v1 && cd pkg/api/foobar/v1
```

## References

- [Generators for kube-like API types](https://github.com/kubernetes/code-generator)
- [Kubernetes Deep Dive: Code Generation for CustomResources](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/)
- [The Kubebuilder Book](https://book.kubebuilder.io/introduction.html)
- [Writing Controllers](https://github.com/kubernetes/community/blob/8decfe4/contributors/devel/controllers.md)
- [浅析 Kubernetes 控制器的工作原理](https://www.yangcs.net/posts/a-deep-dive-into-kubernetes-controllers/)
- [如何使用 CRD 拓展 Kubernetes 集群](https://mp.weixin.qq.com/s/UJ6H686h_r2w_C1XS116Lw)
- [进阶 K8s 高级玩家必备 | Kubebuilder：让编写 CRD 变得更简单](https://mp.weixin.qq.com/s/UzEcj2eXKM0m8f4XzZCYAA)
- [code-generator使用](https://tangxusc.github.io/blog/2019/05/code-generator%E4%BD%BF%E7%94%A8/)
