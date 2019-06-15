## Overview

在容器中做 volume mount，默认情况下都是将容器中的数据 mount 中的 root 目录下。

如果一个 Pod 中存在多个容器，每个容器需要将数据挂载到同一个持久化存储卷的不同目录下，或者不同 Pod 中的不同容器需要将数据保存到同一个持久化存储卷的不同目录下， 就需要用到 subPath。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-using-sub-path
spec:
    containers:
    - name: mysql
      image: mysql:5.7.25
      env:
      - name: MYSQL_ROOT_PASSWORD
        value: "basebit.ai"
      volumeMounts:
      - mountPath: /var/lib/mysql
        name: persistent-local-storage
        # 将数据挂载到 volume 中的 mysql 目录下
        subPath: mysql
    - name: redis
      image: redis:5.0.4
      volumeMounts:
      - mountPath: /data
        name: persistent-local-storage
        # 将数据挂载到 volume 中的 redis 目录下
        subPath: redis
    volumes:
    - name: persistent-local-storage
      persistentVolumeClaim:
        claimName: pvc-test
```

