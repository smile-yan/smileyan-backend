package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/smileyan/backend/config"
	"github.com/smileyan/backend/models"
	"github.com/smileyan/backend/utils"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Dev    string `json:"dev"` // 开发环境跳过认证的标识
	jwt.RegisteredClaims
}

// JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析token
		cfg := config.GetConfig()
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 开发环境：检查 dev token 是否匹配环境变量
		devToken := os.Getenv("SMILEYAN_BACKEND_DEV_TOKEN")
		if devToken != "" && claims.Dev == devToken {
			// 使用 dev token，设置默认管理员用户
			c.Set("user_id", uint(1))
			c.Set("email", "root@smileyan.cn")
			c.Set("role", "admin")
			c.Next()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		// 处理 interface{} 类型
		roleStr, ok := role.(string)
		if !ok || roleStr != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 获取当前用户ID
func GetCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	// 处理 interface{} 类型
	switch v := userID.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	default:
		return 0
	}
}

// 获取当前用户
func GetCurrentUser(c *gin.Context) *models.User {
	userID := GetCurrentUserID(c)
	if userID == 0 {
		return nil
	}

	// 开发环境：如果 user_id 是 1，直接返回默认管理员
	if userID == 1 {
		return &models.User{
			ID:       1,
			Email:    "root@smileyan.cn",
			Nickname: "管理员",
			Role:     models.RoleAdmin,
		}
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		utils.Error("Failed to get current user", utils.Field("error", err))
		return nil
	}
	return &user
}

// 请求日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求开始时间
		utils.Info("Request started",
			utils.Field("method", c.Request.Method),
			utils.Field("path", c.Request.URL.Path),
			utils.Field("ip", c.ClientIP()),
		)

		c.Next()

		// 请求结束
		utils.Info("Request finished",
			utils.Field("method", c.Request.Method),
			utils.Field("path", c.Request.URL.Path),
			utils.Field("status", c.Writer.Status()),
		)
	}
}
