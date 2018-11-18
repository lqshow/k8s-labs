## Overvew
每个节点以 Daemonset 的形式跑一个 Fluentd， 通过 Fluentd 作为 Logging agent 收集日志，并发送到后端的 Elasticsearch， Kibana 从 Elasticsearch 中获取日志进行可视化展示。

## Logging agent (Fluentd)
Fluentd 是一个用于统一日志记录层的开源数据收集器.
我们将使用 Fluentd pods 来收集存储 在Kubernetes 集群中各个节点中的所有日志, 并发送到Elasticsearch.
这些日志可以在集群中的 /var/log/containers 目录下找到.

#### Docker Container logs
Docker logs 默认会保存在宿主机的 /var/lib/docker/containers/{{. 容器 ID}}/{{. 容器 ID}}-json.log 文件里，所以这个目录正是 fluentd 的搜集目标。

#### DaemonSet
对于 Kubernetes 来说，DaemonSet 确保集群内中每一个 Node 上都能启用一个 Pod 副本。

因此，Fluentd 被部署为 DaemonSet，以 DaemonSet 的形式运行在 Kubernetes 集群中，这样就可以保证集群中每个 Node 上都会启动一个 Fluentd， 通过 fluentd 将 Docker 容器里的日志转发到 ElasticSearch 中。

为了使 Fluentd 能够工作，每个 Node 都必须标记 beta.kubernetes.io/fluentd-ds-ready=true。

#### Fluentd configuration

1. 定义了 Fluentd 的 ConfigMap 配置文件，此文件定义了 Fluentd 所获取的日志数据源，以及将这些日志数据输出到 Elasticsearch 中。
    ```bash
    # 创建 ConfigMap
    kubectl apply -f fluentd-es-configmap.yaml
    ```

2. 定义了一个名称为 fluentd-es 的 ServiceAccount，并授予其能够对 namespaces 和 Pods 读取的访问权限， 并以 DaemonSet 类型部署 Fluentd
    ```bash
    # 1. 创建 ServiceAccount
    # 2. 创建 ClusterRole
    # 3. 创建 ClusterRoleBinding
    # 4. 创建 DaemonSet
    #    - 这个 DaemonSet，管理的是一个 fluentd-elasticsearch 镜像的 Pod。
    #    - 这个镜像的功能非常实用：通过 fluentd 将 Docker 容器里的日志转发到 ElasticSearch 中
    #    - 这个容器挂载了两个 hostPath 类型的 Volume，分别对应宿主机的 /var/log 目录和 /var/lib/docker/containers 目录.
    #    - fluentd 启动之后，它会从这两个目录里搜集日志信息，并转发给 ElasticSearch 保存。
    kubectl apply -f fluentd-es-ds.yaml
    ```

## Logging Backend (Elasticsearch)
Elasticsearch 是一种负责存储日志并允许查询的搜索引擎

#### Elasticsearch configuration

1. 定义了一个名称为 elasticsearch-logging 的 ServiceAccount，并授予其能够对 services、namespaces 和endpoints 读取的访问权限；并以 StatefulSet 类型部署 Elasticsearch。
    ```bash
    # 1. 创建 ServiceAccount
    # 2. 创建 ClusterRole
    # 3. 创建 ClusterRoleBinding
    # 4. 创建 StatefulSet
    kubectl apply -f es-statefulset.yaml
    ```
2. 为 Elasticsearch 暴露 9200 端口
    ```bash
    # 创建 Service
    kubectl apply -f es-service.yaml
    ```
3. 检查 Elasticsearch 部署情况
    ```bash
    ➜  lab20-efk git:(master) ✗ kubectl get svc -owide -n kube-system -lk8s-app=elasticsearch-logging
    NAME                    TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE       SELECTOR
    elasticsearch-logging   ClusterIP   10.101.119.167   <none>        9200/TCP   17m       k8s-app=elasticsearch-logging
    ```
    ```bash
    ➜  lab20-efk git:(master) ✗ kubectl exec -it busybox-curl -- curl -XGET '10.101.119.167:9200/_search?pretty'
    {
    "took" : 1,
    "timed_out" : false,
    "_shards" : {
        "total" : 0,
        "successful" : 0,
        "skipped" : 0,
        "failed" : 0
    },
    "hits" : {
        "total" : 0,
        "max_score" : 0.0,
        "hits" : [ ]
    }
    }
    ```
## Logging display (Kibana)
Kibana 是一个图形界面，用于查看和查询存储在 Elasticsearch 中的日志


## References
- [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)
- [fluentd-elasticsearch](https://github.com/kubernetes/kubernetes/tree/master/cluster/addons/fluentd-elasticsearch)
- [Kubernetes Log Analysis with Fluentd, Elasticsearch and Kibana](https://logz.io/blog/kubernetes-log-analysis/)
- [fluentd-kubernetes-daemonset](https://github.com/fluent/fluentd-kubernetes-daemonset)
- [Kubernetes Logging with Fluentd](https://docs.fluentd.org/v0.12/articles/kubernetes-fluentd)
- [Logging Architecture](https://kubernetes.io/docs/concepts/cluster-administration/logging/)

