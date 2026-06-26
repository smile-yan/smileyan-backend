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

// 发送验证码
func SendVerificationCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// 检查是否是管理员邮箱，管理员邮箱跳过频率限制
	cfg := config.GetConfig()
	isAdmin := cfg.IsAdminEmail(req.Email)

	// 非管理员需要检查频率限制
	if !isAdmin {
		ip := c.ClientIP()
		allowed, err := utils.CheckRateLimit(req.Email, ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check rate limit"})
			return
		}
		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, please try again later"})
			return
		}
		// 设置频率限制
		utils.SetRateLimit(req.Email, ip)
	}

	// 生成验证码
	code := utils.GenerateVerificationCode()

	// 保存验证码
	if err := utils.SaveVerificationCode(req.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save verification code"})
		return
	}

	// 发送邮件
	subject := "Smileyan 博客验证码"
	body := "<html><body><h2>您的验证码是：</h2><h1>" + code + "</h1><p>有效期10分钟</p></body></html>"
	if err := utils.SendEmail(req.Email, subject, body); err != nil {
		utils.Error("Failed to send verification email", utils.Field("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
}

// 登录
func Login(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required,len=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 验证验证码
	savedCode, err := utils.GetVerificationCode(req.Email)
	if err != nil || savedCode != req.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid verification code"})
		return
	}

	// 删除已使用的验证码
	utils.DeleteVerificationCode(req.Email)

	// 查找或创建用户
	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		// 创建新用户
		// 检查是否为管理员邮箱
		role := models.RoleReader
		if config.GetConfig().IsAdminEmail(req.Email) {
			role = models.RoleAdmin
		}
		user = models.User{
			Email:    req.Email,
			Nickname: req.Email[:len(req.Email)-10], // 默认昵称
			Role:     role,
			Status:   models.StatusActive,
		}
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// 生成JWT Token
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	})
}

// 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"role":     user.Role,
	})
}

// 更新用户信息
func UpdateUser(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := config.DB.Save(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"email":    user.Email,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"role":     user.Role,
	})
}

// 上传头像
func UploadAvatar(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 检查文件类型
	ext := file.Filename[len(file.Filename)-4:]
	if ext != ".jpg" && ext != ".png" && ext != ".gif" && ext != ".jpeg" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only JPG, PNG, GIF allowed"})
		return
	}

	// 生成文件名
	filename := strconv.Itoa(int(user.ID)) + "_" + file.Filename
	filepath := "./uploads/avatars/" + filename

	// 确保目录存在
	// mkdir -p handled by init

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 更新用户头像
	avatarURL := "/uploads/avatars/" + filename
	user.Avatar = avatarURL
	config.DB.Save(user)

	c.JSON(http.StatusOK, gin.H{"avatar": avatarURL})
}

// 获取用户列表（管理员）
func GetUsers(c *gin.Context) {
	var users []models.User
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	offset := (page - 1) * pageSize
	config.DB.Offset(offset).Limit(pageSize).Find(&users)

	var total int64
	config.DB.Model(&models.User{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"list":  users,
		"total": total,
		"page":  page,
	})
}

// 更新用户角色（管理员）
func UpdateUserRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Role = models.UserRole(req.Role)
	config.DB.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "User role updated"})
}