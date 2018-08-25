# Deploy the app to Kubernetes

```bash
kubectl create -f deployment.yml --record
```

# Check that the Deploy, Replica Set, Pod and Service are created

```bash
kubectl get deploy -o wide
kubectl get rs
kubectl get po -o wide -l app=webapp --show-labels
kubectl get svc -o wide -l app=webapp
```

# Test the app by accessing the NodePort of one of the nodes.
```bash
curl <NODE_IP>:30080
```

