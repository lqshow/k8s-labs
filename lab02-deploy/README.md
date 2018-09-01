# Deploy the app to Kubernetes

```bash
kubectl create -f webapp-deploy.yml --record
kubectl create -f webapp-svc.yml
```

# Check that the Deploy, Replica Set, Pod and Service are created

```bash
kubectl get deploy -o wide
kubectl get rs
kubectl get po -o wide -l app=webapp --show-labels
kubectl get svc -o wide -l app=webapp
```

# Reference
- [k8s 使用 Deployment 部署应用](https://github.com/lqshow/notes/issues/39)

