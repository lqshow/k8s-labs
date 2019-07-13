# Overview

目前所有前端产品的部署流程是这样的，首先由前端项目打包提供纯静态文件，然后基于 nginx 实现静态网页的部署。

现在的做法是将静态文件和 nginx 定制在一个镜像内，由于打包需要 nodejs 的环境，所以最终生产的镜像根据依赖的不同一般都达到了1G 以上。

目前升级 app 的内容，哪怕只更改一个 app 的配置或者调整 nginx 的配置，都需要重新制作一个新的镜像发布，这样的流程非常麻烦。即使事先将配置做成卷映射，但是镜像依然巨大。

## Sidecar pattern

Sidecar 指的就是我们可以在一个 Pod 中，启动一个或多个辅助容器，来完成一些独立于主进程（主容器）之外的工作。

它主要利用在同一 Pod 中的容器可以共享存储空间的能力。

![sidecard](https://user-images.githubusercontent.com/8086910/59548973-fcbbc880-8f88-11e9-964b-cf1f1001cac8.png)

## Frontend Deployment

根据 Sidecar 的容器设计模式，我们可以很容易想到，这里 nginx 应该作为主容器，前端项目作为辅助容器，只需负责提供静态文件即可。

这样部署有以下3个好处

1. 辅助容器可以利用 Docker 多阶段构建来生产出更小的镜像，因为只提供静态文件，使部署更快。(将原先将近1G 的镜像缩减到10M左右)
2. 将 nginx 主容器独立出来，后续升级配置可以统一管控。
3. 解决了 app 中静态文件 和 nginx 之间的耦合关系，做到了每个容器职责分明。

### Step 1: Frontend Dockerfile

通过多阶段构建生产 mini 镜像

```dockerfile
FROM node:latest as builder

# Create app directory
WORKDIR /data/project

# Install app dependencies
COPY package.json ./
RUN npm install

# Bundle app source
COPY ./ ./
RUN npm run build

FROM alpine

WORKDIR /project/dist
COPY --from=builder /data/project/dist  ./
```

```bash
docker build -t lqshow/app-test .
```

### Step 2: Create configmap for Nginx configuration

将 nginx 的配置通过 Configmap 方式注入到容器

```yml
kind: ConfigMap
apiVersion: v1
metadata:
  name: nginx-config
data:
  nginx.conf: |
    user              nginx;
    worker_processes  1;
    error_log  /var/log/nginx/error.log;
    pid        /var/run/nginx.pid;
    events {
        worker_connections  1024;
    }
    http {
        include       /etc/nginx/mime.types;
        default_type  application/octet-stream;
        log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                          '$status $body_bytes_sent "$http_referer" '
                          '"$http_user_agent" "$http_x_forwarded_for"';
        access_log  /var/log/nginx/access.log  main;
        sendfile        on;

        keepalive_timeout  65;

        # Load config files from the /etc/nginx/conf.d directory
        # The default server is in conf.d/default.conf
        include /etc/nginx/conf.d/*.conf;
    }
  default.conf: |
    server {
        listen       80;
        server_name  localhost;
        gzip on;
        gzip_comp_level 9;
        gzip_vary on;
        gzip_static on;
        gzip_types text/plain application/x-javascript text/css application/xml application/json application/javascript application/x-httpd-php image/jpeg image/gif image/png image/svg+xml xml/svg;
        
        # Set nginx to serve files from the shared volume
        root   /usr/share/nginx/data/project;

        location / {
            try_files $uri /index.html;
        }
    }
```

### Step 3: Create deployment

在这个 Pod 模板中，定义了两个容器

1. lqshow/app-test（这个镜像很简单，只提供静态文件，将文件放在 /var/www/html 下）
2. nginx:1.15.2（标准的 nginx 镜像）

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deploy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx

  template:
    metadata:
      labels:
        app: nginx
    spec:
      initContainers:
      - name: app-container
        image: lqshow/app-test
        imagePullPolicy: Always
        command: ["/bin/sh", "-c", "mkdir -p /var/www/html && cp -r /project/dist/. /var/www/html"]
        volumeMounts:
          - name: shared-files
            mountPath: /var/www/html
          - name: app-config-volume
            mountPath: /project/dist/js/apiUrl.js
            subPath: api-url
 
      containers:
      - name: nginx-container
        image: nginx:1.15.2
        imagePullPolicy: Always
        ports:
          - containerPort: 80
        volumeMounts:
          - name: shared-files
            mountPath: /usr/share/nginx/data/project
          - name: nginx-config-volume
            mountPath: /etc/nginx/nginx.conf
            subPath: nginx.conf
          - name: nginx-config-volume
            mountPath: /etc/nginx/conf.d/default.conf
            subPath: default.conf
           
      volumes:
        # 共享文件卷，用于共享 app 静态文件
        - name: shared-files
          emptyDir: {}

        # ConfigMap 向 app 容器注入配置信息
        - name: app-config-volume
          configMap:
            name: app-config

        # ConfigMap 向 ngnix 容器注入配置信息
        - name: nginx-config-volume
          configMap:
            name: nginx-config
```

## Reference

- [Multi-container pods and container communication in Kubernetes](https://www.mirantis.com/blog/multi-container-pods-and-container-communication-in-kubernetes/)
