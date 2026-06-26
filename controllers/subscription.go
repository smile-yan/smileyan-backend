package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/smileyan/backend/config"
	"github.com/smileyan/backend/models"
	"github.com/smileyan/backend/utils"
)

// 订阅博客
func Subscribe(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}

	// 检查是否已订阅
	var sub models.Subscription
	if err := config.DB.Where("email = ?", req.Email).First(&sub).Error; err == nil {
		if sub.IsActive {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Already subscribed"})
			return
		}
		// 重新激活订阅
		sub.IsActive = true
		sub.Token = utils.GenerateUnsubscribeToken(req.Email)
		config.DB.Save(&sub)
		c.JSON(http.StatusOK, gin.H{"message": "Subscribed successfully"})
		return
	}

	// 创建新订阅
	token := utils.GenerateUnsubscribeToken(req.Email)
	sub = models.Subscription{
		Email:    req.Email,
		Token:    token,
		IsActive: true,
	}

	if err := config.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe"})
		return
	}

	// 发送确认邮件
	subject := "欢迎订阅 Smileyan 博客"
	body := "<html><body><h2>感谢您的订阅！</h2>" +
		"<p>您将收到新文章发布通知。</p>" +
		"<p><a href='https://smileyan.cn/unsubscribe?token=" + token + "&email=" + req.Email + "'>退订</a></p>" +
		"</body></html>"
	utils.SendEmail(req.Email, subject, body)

	c.JSON(http.StatusOK, gin.H{"message": "Subscribed successfully"})
}

// 退订
func Unsubscribe(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	var sub models.Subscription
	var err error

	// 如果有 token，同时验证 email 和 token
	if token != "" {
		err = config.DB.Where("email = ? AND token = ?", email, token).First(&sub).Error
	} else {
		// 只根据 email 查找
		err = config.DB.Where("email = ?", email).First(&sub).Error
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	sub.IsActive = false
	config.DB.Save(&sub)

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}

// 获取订阅列表（管理员）
func GetSubscriptions(c *gin.Context) {
	var subs []models.Subscription
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	offset := (parseInt(page) - 1) * parseInt(pageSize)
	config.DB.Offset(offset).Limit(parseInt(pageSize)).Find(&subs)

	var total int64
	config.DB.Model(&models.Subscription{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"list":  subs,
		"total": total,
	})
}

// 删除订阅（管理员）
func DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Subscription{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Subscription deleted"})
}

// 发送新文章通知（管理员）
func NotifySubscribers(c *gin.Context) {
	var req struct {
		PostID uint `json:"post_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 获取文章
	var post models.Post
	if err := config.DB.First(&post, req.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// 获取所有活跃订阅者
	var subs []models.Subscription
	config.DB.Where("is_active = ?", true).Find(&subs)

	// 发送邮件
	for _, sub := range subs {
		subject := "【新文章】" + post.Title
		body := "<html><body>" +
			"<h2>Smileyan 博客有新文章发布</h2>" +
			"<h3>" + post.Title + "</h3>" +
			"<p>" + post.Excerpt + "</p>" +
			"<p><a href='https://smileyan.cn/post/" + post.Slug + "'>阅读全文</a></p>" +
			"<hr>" +
			"<p><a href='https://smileyan.cn/unsubscribe?token=" + sub.Token + "&email=" + sub.Email + "'>退订</a></p>" +
			"</body></html>"
		utils.SendEmail(sub.Email, subject, body)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notifications sent", "count": len(subs)})
}