## Overview

Dashboard 默认的 token 失效时间为 900 秒，即 15 分钟。

```bash
# 通过修改参数 token-ttl 即可调整失效时间
kubectl edit deploy kubernetes-dashboard -n kube-system
```

```yaml
spec:
  containers:
  - args:
    - --auto-generate-certificates
    - --token-ttl=86400
```

## Referneces

- [DefaultTokenTTL](https://github.com/kubernetes/dashboard/blob/master/src/app/backend/auth/api/types.go#L29)
- [Dashboard arguments](https://github.com/kubernetes/dashboard/wiki/Dashboard-arguments)

