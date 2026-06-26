package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smileyan/backend/config"
	"github.com/smileyan/backend/middleware"
	"github.com/smileyan/backend/models"
	"github.com/smileyan/backend/utils"
)

// 获取页面列表
func GetPages(c *gin.Context) {
	var pages []models.Page
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	offset := (parseInt(page) - 1) * parseInt(pageSize)

	config.DB.Offset(offset).Limit(parseInt(pageSize)).Find(&pages)

	var total int64
	config.DB.Model(&models.Page{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"list":  pages,
		"total": total,
	})
}

// 获取页面详情（管理员按 ID）
func GetPageByID(c *gin.Context) {
	id := c.Param("id")
	var page models.Page
	if err := config.DB.First(&page, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}
	c.JSON(http.StatusOK, page)
}

// 获取页面详情
func GetPage(c *gin.Context) {
	slug := c.Param("slug")

	var page models.Page
	if err := config.DB.Where("slug = ?", slug).First(&page).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	// 检查权限
	user := middleware.GetCurrentUser(c)
	if page.Status != models.PostStatusPublished {
		if user == nil || user.Role != models.RoleAdmin {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
			return
		}
	}

	c.JSON(http.StatusOK, page)
}

// 创建页面
func CreatePage(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Slug    string `json:"slug"`
		Content string `json:"content"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Slug == "" {
		req.Slug = generateSlug(req.Title)
	}

	htmlContent := string(utils.MarkdownToHTML([]byte(req.Content)))

	page := models.Page{
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		HTMLContent: htmlContent,
		Status:      models.PostStatusPublished,
	}

	if req.Status != "" {
		page.Status = models.PostStatus(req.Status)
	}

	if err := config.DB.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create page"})
		return
	}

	c.JSON(http.StatusCreated, page)
}

// 更新页面
func UpdatePage(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title   string `json:"title"`
		Slug    string `json:"slug"`
		Content string `json:"content"`
		Status  string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var page models.Page
	if err := config.DB.First(&page, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	if req.Title != "" {
		page.Title = req.Title
	}
	if req.Slug != "" {
		page.Slug = req.Slug
	}
	if req.Content != "" {
		page.Content = req.Content
		page.HTMLContent = string(utils.MarkdownToHTML([]byte(req.Content)))
	}
	if req.Status != "" {
		page.Status = models.PostStatus(req.Status)
	}

	config.DB.Save(&page)
	c.JSON(http.StatusOK, page)
}

// 删除页面
func DeletePage(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Page{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete page"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Page deleted"})
}

func parseInt(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}