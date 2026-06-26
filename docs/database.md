# 数据库表结构文档

## 概述

本文档详细介绍 Smileyan 博客系统的数据库表结构。

## 表列表

1. [users - 用户表](#users-用户表)
2. [categories - 分类表](#categories-分类表)
3. [tags - 标签表](#tags-标签表)
4. [posts - 文章表](#posts-文章表)
5. [pages - 页面表](#pages-页面表)
6. [comments - 评论表](#comments-评论表)
7. [subscriptions - 订阅表](#subscriptions-订阅表)

---

## users - 用户表

用户表存储系统用户信息。

### 表结构

| 字段 | 类型 | 主键 | 自增 | 可为空 | 默认值 | 说明 |
|------|------|------|------|--------|--------|------|
| id | bigint unsigned | 是 | 是 | 否 | 自增 | 用户ID |
| created_at | datetime(3) | - | - | 是 | NULL | 创建时间 |
| updated_at | datetime(3) | - | - | 是 | NULL | 更新时间 |
| deleted_at | datetime(3) | - | - | 是 | NULL | 删除时间（软删除） |
| email | longtext | - | - | 否 | - | 邮箱（唯一索引） |
| nickname | varchar(50) | - | - | 是 | NULL | 昵称 |
| avatar | longtext | - | - | 是 | NULL | 头像URL |
| role | varchar(20) | - | - | 是 | reader | 角色：admin/reader |
| status | varchar(20) | - | - | 是 | active | 状态：active/inactive |

### 索引

| 索引名 | 字段 | 唯一 | 说明 |
|--------|------|------|------|
| idx_users_deleted_at | deleted_at | 否 | 软删除索引 |
| idx_users_email | email | 是 | 邮箱唯一索引 |

### 角色说明

| 角色值 | 说明 |
|--------|------|
| admin | 管理员，拥有所有管理权限 |
| reader | 普通读者，默认角色 |

### 状态说明

| 状态值 | 说明 |
|--------|------|
| active | 活跃 |
| inactive | 未激活 |

---

## categories - 分类表

分类表用于管理文章的分类。

### 表结构

| 字段 | 类型 | 主键 | 自增 | 可为空 | 默认值 | 说明 |
|------|------|------|------|--------|--------|------|
| id | bigint unsigned | 是 | 是 | 否 | 自增 | 分类ID |
| created_at | datetime(3) | - | - | 是 | NULL | 创建时间 |
| updated_at | datetime(3) | - | - | 是 | NULL | 更新时间 |
| deleted_at | datetime(3) | - | - | 是 | NULL | 删除时间（软删除） |
| name | varchar(50) | - | - | 否 | - | 分类名称 |
| slug | varchar(50) | - | - | 是 | NULL | 分类别名（URL友好） |
| description | varchar(200) | - | - | 是 | NULL | 分类描述 |
| sort | bigint | - | - | 是 | 0 | 排序权重 |

### 索引

| 索引名 | 字段 | 唯一 | 说明 |
|--------|------|------|------|
| idx_categories_deleted_at | deleted_at | 否 | 软删除索引 |
| idx_categories_slug | slug | 是 | 别名唯一索引 |

---

## tags - 标签表

标签表用于管理文章的标签。

### 表结构

| 字段 | 类型 | 主键 | 自增 | 可为空 | 默认值 | 说明 |
|------|------|------|------|--------|--------|------|
| id | bigint unsigned | 是 | 是 | 否 | 自增 | 标签ID |
| created_at | datetime(3) | - | - | 是 | NULL | 创建时间 |
| updated_at | datetime(3) | - | - | 是 | NULL | 更新时间 |
| deleted_at | datetime(3) | - | - | 是 | NULL | 删除时间（软删除） |
| name | varchar(50) | - | - | 否 | - | 标签名称 |
| slug | varchar(50) | - | - | 是 | NULL | 标签别名（URL友好） |

### 索引

| 索引名 | 字段 | 唯一 | 说明 |
|--------|------|------|------|
| idx_tags_slug | slug | 是 | 别名唯一索引 |
| idx_tags_deleted_at | deleted_at | 否 | 软删除索引 |

---

## posts - 文章表

文章表存储博客文章内容。

### 表结构

| 字段 | 类型 | 主键 | 自增 | 可为空 | 默认值 | 说明 |
|------|------|------|------|--------|--------|------|
| id | bigint unsigned | 是 | 是 | 否 | 自增 | 文章ID |
| created_at | datetime(3) | - | - | 是 | NULL | 创建时间 |
| updated_at | datetime(3) | - | - | 是 | NULL | 更新时间 |
| deleted_at | datetime(3) | - | - | 是 | NULL | 删除时间（软删除） |
| title | varchar(200) | - | - | 否 | - | 文章标题 |
| slug | varchar(200) | - | - | 是 | NULL | 文章别名（URL友好） |
| content | text | - | - | 是 | NULL | Markdown格式内容 |
| html_content | text | - | - | 是 | NULL | HTML格式内容 |
| excerpt | varchar(500) | - | - | 是 | NULL | 文章摘要 |
| cover_image | longtext | - | - | 是 | NULL | 封面图URL |
| status | varchar(20) | - | - | 是 | draft | 状态：draft/published/hidden |
| view_count | bigint | - | - | 是 | 0 | 阅读数 |
| is_deleted | tinyint(1) | - | - | 是 | false | 软删除标记 |
| category_id | bigint unsigned | - | - | 是 | NULL | 关联分类ID |
| author_id | bigint unsigned | - | - | 是 | NULL | 作者ID |

### 索引

| 索引名 | 字段 | 唯一 | 说明 |
|--------|------|------|------|
| idx_posts_deleted_at | deleted_at | 否 | 软删除索引 |
| idx_posts_slug | slug | 是 | 别名唯一索引 |
| fk_posts_category | category_id | 否 | 外键关联分类 |
| fk_posts_author | author_id | 否 | 外键关联用户 |

### 状态说明

| 状态值 | 说明 |
|--------|------|
| draft | 草稿，未发布 |
| published | 已发布 |
| hidden | 已隐藏 |

### 关联关系

- 多对一：每个文章属于一个分类（Category）
- 多对多：每个文章可以有多个标签（Tag），通过 post_tags 中间表关联
- 多对一：每个文章属于一个作者（User）

---

## pages - 页面表

页面表存储自定义页面内容（如关于页面、联系页面等）。

### 表结构

| 字段 | 类型 | 主键 | 自增 | 可为空 | 默认值 | 说明 |
|------|------|------|------|--------|--------|------|
| id | bigint unsigned | 是 | 是 | 否 | 自增 | 页面ID |
| created_at | datetime(3) | - | - | 是 | NULL | 创建时间 |
| updated_at | datetime(3) | - | - | 是 | NULL | 更新时间 |
| deleted_at | datetime(3) | - | - | 是 | NULL | 删除时间（软删除） |
| title | varchar(200) | - | - | 否 | - | 页面标题 |
| slug | varchar(200) | - | - | 是 | NULL | 页面别名（URL友好） |
| content | text | - | - | 是 | NULL | Markdown格式内容 |
| html_content | text | - | - | 是 | NULL | HTML格式内容 |
| status | varchar(20) | - | - | 是 | published | 状态：draft/published/hidden |

### 索引

| 索引名 | 字段 | 唯一 | 说明 |
|--------|------|------|------|
| idx_pages_deleted_at | deleted_at | 否 | 软删除索引 |
| idx_pages_slug | slug | 是 | 别名唯一索引 |

---

## comments - 评论表

评论表存储用户对文章的评论。

### 表结构

| 字段 | 类型 | 主键 | 自增 | 可为空 | 默认值 | 说明 |
|------|------|------|------|--------|--------|------|
| id | bigint unsigned | 是 | 是 | 否 | 自增 | 评论ID |
| created_at | datetime(3) | - | - | 是 | NULL | 创建时间 |
| updated_at | datetime(3) | - | - | 是 | NULL | 更新时间 |
| deleted_at | datetime(3) | - | - | 是 | NULL | 删除时间（软删除） |
| content | text | - | - | 否 | - | 评论内容 |
| status | varchar(20) | - | - | 是 | pending | 状态：pending/approved/rejected |
| post_id | bigint unsigned | - | - | 是 | NULL | 关联文章ID |
| user_id | bigint unsigned | - | - | 是 | NULL | 评论用户ID |
| parent_id | bigint unsigned | - | - | 是 | NULL | 父评论ID（支持嵌套） |

### 索引

| 索引名 | 字段 | 唯一 | 说明 |
|--------|------|------|------|
| idx_comments_deleted_at | deleted_at | 否 | 软删除索引 |

### 状态说明

| 状态值 | 说明 |
|--------|------|
| pending | 待审核 |
| approved | 已通过 |
| rejected | 已拒绝 |

### 关联关系

- 多对一：每个评论属于一篇文章（Post）
- 多对一：每个评论属于一个用户（User）
- 自关联：支持嵌套评论（ParentID 指向父评论）

---

## subscriptions - 订阅表

订阅表存储邮件订阅用户信息。

### 表结构

| 字段 | 类型 | 主键 | 自增 | 可为空 | 默认值 | 说明 |
|------|------|------|------|--------|--------|------|
| id | bigint unsigned | 是 | 是 | 否 | 自增 | 订阅ID |
| created_at | datetime(3) | - | - | 是 | NULL | 创建时间 |
| updated_at | datetime(3) | - | - | 是 | NULL | 更新时间 |
| deleted_at | datetime(3) | - | - | 是 | NULL | 删除时间（软删除） |
| email | longtext | - | - | 否 | - | 订阅邮箱（唯一索引） |
| token | varchar(64) | - | - | 是 | NULL | 退订Token（唯一索引） |
| is_active | tinyint(1) | - | - | 是 | true | 是否激活 |

### 索引

| 索引名 | 字段 | 唯一 | 说明 |
|--------|------|------|------|
| idx_subscriptions_deleted_at | deleted_at | 否 | 软删除索引 |
| idx_subscriptions_email | email | 是 | 邮箱唯一索引 |
| idx_subscriptions_token | token | 是 | Token唯一索引 |

---

## 关联表

### post_tags - 文章标签关联表

用于实现文章和标签的多对多关系。

| 字段 | 类型 | 说明 |
|------|------|------|
| post_id | bigint unsigned | 文章ID |
| tag_id | bigint unsigned | 标签ID |

---

## 软删除说明

本系统使用 GORM 的软删除功能，通过 `deleted_at` 字段标记记录是否被删除。被标记删除的记录在查询时默认不显示，但数据仍然保留在数据库中。