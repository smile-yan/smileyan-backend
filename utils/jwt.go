package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/smileyan/backend/config"
	"go.uber.org/zap"
)

// Claims 自定义声明
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// 创建JWT Token
func GenerateToken(userID uint, email, role string) (string, error) {
	cfg := config.GetConfig()

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(cfg.JWT.ExpireHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// Field 创建zap字段
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}