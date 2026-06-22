---
name: docker-compose
description: 生成 Go API + PostgreSQL + Redis + Nginx 的 Docker Compose 部署配置
---

# Docker Compose 全栈部署模板

生成标准的 Go Web 应用 Docker Compose 配置。

## 服务组成

| 服务 | 镜像 | 端口 | 说明 |
|------|------|------|------|
| api | 自构建 | 8080 | Go API 服务 |
| postgres | postgres:16-alpine | 5432 | 数据库 |
| redis | redis:7-alpine | 6379 | 缓存 |
| nginx | nginx:alpine | 80, 443 | 反向代理 |

## docker-compose.yml 模板

```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=postgres
      - DB_PASSWORD=postgres
      - DB_NAME={db_name}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    volumes:
      - ./config:/root/config
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: {db_name}
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redisdata:/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
      - ./frontend/dist:/usr/share/nginx/html
    depends_on:
      - api
    restart: unless-stopped

volumes:
  pgdata:
  redisdata:
```

## nginx.conf 模板

```nginx
events {
    worker_connections 1024;
}

http {
    upstream api {
        server api:8080;
    }

    server {
        listen 80;
        server_name {domain};

        location /api/ {
            proxy_pass http://api;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }

        location / {
            root /usr/share/nginx/html;
            try_files $uri $uri/ /index.html;
        }
    }
}
```

## Dockerfile 模板

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/config ./config
EXPOSE 8080
CMD ["./server"]
```

## 使用方式

1. 在项目根目录创建 `docker-compose.yml`
2. 创建 `nginx/nginx.conf`
3. 创建 `Dockerfile`
4. 运行 `docker-compose up -d`
