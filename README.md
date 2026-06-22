# 三甲医院预约挂号平台

基于 Go + Gin + PostgreSQL + Redis 的高并发医院预约挂号系统。

## 项目简介

本项目专注于三甲医院预约挂号核心闭环，支持多院区、百万用户、高并发预约场景。采用微信小程序作为患者端入口，Web 管理后台 + 小程序管理端供医院管理人员使用。

### 核心特性

- **高并发预约**: Redis Lua 脚本原子扣减号源，数据库行锁降级，防止超卖
- **多院区支持**: 一套系统管理多个院区的科室、医生、排班
- **微信小程序登录**: 基于微信 openid 的用户认证
- **预约状态机**: 待支付 → 已支付 → 已就诊 → 已完成，支持取消和爽约
- **就诊人管理**: 支持为家人（父母、子女、配偶）预约
- **预约提醒**: 微信订阅消息 + 短信通知
- **RBAC 权限**: 患者、医生、排班管理员、医院管理员、超级管理员

### 不包含的功能

- EMR 电子病历
- HIS 收费系统
- 药房系统
- 医保系统

## 技术栈

| 组件 | 技术 |
|------|------|
| 后端框架 | Go 1.22+ / Gin |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| ORM | GORM |
| 配置管理 | Viper |
| 认证 | JWT (golang-jwt) |
| 容器化 | Docker / Docker Compose |
| 反向代理 | Nginx |

## 项目结构

```
hos_schedule/
├── cmd/server/main.go              # 启动入口
├── config/config.yaml              # 配置文件
├── internal/
│   ├── config/                     # 配置加载
│   │   ├── config.go               # Viper 配置结构
│   │   ├── database.go             # PostgreSQL 连接
│   │   └── redis.go                # Redis 连接
│   ├── handler/                    # HTTP 处理器
│   │   ├── admin_handler.go        # 管理端 API
│   │   ├── appointment_handler.go  # 预约 API
│   │   ├── auth_handler.go         # 认证 API
│   │   ├── department_handler.go   # 科室 API
│   │   ├── doctor_handler.go       # 医生 API
│   │   ├── hospital_handler.go     # 医院 API
│   │   ├── notification_handler.go # 通知 API
│   │   ├── patient_handler.go      # 就诊人 API
│   │   └── schedule_handler.go     # 排班 API
│   ├── middleware/                  # 中间件
│   │   ├── auth.go                 # JWT 认证
│   │   └── cors.go                 # CORS 跨域
│   ├── model/                      # 数据模型
│   │   ├── appointment.go
│   │   ├── department.go
│   │   ├── doctor.go
│   │   ├── hospital.go
│   │   ├── notification.go
│   │   ├── patient.go
│   │   ├── schedule.go
│   │   └── user.go
│   ├── pkg/                        # 公共包
│   │   ├── jwt/jwt.go              # JWT 工具
│   │   ├── redis/                  # Redis 工具
│   │   │   ├── slot_deduction.lua  # 号源扣减 Lua 脚本
│   │   │   └── slot_manager.go     # 号源管理器
│   │   ├── response/response.go    # 统一响应
│   │   └── wechat/client.go        # 微信 API 客户端
│   ├── repository/                 # 数据访问层
│   │   ├── appointment_repo.go
│   │   ├── department_repo.go
│   │   ├── doctor_repo.go
│   │   ├── hospital_repo.go
│   │   ├── notification_repo.go
│   │   ├── patient_repo.go
│   │   └── schedule_repo.go
│   ├── router/router.go            # 路由注册
│   └── service/                    # 业务逻辑层
│       ├── appointment_service.go
│       ├── auth_service.go
│       ├── department_service.go
│       ├── doctor_service.go
│       ├── hospital_service.go
│       ├── notification_service.go
│       ├── patient_service.go
│       └── schedule_service.go
├── migrations/                     # 数据库迁移
│   ├── 001_create_tables.up.sql
│   └── 001_create_tables.down.sql
├── nginx/nginx.conf                # Nginx 配置
├── docs/                           # 设计文档
│   ├── compose/specs/              # 设计规范
│   └── compose/plans/              # 实现计划
├── Dockerfile
├── docker-compose.yml
├── start.sh
├── go.mod
└── go.sum
```

## 快速开始

### 环境要求

- Go 1.22+
- PostgreSQL 16+
- Redis 7+

### 1. 克隆项目

```bash
git clone https://github.com/hlw422/hos_schedule.git
cd hos_schedule
```

### 2. 配置

编辑 `config/config.yaml`：

```yaml
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
  secret: your-jwt-secret-key  # 请修改为随机密钥
  expire: 24h

wechat:
  appid: your-appid      # 微信小程序 AppID
  secret: your-secret    # 微信小程序 Secret
```

### 3. 启动数据库

```bash
# 使用 Docker Compose 启动 PostgreSQL 和 Redis
docker-compose up -d postgres redis
```

或手动启动 PostgreSQL 和 Redis。

### 4. 运行迁移

```bash
# 使用 psql 执行迁移
psql -h localhost -U postgres -d hospital_schedule -f migrations/001_create_tables.up.sql
```

### 5. 启动服务

```bash
# 直接运行
go run cmd/server/main.go

# 或编译后运行
go build -o server cmd/server/main.go
./server
```

服务将在 `http://localhost:8080` 启动。

### 6. 验证

```bash
# 健康检查
curl http://localhost:8080/api/v1/health

# 预期响应
{"code":0,"message":"success","data":{"status":"ok"}}
```

## Docker 部署

### 一键启动

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f api
```

### 服务访问

| 服务 | 地址 |
|------|------|
| API 服务 | http://localhost:8080 |
| Nginx | http://localhost:80 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6379 |

## API 文档

### 患者端 API

#### 认证

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/auth/login` | 微信登录 | ❌ |
| GET | `/api/v1/me` | 获取当前用户 | ✅ |

#### 医院

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/hospitals` | 医院列表 | ❌ |
| GET | `/api/v1/hospitals/:id` | 医院详情 | ❌ |
| GET | `/api/v1/hospitals/:id/campuses` | 院区列表 | ❌ |

#### 科室

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/departments` | 科室列表 | ❌ |
| GET | `/api/v1/departments/:id` | 科室详情 | ❌ |

#### 医生

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/doctors` | 医生列表 | ❌ |
| GET | `/api/v1/doctors/:id` | 医生详情 | ❌ |

#### 排班

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/schedules` | 排班列表 | ❌ |
| GET | `/api/v1/schedules/:id` | 排班详情 | ❌ |

#### 预约

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/appointments` | 创建预约 | ✅ |
| GET | `/api/v1/appointments` | 我的预约 | ✅ |
| GET | `/api/v1/appointments/:id` | 预约详情 | ✅ |
| PUT | `/api/v1/appointments/:id/cancel` | 取消预约 | ✅ |

#### 就诊人

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/patients` | 就诊人列表 | ✅ |
| POST | `/api/v1/patients` | 添加就诊人 | ✅ |
| PUT | `/api/v1/patients/:id` | 编辑就诊人 | ✅ |
| DELETE | `/api/v1/patients/:id` | 删除就诊人 | ✅ |
| PUT | `/api/v1/patients/:id/default` | 设为默认 | ✅ |

#### 通知

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/v1/notifications/subscribe` | 订阅消息 | ✅ |

### 医生端 API

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/v1/doctor/schedules` | 我的排班 | ✅ |
| GET | `/api/v1/doctor/appointments` | 今日预约 | ✅ |

### 管理端 API

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| PUT | `/api/v1/admin/hospitals/:id` | 更新医院 | ✅ |
| POST | `/api/v1/admin/hospitals/campuses` | 添加院区 | ✅ |
| POST | `/api/v1/admin/departments` | 添加科室 | ✅ |
| PUT | `/api/v1/admin/departments/:id` | 编辑科室 | ✅ |
| DELETE | `/api/v1/admin/departments/:id` | 删除科室 | ✅ |
| POST | `/api/v1/admin/doctors` | 添加医生 | ✅ |
| PUT | `/api/v1/admin/doctors/:id` | 编辑医生 | ✅ |
| PUT | `/api/v1/admin/doctors/:id/status` | 启用/停用 | ✅ |
| POST | `/api/v1/admin/schedules` | 创建排班 | ✅ |
| POST | `/api/v1/admin/schedules/batch` | 批量排班 | ✅ |
| PUT | `/api/v1/admin/schedules/:id` | 编辑排班 | ✅ |
| DELETE | `/api/v1/admin/schedules/:id` | 删除排班 | ✅ |
| GET | `/api/v1/admin/appointments` | 预约列表 | ✅ |
| GET | `/api/v1/admin/appointments/stats` | 预约统计 | ✅ |

## 数据库设计

### 核心表

| 表名 | 说明 |
|------|------|
| users | 用户表 |
| hospitals | 医院表 |
| campuses | 院区表 |
| departments | 科室表 |
| doctors | 医生表 |
| schedules | 排班表 |
| patients | 就诊人表 |
| appointments | 预约表 |
| notifications | 通知表 |

### ER 关系

```
hospitals 1──N campuses
hospitals 1──N departments
campuses 1──N departments
departments 1──N doctors
doctors 1──N schedules
users 1──N patients
patients 1──N appointments
schedules 1──N appointments
```

## 高并发预约设计

### 号源扣减流程

```
用户选择号源
     ↓
Redis Lua 脚本原子扣减
     ↓ 成功
创建预约记录（状态：待支付）
     ↓
支付确认
     ↓
更新状态为已支付
```

### Redis Lua 脚本

```lua
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

### 降级策略

1. **Redis 故障**: 降级到 PostgreSQL 行锁
2. **限流**: 单用户 10 次/分钟预约请求
3. **幂等**: 用户 + 医生 + 日期 + 时段唯一约束
4. **超时释放**: 15 分钟未支付自动释放号源

## 预约状态机

```
PENDING_PAY (待支付)
     ↓
   PAID (已支付)
     ↓
  ┌──────────────┐
  ↓              ↓
COMPLETED    NO_SHOW
(已完成)     (爽约)

PAID → CANCELLED (取消)
```

## 配置说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| server.port | 服务端口 | 8080 |
| server.mode | 运行模式 (debug/release) | debug |
| database.host | 数据库地址 | localhost |
| database.port | 数据库端口 | 5432 |
| database.user | 数据库用户 | postgres |
| database.password | 数据库密码 | postgres |
| database.dbname | 数据库名称 | hospital_schedule |
| redis.host | Redis 地址 | localhost |
| redis.port | Redis 端口 | 6379 |
| jwt.secret | JWT 密钥 | - |
| jwt.expire | Token 过期时间 | 24h |
| wechat.appid | 微信小程序 AppID | - |
| wechat.secret | 微信小程序 Secret | - |

## 开发计划

### 第一阶段（MVP）✅

- [x] 预约挂号核心流程
- [x] 就诊人管理
- [x] 预约提醒（微信订阅消息）
- [x] 医生端 API
- [x] 管理端 API
- [x] Docker 部署

### 第二阶段

- [ ] 在线支付（微信支付）
- [ ] 医院导航
- [ ] 短信通知

### 第三阶段

- [ ] 数据分析看板
- [ ] 智能推荐
- [ ] 微信小程序前端
- [ ] Web 管理后台前端

## 设计文档

- [系统设计文档](docs/compose/specs/2026-06-22-hospital-appointment-design.md)
- [实现计划](docs/compose/plans/2026-06-22-hospital-appointment-plan.md)

## 许可证

MIT License
