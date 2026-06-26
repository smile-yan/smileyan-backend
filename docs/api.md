# API 接口文档

## 概述

Smileyan 博客系统提供 RESTful API 接口。

## 基础信息

- 基础路径：`/api`
- 管理员路径：`/api/admin`
- 认证方式：JWT Token（通过 Header 传递）
- 请求格式：JSON
- 响应格式：JSON

## 通用响应格式

### 成功响应
```json
{
  "message": "操作成功",
  "data": {}
}
```

### 错误响应
```json
{
  "error": "错误信息"
}
```

## 认证接口

### 发送验证码
发送登录验证码到用户邮箱。

- **URL**: `POST /api/send-code`
- **请求体**:
```json
{
  "email": "user@example.com"
}
```
- **响应**:
```json
{
  "message": "Verification code sent"
}
```
- **错误**:
  - `400`: 邮箱格式错误
  - `429`: 请求过于频繁
  - `500`: 发送失败

### 登录
使用邮箱和验证码登录。

- **URL**: `POST /api/login`
- **请求体**:
```json
{
  "email": "user@example.com",
  "code": "123456"
}
```
- **响应**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "nickname": "user",
    "avatar": "/uploads/avatars/1_avatar.jpg",
    "role": "reader"
  }
}
```
- **错误**:
  - `400`: 验证码错误或过期

### 获取当前用户信息
获取已登录用户的信息。

- **URL**: `GET /api/user`
- **认证**: 需要（Bearer Token）
- **响应**:
```json
{
  "id": 1,
  "email": "user@example.com",
  "nickname": "user",
  "avatar": "/uploads/avatars/1_avatar.jpg",
  "role": "reader"
}
```

### 更新用户信息
更新当前用户的信息。

- **URL**: `PUT /api/user`
- **认证**: 需要
- **请求体**:
```json
{
  "nickname": "新昵称",
  "avatar": "/uploads/avatars/new_avatar.jpg"
}
```
- **响应**:
```json
{
  "id": 1,
  "email": "user@example.com",
  "nickname": "新昵称",
  "avatar": "/uploads/avatars/new_avatar.jpg",
  "role": "reader"
}
```

### 上传头像
上传用户头像。

- **URL**: `POST /api/upload-avatar`
- **认证**: 需要
- **Content-Type**: `multipart/form-data`
- **表单字段**: `avatar`（文件）
- **响应**:
```json
{
  "avatar": "/uploads/avatars/1_avatar.jpg"
}
```

---

## 文章接口

### 获取文章列表
获取文章列表，支持分页。

- **URL**: `GET /api/posts`
- **查询参数**:
  - `page`: 页码（默认1）
  - `page_size`: 每页数量（默认10）
  - `category_id`: 分类ID筛选
  - `tag_id`: 标签ID筛选
  - `status`: 状态筛选（published/draft/hidden）
- **响应**:
```json
{
  "list": [
    {
      "id": 1,
      "title": "文章标题",
      "slug": "article-slug",
      "excerpt": "文章摘要",
      "cover_image": "/uploads/cover.jpg",
      "status": "published",
      "view_count": 100,
      "created_at": "2024-01-01T00:00:00Z",
      "category": {
        "id": 1,
        "name": "分类名称"
      },
      "tags": [
        {
          "id": 1,
          "name": "标签"
        }
      ]
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 10
}
```

### 获取单个文章
根据 slug 获取文章详情。

- **URL**: `GET /api/posts/:slug`
- **响应**:
```json
{
  "id": 1,
  "title": "文章标题",
  "slug": "article-slug",
  "content": "# Markdown内容",
  "html_content": "<h1>HTML内容</h1>",
  "excerpt": "文章摘要",
  "cover_image": "/uploads/cover.jpg",
  "status": "published",
  "view_count": 100,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "category": {
    "id": 1,
    "name": "分类名称"
  },
  "tags": [
    {
      "id": 1,
      "name": "标签"
    }
  ],
  "author": {
    "id": 1,
    "nickname": "作者"
  }
}
```

### 创建文章（管理员）
创建新文章。

- **URL**: `POST /api/admin/posts`
- **认证**: 需要（管理员）
- **请求体**:
```json
{
  "title": "文章标题",
  "slug": "article-slug",
  "content": "# Markdown内容",
  "excerpt": "文章摘要",
  "cover_image": "/uploads/cover.jpg",
  "status": "published",
  "category_id": 1,
  "tag_ids": [1, 2, 3]
}
```
- **响应**:
```json
{
  "id": 1,
  "title": "文章标题"
}
```

### 更新文章（管理员）
更新文章。

- **URL**: `PUT /api/admin/posts/:id`
- **认证**: 需要（管理员）
- **请求体**: 同创建

### 删除文章（管理员）
删除文章（软删除）。

- **URL**: `DELETE /api/admin/posts/:id`
- **认证**: 需要（管理员）
- **响应**:
```json
{
  "message": "Post deleted"
}
```

### 恢复文章（管理员）
恢复已删除的文章。

- **URL**: `POST /api/admin/posts/:id/restore`
- **认证**: 需要（管理员）

---

## 分类接口

### 获取分类列表
获取所有分类。

- **URL**: `GET /api/categories`
- **响应**:
```json
{
  "list": [
    {
      "id": 1,
      "name": "分类名称",
      "slug": "category-slug",
      "description": "分类描述",
      "sort": 0
    }
  ]
}
```

### 创建分类（管理员）
- **URL**: `POST /api/admin/categories`
- **认证**: 需要（管理员）
- **请求体**:
```json
{
  "name": "分类名称",
  "slug": "category-slug",
  "description": "分类描述",
  "sort": 0
}
```

### 更新分类（管理员）
- **URL**: `PUT /api/admin/categories/:id`

### 删除分类（管理员）
- **URL**: `DELETE /api/admin/categories/:id`

---

## 标签接口

### 获取标签列表
获取所有标签。

- **URL**: `GET /api/tags`
- **响应**:
```json
{
  "list": [
    {
      "id": 1,
      "name": "标签名称",
      "slug": "tag-slug"
    }
  ]
}
```

### 创建标签（管理员）
- **URL**: `POST /api/admin/tags`
- **认证**: 需要（管理员）
- **请求体**:
```json
{
  "name": "标签名称",
  "slug": "tag-slug"
}
```

### 删除标签（管理员）
- **URL**: `DELETE /api/admin/tags/:id`

---

## 页面接口

### 获取页面列表
获取所有自定义页面。

- **URL**: `GET /api/pages`
- **响应**:
```json
{
  "list": [
    {
      "id": 1,
      "title": "页面标题",
      "slug": "page-slug",
      "status": "published"
    }
  ]
}
```

### 获取单个页面
根据 slug 获取页面详情。

- **URL**: `GET /api/pages/:slug`
- **响应**: 包含 content 和 html_content

### 创建页面（管理员）
- **URL**: `POST /api/admin/pages`
- **请求体**:
```json
{
  "title": "页面标题",
  "slug": "page-slug",
  "content": "# Markdown内容",
  "status": "published"
}
```

### 更新页面（管理员）
- **URL**: `PUT /api/admin/pages/:id`

### 删除页面（管理员）
- **URL**: `DELETE /api/admin/pages/:id`

---

## 评论接口

### 获取评论列表
获取文章的评论列表。

- **URL**: `GET /api/comments`
- **查询参数**:
  - `post_id`: 文章ID（必填）
  - `page`: 页码
  - `page_size`: 每页数量
- **响应**:
```json
{
  "list": [
    {
      "id": 1,
      "content": "评论内容",
      "status": "approved",
      "created_at": "2024-01-01T00:00:00Z",
      "user": {
        "id": 1,
        "nickname": "用户昵称",
        "avatar": "/uploads/avatars/1.jpg"
      }
    }
  ],
  "total": 100
}
```

### 创建评论
发表评论。

- **URL**: `POST /api/comments`
- **认证**: 需要
- **请求体**:
```json
{
  "post_id": 1,
  "content": "评论内容",
  "parent_id": 0
}
```

### 审核评论（管理员）
通过评论。

- **URL**: `POST /api/admin/comments/:id/approve`
- **认证**: 需要（管理员）

### 拒绝评论（管理员）
拒绝评论。

- **URL**: `POST /api/admin/comments/:id/reject`
- **认证**: 需要（管理员）

### 删除评论（管理员）
- **URL**: `DELETE /api/admin/comments/:id`

---

## 搜索接口

### 搜索文章
搜索文章。

- **URL**: `GET /api/search`
- **查询参数**:
  - `keyword`: 关键词（必填）
  - `page`: 页码
  - `page_size`: 每页数量
- **响应**:
```json
{
  "list": [
    {
      "id": "slug",
      "title": "文章标题",
      "excerpt": "摘要",
      "highlights": ["<em>关键词</em>"]
    }
  ],
  "total": 10
}
```

---

## 订阅接口

### 订阅
订阅博客更新通知。

- **URL**: `POST /api/subscribe`
- **请求体**:
```json
{
  "email": "subscriber@example.com"
}
```
- **响应**:
```json
{
  "message": "Subscribed successfully"
}
```

### 退订
退订邮件通知。

- **URL**: `GET /api/unsubscribe`
- **查询参数**:
  - `email`: 邮箱
  - `token`: 退订Token
- **响应**:
```json
{
  "message": "Unsubscribed successfully"
}
```

### 获取订阅列表（管理员）
- **URL**: `GET /api/admin/subscriptions`
- **认证**: 需要（管理员）

### 删除订阅（管理员）
- **URL**: `DELETE /api/admin/subscriptions/:id`

### 发送订阅通知（管理员）
向所有订阅者发送新文章通知。

- **URL**: `POST /api/admin/notify`
- **认证**: 需要（管理员）
- **请求体**:
```json
{
  "post_id": 1
}
```

---

## 用户管理（管理员）

### 获取用户列表
- **URL**: `GET /api/admin/users`
- **认证**: 需要（管理员）
- **查询参数**:
  - `page`: 页码
  - `page_size`: 每页数量

### 更新用户角色
- **URL**: `PUT /api/admin/users/:id/role`
- **认证**: 需要（管理员）
- **请求体**:
```json
{
  "role": "admin"  // admin 或 reader
}
```

---

## 错误码

| 状态码 | 说明 |
|--------|------|
| 400 | 请求参数错误 |
| 401 | 未授权（Token无效或过期） |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 429 | 请求过于频繁 |
| 500 | 服务器内部错误 |