# Deploy the app to Kubernetes

```bash
kubectl create -f nginx-pod.yml
kubectl create -f nginx-svc.yml
```

# Check that the Pods and Service are created

```bash
kubectl get po -o wide -l app=nginx
kubectl get svc -o wide -l app=nginx
```

# Reference
- [k8s 使用 Pod 来部署应用](https://github.com/lqshow/notes/issues/38)

