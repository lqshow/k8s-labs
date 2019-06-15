## Basic Commands

```bash
# 查看 client/server 版本信息
helm version

# 查找可用 Charts
helm search mysql

# 列出指定 repo 下的所有 charts
helm search basebit-chartmuseum/

# 根据指定的 Chart 部署一个Release 到K8s
helm install stable/mariadb

# 查看指定Chart的基本信息
helm inspect stable/mariadb

# 查看 chart 可配置选项
helm inspect values stable/mariadb
```

## Chart

```bash
# 创建一个 chart 骨架
helm create my-chart

# 检验 chart
helm lint

# 通过目录安装 chart
helm install ./my-chart

# 指定 release 名称
helm install ./my-chart --name helm-test

# 测试模板呈现，模拟安装(调试)
helm install --debug --dry-run ./mychart

# 打包 chart
helm package ./my-chart

# 列出指定 chart 的所有版本
helm search -l stable/<some_chart>

# 下载指定版本的 chart
helm fetch stable/postgresql --version=3.7.1
```

## Releases

```bash
# 查看所有已部署的 releases
helm ls

# 查看被删除的 releases
helm ls --deleted

# 查看 release 的清单
helm get manifest <RELEASE_NAME>

# 查看 release 状态（删除了，依然可以查看状态）
helm status <RELEASE_NAME>

# 删除 release
helm delete <RELEASE_NAME>

# 查看 release 历史版本
helm history <RELEASE_NAME>

# Upgrade a release
helm upgrade --set foo=bar --set foo=newbar redis stable/redis

# roll back a release to a previous revision
helm rollback my-release

# 回滚导指定版本
helm rollback <RELEASE_NAME> <VERSION>
```

## Repo

```bash
# 查看当前的仓库配置
helm repo list

# 更新仓库
helm repo update
```
