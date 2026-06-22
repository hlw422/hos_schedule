---
name: go-scaffold
description: 快速创建 Go + Gin + PostgreSQL + Redis 项目脚手架，包含标准分层架构、JWT认证、中间件、Docker部署
---

# Go + Gin 项目脚手架

一键生成标准 Go Web 项目结构，适用于所有 Go 后端项目。

## 生成的项目结构

```
{project_name}/
├── cmd/server/main.go              # 启动入口
├── config/config.yaml              # 配置文件
├── internal/
│   ├── config/
│   │   ├── config.go               # Viper 配置加载
│   │   ├── database.go             # PostgreSQL GORM 连接
│   │   └── redis.go                # Redis 连接
│   ├── middleware/
│   │   ├── auth.go                 # JWT 认证中间件
│   │   └── cors.go                 # CORS 跨域
│   ├── model/                      # GORM 数据模型
│   ├── repository/                 # 数据访问层
│   ├── service/                    # 业务逻辑层
│   ├── handler/                    # HTTP 处理器
│   ├── router/router.go           # 路由注册
│   └── pkg/
│       ├── jwt/jwt.go              # JWT 工具
│       └── response/response.go    # 统一响应
├── migrations/                     # 数据库迁移
├── Dockerfile
├── docker-compose.yml
├── .gitignore
└── go.mod
```

## 使用方式

当用户要求创建 Go 后端项目时，按以下步骤执行：

1. **初始化 Go 模块**
```bash
go mod init {module_name}
```

2. **安装核心依赖**
```bash
go get github.com/gin-gonic/gin
go get github.com/spf13/viper
go get github.com/redis/go-redis/v9
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
```

3. **创建配置文件** — 参考下方模板

4. **创建核心模块** — config → middleware → pkg → model → repository → service → handler → router

5. **创建 Docker 配置** — Dockerfile + docker-compose.yml

6. **验证编译** — `go build ./cmd/server`

## 配置模板

### config/config.yaml
```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: {db_name}
  sslmode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: {random_secret}
  expire: 24h
```

### 统一响应格式
```go
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}
```

### 分层架构规范
- **model**: 纯数据结构，GORM 标签 + JSON 标签
- **repository**: 数据库操作，接收 `*gorm.DB`
- **service**: 业务逻辑，接收 repository
- **handler**: HTTP 处理，接收 service，调用 response 返回

### Docker 模板
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

## 注意事项

- Windows 环境使用 `GOPROXY=https://goproxy.cn,direct`
- Gin v1.10.0 本地缓存稳定，离线时优先使用
- PostgreSQL 连接串格式: `host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`
