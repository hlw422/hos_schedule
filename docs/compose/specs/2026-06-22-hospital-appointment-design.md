# 三甲医院预约挂号平台设计文档

**日期：** 2026-06-22
**版本：** v1.0

---

## [S1] 系统架构设计

### 技术栈
- **后端：** Go + Gin
- **数据库：** PostgreSQL 16
- **缓存：** Redis 7
- **前端：** 微信小程序（患者端 + 医生端 + 管理端）
- **管理后台：** Vue 3 + Element Plus（或 React + Ant Design）
- **部署：** 私有化部署，Docker 容器化

### 架构图

```
┌─────────────────────────────────────────────────────────┐
│                    微信小程序客户端                        │
│         (患者端 + 医生端 + 管理端小程序)                    │
└─────────────────────────┬───────────────────────────────┘
                          │ HTTPS
                          ▼
┌─────────────────────────────────────────────────────────┐
│                    Nginx 反向代理                         │
│              (SSL终止 + 负载均衡)                          │
└─────────────────────────┬───────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                   Go + Gin API 服务                       │
│     ┌──────────┐  ┌──────────┐  ┌──────────┐            │
│     │ 患者模块  │  │ 预约模块  │  │ 排班模块  │            │
│     └──────────┘  └──────────┘  └──────────┘            │
│     ┌──────────┐  ┌──────────┐  ┌──────────┐            │
│     │ 医院模块  │  │ 医生模块  │  │ 用户模块  │            │
│     └──────────┘  └──────────┘  └──────────┘            │
└─────────────────────────┬───────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          ▼               ▼               ▼
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  PostgreSQL  │  │    Redis     │  │  消息服务     │
│   主数据库    │  │   缓存/锁    │  │ (微信+短信)  │
└──────────────┘  └──────────────┘  └──────────────┘
```

### 设计决策
1. **单体应用，模块化结构** — 初期不做微服务，用 Go module 分包，未来可拆分
2. **Redis 缓存号源** — 号源库存放 Redis，Lua 脚本原子扣减，防止超卖
3. **数据库行锁兜底** — Redis 故障时降级到 PostgreSQL 行锁
4. **微信消息模板** — 预约成功、提醒、取消通知

---

## [S2] 数据库设计

### 核心表结构

```sql
-- 医院
CREATE TABLE hospitals (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(20),
    logo VARCHAR(255),
    intro TEXT,
    status SMALLINT DEFAULT 1,  -- 1:启用 0:停用
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 院区
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

-- 科室
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

-- 医生
CREATE TABLE doctors (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT NOT NULL REFERENCES departments(id),
    name VARCHAR(50) NOT NULL,
    avatar VARCHAR(255),
    title VARCHAR(50),  -- 职称：主任医师/副主任医师/主治医师/住院医师
    intro TEXT,
    specialty TEXT,  -- 擅长领域
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 排班
CREATE TABLE schedules (
    id BIGSERIAL PRIMARY KEY,
    doctor_id BIGINT NOT NULL REFERENCES doctors(id),
    campus_id BIGINT NOT NULL REFERENCES campuses(id),
    date DATE NOT NULL,
    time_period VARCHAR(20) NOT NULL,  -- MORNING/AFTERNOON/EVENING
    total_count INT NOT NULL,  -- 总号源
    used_count INT DEFAULT 0,  -- 已用号源
    remain_count INT NOT NULL,  -- 剩余号源
    status SMALLINT DEFAULT 1,  -- 1:正常 0:停诊
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(doctor_id, date, time_period)
);

-- 就诊人
CREATE TABLE patients (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(50) NOT NULL,
    id_card VARCHAR(18),  -- 身份证号
    phone VARCHAR(20),
    relation VARCHAR(20),  -- 本人/父母/子女/配偶/其他
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 用户
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    openid VARCHAR(100) UNIQUE NOT NULL,
    unionid VARCHAR(100),
    phone VARCHAR(20),
    nickname VARCHAR(50),
    avatar VARCHAR(255),
    role VARCHAR(20) DEFAULT 'PATIENT',  -- PATIENT/DOCTOR/SCHEDULER/HOSPITAL_ADMIN/SUPER_ADMIN
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 预约
CREATE TABLE appointments (
    id BIGSERIAL PRIMARY KEY,
    patient_id BIGINT NOT NULL REFERENCES patients(id),
    doctor_id BIGINT NOT NULL REFERENCES doctors(id),
    schedule_id BIGINT NOT NULL REFERENCES schedules(id),
    campus_id BIGINT NOT NULL REFERENCES campuses(id),
    date DATE NOT NULL,
    time_period VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'PENDING_PAY',  -- PENDING_PAY/PAID/CANCELLED/COMPLETED/NO_SHOW
    pay_type VARCHAR(20),  -- ONLINE/ONSITE
    pay_amount DECIMAL(10, 2),
    cancel_reason VARCHAR(255),
    visit_no VARCHAR(50),  -- 就诊号
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 通知记录
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    type VARCHAR(50),  -- APPOINTMENT_SUCCESS/REMINDER_1DAY/REMINDER_2HOUR/CANCEL
    template_id VARCHAR(100),
    content TEXT,
    status VARCHAR(20),  -- PENDING/SENT/FAILED
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 关键索引
- `schedules`: `(doctor_id, date, time_period)` 联合索引
- `appointments`: `(user_id, status)`, `(doctor_id, date)`
- `patients`: `(user_id)`

---

## [S3] API 设计

### 患者端 API

```
/api/v1/
├── auth/
│   ├── POST /login          # 微信登录
│   └── POST /refresh        # 刷新token
│
├── hospitals/
│   ├── GET /                # 医院列表
│   ├── GET /:id             # 医院详情
│   └── GET /:id/campuses    # 院区列表
│
├── departments/
│   ├── GET /                # 科室列表
│   └── GET /:id             # 科室详情
│
├── doctors/
│   ├── GET /                # 医生列表
│   ├── GET /:id             # 医生详情
│   └── GET /:id/schedules   # 医生排班
│
├── schedules/
│   ├── GET /                # 排班列表（按日期/科室）
│   └── GET /:id/slots       # 号源详情
│
├── appointments/
│   ├── POST /               # 创建预约
│   ├── GET /                # 我的预约列表
│   ├── GET /:id             # 预约详情
│   ├── PUT /:id/cancel      # 取消预约
│   └── PUT /:id/confirm     # 确认就诊
│
├── patients/
│   ├── GET /                # 就诊人列表
│   ├── POST /               # 添加就诊人
│   ├── PUT /:id             # 编辑就诊人
│   └── DELETE /:id          # 删除就诊人
│
└── notifications/
    └── POST /subscribe      # 订阅消息授权
```

### 医生端 API

```
/api/v1/doctor/
├── GET /schedules           # 我的排班
├── GET /appointments        # 今日预约
├── POST /leave              # 申请停诊
└── GET /stats               # 统计数据
```

### 管理端 API

```
/api/v1/admin/
├── dashboard/
│   └── GET /stats           # 数据看板
│
├── hospitals/
│   ├── PUT /:id             # 更新医院信息
│   └── POST /campuses       # 添加院区
│
├── departments/
│   ├── POST /               # 添加科室
│   ├── PUT /:id             # 编辑科室
│   └── DELETE /:id          # 删除科室
│
├── doctors/
│   ├── POST /               # 添加医生
│   ├── PUT /:id             # 编辑医生
│   └── PUT /:id/status      # 启用/停用医生
│
├── schedules/
│   ├── POST /               # 创建排班
│   ├── POST /batch          # 批量排班
│   ├── PUT /:id             # 编辑排班
│   └── DELETE /:id          # 删除排班
│
└── appointments/
    ├── GET /                # 预约列表
    └── GET /stats           # 预约统计
```

---

## [S4] 高并发预约设计

### 号源库存模型

Redis 存储号源库存：
- Key: `schedule:{schedule_id}:remain`
- Value: 剩余号源数量

### Lua 脚本原子扣减

```lua
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

### 预约流程（防超卖）

```
1. 用户选择号源
   ↓
2. Redis Lua 脚本原子扣减
   ↓ 成功
3. 创建预约记录（状态：待支付）
   ↓
4. 发起支付（如选择在线支付）
   ↓
5. 支付成功 → 更新状态为已支付
   ↓
6. 超时未支付 → 释放号源（Redis + DB）
```

### 降级策略
- **Redis 故障：** 降级到 PostgreSQL 行锁
  ```sql
  UPDATE schedules SET remain_count = remain_count - 1 
  WHERE id = ? AND remain_count > 0
  ```
- **限流：** 单用户 10 次/分钟预约请求
- **幂等：** 用户 + 医生 + 日期 + 时段 唯一约束

### 超时释放
- 定时任务每分钟扫描超时预约（15分钟未支付）
- 释放号源：Redis incr + DB update

---

## [S5] 微信小程序设计

### 页面结构

```
小程序页面
├── 首页 (pages/index)
│   ├── 搜索框
│   ├── 轮播图（公告）
│   ├── 快捷入口
│   ├── 热门科室
│   └── 推荐专家
│
├── 医院列表 (pages/hospitals/list)
├── 医院详情 (pages/hospitals/detail)
│   ├── 院区切换
│   ├── 科室列表
│   └── 医生列表
│
├── 科室详情 (pages/departments/detail)
│   ├── 排班日历
│   └── 医生列表
│
├── 医生详情 (pages/doctors/detail)
│   ├── 医生信息
│   ├── 排班日历
│   └── 立即预约按钮
│
├── 预约确认 (pages/appointment/confirm)
│   ├── 选择就诊人
│   ├── 预约信息确认
│   └── 提交预约
│
├── 预约成功 (pages/appointment/success)
│
├── 我的预约 (pages/appointment/list)
│   ├── 待就诊
│   ├── 已完成
│   └── 已取消
│
├── 就诊人管理 (pages/patients/list)
│   ├── 添加就诊人
│   └── 编辑就诊人
│
└── 我的 (pages/profile)
    ├── 个人信息
    ├── 就诊人管理
    └── 设置
```

### UI 设计规范
- **主色：** #1677FF（医疗蓝）
- **辅助色：** #52C41A（成功绿）
- **背景：** #F7F8FA
- **设计风格：** 微信官方设计规范，Apple Health 医疗科技感

---

## [S6] Web 管理后台设计

### 模块结构

```
Web 管理后台
├── 数据看板
│   ├── 今日预约量/取消量/到诊率
│   ├── 本月预约趋势
│   └── 热门科室/医生
│
├── 医院管理
│   ├── 医院信息
│   └── 院区管理
│
├── 科室管理
│   ├── 科室列表
│   └── 科室编辑
│
├── 医生管理
│   ├── 医生列表
│   ├── 医生详情
│   └── 医生停用/启用
│
├── 排班管理
│   ├── 排班日历视图
│   ├── 单次排班
│   ├── 周期排班（批量）
│   └── 号源管理
│
├── 预约管理
│   ├── 预约列表
│   ├── 预约详情
│   └── 预约统计
│
└── 系统管理
    ├── 用户管理
    ├── 角色管理
    └── 操作日志
```

### 技术方案
- **前端：** Vue 3 + Element Plus（或 React + Ant Design）
- **部署：** 与后端同服务器，Nginx 静态托管

---

## [S7] 私有化部署方案

### 服务器规划

```
服务器配置（最低）：
├── 应用服务器 × 2
│   ├── Go API 服务
│   ├── Nginx
│   └── Redis
│
├── 数据库服务器 × 1
│   ├── PostgreSQL 16
│   └── 定时备份
│
└── 文件服务器 × 1（可选）
    └── 图片/文档存储
```

### Docker 部署

```yaml
version: '3.8'
services:
  api:
    build: .
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
  
  postgres:
    image: postgres:16
    volumes:
      - pgdata:/var/lib/postgresql/data
    environment:
      - POSTGRES_DB=hospital
      - POSTGRES_USER=admin
      - POSTGRES_PASSWORD=xxx
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
  
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./frontend/dist:/usr/share/nginx/html

volumes:
  pgdata:
```

### 备份策略
- 数据库：每日全量备份 + 增量备份
- 保留最近 30 天备份
- 异地备份（可选）

---

## [S8] MVP 范围

### 第一阶段（MVP）
- ✅ 预约挂号核心流程
- ✅ 就诊人管理
- ✅ 预约提醒（微信订阅消息 + 短信）
- ✅ 医生端小程序

### 第二阶段
- ⏳ 在线支付（微信支付）
- ⏳ 医院导航

### 第三阶段
- ⏳ 数据分析看板
- ⏳ 智能推荐
