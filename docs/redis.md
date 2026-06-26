# Redis 使用文档

## 概述

Smileyan 博客系统使用 Redis 作为缓存和状态存储，用于验证码存储和发送频率限制。

## 配置

Redis 配置位于 `config.yaml` 文件中：

```yaml
redis:
  host: localhost      # Redis 地址
  port: 6379           # Redis 端口
  username: ""         # Redis 用户名（可选）
  password: ""         # Redis 密码
  db: 0                # Redis 数据库编号
```

可通过环境变量覆盖：

| 环境变量 | 说明 |
|---------|------|
| SMILEYAN_BACKEND_REDIS_HOST | Redis 服务器地址 |
| SMILEYAN_BACKEND_REDIS_PORT | Redis 端口 |
| SMILEYAN_BACKEND_REDIS_USERNAME | Redis 用户名 |
| SMILEYAN_BACKEND_REDIS_PASSWORD | Redis 密码 |

## 使用场景

### 1. 邮箱验证码存储

用于存储用户注册/登录时发送的邮箱验证码。

**Redis Key 格式**: `verification_code:{email}`

**数据结构**: String

**过期时间**: 10 分钟

**操作函数**:
- `SaveVerificationCode(email, code)` - 保存验证码
- `GetVerificationCode(email)` - 获取验证码
- `DeleteVerificationCode(email)` - 删除验证码

### 2. 发送频率限制

防止恶意刷邮件攻击，保护邮件发送配额。

#### 邮箱频率限制

**Redis Key 格式**: `rate_limit:email:{email}`

**数据结构**: String (值为 "1")

**过期时间**: 600 秒 (10 分钟)

**限制规则**: 同一邮箱 10 分钟内只能发送一次验证码

#### IP 频率限制

**Redis Key 格式**: `rate_limit:ip:{ip}`

**数据结构**: Integer

**过期时间**: 3600 秒 (1 小时)

**限制规则**: 同一 IP 每小时最多发送 100 次验证码

**操作函数**:
- `CheckRateLimit(email, ip)` - 检查是否超过限制
- `SetRateLimit(email, ip)` - 设置频率限制

## Redis Key 汇总

| Key 模式 | 类型 | 过期时间 | 说明 |
|---------|------|---------|------|
| `verification_code:{email}` | String | 10 分钟 | 邮箱验证码 |
| `rate_limit:email:{email}` | String | 10 分钟 | 邮箱发送频率限制 |
| `rate_limit:ip:{ip}` | Integer | 1 小时 | IP 发送频率限制 |

## 代码位置

- Redis 初始化: [config/config.go](../config/config.go)
- Redis 操作函数: [utils/email.go](../utils/email.go)
