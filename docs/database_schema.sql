-- Smileyan Blog Database Schema
-- Generated at: 2026-04-20
-- Database Engine: MySQL InnoDB

-- =============================================================================
-- posts - 文章表
-- =============================================================================
DROP TABLE IF EXISTS posts;
CREATE TABLE posts (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at      DATETIME(3) DEFAULT NULL,
    updated_at      DATETIME(3) DEFAULT NULL,
    deleted_at      DATETIME(3) DEFAULT NULL,
    title           VARCHAR(200) NOT NULL,
    slug            VARCHAR(200) NOT NULL,
    content         TEXT,
    html_content    TEXT,
    excerpt         VARCHAR(500) DEFAULT NULL,
    cover_image     LONGTEXT,
    status          VARCHAR(20) DEFAULT 'draft',
    view_count      BIGINT DEFAULT '0',
    is_deleted      TINYINT(1) DEFAULT '0',
    category_id     BIGINT UNSIGNED DEFAULT NULL,
    author_id       BIGINT UNSIGNED DEFAULT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_posts_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =============================================================================
-- categories - 分类表
-- =============================================================================
DROP TABLE IF EXISTS categories;
CREATE TABLE categories (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at      DATETIME(3) DEFAULT NULL,
    updated_at      DATETIME(3) DEFAULT NULL,
    deleted_at      DATETIME(3) DEFAULT NULL,
    name            VARCHAR(50) NOT NULL,
    slug            VARCHAR(50) NOT NULL,
    description     VARCHAR(200) DEFAULT NULL,
    sort            BIGINT DEFAULT '0',
    PRIMARY KEY (id),
    UNIQUE INDEX idx_categories_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =============================================================================
-- tags - 标签表
-- =============================================================================
DROP TABLE IF EXISTS tags;
CREATE TABLE tags (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at      DATETIME(3) DEFAULT NULL,
    updated_at      DATETIME(3) DEFAULT NULL,
    deleted_at      DATETIME(3) DEFAULT NULL,
    name            VARCHAR(50) NOT NULL,
    slug            VARCHAR(50) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE INDEX idx_tags_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =============================================================================
-- pages - 页面表
-- =============================================================================
DROP TABLE IF EXISTS pages;
CREATE TABLE pages (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at      DATETIME(3) DEFAULT NULL,
    updated_at      DATETIME(3) DEFAULT NULL,
    deleted_at      DATETIME(3) DEFAULT NULL,
    title           VARCHAR(200) NOT NULL,
    slug            VARCHAR(200) NOT NULL,
    content         TEXT,
    html_content    TEXT,
    status          VARCHAR(20) DEFAULT 'published',
    PRIMARY KEY (id),
    UNIQUE INDEX idx_pages_slug (slug)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =============================================================================
-- users - 用户表
-- =============================================================================
DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at      DATETIME(3) DEFAULT NULL,
    updated_at      DATETIME(3) DEFAULT NULL,
    deleted_at      DATETIME(3) DEFAULT NULL,
    email           VARCHAR(255) NOT NULL,
    nickname        VARCHAR(50) DEFAULT NULL,
    avatar          LONGTEXT,
    role            VARCHAR(20) DEFAULT 'reader',
    status          VARCHAR(20) DEFAULT 'active',
    PRIMARY KEY (id),
    UNIQUE INDEX idx_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =============================================================================
-- comments - 评论表
-- =============================================================================
DROP TABLE IF EXISTS comments;
CREATE TABLE comments (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at      DATETIME(3) DEFAULT NULL,
    updated_at      DATETIME(3) DEFAULT NULL,
    deleted_at      DATETIME(3) DEFAULT NULL,
    content         TEXT NOT NULL,
    status          VARCHAR(20) DEFAULT 'pending',
    post_id         BIGINT UNSIGNED DEFAULT NULL,
    user_id         BIGINT UNSIGNED DEFAULT NULL,
    parent_id       BIGINT UNSIGNED DEFAULT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =============================================================================
-- subscriptions - 订阅表
-- =============================================================================
DROP TABLE IF EXISTS subscriptions;
CREATE TABLE subscriptions (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at      DATETIME(3) DEFAULT NULL,
    updated_at      DATETIME(3) DEFAULT NULL,
    deleted_at      DATETIME(3) DEFAULT NULL,
    email           VARCHAR(255) NOT NULL,
    token           VARCHAR(64) DEFAULT NULL,
    is_active       TINYINT(1) DEFAULT '1',
    PRIMARY KEY (id),
    UNIQUE INDEX idx_subscriptions_email (email),
    UNIQUE INDEX idx_subscriptions_token (token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =============================================================================
-- post_tags - 文章标签关联表
-- =============================================================================
DROP TABLE IF EXISTS post_tags;
CREATE TABLE post_tags (
    post_id         BIGINT UNSIGNED NOT NULL,
    tag_id          BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (post_id, tag_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
