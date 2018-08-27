# Build a Docker image from existing Golang source code and push it to Docker Hub

> multi-container-pod has a golang-server-container container and a redis-container container
```bash
cd src
docker build -t lqshow/k8s-multi-container-pod:0.0.1 .
docker push lqshow/k8s-multi-container-pod:0.0.1
```

# Deploy the app to Kubernetes
```bash
cd kubernetes

kubectl create -f deployment.yml --record
```

# Connecting Applications with Services via environment variables(localhost)

### Get enviroment variables
```go
redisHost := os.Getenv("REDIS_HOST")
redisPort := os.Getenv("REDIS_PORT")
```

# Test the app by accessing the NodePort of one of the nodes.
```bash
curl localhost:30083
```
