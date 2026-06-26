package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smileyan/backend/config"
	"github.com/smileyan/backend/middleware"
	"github.com/smileyan/backend/models"
	"github.com/smileyan/backend/utils"
)

// 获取文章评论（公开接口）
func GetComments(c *gin.Context) {
	postID := c.Query("post_id")

	var comments []models.Comment
	query := config.DB.Preload("User").Preload("Children.User").Order("created_at DESC")

	if postID != "" {
		query = query.Where("post_id = ?", postID)
	}

	// 非管理员只能看到已审核的评论
	user := middleware.GetCurrentUser(c)
	if user != nil && user.Role == models.RoleAdmin {
		// 管理员：按 status 筛选，如果不传则显示所有
		if status := c.Query("status"); status != "" {
			query = query.Where("status = ?", models.CommentStatus(status))
		}
	} else {
		query = query.Where("status = ?", models.CommentStatusApproved)
	}

	query.Find(&comments)

	c.JSON(http.StatusOK, comments)
}

// 获取评论列表（管理员专用接口，支持分页）
func GetAdminComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	offset := (page - 1) * pageSize

	query := config.DB.Model(&models.Comment{}).Preload("User").Preload("Post")

	if status != "" {
		query = query.Where("status = ?", models.CommentStatus(status))
	}

	var total int64
	query.Count(&total)

	var comments []models.Comment
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"list":  comments,
		"total": total,
	})
}

// 创建评论
func CreateComment(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please login first"})
		return
	}

	var req struct {
		PostID   uint   `json:"post_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 检查文章是否存在
	var post models.Post
	if err := config.DB.First(&post, req.PostID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	comment := models.Comment{
		PostID:   req.PostID,
		UserID:   user.ID,
		Content:  req.Content,
		ParentID: req.ParentID,
		Status:   models.CommentStatusPending, // 默认待审核
	}

	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	// 如果是回复，发送邮件通知
	if req.ParentID != nil {
		var parentComment models.Comment
		if err := config.DB.Preload("User").First(&parentComment, req.ParentID).Error; err == nil {
			if parentComment.User.Email != user.Email {
				subject := "您在 Smileyan 博客的评论有新回复"
				body := "<html><body><p>您好，" + parentComment.User.Nickname + "</p>" +
					"<p>您有一条新的回复：</p>" +
					"<blockquote>" + req.Content + "</blockquote>" +
					"<p><a href='/post/" + post.Slug + "#comment-" + strconv.Itoa(int(comment.ID)) + "'>查看回复</a></p>" +
					"</body></html>"
				utils.SendEmail(parentComment.User.Email, subject, body)
			}
		}
	}

	// 加载用户信息
	config.DB.Preload("User").First(&comment, comment.ID)

	c.JSON(http.StatusCreated, comment)
}

// 审核评论（管理员）
func ApproveComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var comment models.Comment
	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	comment.Status = models.CommentStatusApproved
	config.DB.Save(&comment)

	c.JSON(http.StatusOK, gin.H{"message": "Comment approved"})
}

// 拒绝评论（管理员）
func RejectComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var comment models.Comment
	if err := config.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	comment.Status = models.CommentStatusRejected
	config.DB.Save(&comment)

	c.JSON(http.StatusOK, gin.H{"message": "Comment rejected"})
}

// 删除评论（管理员）
func DeleteComment(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := config.DB.Delete(&models.Comment{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}