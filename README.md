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

## 部署

推送 `v*.*.*` 形式的 tag 会自动触发 `.github/workflows/release.yml`：构建 `linux/amd64` 静态二进制、发布到 GitHub Release、并通过 SSH 部署到 `SERVERS_CONFIG` 列出的每一台服务器（并行）。`v*.*.*-rc` / `-alpha` / `-beta` 形式的 tag 只发 release 不部署，并自动标记为 GitHub prerelease。

### 必需的 GitHub 配置

**Secrets**（`Settings -> Secrets and variables -> Actions -> Secrets`）：

| Secret | 说明 |
|--------|------|
| `SMILEYAN_BACKEND_DB_HOST` / `_USER` / `_PASSWORD` / `_NAME` | MySQL 连接信息（**所有服务器共用**） |
| `SMILEYAN_BACKEND_REDIS_HOST` / `_USERNAME` / `_PASSWORD` | Redis 连接信息（**所有服务器共用**） |
| `SMILEYAN_BACKEND_EMAIL_PASSWORD` | SMTP 邮箱密码（**所有服务器共用**） |
| `SMILEYAN_BACKEND_JWT_SECRET` | JWT 签名密钥（**所有服务器共用**） |
| `DEPLOY_SSH_KEY` | 部署用 SSH 私钥（**所有服务器共用同一把**，需提前把对应公钥加入每台机器的 `authorized_keys`） |

**Repository Variable**（`Settings -> Secrets and variables -> Actions -> Variables`）：

| Variable | 说明 |
|----------|------|
| `SERVERS_CONFIG` | 部署目标服务器列表（JSON 数组，见下） |

> **为什么不是 Secret**：GitHub Actions 拒绝在 `strategy.matrix` 表达式里使用 `secrets.*`（matrix 值会写入 workflow run 日志，泄露 secret 风险）。Repository Variable 加密存储但允许出现在 matrix 中，地址/端口/用户/部署路径也不属于真正的敏感信息（真正的私钥仍在 `DEPLOY_SSH_KEY` 中）。

### `SERVERS_CONFIG`

JSON 数组，每项描述一台目标服务器。`name` 在数组内必须唯一；`port` 接受字符串或整数；其余字段为字符串。

```json
[
  {
    "name": "prod-1",
    "addr": "1.2.3.4",
    "port": "22",
    "user": "deploy",
    "deploy_path": "/opt/smileyan-backend"
  },
  {
    "name": "prod-2",
    "addr": "5.6.7.8",
    "port": "2222",
    "user": "ubuntu",
    "deploy_path": "/srv/smileyan-backend"
  }
]
```

| 字段 | 说明 |
|------|------|
| `name` | 仅用于日志标识 / Actions UI 区分多机 job，必须在数组内唯一 |
| `addr` | 主机名或 IP |
| `port` | SSH 端口 |
| `user` | SSH 登录用户名（须对 `deploy_path` 有写权限） |
| `deploy_path` | 服务器上的部署根目录，脚本会创建 `releases/<tag>/` 和 `shared/.env` |

添加 / 删除部署目标**只需编辑这一个 secret**，不需要改 workflow 文件。每次 tag 触发时所有机器**并行部署**（`strategy.fail-fast: false`），单台失败不会取消其他机器 —— 想看完整结果再介入时这个行为很关键。

> 启动 `Diagnose secret availability` 步骤会预先校验 `SERVERS_CONFIG` 的结构（数组非空、每项字段齐全、端口为数字、名字不重复），配置写错会立即以具体错误信息失败，不会让 `fromJSON` 抛一个模糊的报错。

### 已废弃的 secrets

以下旧版 workflow 用过的 secret 不再被读取，可以从仓库设置中删除（保留也完全无害）：

- `SSH_INSTANCE_ADDR`、`SSH_INSTANCE_PORT`、`SSH_INSTANCE_USER`
- `SERVER_DEPLOY_PATH`

## 许可证

MIT