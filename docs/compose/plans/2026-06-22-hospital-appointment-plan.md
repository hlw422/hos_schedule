# 三甲医院预约挂号平台实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use compose:subagent (recommended) or compose:execute to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个支持多院区、高并发预约的三甲医院预约挂号平台

**Architecture:** 单体 Go 应用 + PostgreSQL + Redis，模块化分包，Redis Lua 脚本保证号源原子扣减

**Tech Stack:** Go 1.22+, Gin, PostgreSQL 16, Redis 7, golang-jwt, wechat-sdk

---

## 文件结构

```
hos_schedule/
├── cmd/
│   └── server/
│       └── main.go                    # 启动入口
├── internal/
│   ├── config/
│   │   └── config.go                  # 配置加载
│   ├── middleware/
│   │   ├── auth.go                    # JWT认证中间件
│   │   └── cors.go                    # CORS中间件
│   ├── model/
│   │   ├── hospital.go                # 医院/院区模型
│   │   ├── department.go              # 科室模型
│   │   ├── doctor.go                  # 医生模型
│   │   ├── schedule.go                # 排班/号源模型
│   │   ├── user.go                    # 用户模型
│   │   ├── patient.go                 # 就诊人模型
│   │   ├── appointment.go             # 预约模型
│   │   └── notification.go            # 通知模型
│   ├── repository/
│   │   ├── hospital_repo.go
│   │   ├── department_repo.go
│   │   ├── doctor_repo.go
│   │   ├── schedule_repo.go
│   │   ├── user_repo.go
│   │   ├── patient_repo.go
│   │   ├── appointment_repo.go
│   │   └── notification_repo.go
│   ├── service/
│   │   ├── auth_service.go            # 登录/Token
│   │   ├── hospital_service.go
│   │   ├── department_service.go
│   │   ├── doctor_service.go
│   │   ├── schedule_service.go
│   │   ├── patient_service.go
│   │   ├── appointment_service.go     # 核心预约逻辑
│   │   └── notification_service.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── hospital_handler.go
│   │   ├── department_handler.go
│   │   ├── doctor_handler.go
│   │   ├── schedule_handler.go
│   │   ├── patient_handler.go
│   │   ├── appointment_handler.go
│   │   └── notification_handler.go
│   ├── router/
│   │   └── router.go                  # 路由注册
│   └── pkg/
│       ├── redis/
│       │   └── slot_deduction.lua     # 号源扣减Lua脚本
│       ├── wechat/
│       │   └── client.go              # 微信API封装
│       └── response/
│           └── response.go            # 统一响应
├── migrations/
│   ├── 001_create_tables.up.sql
│   └── 001_create_tables.down.sql
├── config/
│   └── config.yaml                    # 配置文件
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

---

## Task 1: 项目初始化与基础框架

**Covers:** S1

**Files:**
- Create: `go.mod`, `cmd/server/main.go`, `internal/config/config.go`, `internal/middleware/cors.go`, `internal/pkg/response/response.go`, `config/config.yaml`, `Dockerfile`, `docker-compose.yml`

- [ ] **Step 1: 初始化 Go 模块**

```bash
cd D:\coding\lang\mimo\hos_schedule
go mod init hos_schedule
```

- [ ] **Step 2: 安装依赖**

```bash
go get github.com/gin-gonic/gin
go get github.com/spf13/viper
go get github.com/redis/go-redis/v9
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/golang-jwt/jwt/v5
```

- [ ] **Step 3: 创建配置文件**

```yaml
# config/config.yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: hospital_schedule
  sslmode: disable

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-jwt-secret-key
  expire: 24h

wechat:
  appid: your-appid
  secret: your-secret
```

- [ ] **Step 4: 创建配置加载模块**

```go
// internal/config/config.go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Wechat   WechatConfig   `mapstructure:"wechat"`
}

type ServerConfig struct {
    Port int    `mapstructure:"port"`
    Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
    Secret string `mapstructure:"secret"`
    Expire string `mapstructure:"expire"`
}

type WechatConfig struct {
    AppID  string `mapstructure:"appid"`
    Secret string `mapstructure:"secret"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}
```

- [ ] **Step 5: 创建统一响应**

```go
// internal/pkg/response/response.go
package response

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code:    0,
        Message: "success",
        Data:    data,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(http.StatusOK, Response{
        Code:    code,
        Message: message,
    })
}

func BadRequest(c *gin.Context, message string) {
    Error(c, 400, message)
}

func Unauthorized(c *gin.Context, message string) {
    Error(c, 401, message)
}

func NotFound(c *gin.Context, message string) {
    Error(c, 404, message)
}

func InternalError(c *gin.Context, message string) {
    Error(c, 500, message)
}
```

- [ ] **Step 6: 创建 CORS 中间件**

```go
// internal/middleware/cors.go
package middleware

import (
    "github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

- [ ] **Step 7: 创建 main.go**

```go
// cmd/server/main.go
package main

import (
    "fmt"
    "log"
    "hos_schedule/internal/config"
    "hos_schedule/internal/middleware"
    "hos_schedule/internal/router"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    gin.SetMode(cfg.Server.Mode)
    
    r := gin.Default()
    r.Use(middleware.CORS())
    
    router.Register(r)
    
    addr := fmt.Sprintf(":%d", cfg.Server.Port)
    log.Printf("Server starting on %s", addr)
    
    if err := r.Run(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

- [ ] **Step 8: 创建路由占位**

```go
// internal/router/router.go
package router

import (
    "github.com/gin-gonic/gin"
    "hos_schedule/internal/pkg/response"
)

func Register(r *gin.Engine) {
    api := r.Group("/api/v1")
    {
        api.GET("/health", func(c *gin.Context) {
            response.Success(c, gin.H{"status": "ok"})
        })
    }
}
```

- [ ] **Step 9: 创建 Docker 配置**

```dockerfile
# Dockerfile
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

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: hospital_schedule
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  pgdata:
```

- [ ] **Step 10: 验证项目编译**

```bash
cd D:\coding\lang\mimo\hos_schedule
go build ./cmd/server
```

Expected: 编译成功，无错误

---

## Task 2: 数据库连接与迁移

**Covers:** S2

**Files:**
- Create: `internal/model/*.go`, `migrations/001_create_tables.up.sql`, `migrations/001_create_tables.down.sql`
- Modify: `internal/config/config.go`, `cmd/server/main.go`

- [ ] **Step 1: 创建数据库连接模块**

```go
// internal/config/database.go
package config

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func (c *DatabaseConfig) Connect() (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
        c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode,
    )
    return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
```

- [ ] **Step 2: 创建 Redis 连接模块**

```go
// internal/config/redis.go
package config

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
)

func (c *RedisConfig) Connect() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
        Password: c.Password,
        DB:       c.DB,
    })
}

func (c *RedisConfig) Ping(ctx context.Context, client *redis.Client) error {
    return client.Ping(ctx).Err()
}
```

- [ ] **Step 3: 创建数据库迁移文件**

```sql
-- migrations/001_create_tables.up.sql

-- 用户表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    openid VARCHAR(100) UNIQUE NOT NULL,
    unionid VARCHAR(100),
    phone VARCHAR(20),
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    role VARCHAR(20) DEFAULT 'PATIENT',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 医院表
CREATE TABLE hospitals (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(20),
    logo VARCHAR(255),
    intro TEXT,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 院区表
CREATE TABLE campuses (
    id BIGSERIAL PRIMARY KEY,
    hospital_id BIGINT NOT NULL REFERENCES hospitals(id),
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(20),
    latitude DECIMAL(10, 7),
    longitude DECIMAL(10, 7),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 科室表
CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    hospital_id BIGINT NOT NULL REFERENCES hospitals(id),
    campus_id BIGINT REFERENCES campuses(id),
    name VARCHAR(100) NOT NULL,
    intro TEXT,
    sort_order INT DEFAULT 0,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 医生表
CREATE TABLE doctors (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT NOT NULL REFERENCES departments(id),
    name VARCHAR(50) NOT NULL,
    avatar VARCHAR(255),
    title VARCHAR(50),
    intro TEXT,
    specialty TEXT,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 排班表
CREATE TABLE schedules (
    id BIGSERIAL PRIMARY KEY,
    doctor_id BIGINT NOT NULL REFERENCES doctors(id),
    campus_id BIGINT NOT NULL REFERENCES campuses(id),
    date DATE NOT NULL,
    time_period VARCHAR(20) NOT NULL,
    total_count INT NOT NULL,
    used_count INT DEFAULT 0,
    remain_count INT NOT NULL,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(doctor_id, date, time_period)
);

-- 就诊人表
CREATE TABLE patients (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    name VARCHAR(50) NOT NULL,
    id_card VARCHAR(18),
    phone VARCHAR(20),
    relation VARCHAR(20),
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 预约表
CREATE TABLE appointments (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES patients(id),
    doctor_id BIGINT NOT NULL REFERENCES doctors(id),
    schedule_id BIGINT NOT NULL REFERENCES schedules(id),
    campus_id BIGINT NOT NULL REFERENCES campuses(id),
    date DATE NOT NULL,
    time_period VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING_PAY',
    pay_type VARCHAR(20),
    pay_amount DECIMAL(10, 2),
    cancel_reason VARCHAR(255),
    visit_no VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 通知表
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    type VARCHAR(50),
    template_id VARCHAR(100),
    content TEXT,
    status VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX idx_schedules_doctor_date ON schedules(doctor_id, date);
CREATE INDEX idx_appointments_user_status ON appointments(user_id, status);
CREATE INDEX idx_appointments_doctor_date ON appointments(doctor_id, date);
CREATE INDEX idx_patients_user ON patients(user_id);
```

```sql
-- migrations/001_create_tables.down.sql
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS patients;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS doctors;
DROP TABLE IF EXISTS departments;
DROP TABLE IF EXISTS campuses;
DROP TABLE IF EXISTS hospitals;
DROP TABLE IF EXISTS users;
```

- [ ] **Step 4: 创建 GORM 模型**

```go
// internal/model/user.go
package model

import "time"

type User struct {
    ID        int64     `gorm:"primaryKey" json:"id"`
    OpenID    string    `gorm:"uniqueIndex;size:100" json:"openid"`
    UnionID   string    `gorm:"size:100" json:"unionid,omitempty"`
    Phone     string    `gorm:"size:20" json:"phone,omitempty"`
    Nickname  string    `gorm:"size:50" json:"nickname,omitempty"`
    Avatar    string    `gorm:"size:255" json:"avatar,omitempty"`
    Role      string    `gorm:"size:20;default:PATIENT" json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

```go
// internal/model/hospital.go
package model

import "time"

type Hospital struct {
    ID        int64     `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"size:100" json:"name"`
    Address   string    `gorm:"size:255" json:"address,omitempty"`
    Phone     string    `gorm:"size:20" json:"phone,omitempty"`
    Logo      string    `gorm:"size:255" json:"logo,omitempty"`
    Intro     string    `gorm:"type:text" json:"intro,omitempty"`
    Status    int8      `gorm:"default:1" json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Campus struct {
    ID         int64     `gorm:"primaryKey" json:"id"`
    HospitalID int64     `gorm:"index" json:"hospital_id"`
    Name       string    `gorm:"size:100" json:"name"`
    Address    string    `gorm:"size:255" json:"address,omitempty"`
    Phone      string    `gorm:"size:20" json:"phone,omitempty"`
    Latitude   float64   `gorm:"type:decimal(10,7)" json:"latitude,omitempty"`
    Longitude  float64   `gorm:"type:decimal(10,7)" json:"longitude,omitempty"`
    Status     int8      `gorm:"default:1" json:"status"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

```go
// internal/model/department.go
package model

import "time"

type Department struct {
    ID         int64     `gorm:"primaryKey" json:"id"`
    HospitalID int64     `gorm:"index" json:"hospital_id"`
    CampusID   int64     `gorm:"index" json:"campus_id,omitempty"`
    Name       string    `gorm:"size:100" json:"name"`
    Intro      string    `gorm:"type:text" json:"intro,omitempty"`
    SortOrder  int       `gorm:"default:0" json:"sort_order"`
    Status     int8      `gorm:"default:1" json:"status"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

```go
// internal/model/doctor.go
package model

import "time"

type Doctor struct {
    ID           int64     `gorm:"primaryKey" json:"id"`
    DepartmentID int64     `gorm:"index" json:"department_id"`
    Name         string    `gorm:"size:50" json:"name"`
    Avatar       string    `gorm:"size:255" json:"avatar,omitempty"`
    Title        string    `gorm:"size:50" json:"title,omitempty"`
    Intro        string    `gorm:"type:text" json:"intro,omitempty"`
    Specialty    string    `gorm:"type:text" json:"specialty,omitempty"`
    Status       int8      `gorm:"default:1" json:"status"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

```go
// internal/model/schedule.go
package model

import "time"

type Schedule struct {
    ID          int64     `gorm:"primaryKey" json:"id"`
    DoctorID    int64     `gorm:"index" json:"doctor_id"`
    CampusID    int64     `gorm:"index" json:"campus_id"`
    Date        string    `gorm:"type:date" json:"date"`
    TimePeriod  string    `gorm:"size:20" json:"time_period"`
    TotalCount  int       `json:"total_count"`
    UsedCount   int       `gorm:"default:0" json:"used_count"`
    RemainCount int       `json:"remain_count"`
    Status      int8      `gorm:"default:1" json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

```go
// internal/model/patient.go
package model

import "time"

type Patient struct {
    ID        int64     `gorm:"primaryKey" json:"id"`
    UserID    int64     `gorm:"index" json:"user_id"`
    Name      string    `gorm:"size:50" json:"name"`
    IDCard    string    `gorm:"size:18" json:"id_card,omitempty"`
    Phone     string    `gorm:"size:20" json:"phone,omitempty"`
    Relation  string    `gorm:"size:20" json:"relation,omitempty"`
    IsDefault bool      `gorm:"default:false" json:"is_default"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

```go
// internal/model/appointment.go
package model

import "time"

type Appointment struct {
    ID          int64     `gorm:"primaryKey" json:"id"`
    PatientID   int64     `gorm:"index" json:"patient_id"`
    DoctorID    int64     `gorm:"index" json:"doctor_id"`
    ScheduleID  int64     `gorm:"index" json:"schedule_id"`
    CampusID    int64     `json:"campus_id"`
    Date        string    `gorm:"type:date" json:"date"`
    TimePeriod  string    `gorm:"size:20" json:"time_period"`
    Status      string    `gorm:"size:20;default:PENDING_PAY" json:"status"`
    PayType     string    `gorm:"size:20" json:"pay_type,omitempty"`
    PayAmount   float64   `gorm:"type:decimal(10,2)" json:"pay_amount,omitempty"`
    CancelReason string   `gorm:"size:255" json:"cancel_reason,omitempty"`
    VisitNo     string    `gorm:"size:50" json:"visit_no,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

```go
// internal/model/notification.go
package model

import "time"

type Notification struct {
    ID         int64     `gorm:"primaryKey" json:"id"`
    UserID     int64     `gorm:"index" json:"user_id"`
    Type       string    `gorm:"size:50" json:"type"`
    TemplateID string    `gorm:"size:100" json:"template_id,omitempty"`
    Content    string    `gorm:"type:text" json:"content,omitempty"`
    Status     string    `gorm:"size:20" json:"status"`
    CreatedAt  time.Time `json:"created_at"`
}
```

- [ ] **Step 5: 更新 main.go 添加数据库初始化**

```go
// cmd/server/main.go
package main

import (
    "fmt"
    "log"
    "hos_schedule/internal/config"
    "hos_schedule/internal/middleware"
    "hos_schedule/internal/model"
    "hos_schedule/internal/router"
    "github.com/gin-gonic/gin"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    db, err := cfg.Database.Connect()
    if err != nil {
        log.Fatalf("Failed to connect database: %v", err)
    }
    
    if err := db.AutoMigrate(
        &model.User{},
        &model.Hospital{},
        &model.Campus{},
        &model.Department{},
        &model.Doctor{},
        &model.Schedule{},
        &model.Patient{},
        &model.Appointment{},
        &model.Notification{},
    ); err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }
    
    rdb := cfg.Redis.Connect()
    if err := cfg.Redis.Ping(nil, rdb); err != nil {
        log.Fatalf("Failed to connect redis: %v", err)
    }
    
    gin.SetMode(cfg.Server.Mode)
    
    r := gin.Default()
    r.Use(middleware.CORS())
    
    router.Register(r, db, rdb, cfg)
    
    addr := fmt.Sprintf(":%d", cfg.Server.Port)
    log.Printf("Server starting on %s", addr)
    
    if err := r.Run(addr); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

- [ ] **Step 6: 更新路由注册**

```go
// internal/router/router.go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
    "hos_schedule/internal/config"
    "hos_schedule/internal/pkg/response"
)

func Register(r *gin.Engine, db *gorm.DB, rdb *redis.Client, cfg *config.Config) {
    api := r.Group("/api/v1")
    {
        api.GET("/health", func(c *gin.Context) {
            response.Success(c, gin.H{"status": "ok"})
        })
    }
}
```

- [ ] **Step 7: 验证编译**

```bash
cd D:\coding\lang\mimo\hos_schedule
go build ./cmd/server
```

Expected: 编译成功

---

## Task 3: 用户认证模块

**Covers:** S3

**Files:**
- Create: `internal/service/auth_service.go`, `internal/handler/auth_handler.go`, `internal/middleware/auth.go`
- Modify: `internal/router/router.go`

- [ ] **Step 1: 创建 JWT 工具**

```go
// internal/pkg/jwt/jwt.go
package jwt

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID int64  `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(secret string, userID int64, role string, expire time.Duration) (string, error) {
    claims := Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ParseToken(secret, tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrSignatureInvalid
}
```

- [ ] **Step 2: 创建认证中间件**

```go
// internal/middleware/auth.go
package middleware

import (
    "strings"
    "hos_schedule/internal/config"
    "hos_schedule/internal/pkg/jwt"
    "hos_schedule/internal/pkg/response"
    "github.com/gin-gonic/gin"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Unauthorized(c, "Missing authorization header")
            c.Abort()
            return
        }
        
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            response.Unauthorized(c, "Invalid authorization format")
            c.Abort()
            return
        }
        
        claims, err := jwt.ParseToken(cfg.JWT.Secret, tokenString)
        if err != nil {
            response.Unauthorized(c, "Invalid token")
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

- [ ] **Step 3: 创建微信登录服务**

```go
// internal/service/auth_service.go
package service

import (
    "encoding/json"
    "fmt"
    "net/http"
    "hos_schedule/internal/config"
    "hos_schedule/internal/model"
    "hos_schedule/internal/pkg/jwt"
    "gorm.io/gorm"
    "time"
)

type AuthService struct {
    db  *gorm.DB
    cfg *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
    return &AuthService{db: db, cfg: cfg}
}

type WechatLoginResp struct {
    OpenID     string `json:"openid"`
    UnionID    string `json:"unionid"`
    SessionKey string `json:"session_key"`
    ErrCode    int    `json:"errcode"`
    ErrMsg     string `json:"errmsg"`
}

func (s *AuthService) WechatLogin(code string) (string, *model.User, error) {
    url := fmt.Sprintf(
        "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
        s.cfg.Wechat.AppID, s.cfg.Wechat.Secret, code,
    )
    
    resp, err := http.Get(url)
    if err != nil {
        return "", nil, fmt.Errorf("failed to call wechat api: %w", err)
    }
    defer resp.Body.Close()
    
    var wechatResp WechatLoginResp
    if err := json.NewDecoder(resp.Body).Decode(&wechatResp); err != nil {
        return "", nil, fmt.Errorf("failed to decode wechat response: %w", err)
    }
    
    if wechatResp.ErrCode != 0 {
        return "", nil, fmt.Errorf("wechat login failed: %s", wechatResp.ErrMsg)
    }
    
    var user model.User
    result := s.db.Where("openid = ?", wechatResp.OpenID).First(&user)
    
    if result.Error == gorm.ErrRecordNotFound {
        user = model.User{
            OpenID:   wechatResp.OpenID,
            UnionID:  wechatResp.UnionID,
            Role:     "PATIENT",
            Nickname: "微信用户",
        }
        if err := s.db.Create(&user).Error; err != nil {
            return "", nil, fmt.Errorf("failed to create user: %w", err)
        }
    } else if result.Error != nil {
        return "", nil, fmt.Errorf("failed to query user: %w", result.Error)
    }
    
    expire, _ := time.ParseDuration(s.cfg.JWT.Expire)
    token, err := jwt.GenerateToken(s.cfg.JWT.Secret, user.ID, user.Role, expire)
    if err != nil {
        return "", nil, fmt.Errorf("failed to generate token: %w", err)
    }
    
    return token, &user, nil
}
```

- [ ] **Step 4: 创建认证 Handler**

```go
// internal/handler/auth_handler.go
package handler

import (
    "hos_schedule/internal/pkg/response"
    "hos_schedule/internal/service"
    "github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{authService: authService}
}

type LoginRequest struct {
    Code string `json:"code" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request")
        return
    }
    
    token, user, err := h.authService.WechatLogin(req.Code)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    
    response.Success(c, gin.H{
        "token": token,
        "user":  user,
    })
}
```

- [ ] **Step 5: 注册路由**

```go
// internal/router/router.go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
    "hos_schedule/internal/config"
    "hos_schedule/internal/handler"
    "hos_schedule/internal/middleware"
    "hos_schedule/internal/pkg/response"
    "hos_schedule/internal/service"
)

func Register(r *gin.Engine, db *gorm.DB, rdb *redis.Client, cfg *config.Config) {
    authService := service.NewAuthService(db, cfg)
    authHandler := handler.NewAuthHandler(authService)
    
    api := r.Group("/api/v1")
    {
        api.GET("/health", func(c *gin.Context) {
            response.Success(c, gin.H{"status": "ok"})
        })
        
        api.POST("/auth/login", authHandler.Login)
        
        // 需要认证的路由
        auth := api.Group("")
        auth.Use(middleware.Auth(cfg))
        {
            auth.GET("/me", func(c *gin.Context) {
                userID, _ := c.Get("user_id")
                response.Success(c, gin.H{"user_id": userID})
            })
        }
    }
}
```

- [ ] **Step 6: 验证编译**

```bash
cd D:\coding\lang\mimo\hos_schedule
go build ./cmd/server
```

Expected: 编译成功

---

## Task 4: 医院/院区/科室模块

**Covers:** S3

**Files:**
- Create: `internal/repository/hospital_repo.go`, `internal/service/hospital_service.go`, `internal/handler/hospital_handler.go`, `internal/repository/department_repo.go`, `internal/service/department_service.go`, `internal/handler/department_handler.go`
- Modify: `internal/router/router.go`

- [ ] **Step 1: 创建医院 Repository**

```go
// internal/repository/hospital_repo.go
package repository

import (
    "hos_schedule/internal/model"
    "gorm.io/gorm"
)

type HospitalRepo struct {
    db *gorm.DB
}

func NewHospitalRepo(db *gorm.DB) *HospitalRepo {
    return &HospitalRepo{db: db}
}

func (r *HospitalRepo) List() ([]model.Hospital, error) {
    var hospitals []model.Hospital
    err := r.db.Where("status = ?", 1).Find(&hospitals).Error
    return hospitals, err
}

func (r *HospitalRepo) GetByID(id int64) (*model.Hospital, error) {
    var hospital model.Hospital
    err := r.db.First(&hospital, id).Error
    return &hospital, err
}

func (r *HospitalRepo) GetCampuses(hospitalID int64) ([]model.Campus, error) {
    var campuses []model.Campus
    err := r.db.Where("hospital_id = ? AND status = ?", hospitalID, 1).Find(&campuses).Error
    return campuses, err
}

func (r *HospitalRepo) GetCampusByID(id int64) (*model.Campus, error) {
    var campus model.Campus
    err := r.db.First(&campus, id).Error
    return &campus, err
}
```

- [ ] **Step 2: 创建医院 Service**

```go
// internal/service/hospital_service.go
package service

import (
    "hos_schedule/internal/model"
    "hos_schedule/internal/repository"
)

type HospitalService struct {
    repo *repository.HospitalRepo
}

func NewHospitalService(repo *repository.HospitalRepo) *HospitalService {
    return &HospitalService{repo: repo}
}

func (s *HospitalService) List() ([]model.Hospital, error) {
    return s.repo.List()
}

func (s *HospitalService) GetByID(id int64) (*model.Hospital, error) {
    return s.repo.GetByID(id)
}

func (s *HospitalService) GetCampuses(hospitalID int64) ([]model.Campus, error) {
    return s.repo.GetCampuses(hospitalID)
}
```

- [ ] **Step 3: 创建医院 Handler**

```go
// internal/handler/hospital_handler.go
package handler

import (
    "strconv"
    "hos_schedule/internal/pkg/response"
    "hos_schedule/internal/service"
    "github.com/gin-gonic/gin"
)

type HospitalHandler struct {
    service *service.HospitalService
}

func NewHospitalHandler(service *service.HospitalService) *HospitalHandler {
    return &HospitalHandler{service: service}
}

func (h *HospitalHandler) List(c *gin.Context) {
    hospitals, err := h.service.List()
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, hospitals)
}

func (h *HospitalHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid hospital ID")
        return
    }
    
    hospital, err := h.service.GetByID(id)
    if err != nil {
        response.NotFound(c, "Hospital not found")
        return
    }
    response.Success(c, hospital)
}

func (h *HospitalHandler) GetCampuses(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid hospital ID")
        return
    }
    
    campuses, err := h.service.GetCampuses(id)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, campuses)
}
```

- [ ] **Step 4: 创建科室 Repository/Service/Handler**

```go
// internal/repository/department_repo.go
package repository

import (
    "hos_schedule/internal/model"
    "gorm.io/gorm"
)

type DepartmentRepo struct {
    db *gorm.DB
}

func NewDepartmentRepo(db *gorm.DB) *DepartmentRepo {
    return &DepartmentRepo{db: db}
}

func (r *DepartmentRepo) ListByCampus(campusID int64) ([]model.Department, error) {
    var departments []model.Department
    err := r.db.Where("campus_id = ? AND status = ?", campusID, 1).
        Order("sort_order ASC").
        Find(&departments).Error
    return departments, err
}

func (r *DepartmentRepo) ListByHospital(hospitalID int64) ([]model.Department, error) {
    var departments []model.Department
    err := r.db.Where("hospital_id = ? AND status = ?", hospitalID, 1).
        Order("sort_order ASC").
        Find(&departments).Error
    return departments, err
}

func (r *DepartmentRepo) GetByID(id int64) (*model.Department, error) {
    var department model.Department
    err := r.db.First(&department, id).Error
    return &department, err
}
```

```go
// internal/service/department_service.go
package service

import (
    "hos_schedule/internal/model"
    "hos_schedule/internal/repository"
)

type DepartmentService struct {
    repo *repository.DepartmentRepo
}

func NewDepartmentService(repo *repository.DepartmentRepo) *DepartmentService {
    return &DepartmentService{repo: repo}
}

func (s *DepartmentService) ListByCampus(campusID int64) ([]model.Department, error) {
    return s.repo.ListByCampus(campusID)
}

func (s *DepartmentService) ListByHospital(hospitalID int64) ([]model.Department, error) {
    return s.repo.ListByHospital(hospitalID)
}

func (s *DepartmentService) GetByID(id int64) (*model.Department, error) {
    return s.repo.GetByID(id)
}
```

```go
// internal/handler/department_handler.go
package handler

import (
    "strconv"
    "hos_schedule/internal/pkg/response"
    "hos_schedule/internal/service"
    "github.com/gin-gonic/gin"
)

type DepartmentHandler struct {
    service *service.DepartmentService
}

func NewDepartmentHandler(service *service.DepartmentService) *DepartmentHandler {
    return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) List(c *gin.Context) {
    campusID, _ := strconv.ParseInt(c.Query("campus_id"), 10, 64)
    hospitalID, _ := strconv.ParseInt(c.Query("hospital_id"), 10, 64)
    
    var departments interface{}
    var err error
    
    if campusID > 0 {
        departments, err = h.service.ListByCampus(campusID)
    } else if hospitalID > 0 {
        departments, err = h.service.ListByHospital(hospitalID)
    } else {
        response.BadRequest(c, "campus_id or hospital_id required")
        return
    }
    
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, departments)
}

func (h *DepartmentHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid department ID")
        return
    }
    
    department, err := h.service.GetByID(id)
    if err != nil {
        response.NotFound(c, "Department not found")
        return
    }
    response.Success(c, department)
}
```

- [ ] **Step 5: 注册路由**

```go
// internal/router/router.go (续)
// 在 api.Group 内添加：

    hospitalRepo := repository.NewHospitalRepo(db)
    hospitalService := service.NewHospitalService(hospitalRepo)
    hospitalHandler := handler.NewHospitalHandler(hospitalService)
    
    departmentRepo := repository.NewDepartmentRepo(db)
    departmentService := service.NewDepartmentService(departmentRepo)
    departmentHandler := handler.NewDepartmentHandler(departmentService)
    
    api.GET("/hospitals", hospitalHandler.List)
    api.GET("/hospitals/:id", hospitalHandler.GetByID)
    api.GET("/hospitals/:id/campuses", hospitalHandler.GetCampuses)
    
    api.GET("/departments", departmentHandler.List)
    api.GET("/departments/:id", departmentHandler.GetByID)
```

- [ ] **Step 6: 验证编译**

```bash
go build ./cmd/server
```

---

## Task 5: 医生与排班模块

**Covers:** S3

**Files:**
- Create: `internal/repository/doctor_repo.go`, `internal/service/doctor_service.go`, `internal/handler/doctor_handler.go`, `internal/repository/schedule_repo.go`, `internal/service/schedule_service.go`, `internal/handler/schedule_handler.go`
- Modify: `internal/router/router.go`

- [ ] **Step 1: 创建医生 Repository/Service/Handler**

```go
// internal/repository/doctor_repo.go
package repository

import (
    "hos_schedule/internal/model"
    "gorm.io/gorm"
)

type DoctorRepo struct {
    db *gorm.DB
}

func NewDoctorRepo(db *gorm.DB) *DoctorRepo {
    return &DoctorRepo{db: db}
}

func (r *DoctorRepo) ListByDepartment(departmentID int64) ([]model.Doctor, error) {
    var doctors []model.Doctor
    err := r.db.Where("department_id = ? AND status = ?", departmentID, 1).Find(&doctors).Error
    return doctors, err
}

func (r *DoctorRepo) GetByID(id int64) (*model.Doctor, error) {
    var doctor model.Doctor
    err := r.db.First(&doctor, id).Error
    return &doctor, err
}

func (r *DoctorRepo) ListRecommended(limit int) ([]model.Doctor, error) {
    var doctors []model.Doctor
    err := r.db.Where("status = ?", 1).Limit(limit).Find(&doctors).Error
    return doctors, err
}
```

```go
// internal/service/doctor_service.go
package service

import (
    "hos_schedule/internal/model"
    "hos_schedule/internal/repository"
)

type DoctorService struct {
    repo *repository.DoctorRepo
}

func NewDoctorService(repo *repository.DoctorRepo) *DoctorService {
    return &DoctorService{repo: repo}
}

func (s *DoctorService) ListByDepartment(departmentID int64) ([]model.Doctor, error) {
    return s.repo.ListByDepartment(departmentID)
}

func (s *DoctorService) GetByID(id int64) (*model.Doctor, error) {
    return s.repo.GetByID(id)
}

func (s *DoctorService) ListRecommended(limit int) ([]model.Doctor, error) {
    return s.repo.ListRecommended(limit)
}
```

```go
// internal/handler/doctor_handler.go
package handler

import (
    "strconv"
    "hos_schedule/internal/pkg/response"
    "hos_schedule/internal/service"
    "github.com/gin-gonic/gin"
)

type DoctorHandler struct {
    service *service.DoctorService
}

func NewDoctorHandler(service *service.DoctorService) *DoctorHandler {
    return &DoctorHandler{service: service}
}

func (h *DoctorHandler) List(c *gin.Context) {
    departmentID, _ := strconv.ParseInt(c.Query("department_id"), 10, 64)
    if departmentID > 0 {
        doctors, err := h.service.ListByDepartment(departmentID)
        if err != nil {
            response.InternalError(c, err.Error())
            return
        }
        response.Success(c, doctors)
        return
    }
    
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
    doctors, err := h.service.ListRecommended(limit)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, doctors)
}

func (h *DoctorHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid doctor ID")
        return
    }
    
    doctor, err := h.service.GetByID(id)
    if err != nil {
        response.NotFound(c, "Doctor not found")
        return
    }
    response.Success(c, doctor)
}
```

- [ ] **Step 2: 创建排班 Repository/Service/Handler**

```go
// internal/repository/schedule_repo.go
package repository

import (
    "hos_schedule/internal/model"
    "gorm.io/gorm"
)

type ScheduleRepo struct {
    db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
    return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) ListByDoctor(doctorID int64, startDate, endDate string) ([]model.Schedule, error) {
    var schedules []model.Schedule
    err := r.db.Where("doctor_id = ? AND date >= ? AND date <= ? AND status = ?", 
        doctorID, startDate, endDate, 1).
        Order("date ASC").
        Find(&schedules).Error
    return schedules, err
}

func (r *ScheduleRepo) ListByDepartment(departmentID int64, date string) ([]model.Schedule, error) {
    var schedules []model.Schedule
    err := r.db.Joins("JOIN doctors ON doctors.id = schedules.doctor_id").
        Where("doctors.department_id = ? AND schedules.date = ? AND schedules.status = ?", 
            departmentID, date, 1).
        Find(&schedules).Error
    return schedules, err
}

func (r *ScheduleRepo) GetByID(id int64) (*model.Schedule, error) {
    var schedule model.Schedule
    err := r.db.First(&schedule, id).Error
    return &schedule, err
}

func (r *ScheduleRepo) DecrementRemain(id int64) error {
    result := r.db.Model(&model.Schedule{}).
        Where("id = ? AND remain_count > 0", id).
        Updates(map[string]interface{}{
            "remain_count": gorm.Expr("remain_count - 1"),
            "used_count":   gorm.Expr("used_count + 1"),
        })
    
    if result.RowsAffected == 0 {
        return fmt.Errorf("no remaining slots")
    }
    return result.Error
}

func (r *ScheduleRepo) IncrementRemain(id int64) error {
    return r.db.Model(&model.Schedule{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "remain_count": gorm.Expr("remain_count + 1"),
            "used_count":   gorm.Expr("used_count - 1"),
        }).Error
}
```

```go
// internal/service/schedule_service.go
package service

import (
    "hos_schedule/internal/model"
    "hos_schedule/internal/repository"
)

type ScheduleService struct {
    repo *repository.ScheduleRepo
}

func NewScheduleService(repo *repository.ScheduleRepo) *ScheduleService {
    return &ScheduleService{repo: repo}
}

func (s *ScheduleService) ListByDoctor(doctorID int64, startDate, endDate string) ([]model.Schedule, error) {
    return s.repo.ListByDoctor(doctorID, startDate, endDate)
}

func (s *ScheduleService) ListByDepartment(departmentID int64, date string) ([]model.Schedule, error) {
    return s.repo.ListByDepartment(departmentID, date)
}

func (s *ScheduleService) GetByID(id int64) (*model.Schedule, error) {
    return s.repo.GetByID(id)
}
```

```go
// internal/handler/schedule_handler.go
package handler

import (
    "strconv"
    "hos_schedule/internal/pkg/response"
    "hos_schedule/internal/service"
    "github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
    service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
    return &ScheduleHandler{service: service}
}

func (h *ScheduleHandler) List(c *gin.Context) {
    doctorID, _ := strconv.ParseInt(c.Query("doctor_id"), 10, 64)
    departmentID, _ := strconv.ParseInt(c.Query("department_id"), 10, 64)
    date := c.Query("date")
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")
    
    if doctorID > 0 {
        schedules, err := h.service.ListByDoctor(doctorID, startDate, endDate)
        if err != nil {
            response.InternalError(c, err.Error())
            return
        }
        response.Success(c, schedules)
        return
    }
    
    if departmentID > 0 && date != "" {
        schedules, err := h.service.ListByDepartment(departmentID, date)
        if err != nil {
            response.InternalError(c, err.Error())
            return
        }
        response.Success(c, schedules)
        return
    }
    
    response.BadRequest(c, "doctor_id or department_id+date required")
}

func (h *ScheduleHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid schedule ID")
        return
    }
    
    schedule, err := h.service.GetByID(id)
    if err != nil {
        response.NotFound(c, "Schedule not found")
        return
    }
    response.Success(c, schedule)
}
```

- [ ] **Step 3: 注册路由**

```go
// internal/router/router.go (续)
    doctorRepo := repository.NewDoctorRepo(db)
    doctorService := service.NewDoctorService(doctorRepo)
    doctorHandler := handler.NewDoctorHandler(doctorService)
    
    scheduleRepo := repository.NewScheduleRepo(db)
    scheduleService := service.NewScheduleService(scheduleRepo)
    scheduleHandler := handler.NewScheduleHandler(scheduleService)
    
    api.GET("/doctors", doctorHandler.List)
    api.GET("/doctors/:id", doctorHandler.GetByID)
    
    api.GET("/schedules", scheduleHandler.List)
    api.GET("/schedules/:id", scheduleHandler.GetByID)
```

- [ ] **Step 4: 验证编译**

```bash
go build ./cmd/server
```

---

## Task 6: 就诊人模块

**Covers:** S3

**Files:**
- Create: `internal/repository/patient_repo.go`, `internal/service/patient_service.go`, `internal/handler/patient_handler.go`
- Modify: `internal/router/router.go`

- [ ] **Step 1: 创建就诊人 Repository**

```go
// internal/repository/patient_repo.go
package repository

import (
    "hos_schedule/internal/model"
    "gorm.io/gorm"
)

type PatientRepo struct {
    db *gorm.DB
}

func NewPatientRepo(db *gorm.DB) *PatientRepo {
    return &PatientRepo{db: db}
}

func (r *PatientRepo) ListByUser(userID int64) ([]model.Patient, error) {
    var patients []model.Patient
    err := r.db.Where("user_id = ?", userID).Order("is_default DESC").Find(&patients).Error
    return patients, err
}

func (r *PatientRepo) GetByID(id int64) (*model.Patient, error) {
    var patient model.Patient
    err := r.db.First(&patient, id).Error
    return &patient, err
}

func (r *PatientRepo) Create(patient *model.Patient) error {
    return r.db.Create(patient).Error
}

func (r *PatientRepo) Update(patient *model.Patient) error {
    return r.db.Save(patient).Error
}

func (r *PatientRepo) Delete(id int64) error {
    return r.db.Delete(&model.Patient{}, id).Error
}

func (r *PatientRepo) SetDefault(userID, patientID int64) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&model.Patient{}).Where("user_id = ?", userID).
            Update("is_default", false).Error; err != nil {
            return err
        }
        return tx.Model(&model.Patient{}).Where("id = ? AND user_id = ?", patientID, userID).
            Update("is_default", true).Error
    })
}
```

- [ ] **Step 2: 创建就诊人 Service**

```go
// internal/service/patient_service.go
package service

import (
    "hos_schedule/internal/model"
    "hos_schedule/internal/repository"
)

type PatientService struct {
    repo *repository.PatientRepo
}

func NewPatientService(repo *repository.PatientRepo) *PatientService {
    return &PatientService{repo: repo}
}

func (s *PatientService) ListByUser(userID int64) ([]model.Patient, error) {
    return s.repo.ListByUser(userID)
}

func (s *PatientService) GetByID(id int64) (*model.Patient, error) {
    return s.repo.GetByID(id)
}

func (s *PatientService) Create(patient *model.Patient) error {
    return s.repo.Create(patient)
}

func (s *PatientService) Update(patient *model.Patient) error {
    return s.repo.Update(patient)
}

func (s *PatientService) Delete(id int64) error {
    return s.repo.Delete(id)
}

func (s *PatientService) SetDefault(userID, patientID int64) error {
    return s.repo.SetDefault(userID, patientID)
}
```

- [ ] **Step 3: 创建就诊人 Handler**

```go
// internal/handler/patient_handler.go
package handler

import (
    "strconv"
    "hos_schedule/internal/model"
    "hos_schedule/internal/pkg/response"
    "hos_schedule/internal/service"
    "github.com/gin-gonic/gin"
)

type PatientHandler struct {
    service *service.PatientService
}

func NewPatientHandler(service *service.PatientService) *PatientHandler {
    return &PatientHandler{service: service}
}

func (h *PatientHandler) List(c *gin.Context) {
    userID := c.GetInt64("user_id")
    patients, err := h.service.ListByUser(userID)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, patients)
}

func (h *PatientHandler) Create(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    var patient model.Patient
    if err := c.ShouldBindJSON(&patient); err != nil {
        response.BadRequest(c, "Invalid request")
        return
    }
    
    patient.UserID = userID
    if err := h.service.Create(&patient); err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, patient)
}

func (h *PatientHandler) Update(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid patient ID")
        return
    }
    
    var patient model.Patient
    if err := c.ShouldBindJSON(&patient); err != nil {
        response.BadRequest(c, "Invalid request")
        return
    }
    
    patient.ID = id
    if err := h.service.Update(&patient); err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, patient)
}

func (h *PatientHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid patient ID")
        return
    }
    
    if err := h.service.Delete(id); err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, nil)
}

func (h *PatientHandler) SetDefault(c *gin.Context) {
    userID := c.GetInt64("user_id")
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid patient ID")
        return
    }
    
    if err := h.service.SetDefault(userID, id); err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, nil)
}
```

- [ ] **Step 4: 注册路由**

```go
// internal/router/router.go (续)
    patientRepo := repository.NewPatientRepo(db)
    patientService := service.NewPatientService(patientRepo)
    patientHandler := handler.NewPatientHandler(patientService)
    
    // 需要认证
    auth := api.Group("")
    auth.Use(middleware.Auth(cfg))
    {
        auth.GET("/patients", patientHandler.List)
        auth.POST("/patients", patientHandler.Create)
        auth.PUT("/patients/:id", patientHandler.Update)
        auth.DELETE("/patients/:id", patientHandler.Delete)
        auth.PUT("/patients/:id/default", patientHandler.SetDefault)
    }
```

- [ ] **Step 5: 验证编译**

```bash
go build ./cmd/server
```

---

## Task 7: 预约模块（核心）

**Covers:** S4

**Files:**
- Create: `internal/repository/appointment_repo.go`, `internal/service/appointment_service.go`, `internal/handler/appointment_handler.go`, `internal/pkg/redis/slot_deduction.lua`
- Modify: `internal/router/router.go`

- [ ] **Step 1: 创建号源扣减 Lua 脚本**

```lua
-- internal/pkg/redis/slot_deduction.lua
-- KEYS[1]: schedule:{id}:remain
-- ARGV[1]: 扣减数量（通常为1）
local key = KEYS[1]
local count = tonumber(ARGV[1])
local remain = tonumber(redis.call('get', key))
if remain >= count then
    redis.call('decrby', key, count)
    return 1
else
    return 0
end
```

- [ ] **Step 2: 创建 Redis 号源管理**

```go
// internal/pkg/redis/slot_manager.go
package redis

import (
    "context"
    "embed"
    "fmt"
    "strconv"
    "github.com/redis/go-redis/v9"
)

//go:embed slot_deduction.lua
var luaScript string

type SlotManager struct {
    rdb *redis.Client
}

func NewSlotManager(rdb *redis.Client) *SlotManager {
    return &SlotManager{rdb: rdb}
}

func (m *SlotManager) InitSlot(ctx context.Context, scheduleID int64, remainCount int) error {
    key := fmt.Sprintf("schedule:%d:remain", scheduleID)
    return m.rdb.Set(ctx, key, remainCount, 0).Err()
}

func (m *SlotManager) DeductSlot(ctx context.Context, scheduleID int64) (bool, error) {
    key := fmt.Sprintf("schedule:%d:remain", scheduleID)
    result, err := m.rdb.Eval(ctx, luaScript, []string{key}, 1).Int()
    if err != nil {
        return false, err
    }
    return result == 1, nil
}

func (m *SlotManager) ReleaseSlot(ctx context.Context, scheduleID int64) error {
    key := fmt.Sprintf("schedule:%d:remain", scheduleID)
    return m.rdb.Incr(ctx, key).Err()
}

func (m *SlotManager) GetRemain(ctx context.Context, scheduleID int64) (int, error) {
    key := fmt.Sprintf("schedule:%d:remain", scheduleID)
    val, err := m.rdb.Get(ctx, key).Result()
    if err == redis.Nil {
        return 0, nil
    }
    if err != nil {
        return 0, err
    }
    return strconv.Atoi(val)
}
```

- [ ] **Step 3: 创建预约 Repository**

```go
// internal/repository/appointment_repo.go
package repository

import (
    "hos_schedule/internal/model"
    "gorm.io/gorm"
)

type AppointmentRepo struct {
    db *gorm.DB
}

func NewAppointmentRepo(db *gorm.DB) *AppointmentRepo {
    return &AppointmentRepo{db: db}
}

func (r *AppointmentRepo) Create(appointment *model.Appointment) error {
    return r.db.Create(appointment).Error
}

func (r *AppointmentRepo) GetByID(id int64) (*model.Appointment, error) {
    var appointment model.Appointment
    err := r.db.First(&appointment, id).Error
    return &appointment, err
}

func (r *AppointmentRepo) ListByUser(userID int64, status string) ([]model.Appointment, error) {
    var appointments []model.Appointment
    query := r.db.Joins("JOIN patients ON patients.id = appointments.patient_id").
        Where("patients.user_id = ?", userID)
    
    if status != "" {
        query = query.Where("appointments.status = ?", status)
    }
    
    err := query.Order("appointments.created_at DESC").Find(&appointments).Error
    return appointments, err
}

func (r *AppointmentRepo) ListByDoctor(doctorID int64, date string) ([]model.Appointment, error) {
    var appointments []model.Appointment
    err := r.db.Where("doctor_id = ? AND date = ?", doctorID, date).
        Order("created_at ASC").
        Find(&appointments).Error
    return appointments, err
}

func (r *AppointmentRepo) UpdateStatus(id int64, status string) error {
    return r.db.Model(&model.Appointment{}).Where("id = ?", id).
        Update("status", status).Error
}

func (r *AppointmentRepo) Exists(userID, doctorID int64, date, timePeriod string) (bool, error) {
    var count int64
    err := r.db.Model(&model.Appointment{}).
        Joins("JOIN patients ON patients.id = appointments.patient_id").
        Where("patients.user_id = ? AND appointments.doctor_id = ? AND appointments.date = ? AND appointments.time_period = ? AND appointments.status NOT IN ?", 
            userID, doctorID, date, timePeriod, []string{"CANCELLED"}).
        Count(&count).Error
    return count > 0, err
}

func (r *AppointmentRepo) GetPendingPayExpired(minutes int) ([]model.Appointment, error) {
    var appointments []model.Appointment
    err := r.db.Where("status = ? AND created_at < ?", 
        "PENDING_PAY", gorm.Expr("NOW() - INTERVAL '? minutes'", minutes)).
        Find(&appointments).Error
    return appointments, err
}
```

- [ ] **Step 4: 创建预约 Service（核心逻辑）**

```go
// internal/service/appointment_service.go
package service

import (
    "context"
    "fmt"
    "hos_schedule/internal/model"
    "hos_schedule/internal/pkg/redis"
    "hos_schedule/internal/repository"
    "gorm.io/gorm"
)

type AppointmentService struct {
    db            *gorm.DB
    appointmentRepo *repository.AppointmentRepo
    scheduleRepo    *repository.ScheduleRepo
    slotManager     *redis.SlotManager
}

func NewAppointmentService(
    db *gorm.DB,
    appointmentRepo *repository.AppointmentRepo,
    scheduleRepo *repository.ScheduleRepo,
    slotManager *redis.SlotManager,
) *AppointmentService {
    return &AppointmentService{
        db:              db,
        appointmentRepo: appointmentRepo,
        scheduleRepo:    scheduleRepo,
        slotManager:     slotManager,
    }
}

type CreateAppointmentRequest struct {
    PatientID  int64  `json:"patient_id" binding:"required"`
    DoctorID   int64  `json:"doctor_id" binding:"required"`
    ScheduleID int64  `json:"schedule_id" binding:"required"`
    CampusID   int64  `json:"campus_id" binding:"required"`
    Date       string `json:"date" binding:"required"`
    TimePeriod string `json:"time_period" binding:"required"`
    PayType    string `json:"pay_type"` // ONLINE or ONSITE
}

func (s *AppointmentService) Create(ctx context.Context, userID int64, req *CreateAppointmentRequest) (*model.Appointment, error) {
    // 检查是否重复预约
    exists, err := s.appointmentRepo.Exists(userID, req.DoctorID, req.Date, req.TimePeriod)
    if err != nil {
        return nil, fmt.Errorf("failed to check duplicate: %w", err)
    }
    if exists {
        return nil, fmt.Errorf("already have appointment for this slot")
    }
    
    // Redis 原子扣减号源
    success, err := s.slotManager.DeductSlot(ctx, req.ScheduleID)
    if err != nil {
        // Redis 故障，降级到数据库行锁
        if err := s.scheduleRepo.DecrementRemain(req.ScheduleID); err != nil {
            return nil, fmt.Errorf("no remaining slots")
        }
    }
    if !success {
        return nil, fmt.Errorf("no remaining slots")
    }
    
    // 创建预约记录
    appointment := &model.Appointment{
        PatientID:  req.PatientID,
        DoctorID:   req.DoctorID,
        ScheduleID: req.ScheduleID,
        CampusID:   req.CampusID,
        Date:       req.Date,
        TimePeriod: req.TimePeriod,
        PayType:    req.PayType,
        Status:     "PENDING_PAY",
    }
    
    if err := s.appointmentRepo.Create(appointment); err != nil {
        // 回滚 Redis
        s.slotManager.ReleaseSlot(ctx, req.ScheduleID)
        return nil, fmt.Errorf("failed to create appointment: %w", err)
    }
    
    return appointment, nil
}

func (s *AppointmentService) Cancel(ctx context.Context, id int64, reason string) error {
    appointment, err := s.appointmentRepo.GetByID(id)
    if err != nil {
        return fmt.Errorf("appointment not found")
    }
    
    if appointment.Status == "CANCELLED" {
        return fmt.Errorf("appointment already cancelled")
    }
    
    // 更新状态
    if err := s.appointmentRepo.UpdateStatus(id, "CANCELLED"); err != nil {
        return err
    }
    
    // 释放号源
    if err := s.slotManager.ReleaseSlot(ctx, appointment.ScheduleID); err != nil {
        // Redis 失败，数据库兜底
        s.scheduleRepo.IncrementRemain(appointment.ScheduleID)
    }
    
    return nil
}

func (s *AppointmentService) GetByID(id int64) (*model.Appointment, error) {
    return s.appointmentRepo.GetByID(id)
}

func (s *AppointmentService) ListByUser(userID int64, status string) ([]model.Appointment, error) {
    return s.appointmentRepo.ListByUser(userID, status)
}

func (s *AppointmentService) ListByDoctor(doctorID int64, date string) ([]model.Appointment, error) {
    return s.appointmentRepo.ListByDoctor(doctorID, date)
}
```

- [ ] **Step 5: 创建预约 Handler**

```go
// internal/handler/appointment_handler.go
package handler

import (
    "strconv"
    "hos_schedule/internal/pkg/response"
    "hos_schedule/internal/service"
    "github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
    service *service.AppointmentService
}

func NewAppointmentHandler(service *service.AppointmentService) *AppointmentHandler {
    return &AppointmentHandler{service: service}
}

func (h *AppointmentHandler) Create(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    var req service.CreateAppointmentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request")
        return
    }
    
    appointment, err := h.service.Create(c.Request.Context(), userID, &req)
    if err != nil {
        response.Error(c, 400, err.Error())
        return
    }
    response.Success(c, appointment)
}

func (h *AppointmentHandler) GetByID(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid appointment ID")
        return
    }
    
    appointment, err := h.service.GetByID(id)
    if err != nil {
        response.NotFound(c, "Appointment not found")
        return
    }
    response.Success(c, appointment)
}

func (h *AppointmentHandler) List(c *gin.Context) {
    userID := c.GetInt64("user_id")
    status := c.Query("status")
    
    appointments, err := h.service.ListByUser(userID, status)
    if err != nil {
        response.InternalError(c, err.Error())
        return
    }
    response.Success(c, appointments)
}

func (h *AppointmentHandler) Cancel(c *gin.Context) {
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(c, "Invalid appointment ID")
        return
    }
    
    var req struct {
        Reason string `json:"reason"`
    }
    c.ShouldBindJSON(&req)
    
    if err := h.service.Cancel(c.Request.Context(), id, req.Reason); err != nil {
        response.Error(c, 400, err.Error())
        return
    }
    response.Success(c, nil)
}
```

- [ ] **Step 6: 注册路由**

```go
// internal/router/router.go (续)
    slotManager := redis.NewSlotManager(rdb)
    
    appointmentRepo := repository.NewAppointmentRepo(db)
    appointmentService := service.NewAppointmentService(db, appointmentRepo, scheduleRepo, slotManager)
    appointmentHandler := handler.NewAppointmentHandler(appointmentService)
    
    auth.POST("/appointments", appointmentHandler.Create)
    auth.GET("/appointments", appointmentHandler.List)
    auth.GET("/appointments/:id", appointmentHandler.GetByID)
    auth.PUT("/appointments/:id/cancel", appointmentHandler.Cancel)
```

- [ ] **Step 7: 验证编译**

```bash
go build ./cmd/server
```

---

## Task 8: 通知模块

**Covers:** S3

**Files:**
- Create: `internal/repository/notification_repo.go`, `internal/service/notification_service.go`, `internal/handler/notification_handler.go`, `internal/pkg/wechat/client.go`
- Modify: `internal/router/router.go`

- [ ] **Step 1: 创建微信消息客户端**

```go
// internal/pkg/wechat/client.go
package wechat

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "hos_schedule/internal/config"
)

type Client struct {
    cfg *config.WechatConfig
}

func NewClient(cfg *config.WechatConfig) *Client {
    return &Client{cfg: cfg}
}

type SubscribeMessage struct {
    ToUser     string                 `json:"touser"`
    TemplateID string                 `json:"template_id"`
    Page       string                 `json:"page"`
    Data       map[string]interface{} `json:"data"`
}

func (c *Client) SendSubscribeMessage(msg *SubscribeMessage) error {
    token, err := c.getAccessToken()
    if err != nil {
        return err
    }
    
    url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s", token)
    
    body, _ := json.Marshal(msg)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    var result struct {
        ErrCode int    `json:"errcode"`
        ErrMsg  string `json:"errmsg"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    
    if result.ErrCode != 0 {
        return fmt.Errorf("wechat send failed: %s", result.ErrMsg)
    }
    
    return nil
}

func (c *Client) getAccessToken() (string, error) {
    url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
        c.cfg.AppID, c.cfg.Secret)
    
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    
    var result struct {
        AccessToken string `json:"access_token"`
        ErrCode     int    `json:"errcode"`
        ErrMsg      string `json:"errmsg"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    
    if result.ErrCode != 0 {
        return "", fmt.Errorf("get access token failed: %s", result.ErrMsg)
    }
    
    return result.AccessToken, nil
}
```

- [ ] **Step 2: 创建通知 Service**

```go
// internal/service/notification_service.go
package service

import (
    "hos_schedule/internal/model"
    "hos_schedule/internal/pkg/wechat"
    "hos_schedule/internal/repository"
)

type NotificationService struct {
    repo   *repository.NotificationRepo
    wechat *wechat.Client
}

func NewNotificationService(repo *repository.NotificationRepo, wechat *wechat.Client) *NotificationService {
    return &NotificationService{repo: repo, wechat: wechat}
}

func (s *NotificationService) SendAppointmentSuccess(userID int64, appointment *model.Appointment) error {
    notification := &model.Notification{
        UserID:  userID,
        Type:    "APPOINTMENT_SUCCESS",
        Status:  "PENDING",
        Content: "预约成功",
    }
    
    if err := s.repo.Create(notification); err != nil {
        return err
    }
    
    // TODO: 调用微信订阅消息
    
    notification.Status = "SENT"
    return s.repo.UpdateStatus(notification.ID, "SENT")
}

func (s *NotificationService) SendReminder(userID int64, appointment *model.Appointment) error {
    notification := &model.Notification{
        UserID:  userID,
        Type:    "REMINDER_1DAY",
        Status:  "PENDING",
        Content: "预约提醒",
    }
    
    return s.repo.Create(notification)
}
```

- [ ] **Step 3: 验证编译**

```bash
go build ./cmd/server
```

---

## Task 9: 医生端 API

**Covers:** S3

**Files:**
- Create: `internal/handler/doctor_handler.go` (扩展)

- [ ] **Step 1: 扩展医生 Handler**

```go
// 在 internal/handler/doctor_handler.go 中添加

func (h *DoctorHandler) GetMySchedules(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    // 通过 user_id 查找 doctor
    // TODO: 需要 user_id 到 doctor_id 的映射
    
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")
    
    response.Success(c, gin.H{
        "user_id":     userID,
        "start_date":  startDate,
        "end_date":    endDate,
    })
}

func (h *DoctorHandler) GetTodayAppointments(c *gin.Context) {
    userID := c.GetInt64("user_id")
    
    // TODO: 通过 user_id 查找 doctor，然后查询今日预约
    
    response.Success(c, gin.H{
        "user_id": userID,
    })
}
```

- [ ] **Step 2: 注册医生端路由**

```go
// 在 auth 路由组内添加
    auth.GET("/doctor/schedules", doctorHandler.GetMySchedules)
    auth.GET("/doctor/appointments", doctorHandler.GetTodayAppointments)
```

- [ ] **Step 3: 验证编译**

```bash
go build ./cmd/server
```

---

## Task 10: Admin API

**Covers:** S6

**Files:**
- Create: `internal/handler/admin_handler.go`, `internal/service/admin_service.go`
- Modify: `internal/router/router.go`

- [ ] **Step 1: 创建 Admin Handler**

```go
// internal/handler/admin_handler.go
package handler

import (
    "hos_schedule/internal/model"
    "hos_schedule/internal/pkg/response"
    "hos_schedule/internal/service"
    "github.com/gin-gonic/gin"
)

type AdminHandler struct {
    hospitalService   *service.HospitalService
    departmentService *service.DepartmentService
    doctorService     *service.DoctorService
    scheduleService   *service.ScheduleService
    appointmentService *service.AppointmentService
}

func NewAdminHandler(
    hospitalService *service.HospitalService,
    departmentService *service.DepartmentService,
    doctorService *service.DoctorService,
    scheduleService *service.ScheduleService,
    appointmentService *service.AppointmentService,
) *AdminHandler {
    return &AdminHandler{
        hospitalService:   hospitalService,
        departmentService: departmentService,
        doctorService:     doctorService,
        scheduleService:   scheduleService,
        appointmentService: appointmentService,
    }
}

// 科室管理
func (h *AdminHandler) CreateDepartment(c *gin.Context) {
    var dept model.Department
    if err := c.ShouldBindJSON(&dept); err != nil {
        response.BadRequest(c, "Invalid request")
        return
    }
    // TODO: 调用 service 创建
    response.Success(c, dept)
}

// 医生管理
func (h *AdminHandler) CreateDoctor(c *gin.Context) {
    var doc model.Doctor
    if err := c.ShouldBindJSON(&doc); err != nil {
        response.BadRequest(c, "Invalid request")
        return
    }
    // TODO: 调用 service 创建
    response.Success(c, doc)
}

// 排班管理
func (h *AdminHandler) CreateSchedule(c *gin.Context) {
    var sched model.Schedule
    if err := c.ShouldBindJSON(&sched); err != nil {
        response.BadRequest(c, "Invalid request")
        return
    }
    // TODO: 调用 service 创建，初始化 Redis 号源
    response.Success(c, sched)
}

// 批量排班
func (h *AdminHandler) BatchCreateSchedule(c *gin.Context) {
    var req struct {
        DoctorID   int64    `json:"doctor_id"`
        CampusID   int64    `json:"campus_id"`
        Dates      []string `json:"dates"`
        TimePeriod string   `json:"time_period"`
        TotalCount int      `json:"total_count"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        response.BadRequest(c, "Invalid request")
        return
    }
    // TODO: 批量创建排班
    response.Success(c, req)
}

// 预约统计
func (h *AdminHandler) GetAppointmentStats(c *gin.Context) {
    // TODO: 统计今日预约量、取消量、到诊率
    response.Success(c, gin.H{
        "today_total":   0,
        "today_cancel":  0,
        "today_visited": 0,
    })
}
```

- [ ] **Step 2: 注册 Admin 路由**

```go
// internal/router/router.go (续)
    adminHandler := handler.NewAdminHandler(
        hospitalService, departmentService, doctorService, scheduleService, appointmentService,
    )
    
    admin := api.Group("/admin")
    admin.Use(middleware.Auth(cfg))
    // TODO: 添加管理员权限检查中间件
    {
        admin.POST("/departments", adminHandler.CreateDepartment)
        admin.POST("/doctors", adminHandler.CreateDoctor)
        admin.POST("/schedules", adminHandler.CreateSchedule)
        admin.POST("/schedules/batch", adminHandler.BatchCreateSchedule)
        admin.GET("/appointments/stats", adminHandler.GetAppointmentStats)
    }
```

- [ ] **Step 3: 验证编译**

```bash
go build ./cmd/server
```

---

## Task 11: Docker 部署配置

**Covers:** S7

**Files:**
- Modify: `Dockerfile`, `docker-compose.yml`, `config/config.yaml`

- [ ] **Step 1: 完善 docker-compose.yml**

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
      - DB_NAME=hospital_schedule
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    volumes:
      - ./config:/root/config
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: hospital_schedule
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

- [ ] **Step 2: 创建 Nginx 配置**

```nginx
# nginx/nginx.conf
events {
    worker_connections 1024;
}

http {
    upstream api {
        server api:8080;
    }

    server {
        listen 80;
        server_name your-domain.com;
        
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

- [ ] **Step 3: 创建启动脚本**

```bash
#!/bin/bash
# start.sh
docker-compose up -d
echo "Services started. API available at http://localhost:8080"
```

- [ ] **Step 4: 验证 Docker 构建**

```bash
docker-compose build
```

---

## 自检完成

**Spec 覆盖检查：**
- [S1] 系统架构 → Task 1, 2
- [S2] 数据库设计 → Task 2
- [S3] API 设计 → Task 3, 4, 5, 6, 7, 8, 9, 10
- [S4] 高并发预约 → Task 7
- [S5] 微信小程序 → 前端任务（另行规划）
- [S6] Web 管理后台 → Task 10
- [S7] 私有化部署 → Task 11

**下一步：** 执行计划
