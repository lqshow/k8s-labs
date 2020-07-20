## 更改 VolumeMount 文件权限

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-chmod-perm
spec:
  initContainers:
  - command:
    - sh
    - -c
    - chmod -R 777 /tmp/docker
    image: lqshow/busybox-curl:1.28
    name: volume-mount-hack
    volumeMounts:
    - mountPath: /tmp/docker
      name: workspace
  containers:
  - command:
    - /bin/sh
    image: lqshow/busybox-curl:1.28
    imagePullPolicy: Always
    name: user
    stdin: true
    tty: true
    volumeMounts:
    - mountPath: /root/attached_storage
      mountPropagation: HostToContainer
      name: workspace
  volumes:
  - hostPath:
      path: /workspace/test
      type: DirectoryOrCreate
    name: workspace
```

## 同步 pod 时区，与 host 主机保持同步

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sync-host-tz
  labels:
    app: sync-host-tz
spec:
  selector:
    matchLabels:
      app: sync-host-tz

  replicas: 1
  
  template:
    metadata:
      labels:
        app: sync-host-tz
    spec:
      containers:
      - name: busybox
        image: lqshow/busybox-curl:1.28
        stdin: true
        tty: true
        volumeMounts:
        - name: host-tz
          mountPath: /etc/localtime
      volumes:
      - name: host-tz
        hostPath:
          path: /etc/localtime
```

