# Overview

Helm2 to  Helm 3 一些变化

1. 移除了 Tiller 服务器端组件，更简单，更安全
   - Helm 3 使用与 kubectl 上下文相同的访问权限
   - Helm 3 允许同一个 release name 存在不同的 namespace 中，即 release 限制在 namespace 范围内， 不再是全局资源
   - Helm 3 release 区分命名空间，需要带上 -A 参数，显示所有命名空间
   - Helm 3 同时删除了 helm init 和 home，无需再使用 helm init 来初始化 Helm
2. 如果希望 chart 既可以用于 Helm 2，也可以用于 Helm 3，请确保它们创建了命名空间并同时使用crd install钩子和crds/目录；Helm 3将忽略该钩子，并发出警告。
3. charts 变化
    - 把 requirements.yaml 合并成 Chart.yaml
    - 把 Chart.yaml 配置中 apiVersion: v1 修改成 v2
4. helm cli命令重命名
    ```bash
   # v2中删除项目需要提供 --purge 参数，v3默认情况下启用此功能。要保留以前的行为，请使用 helm uninstall --keep-history
    helm delete  重命名为 helm uninstall
    helm fetch   重命名为 helm pull
    helm inspect 重命名为 helm show
    ```

## 下载 Helm 3 二进制文件

1. 下载 [Helm 3](https://github.com/helm/helm/releases) 二进制文件

    ```bash
    # 将现有的v2二进制文件重命名为 helm2
    mv /usr/local/bin/helm /usr/local/bin/helm2

    # 将最新版本重命名为 helm3
    tar -zxvf helm-v3.3.4-linux-amd64.tar.gz
    mv linux-amd64/helm /usr/local/bin/helm3
    ```
    
## 迁移 Helm 3
> 已安装的 Kubernetes 对象将不会被修改或删除，迁移过程中不会影响线上运行的服务

1. 安装 helm-2to3 插件
    ```bash
    helm3 plugin install https://github.com/helm/helm-2to3

    # 验证 2to3 插件是否安装成功
    ➜ helm3 plugin list
    NAME    VERSION DESCRIPTION
    2to3    0.7.0   migrate and cleanup Helm v2 configuration and releases in-place to Helm v3
    ```
    
2. 迁移 Helm V2 configuration
    
    > 主要迁移 helm repo 以及 plugin
    
    ```bash
    # 迁移前的 repo 列表
    ➜ helm3 repo list
    NAME    URL
    stable  https://kubernetes-charts.storage.googleapis.com/
 
    # 执行迁移命令
    helm3 2to3 move config
    
    # 迁移后的 repo 列表
    ➜ helm3 repo list
    NAME                    URL
    stable                  https://kubernetes-charts.storage.googleapis.com
    local                   http://127.0.0.1:8879/charts
    basebit-chartmuseum     http://chartmuseum.basebit.me
    harbor                  https://helm.goharbor.io
    ingress-nginx           https://kubernetes.github.io/ingress-nginx/
    ```
    
 3. 迁移 Helm V2 Release
 
    ```bash
    # 模拟迁移，可测试下效果
    helm3 2to3 convert <RELEASE NAME> --dry-run

    # 执行迁移命令
    helm3 2to3 convert <RELEASE NAME>

    # 验证迁移结果
    helm3 list
    ```   
    
    ```bash
    # List Helm 2 Releases
    RELEASES=$(helm2 list -aq)
    # Loop through releases and, for each one, test conversion
    while IFS= read -r release; do
      helm3 2to3 convert $release --dry-run
    done <<< "$RELEASES"
    ```
    
4. 清理 Helm v2 data
    > note: 执行 cleanup 后，Tiller Pod 会被删除，并且 kube-system 命名空间中 configmaps 历史版本信息也会被清理
    
      ```bash
      # 模拟清理
      helm3 2to3 cleanup --dry-run
      
      # 执行清理命令
      helm3 2to3 cleanup
      ```  
## References

- [Helm 2和Helm 3的主要区别是什么?](https://helm.sh/zh/docs/faq/)
- [How to migrate from Helm v2 to Helm v3](https://helm.sh/blog/migrate-from-helm-v2-to-helm-v3/)
- [Migrating to Helm 3: What You Need to Know](https://itnext.io/additional-tips-for-migrating-to-helm-3-304c9d50f1b4)
- [Breaking Changes in Helm 3 (and How to Fix Them)](https://itnext.io/breaking-changes-in-helm-3-and-how-to-fix-them-39fea23e06ff)