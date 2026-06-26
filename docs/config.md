# 配置文件说明

## 概述

本文档说明 `config.yaml` 配置文件中各配置项的含义。

## 配置项说明

### server - 服务器配置

```yaml
server:
  port: 8080      # 服务端口
  mode: debug     # 运行模式：debug/release
```

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| port | int | 8080 | HTTP 服务监听端口 |
| mode | string | debug | 运行模式：debug（调试）或 release（生产） |

---

### database - 数据库配置

```yaml
database:
  host: localhost       # 数据库地址
  port: 3306           # 端口
  user: root           # 用户名
  password: ""         # 密码（建议通过环境变量设置）
  dbname: smileyan     # 数据库名称
  charset: utf8mb4     # 字符编码
  max_idle_conns: 10   # 最大空闲连接数
  max_open_conns: 100  # 最大打开连接数
```

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| host | string | localhost | MySQL 服务器地址 |
| port | int | 3306 | MySQL 端口 |
| user | string | root | 数据库用户名 |
| password | string | - | 数据库密码（敏感信息，建议通过环境变量设置） |
| dbname | string | - | 数据库名称 |
| charset | string | utf8mb4 | 字符编码 |
| max_idle_conns | int | 10 | 最大空闲连接数 |
| max_open_conns | int | 100 | 最大打开连接数 |

**环境变量覆盖**：

| 环境变量 | 说明 |
|----------|------|
| SMILEYAN_BACKEND_DB_HOST | 数据库地址 |
| SMILEYAN_BACKEND_DB_USER | 数据库用户名 |
| SMILEYAN_BACKEND_DB_PASSWORD | 数据库密码 |
| SMILEYAN_BACKEND_DB_NAME | 数据库名称 |

---

### redis - Redis 配置

```yaml
redis:
  host: localhost    # Redis 地址
  port: 6379         # 端口
  username: ""       # 用户名（可选）
  password: ""       # 密码（建议通过环境变量设置）
  db: 0              # 数据库编号
```

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| host | string | localhost | Redis 服务器地址 |
| port | int | 6379 | Redis 端口 |
| username | string | - | Redis 用户名（部分 Redis 服务需要） |
| password | string | - | Redis 密码（敏感信息） |
| db | int | 0 | Redis 数据库编号 |

**环境变量覆盖**：

| 环境变量 | 说明 |
|----------|------|
| SMILEYAN_BACKEND_REDIS_HOST | Redis 地址 |
| SMILEYAN_BACKEND_REDIS_USERNAME | Redis 用户名 |
| SMILEYAN_BACKEND_REDIS_PASSWORD | Redis 密码 |

---

### email - 邮件配置

```yaml
email:
  host: smtp.exmail.qq.com  # SMTP 服务器地址
  port: 465                 # SMTP 端口
  username: root@smileyan.cn # 邮箱地址
  password: ""              # 邮箱密码（建议通过环境变量设置）
  use_ssl: true             # 是否使用 SSL
```

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| host | string | - | SMTP 服务器地址 |
| port | int | - | SMTP 端口 |
| username | string | - | 发送邮件的邮箱地址 |
| password | string | - | 邮箱密码或授权码（敏感信息） |
| use_ssl | bool | true | 是否使用 SSL 连接 |

**常用 SMTP 配置示例**：

- QQ 企业邮箱：`smtp.exmail.qq.com`（端口 465/587）
- 163 邮箱：`smtp.163.com`（端口 465）
- Gmail：`smtp.gmail.com`（端口 587）

**环境变量覆盖**：

| 环境变量 | 说明 |
|----------|------|
| SMILEYAN_BACKEND_EMAIL_PASSWORD | 邮箱密码 |

---

### jwt - JWT 配置

```yaml
jwt:
  secret: ""           # JWT 密钥（建议通过环境变量设置）
  expire_hours: 168    # Token 过期时间（小时），默认 7 天
```

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| secret | string | - | JWT 签名密钥（敏感信息，建议使用复杂的随机字符串） |
| expire_hours | int | 168 | Token 有效期（小时），168 = 7 天 |

**环境变量覆盖**：

| 环境变量 | 说明 |
|----------|------|
| SMILEYAN_BACKEND_JWT_SECRET | JWT 密钥 |

---

### admin_emails - 管理员邮箱列表

```yaml
admin_emails:
  - admin@example.com
  - editor@example.com
```

| 配置项 | 类型 | 说明 |
|--------|------|------|
| admin_emails | []string | 管理员邮箱列表 |

在列表中的邮箱首次登录时将自动获得管理员角色（admin），普通用户无法通过注册获得管理员权限。

---

### upload - 上传配置

```yaml
upload:
  path: ./uploads          # 上传文件存储根目录
  avatar_path: ./uploads/avatars  # 头像存储目录
  max_size: 5              # 最大上传文件大小（MB）
```

| 配置项 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| path | string | ./uploads | 上传文件存储根目录 |
| avatar_path | string | ./uploads/avatars | 头像存储目录 |
| max_size | int | 5 | 单个文件最大大小（MB） |

---

## 环境变量

所有敏感配置都支持通过环境变量覆盖，优先级：**环境变量 > .env 文件 > config.yaml**。

### 配置环境变量的方式

#### 方式1：使用 .env 文件（推荐）

在项目根目录创建 `.env` 文件：

```bash
# 复制模板
cp .env.example .env

# 编辑 .env 文件
vim .env
```

`.env` 文件内容示例：

```env
SMILEYAN_BACKEND_DB_HOST=rm-xxx.mysql.rds.aliyuncs.com
SMILEYAN_BACKEND_DB_USER=smileyan
SMILEYAN_BACKEND_DB_PASSWORD=your_password
SMILEYAN_BACKEND_DB_NAME=smileyan

SMILEYAN_BACKEND_REDIS_HOST=r-xxx.redis.rds.aliyuncs.com
SMILEYAN_BACKEND_REDIS_USERNAME=your_username
SMILEYAN_BACKEND_REDIS_PASSWORD=your_password

SMILEYAN_BACKEND_EMAIL_PASSWORD=your_email_password
SMILEYAN_BACKEND_JWT_SECRET=your_jwt_secret
```

#### 方式2：直接在终端设置

```bash
export SMILEYAN_BACKEND_DB_PASSWORD=your_password
export SMILEYAN_BACKEND_JWT_SECRET=your_jwt_secret
go run main.go
```

---

## 配置示例

### 开发环境

```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  user: root
  password: ""
  dbname: smileyan_dev
  charset: utf8mb4

redis:
  host: localhost
  port: 6379
  username: ""
  password: ""
  db: 0

email:
  host: smtp.example.com
  port: 465
  username: dev@example.com
  password: ""
  use_ssl: true

jwt:
  secret: ""
  expire_hours: 24

admin_emails:
  - dev@example.com

upload:
  path: ./uploads
  avatar_path: ./uploads/avatars
  max_size: 5
```

### 生产环境

```yaml
server:
  port: 8080
  mode: release

database:
  host: rm-xxx.mysql.rds.aliyuncs.com
  port: 3306
  user: your_user
  password: ""
  dbname: smileyan
  charset: utf8mb4
  max_idle_conns: 10
  max_open_conns: 100

redis:
  host: r-xxx.redis.rds.aliyuncs.com
  port: 6379
  username: ""
  password: ""
  db: 0

email:
  host: smtp.exmail.qq.com
  port: 465
  username: notify@yourdomain.com
  password: ""
  use_ssl: true

jwt:
  secret: ""
  expire_hours: 168

admin_emails:
  - admin@yourdomain.com
  - editor@yourdomain.com

upload:
  path: /var/www/smileyan/uploads
  avatar_path: /var/www/smileyan/uploads/avatars
  max_size: 10
```