package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := "smileyan-jwt-secret-key-2024"
	if s := os.Getenv("SMILEYAN_BACKEND_JWT_SECRET"); s != "" {
		secret = s
	}

	// 方式1: 构造 dev token（需要服务器设置了 SMILEYAN_BACKEND_DEV_TOKEN 环境变量）
	devTokenValue := "dev-skip-auth-token-2024"
	if d := os.Getenv("SMILEYAN_BACKEND_DEV_TOKEN"); d != "" {
		devTokenValue = d
	}

	devClaims := jwt.MapClaims{
		"dev": devTokenValue,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(365 * 24 * time.Hour).Unix(),
	}
	devJWT, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, devClaims).SignedString([]byte(secret))

	// 方式2: 直接构造管理员 token（绕过 dev token 机制）
	// admin_emails 列表中有 i@ccyan.cn，构造一个该邮箱的管理员 JWT
	adminClaims := jwt.MapClaims{
		"user_id": float64(1),
		"email":   "i@ccyan.cn",
		"role":    "admin",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(365 * 24 * time.Hour).Unix(),
	}
	adminJWT, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, adminClaims).SignedString([]byte(secret))

	fmt.Println("=== 方式1: Dev Token（需服务器设置 DEV_TOKEN 环境变量）===")
	fmt.Println(devJWT)
	fmt.Println()
	fmt.Println("=== 方式2: 管理员 Token（直接声明 admin 角色，推荐）===")
	fmt.Println(adminJWT)
	fmt.Println()
	fmt.Println("使用方法：")
	fmt.Println("  curl -H 'Authorization: Bearer " + adminJWT + "' https://bigbigpig.cn/api/admin/posts")
}
