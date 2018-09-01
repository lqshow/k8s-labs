# Build a Docker image from existing Golang source code and push it to Docker Hub

```bash
cd src
docker build -t lqshow/golang-redis:0.0.1 .
docker push lqshow/golang-redis:0.0.1
```

# Deploy the app to Kubernetes
```bash
cd kubernetes

kubectl create -f redis-deploy.yml --record
kubectl create -f app-deploy.yml --record
```

# Connecting Applications with Services via environment variables
```bash
# Check Environment Variables
kubectl exec redis-deploy-96bdf866c-jwsj5 -- printenv | grep SERVICE
```
### Get enviroment variables
```go
redisHost := os.Getenv("REDIS_SVC_SERVICE_HOST")
redisPort := os.Getenv("REDIS_SVC_SERVICE_PORT")
```

# Test the app by accessing the NodePort of one of the nodes.
```bash
curl localhost:30081
```
![service](https://user-images.githubusercontent.com/8086910/44656248-c37e9000-aa2a-11e8-9bc2-87e26da09bd6.gif)
