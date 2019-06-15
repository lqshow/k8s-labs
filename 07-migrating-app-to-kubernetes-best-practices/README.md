# Overview

当前 kubernetes 盛行，相信大家对如何将现有的传统应用迁移到 kubernetes 中比较关心。

改造一般涉及以下几个方面

1. 应用程序容器化
2. 配置文件与代码解耦，将配置信息作为 ConfigMap 资源注入到 kubernetes 中
3. 把日志当做事件流
4. 网络访问

## Containerize your application

首选需要将应用容器化，使用 Docker 定义一个包含安装步骤的 Dockerfile 文件，最后构建镜像。

## Configuration files

使用 ConfigMap 来将配置和代码解耦， 在 Kubernetes 中通过 ConfigMap 资源将应用的配置文件加载到集群中。

kubernetes 提供两种方式将  ConfigMap 注入到应用程序。

1. 环境变量
2. 存储卷(volume)

### Store config in the environment (The Twelve Factors)

> **12-Factor推荐将应用的配置存储于 环境变量 中**（ *env vars*, *env* ）。环境变量可以非常方便地在不同的部署间做修改，却不动一行代码。
>
> 缺点：不支持热更新

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
  namespace: dev
data:
  NODE_ENV: production
  DEV_MODE: "false"
  PORT: "3000"
  ENABLE_SWAGGER_UI: "true"
  
---
apiVersion: v1
kind: Pod
...
spec:
  containers:
  - name: app-container
    # 将 ConfigMap 中的所有键-值对配置为容器环境变量
    envFrom:
      - configMapRef:
          name: app-config
    env:
      # 为容器定义环境变量
      - name: ENV_NAME
        value: "ENV_VALUE"
      # 通过环境变量向容器公开Pod信息
      - name: SPARK_DRIVER_HOST
    	valueFrom:
    	  fieldRef:
    	    fieldPath: status.hostIP
      # 将 ConfigMap 中指定的键配置到环境变量中
      - name: NODE_ENV
      	valueFrom:
      	  configMapKeyRef:
      	    name: app-config
      	    key: NODE_ENV
...
```

### Use config in volume

> 用作 Pod 启动时挂载的 volume
>
> 优点：支持热更新

1. 挂载文件

   比如项目不同目录下有两个配置文件需要挂载

   - /project/service/config.py
   - /project/deploy/supervisord.conf

   ```yaml
   apiVersion: v1
   kind: Pod
   ...
   spec:
     containers:
     - name: app-container
       volumeMounts:
         - name: app-config-volume
           mountPath: /project/service/config.py
           subPath: config.py
         - name: app-config-volume
           mountPath: /project/deploy/supervisord.conf
           subPath: supervisord.conf
     volumes:
       - name: app-config-volume
         configMap:
           name: app-config
   ...
   ```

2. 挂载目录

   - 如果一个目录下的文件都是需要配置的，通过挂载整个目录的方式来实现

   - 生成 configmap

     ```bash
     kubectl create configmap hadoop-conf --from-file=conf/hadoop-conf
     ```

   - 需要注意的是，挂载目录会抹掉目录下原有的文件

   ```yaml
   apiVersion: v1
   kind: Pod
   ...
   spec:
     containers:
     - name: app-container
       volumeMounts:
         - name: hadoop-config-volume
           mountPath: /etc/hadoop/conf
     volumes:
       - name: hadoop-config-volume
         configMap:
           defaultMode: 0744
           name: hadoop-conf
   ...
   ```

## Logs

> 日志输出到 stdout 和 stderr

1. 日志应该是 [事件流](https://adam.herokuapp.com/past/2011/4/1/logs_are_streams_not_files/) 的汇总，将所有运行中进程和后端服务的输出流按照时间顺序收集起来。日志最原始的格式是一个事件一行。
2. 应用本身不应该考虑存储自己的输出流。这种方式难以维护管理，运维人员不可能一一去了解每个应用，将日志目录挂载出来。
3. 应用容器化后的日志默认会保存在宿主机的 /var/lib/docker/containers/{{. 容器 ID}}/{{. 容器 ID}}-json.log 文件里，这个目录正是类似 [Logplex](https://github.com/heroku/logplex) 和 [Fluentd](https://github.com/fluent/fluentd) 等开源工具搜集的目标。

kubernets 容器集群对日志收集的解决方案详见: [Kubernetes Log Analysis with Fluentd, Elasticsearch and Kibana](https://github.com/lqshow/k8s-labs/tree/master/lab20-efk)

## Networking

1. 一个Pod一个IP，Pod内的容器之间通信，直接通过 localhost 访问。
2. 集群内访问走Service
   - <自定义的访问方式名称>.<工作负载所在命名空间> （例如：redis-svc.default）

   - <自定义的访问方式名称>.<工作负载所在命名空间>.svc.cluster.local（例如：redis-svc.default.svc.cluster.local）
3. 集群外访问走Ingress

## Other

1. 持久化的数据写入到分布式存储卷
2. 控制输出到stdout和stderr的日志写入量
3. 使用 Secret 来保护API Key这样的隐私数据
4. 使用 liveness 和 readiness 探针来实现健康检查

## References

- [The Twelve-Factor App](https://12factor.net/)
- [How do you build 12-factor apps using Kubernetes?](https://www.mirantis.com/blog/how-do-you-build-12-factor-apps-using-kubernetes/)
- [Getting started with Kubernetes for your SaaS](https://medium.freecodecamp.org/getting-started-with-kubernetes-for-your-saas-91e91116dd7d)
- [Kubernetes Log Analysis with Fluentd, Elasticsearch and Kibana](https://github.com/lqshow/k8s-labs/tree/master/lab20-efk)

- [Five tips to move your project to Kubernetes](https://blog.alexellis.io/move-your-project-to-kubernetes/)
- [Cloud Migration Best Practices: How to Move Your Project to Kubernetes](https://dzone.com/articles/cloud-migration-best-practices-how-to-move-your-pr)
- [Moving a VM based app to Kubernetes](https://cloud.ibm.com/docs/tutorials?topic=solution-tutorials-vm-to-containers-and-kubernetes)
- [How To Migrate a Docker Compose Workflow to Kubernetes](https://www.digitalocean.com/community/tutorials/how-to-migrate-a-docker-compose-workflow-to-kubernetes)
- [Migrating to Kubernetes with zero downtime: the why and how](https://www.manifold.co/blog/migrating-to-kubernetes-with-zero-downtime-the-why-and-how-d64ba9a92619)
- [Migrating To Kubernetes](https://gravitational.com/blog/migrating-to-kubernetes/)
- [Best Practices for Migrating Your Apps to Containers and Kubernetes](https://mapr.com/blog/best-practices-for-migrating-your-apps-to-containers-and-kubernetes/)
- [Migrating a Spring Boot service to Kubernetes in 5 steps](https://itnext.io/migrating-a-spring-boot-service-to-kubernetes-in-5-steps-7c1702da81b6)
- [Migrating a monolithic / legacy app and DB to Docker and Kubernetes](https://medium.com/faun/migrating-a-monolithic-legacy-app-and-db-to-docker-and-kubernetes-efb314af6656)

