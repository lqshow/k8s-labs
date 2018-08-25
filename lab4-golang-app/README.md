# Build a Docker image from existing Go source code and push it to Docker Hub

```bash
cd src
docker build -t lqshow/golang-app:0.0.1 .
docker push lqshow/golang-app:0.0.1
```

# Launch the app with Docker Compose
```bash
docker-compose up -d
```

# Test the app
```bash
curl localhost:3000
```

# Deploy the app to Kubernetes
```bash
cd kubernetes

kubectl create -f deploy.yml --record
```

# Check that the Pods and Services are created
```bash
kubectl get deploy -o wide
kubectl get rs
kubectl get po -o wide -l app=golang-app --show-labels
kubectl get svc -o wide -l app=golang-app
```

# Test the app by accessing the NodePort of one of the nodes.
```bash
curl localhost:30081
```
![deploy2](https://user-images.githubusercontent.com/8086910/44620267-b004e500-a8c3-11e8-9373-c8afe36ba256.gif)


# Reference
- [Creating a RESTful API With Golang](https://tutorialedge.net/golang/creating-restful-api-with-golang/)

