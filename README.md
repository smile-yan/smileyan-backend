# Smileyan Backend

基于 Gin 框架开发的后端服务，提供博客系统的核心 API 功能。

## 技术栈

- **Web 框架**: Gin
- **数据库**: MySQL (GORM)
- **缓存**: Redis
- **搜索引擎**: Bleve
- **配置管理**: Viper

## 项目结构

```
.
├── config/          # 配置加载模块
├── controllers/    # 控制器层
├── middleware/      # 中间件
├── models/          # 数据模型
├── routes/          # 路由定义
├── services/        # 业务逻辑层
├── utils/           # 工具函数
├── main.go          # 入口文件
└── config.yaml      # 配置文件
```

## 快速开始

### 1. 克隆项目

```bash
git clone <your-repo-url>
cd smileyan/backend
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置环境变量

复制 `.env.example` 为 `.env` 并填写敏感配置：

```bash
cp .env.example .env
```

然后编辑 `.env` 文件，填入你的配置：

```env
# 数据库配置
SMILEYAN_BACKEND_DB_HOST=你的数据库地址
SMILEYAN_BACKEND_DB_USER=数据库用户名
SMILEYAN_BACKEND_DB_PASSWORD=数据库密码
SMILEYAN_BACKEND_DB_NAME=数据库名称

# Redis 配置
SMILEYAN_BACKEND_REDIS_HOST=Redis 地址
SMILEYAN_BACKEND_REDIS_USERNAME=Redis 用户名（可选）
SMILEYAN_BACKEND_REDIS_PASSWORD=Redis 密码

# 邮箱配置
SMILEYAN_BACKEND_EMAIL_PASSWORD=邮箱密码

# JWT 配置
SMILEYAN_BACKEND_JWT_SECRET=JWT 密钥
```

> **注意**: `.env` 文件包含敏感信息，已添加到 `.gitignore`，不会提交到代码仓库。

### 4. 运行项目

```bash
go run main.go
```

## 配置说明

项目使用 `config.yaml` 管理非敏感配置，敏感信息通过环境变量读取：

| 配置项 | 环境变量 | 说明 |
|--------|----------|------|
| 数据库地址 | `SMILEYAN_BACKEND_DB_HOST` | MySQL 服务器地址 |
| 数据库用户 | `SMILEYAN_BACKEND_DB_USER` | MySQL 用户名 |
| 数据库密码 | `SMILEYAN_BACKEND_DB_PASSWORD` | MySQL 密码 |
| 数据库名称 | `SMILEYAN_BACKEND_DB_NAME` | 数据库名 |
| Redis 地址 | `SMILEYAN_BACKEND_REDIS_HOST` | Redis 服务器地址 |
| Redis 用户名 | `SMILEYAN_BACKEND_REDIS_USERNAME` | Redis 用户名 |
| Redis 密码 | `SMILEYAN_BACKEND_REDIS_PASSWORD` | Redis 密码 |
| 邮箱密码 | `SMILEYAN_BACKEND_EMAIL_PASSWORD` | SMTP 邮箱密码 |
| JWT 密钥 | `SMILEYAN_BACKEND_JWT_SECRET` | JWT 签名密钥 |

## API 文档

服务启动后访问 `http://localhost:8080` 查看相关接口。

## 许可证

MIT