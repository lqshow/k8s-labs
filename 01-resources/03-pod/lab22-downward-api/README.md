## Overview
Downward API 让 Pod 里的容器能够直接获取到这个 Pod API 对象本身的信息。
Downward API 能够获取到的信息，一定是 Pod 里的容器进程启动之前就能够确定下来的信息


## Use Pod fields as values for environment variables
| fieldRef                      | desc                          |
| ----------------------------- | ----------------------------- |
| metadata.uid                  | Pod 的 UID                    |
| metadata.namespace            | Pod 的 Namespace              |
| metadata.name                 | Pod 的名字                    |
| metadata.labels               | Pod 的所有 Label              |
| metadata.annotations          | Pod 的所有 Annotation         |
| metadata.labels['<KEY>']      | 指定 <KEY> 的 Label 值        |
| metadata.annotations['<KEY>'] | 指定 <KEY> 的 Annotation 值   |
| spec.nodeName                 | 宿主机名字                    |
| spec.serviceAccountName       | Pod 的 Service Account 的名字 |
| status.hostIP                 | 宿主机 IP                     |
| status.podIP                  | Pod 的 IP                     |

## Use Container fields as values for environment variables
| resourceFieldRef | desc                  |
| ---------------- | --------------------- |
| requests.cpu     | 容器的 CPU request    |
| limits.cpu       | 容器的 CPU limit      |
| requests.memory  | 容器的 memory request |
| limits.memory    | 容器的 memory limit   |


### DownwardAPI VolumeFiles

## 参考
- [Expose Pod Information to Containers Through Environment Variables](https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/)
- [Expose Pod Information to Containers Through Files](https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/)
- [深入解析Pod对象（二）：使用进阶](https://time.geekbang.org/column/article/40466)