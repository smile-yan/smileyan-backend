# 博客 API 接口文档

本文档记录了博客创建相关的 curl 请求命令。

## 环境配置

```bash
export SMILEYAN_BACKEND_JWT_SECRET="smileyan-jwt-secret-key-2024"
export SMILEYAN_BACKEND_DEV_TOKEN="dev-skip-auth-token-2024"
```

## 生成 JWT Token

```bash
cd backend && go run cmd/genjwt/main.go
```

获取到的 Token 用于后续接口的 Authorization 头。

## 创建博客

### 1. Smileyan 博客系统 - 技术栈详解

```bash
curl -X POST http://localhost:8080/api/admin/posts \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Smileyan 博客系统 - 技术栈详解",
    "slug": "smileyan-tech-stack",
    "content": "# Smileyan 博客系统技术栈\n\nSmileyan 是一个现代化的博客系统，采用前后端分离架构构建。\n\n## 技术架构\n\n### 后端\n\n- **Gin**: Go 语言高性能 Web 框架\n- **GORM**: Go 语言的 ORM 库\n- **MySQL**: 主数据库存储\n- **Redis**: 缓存和验证码存储\n\n### 前端\n\n- **Vue 3**: 渐进式 JavaScript 框架\n- **Pinia**: Vue 3 状态管理\n- **Vue Router**: 路由管理\n- **Markdown-it**: Markdown 解析\n\n## 核心功能\n\n1. **用户认证**: 邮箱验证码登录\n2. **文章管理**: Markdown 文章支持\n3. **分类标签**: 灵活的内容组织\n4. **评论系统**: 嵌套评论\n5. **订阅通知**: 邮件订阅更新",
    "status": "published"
  }'
```

### 2. Go语言 Gin 框架实战

```bash
curl -X POST http://localhost:8080/api/admin/posts \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Go语言 Gin 框架实战：构建高性能博客后端",
    "slug": "go-gin-blog-backend",
    "content": "# Go语言 Gin 框架实战：构建高性能博客后端\n\n本文将介绍如何使用 Go 语言的 Gin 框架构建一个高性能的博客后端系统。\n\n## 为什么选择 Gin？\n\n- **高性能**：基于 Radix 树的路由\n- **中间件支持**：灵活的中间件机制\n- **JSON 支持**：内置强大的 JSON 解析",
    "status": "published"
  }'
```

### 3. Redis 缓存实战

```bash
curl -X POST http://localhost:8080/api/admin/posts \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Redis 缓存实战：提升应用性能的终极指南",
    "slug": "redis-cache-guide",
    "content": "# Redis 缓存实战\n\nRedis 是最流行的内存数据库之一，本文将介绍如何在博客系统中高效使用 Redis。\n\n## 在博客系统中的应用\n\n### 1. 验证码存储\n### 2. 频率限制\n### 3. 会话缓存",
    "status": "published"
  }'
```

### 4. 开发效率提升之道

```bash
curl -X POST http://localhost:8080/api/admin/posts \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "我的开发效率提升之道：工具与习惯",
    "slug": "dev-efficiency-tips",
    "content": "# 我的开发效率提升之道\n\n分享一些提升开发效率的实用工具和良好习惯。\n\n## 效率工具\n- VS Code\n- iTerm2\n- tmux\n\n## 开发习惯\n- 代码规范\n- 版本控制\n- 自动化",
    "status": "published"
  }'
```

### 5. 周末coding

```bash
curl -X POST http://localhost:8080/api/admin/posts \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "周末coding：一个人的技术狂欢",
    "slug": "weekend-coding",
    "content": "# 周末coding：一个人的技术狂欢\n\n每个周末，我都会抽出时间Coding，这已经成为我生活中不可或缺的仪式。\n\n## 为什么选择周末？\n- 解决复杂的技术难题\n- 学习新的技术栈\n- 实践项目创意",
    "status": "published"
  }'
```

### 6. Vue 3 组合式 API

```bash
curl -X POST http://localhost:8080/api/admin/posts \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Vue 3 组合式 API 实战指南",
    "slug": "vue3-composition-api",
    "content": "# Vue 3 组合式 API 实战指南\n\nVue 3 引入的组合式 API 是近年来前端领域最重要的变革之一。\n\n## 核心 API\n- setup()\n- ref 和 reactive\n- 生命周期钩子",
    "status": "published"
  }'
```

### 7. Markdown 完全指南

```bash
curl -X POST http://localhost:8080/api/admin/posts \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Markdown 完全指南：让你的博客文章更精彩",
    "slug": "markdown-guide",
    "content": "# Markdown 完全指南\n\nMarkdown 是一种轻量级标记语言，让我们能够用简洁的语法编写格式丰富的文本。\n\n## 1. 标题\n## 2. 文本格式\n## 3. 代码块\n## 4. 列表\n## 5. 表格",
    "status": "published"
  }'
```

## 其他常用接口

### 获取所有博客

```bash
curl http://localhost:8080/api/posts
```

### 获取单个博客

```bash
curl http://localhost:8080/api/posts/<slug>
```

### 更新博客

```bash
curl -X PUT http://localhost:8080/api/admin/posts/<id> \
  -H "Authorization: Bearer <YOUR_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "新标题",
    "content": "新内容",
    "status": "published"
  }'
```

### 删除博客

```bash
curl -X DELETE http://localhost:8080/api/admin/posts/<id> \
  -H "Authorization: Bearer <YOUR_TOKEN>"
```
