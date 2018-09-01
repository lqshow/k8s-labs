# Build a Docker image from existing Golang source code and push it to Docker Hub

```bash
cd src
docker build -t lqshow/k8s-skaffold-test:0.0.1 .
docker push lqshow/k8s-skaffold-test:0.0.1
```