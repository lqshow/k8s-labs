# Overview

对于容器来说，我们一般建议应用把日志输出到 stdout 和 stderr 中，Docker logs 默认会保存在宿主机的 /var/lib/docker/containers/{{. 容器 ID}}/{{. 容器 ID}}-json.log 文件里。然后再由 [EFK](../../08-efk-logging-stack) 统一去收集日志。

但是有些老项目并没有遵守这个规则，直接将日志输出到了容器里的某个文件中，我们可以通过 Sidecar 容器将这些日志文件重新输出到 Sidecar 容器的 stdout 和 stderr 中，这样我们就可以由 [EFK](../../08-efk-logging-stack) 统一去收集日志，同时能通过 kubectl logs 命令方便的查看容器的日志。

![sidecar](https://user-images.githubusercontent.com/8086910/61169726-74364500-a593-11e9-9df1-e43b9900a13a.png)


## Pod manifest

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: counter
spec:
  containers:
  - name: count
    image: busybox
    args:
    - /bin/sh
    - -c
    - >
      i=0;
      while true;
      do
        echo "$(date) INFO $i" >> /var/log/date.log;
        i=$((i+1));
        sleep 1;
      done
    volumeMounts:
    - name: varlog
      mountPath: /var/log
  - name: log
    image: busybox
    args: [/bin/sh, -c, 'tail -n+1 -f /var/log/date.log']
    volumeMounts:
    - name: varlog
      mountPath: /var/log
  volumes:
  - name: varlog
    emptyDir: {}
```

```bash
# 查看日志命令

kubectl logs -f counter log
```

## Notes

当然这种不是最合适的处理方式，因为这样会导致宿主机上存在两份相同的日志文件，如果应用流量特别大，对磁盘会造成很大的浪费，建议还是修改应用容器的日志输出方式才是上策。

## References

- [让日志无处可逃：容器日志收集与管理](https://time.geekbang.org/column/article/73156)
