### Headless Service 
Headless 的含义: clusterIP 字段的值是：None，即：这个 Service，没有一个 VIP 作为“头”。
不需要分配一个 Virtual IP，而是可以直接以 DNS 记录的方式解析出被代理 Pod 的 IP 地址。

### StatefulSet（有状态应用）
> 与 Deployment 的唯一区别，就是多了一个 serviceName 字段。

用于解决各个pod实例独立生命周期管理，提供各个实例的启动顺序和唯一性
- 稳定，唯一的网络标识符。
- 稳定，持久存储。
- 有序的，优雅的部署和扩展。
- 有序，优雅的删除和终止。
- 有序的自动滚动更新.


web-0.nginx-headless-svc.default.svc.cluster.local

通过 DNS 名称验证服务访问是否正常
```bash
kubectl exec -it net-test-dff6845bb-jc42c curl web-0.nginx-headless-svc
kubectl exec -it net-test-dff6845bb-jc42c curl web-1.nginx-headless-svc
```