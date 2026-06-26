package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/smileyan/backend/config"
	"github.com/smileyan/backend/middleware"
	"github.com/smileyan/backend/models"
	"github.com/smileyan/backend/services"
	"github.com/smileyan/backend/utils"
)

// 获取文章列表（管理员用，支持筛选状态）
func GetAdminPosts(c *gin.Context) {
	var posts []models.Post
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	categoryID := c.Query("category_id")
	tagID := c.Query("tag_id")
	keyword := c.Query("keyword")

	offset := (page - 1) * pageSize

	// 先构建查询条件，不应用分页
	db := config.DB.Model(&models.Post{}).Order("created_at DESC")

	// 支持按状态筛选
	if status != "" {
		db = db.Where("status = ?", models.PostStatus(status))
	}

	// 支持按分类筛选
	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	// 支持按标签筛选
	if tagID != "" {
		db = db.Joins("JOIN post_tags ON posts.id = post_tags.post_id").Where("post_tags.tag_id = ?", tagID)
	}

	// 支持按标题关键字搜索
	if keyword != "" {
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}

	// 先计算总数（不带 offset/limit）
	var total int64
	db.Count(&total)

	// 再查询数据（带 offset/limit）
	db.Offset(offset).Limit(pageSize).Preload("Category").Preload("Tags").Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"list":      posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// 获取文章列表（前台）
func GetPosts(c *gin.Context) {
	var posts []models.Post
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	categoryID := c.Query("category_id")
	tagID := c.Query("tag_id")
	keyword := c.Query("keyword")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	offset := (page - 1) * pageSize

	db := config.DB.Model(&models.Post{})

	// 管理员可以看到所有文章，非管理员只能看已发布的
	user := middleware.GetCurrentUser(c)
	isAdmin := user != nil && user.Role == models.RoleAdmin
	if !isAdmin {
		db = db.Where("status = ?", models.PostStatusPublished)
	} else if status := c.Query("status"); status != "" {
		db = db.Where("status = ?", models.PostStatus(status))
	}

	if categoryID != "" {
		db = db.Where("category_id = ?", categoryID)
	}

	if tagID != "" {
		db = db.Joins("JOIN post_tags ON post_tags.post_id = posts.id").Where("post_tags.tag_id = ?", tagID)
	}

	if keyword != "" {
		db = db.Where("title LIKE ?", "%"+keyword+"%")
	}

	// 支持按日期范围筛选
	if startDate != "" {
		db = db.Where("DATE(created_at) >= ?", startDate)
	}
	if endDate != "" {
		db = db.Where("DATE(created_at) <= ?", endDate)
	}

	var total int64
	db.Count(&total)

	db.Preload("Category").Preload("Tags").Preload("Author").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts)

	// 处理摘要
	for i := range posts {
		if posts[i].Excerpt == "" && posts[i].Content != "" {
			posts[i].Excerpt = getExcerpt(posts[i].Content, 200)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"list":  posts,
		"total": total,
		"page":  page,
	})
}

// 获取文章详情
func GetPost(c *gin.Context) {
	slug := c.Param("slug")

	var post models.Post
	if err := config.DB.Preload("Category").Preload("Tags").Preload("Author").
		Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 检查权限：管理员可以看到所有文章，非管理员只能看已发布的
	user := middleware.GetCurrentUser(c)
	isAdmin := user != nil && user.Role == models.RoleAdmin
	// 管理员在编辑模式(isAdmin && edit=true)下可以查看所有状态的文章
	if post.Status != models.PostStatusPublished && !(isAdmin && c.Query("edit") == "true") {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 编辑模式不增加阅读量
	if c.Query("edit") != "true" {
		config.DB.Model(&post).Update("view_count", post.ViewCount+1)
	}

	c.JSON(http.StatusOK, post)
}

// 获取文章详情（管理员用，通过ID）
func GetPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var post models.Post
	if err := config.DB.Preload("Category").Preload("Tags").Preload("Author").
		First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// 获取文章详情（管理员用，通过Slug）
func GetAdminPostBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var post models.Post
	if err := config.DB.Preload("Category").Preload("Tags").Preload("Author").
		Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// 创建文章
func CreatePost(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Title       string   `json:"title" binding:"required"`
		Slug        string   `json:"slug"`
		Content     string   `json:"content" binding:"required"`
		Excerpt     string   `json:"excerpt"`
		CoverImage  string   `json:"cover_image"`
		CategoryID  *uint    `json:"category_id"`
		TagIDs      []uint   `json:"tag_ids"`
		Status      string   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 生成slug
	if req.Slug == "" {
		req.Slug = generateSlug(req.Title)
	}

	// 检查slug是否已存在（包括已删除的文章）
	var existingPost models.Post
	if err := config.DB.Unscoped().Where("slug = ?", req.Slug).First(&existingPost).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug already exists"})
		return
	}

	// Markdown转HTML
	htmlContent := string(utils.MarkdownToHTML([]byte(req.Content)))

	post := models.Post{
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		HTMLContent: htmlContent,
		Excerpt:     req.Excerpt,
		CoverImage:  req.CoverImage,
		CategoryID:  req.CategoryID,
		Status:      models.PostStatus(req.Status),
		AuthorID:    user.ID,
	}

	// 保存文章
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		config.DB.Find(&tags, req.TagIDs)
		post.Tags = tags
		config.DB.Save(&post)
	}

	// 更新搜索索引
	services.AddToSearchIndex(post)

	c.JSON(http.StatusCreated, post)
}

// 更新文章
func UpdatePost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Title       string   `json:"title"`
		Slug        string   `json:"slug"`
		Content     string   `json:"content"`
		Excerpt     string   `json:"excerpt"`
		CoverImage  string   `json:"cover_image"`
		CategoryID  *uint    `json:"category_id"`
		TagIDs      []uint   `json:"tag_ids"`
		Status      string   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 更新字段
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.Content != "" {
		post.Content = req.Content
		post.HTMLContent = string(utils.MarkdownToHTML([]byte(req.Content)))
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	if req.CoverImage != "" {
		post.CoverImage = req.CoverImage
	}
	if req.CategoryID != nil {
		post.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		post.Status = models.PostStatus(req.Status)
	}

	config.DB.Save(&post)

	// 更新标签：先清空旧关联，再写入新关联
	if req.TagIDs != nil {
		var tags []models.Tag
		config.DB.Find(&tags, req.TagIDs)
		config.DB.Model(&post).Association("Tags").Replace(tags)
	}

	// 更新搜索索引
	services.UpdateSearchIndex(post)

	c.JSON(http.StatusOK, post)
}

// 删除文章（隐藏状态则硬删除，其他状态软删除）
func DeletePost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 检查权限：只有管理员可以删除文章
	user := middleware.GetCurrentUser(c)
	if user == nil || user.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// 隐藏状态的文章进行硬删除
	if post.Status == models.PostStatusHidden {
		config.DB.Unscoped().Delete(&post)
		// 从搜索索引中删除
		services.RemoveFromSearchIndex(post.ID)
		c.JSON(http.StatusOK, gin.H{"message": "Post permanently deleted"})
		return
	}

	// 其他状态执行软删除
	post.IsDeleted = true
	post.Status = models.PostStatusHidden
	post.Slug = utils.GenerateRandomString(8) // 修改slug为8位随机码
	config.DB.Save(&post)

	// 从搜索索引中删除
	services.RemoveFromSearchIndex(post.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}

// 恢复文章
func RestorePost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var post models.Post
	// 使用 Unscoped 查询已软删除的记录
	if err := config.DB.Unscoped().First(&post, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 检查权限：只有管理员可以恢复文章
	user := middleware.GetCurrentUser(c)
	if user == nil || user.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	post.IsDeleted = false
	post.Status = models.PostStatusDraft
	config.DB.Save(&post)

	// 重新添加到搜索索引
	services.AddToSearchIndex(post)

	c.JSON(http.StatusOK, gin.H{"message": "Post restored"})
}

// 获取分类列表
func GetCategories(c *gin.Context) {
	var categories []models.Category
	config.DB.Order("sort ASC, id ASC").Find(&categories)

	// 计算每个分类的文章数量（只统计已发布的）
	type categoryWithCount struct {
		models.Category
		PostCount int64 `json:"post_count"`
	}
	result := make([]categoryWithCount, len(categories))
	for i, cat := range categories {
		result[i].Category = cat
		config.DB.Model(&models.Post{}).
			Where("category_id = ? AND status = ? AND is_deleted = false", cat.ID, models.PostStatusPublished).
			Count(&result[i].PostCount)
	}

	c.JSON(http.StatusOK, result)
}

// 创建分类
func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	category := models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Sort:        req.Sort,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// 更新分类
func UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		Description string `json:"description"`
		Sort        int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.Sort > 0 {
		category.Sort = req.Sort
	}

	config.DB.Save(&category)
	c.JSON(http.StatusOK, category)
}

// 删除分类
func DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := config.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

// 获取标签列表
func GetTags(c *gin.Context) {
	var tags []models.Tag
	config.DB.Find(&tags)

	// 计算每个标签的文章数量（只统计已发布的）
	type tagWithCount struct {
		models.Tag
		PostCount int64 `json:"post_count"`
	}
	result := make([]tagWithCount, len(tags))
	for i, tag := range tags {
		result[i].Tag = tag
		config.DB.Model(&models.Post{}).
			Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Where("post_tags.tag_id = ? AND posts.status = ? AND posts.is_deleted = false", tag.ID, models.PostStatusPublished).
			Count(&result[i].PostCount)
	}

	c.JSON(http.StatusOK, result)
}

// 创建标签
func CreateTag(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
		Slug string `json:"slug"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Slug == "" {
		req.Slug = generateSlug(req.Name)
	}

	tag := models.Tag{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := config.DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// 删除标签
func DeleteTag(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := config.DB.Delete(&models.Tag{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tag"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted"})
}

// 辅助函数
func getExcerpt(content string, length int) string {
	// 去除markdown语法
	content = strings.ReplaceAll(content, "#", "")
	content = strings.ReplaceAll(content, "*", "")
	content = strings.ReplaceAll(content, "_", "")
	content = strings.ReplaceAll(content, "`", "")
	content = strings.ReplaceAll(content, "\n", " ")
	content = strings.ReplaceAll(content, "  ", " ")

	if len(content) > length {
		return content[:length] + "..."
	}
	return content
}

func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	// 保留中文字符
	var result []rune
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r >= 0x4E00 {
			result = append(result, r)
		}
	}
	return string(result)
}