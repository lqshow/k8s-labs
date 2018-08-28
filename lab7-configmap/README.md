> ConfigMap 可以将配置信息与 docker 镜像文件解耦，ConfigMap提供了向容器注入配置数据的机制，可用于存储细粒度信息如单个属性，或粗粒度信息如整个配置文件。

# 使用来自多个ConfigMaps的数据定义容器环境变量
```bash
kubectl create -f using-multiple-configmap.yml
kubectl exec $(kubectl get pods -l run=using-multiple-configmap-pod  -o=name|cut -d "/" -f2) env |grep -E 'MYSQL|LOG'
kubectl delete -f using-multiple-configmap.yml
```

# 将ConfigMap中的所有键-值对配置为容器环境变量
```bash
kubectl create -f using-all-key-value-pairs.yml
kubectl exec $(kubectl get pods -l run=using-all-key-value-pairs-in-configmap-pod  -o=name|cut -d "/" -f2) env |grep redis
kubectl delete -f using-all-key-value-pairs.yml
```

# 通过数据卷插件使用ConfigMap
```bash
kubectl create -f add-configmap-to-volume.yml
kubectl exec $(kubectl get pods -l run=add-configmap-to-volume-pod  -o=name|cut -d "/" -f2) cat /etc/config/game
kubectl delete -f add-configmap-to-volume.yml
```
### 挂载到 Volume 的配置文件支持热更新
```bash
# 更改lives对应的值为4后，过10来秒后，pod 中的配置信息会同步更新
kubectl edit configmap game-config
```

![configmap](https://user-images.githubusercontent.com/8086910/44737433-8dbdd200-ab24-11e8-8522-fe254c076220.gif)



# Notes

1. 当ConfigMap以数据卷的形式挂载进Pod时，更新ConfigMap（或删掉重建ConfigMap），Pod内挂载的配置信息会热更新，但使用环境变量方式加载到pod，则不会自动更新。
2. ConfigMaps 和 Pod 必须在同一个 namespace 下才能引用的到。