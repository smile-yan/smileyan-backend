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

## 部署环境变量（CI/CD）

`release.yml` 的 `deploy` job 引用 **10 个 Secret** + **1 个 Variable**：

> **关于 Secret vs Variable 的取舍**：  
> - **Secret** 适合**真正的凭据**（密码、私钥、token），运行时不会出现在 workflow 日志里（runner 会 mask）。  
> - **Variable** 适合**非敏感的部署拓扑**（服务器地址、端口、路径、用户）。明文存储，但**可以**用在 `strategy.matrix` 等 contexts 里（Secret 不行）。  
>   详见 GitHub 官方 [Encrypted secrets](https://docs.github.com/zh/actions/security-for-github-actions/security-guides/using-secrets-in-github-actions) 和 [Variables](https://docs.github.com/zh/actions/learn-github-actions/variables) 文档。

### Secret（10 个，`Settings → Secrets and variables → Actions → Secrets` tab）

应用层 9 个（被 Go 二进制 `config.go` 读，并经 `deploy-remote.sh` 转发到服务器的 `shared/.env`）：

| 名称 | 用途 |
|------|------|
| `SMILEYAN_BACKEND_DB_HOST` | MySQL 主机 / IP |
| `SMILEYAN_BACKEND_DB_USER` | MySQL 用户名 |
| `SMILEYAN_BACKEND_DB_PASSWORD` | MySQL 密码 |
| `SMILEYAN_BACKEND_DB_NAME` | MySQL 数据库名 |
| `SMILEYAN_BACKEND_REDIS_HOST` | Redis 主机 / IP |
| `SMILEYAN_BACKEND_REDIS_USERNAME` | Redis 用户名（未启用 ACL 时填占位如 `default`） |
| `SMILEYAN_BACKEND_REDIS_PASSWORD` | Redis 密码 |
| `SMILEYAN_BACKEND_EMAIL_PASSWORD` | SMTP 邮箱密码 |
| `SMILEYAN_BACKEND_JWT_SECRET` | JWT 签名密钥 |

部署层 1 个（被 `release.yml` 的 `webfactory/ssh-agent` 用，**所有目标服务器共用同一把私钥**）：

| 名称 | 用途 |
|------|------|
| `DEPLOY_SSH_KEY` | 部署用 SSH 私钥（需提前把对应公钥加入每台目标机器的 `authorized_keys`） |

> **关于 `SMILEYAN_BACKEND_REDIS_USERNAME`**：CLAUDE.md 注明该变量在本地运行时是 optional（Redis 未启用 ACL 时可空），但当前 `release.yml` 的 `Check required env vars` 步骤会把它当作必填。**线上部署时请给它一个占位值（如 `default`），否则第一步就会被拦下来。**

### Variable（1 个，`Settings → Secrets and variables → Actions → Variables` tab）

| 名称 | 用途 |
|------|------|
| `SERVERS_CONFIG` | 部署目标服务器列表（JSON 数组，见下） |

#### `SERVERS_CONFIG` 格式

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

| 字段 | 必填 | 说明 |
|------|------|------|
| `name` | ✓ | 仅用于日志标识 / Actions UI 区分多机 job，必须在数组内唯一 |
| `addr` | ✓ | 主机名或 IP |
| `port` | ✓ | SSH 端口（字符串或整数都行，校验时只看是否数字） |
| `user` | ✓ | SSH 登录用户名（须对 `deploy_path` 有写权限） |
| `deploy_path` | ✓ | 服务器上的部署根目录，脚本会创建 `releases/<tag>/` 和 `shared/.env` |

> **为什么 `SERVERS_CONFIG` 不是 Secret**：GitHub Actions 拒绝在 `strategy.matrix` 表达式里使用 `secrets.*`（matrix 值会写入 workflow run 日志）。Repository Variable 加密存储、允许出现在 matrix 中 —— 部署目标本来也不算敏感信息（真正的私钥在 `DEPLOY_SSH_KEY` 这个 Secret 里）。

**添加 / 删除部署目标只需编辑 `SERVERS_CONFIG` 这一个 Variable**，不需要改 `release.yml`。每次打 tag 触发后所有机器**并行部署**（`strategy.fail-fast: false`），单台失败不会取消其他机器。

如果 `SERVERS_CONFIG` 未配置或为空数组，deploy job 会被静默 skip（不报错）；同时 `release` job 的 `Print SERVERS_CONFIG (debug)` 步骤会打印 runner 实际看到的原始值（长度、行数、Python repr），方便排查「配了但没读到」类问题。

### 调试：跑一次失败的 deploy 怎么排错

| 失败位置 | 看什么 |
|---|---|
| `Check required env vars` 步骤 | 日志列出哪几项 `EMPTY (not accessible to this workflow)`，照名字补上即可。常见原因：拼错大小写、配到了 `Settings/Keys`（Deploy keys）、漏配。 |
| Strategy 求值阶段就报 `fromJSON: empty input` | `vars.SERVERS_CONFIG` 未配。`Print SERVERS_CONFIG (debug)` 步骤会显示原始值。 |
| `Diagnose` 步骤报 `SERVERS_CONFIG_VALIDATION` | `SERVERS_CONFIG` 存在但 JSON 不合法（不是数组、缺字段、端口非数字、有重名），stderr 上有具体哪一项的哪个字段。 |
| SSH / scp 步骤 | `Add server to known_hosts`、`Upload tarball` 等步骤会用 `matrix.server.X` 拼出目标。如果某台机器 `addr` 填错、SSH 端口不通、或 `DEPLOY_SSH_KEY` 公钥没加到目标机器，会在这里挂掉，**但不影响其他机器**（`fail-fast: false`）。 |

## API 文档

服务启动后访问 `http://localhost:8080` 查看相关接口。

## 许可证

MIT