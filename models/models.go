package models

import (
	"time"

	"gorm.io/gorm"
)

// 用户角色
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleReader UserRole = "reader"
)

// 用户状态
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
)

// 用户
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `json:"avatar"`
	Role      UserRole       `gorm:"size:20;default:reader" json:"role"`
	Status    UserStatus     `gorm:"size:20;default:active" json:"status"`
}

func (User) TableName() string {
	return "users"
}

// 分类
type Category struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Slug      string         `gorm:"size:50;uniqueIndex" json:"slug"`
	Description string       `gorm:"size:200" json:"description"`
	Sort      int            `gorm:"default:0" json:"sort"`
}

func (Category) TableName() string {
	return "categories"
}

// 标签
type Tag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Slug      string         `gorm:"size:50;uniqueIndex" json:"slug"`
}

func (Tag) TableName() string {
	return "tags"
}

// 文章状态
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusHidden    PostStatus = "hidden"
)

// 文章
type Post struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Slug        string         `gorm:"size:200;uniqueIndex" json:"slug"`
	Content     string         `gorm:"type:text" json:"content"`          // Markdown
	HTMLContent string         `gorm:"type:text" json:"html_content"`     // HTML
	Excerpt     string         `gorm:"size:500" json:"excerpt"`           // 摘要
	CoverImage  string         `json:"cover_image"`                       // 封面图
	Status      PostStatus     `gorm:"size:20;default:draft;index" json:"status"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`       // 阅读数
	IsDeleted   bool           `gorm:"default:false;index" json:"is_deleted"` // 软删除标记

	// 关联
	CategoryID *uint      `gorm:"index" json:"category_id"`
	Category   Category   `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags       []Tag      `gorm:"many2many:post_tags;" json:"tags"`
	AuthorID   uint       `gorm:"index" json:"author_id"`
	Author     User       `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
}

func (Post) TableName() string {
	return "posts"
}

// 自定义页面
type Page struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Slug        string         `gorm:"size:200;uniqueIndex" json:"slug"`
	Content     string         `gorm:"type:text" json:"content"`
	HTMLContent string         `gorm:"type:text" json:"html_content"`
	Status      PostStatus     `gorm:"size:20;default:published" json:"status"`
}

func (Page) TableName() string {
	return "pages"
}

// 评论状态
type CommentStatus string

const (
	CommentStatusPending  CommentStatus = "pending"
	CommentStatusApproved CommentStatus = "approved"
	CommentStatusRejected CommentStatus = "rejected"
)

// 评论
type Comment struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Status    CommentStatus `gorm:"size:20;default:pending" json:"status"`

	// 关联
	PostID    uint   `json:"post_id"`
	Post      Post   `gorm:"foreignKey:PostID" json:"post,omitempty"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ParentID  *uint  `json:"parent_id"`  // 父评论ID，支持嵌套
	Parent    *Comment `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children  []Comment `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

// 订阅
type Subscription struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Token     string         `gorm:"size:64;uniqueIndex" json:"token"` // 退订token
	IsActive  bool           `gorm:"default:true" json:"is_active"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}