package utils

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-mail/mail/v2"
	"github.com/smileyan/backend/config"
)

// 发送邮件
func SendEmail(to, subject, body string) error {
	cfg := config.GetConfig()

	m := mail.NewMessage()
	m.SetHeader("From", cfg.Email.Username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dialer := mail.NewDialer(
		cfg.Email.Host,
		cfg.Email.Port,
		cfg.Email.Username,
		cfg.Email.Password,
	)
	dialer.SSL = cfg.Email.UseSSL

	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// 生成验证码
func GenerateVerificationCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := ""
	for i := 0; i < 6; i++ {
		code += fmt.Sprintf("%d", r.Intn(10))
	}
	return code
}

// 保存验证码到Redis
func SaveVerificationCode(email, code string) error {
	ctx := context.Background()
	key := fmt.Sprintf("verification_code:%s", email)

	// 10分钟有效期
	return config.Redis.Set(ctx, key, code, 10*time.Minute).Err()
}

// 获取验证码
func GetVerificationCode(email string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("verification_code:%s", email)

	code, err := config.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return code, nil
}

// 删除验证码
func DeleteVerificationCode(email string) error {
	ctx := context.Background()
	key := fmt.Sprintf("verification_code:%s", email)

	return config.Redis.Del(ctx, key).Err()
}

// 检查发送频率限制
func CheckRateLimit(email, ip string) (bool, error) {
	ctx := context.Background()

	// 检查同一邮箱10分钟（600秒）限制
	emailKey := fmt.Sprintf("rate_limit:email:%s", email)
	exists, err := config.Redis.Exists(ctx, emailKey).Result()
	if err != nil {
		return false, err
	}
	if exists > 0 {
		return false, nil
	}

	// 检查同一IP每小时限制（100次）
	ipKey := fmt.Sprintf("rate_limit:ip:%s", ip)
	count, err := config.Redis.Get(ctx, ipKey).Result()
	if err != nil && err.Error() != "redis: nil" {
		return false, err
	}

	if count != "" {
		var c int
		fmt.Sscanf(count, "%d", &c)
		if c >= 100 {
			return false, nil
		}
	}

	return true, nil
}

// 设置发送频率限制
func SetRateLimit(email, ip string) error {
	ctx := context.Background()

	// 设置邮箱10分钟（600秒）限制
	emailKey := fmt.Sprintf("rate_limit:email:%s", email)
	config.Redis.Set(ctx, emailKey, "1", 600*time.Second)

	// 设置IP每小时限制（100次）
	ipKey := fmt.Sprintf("rate_limit:ip:%s", ip)
	config.Redis.Incr(ctx, ipKey)
	config.Redis.Expire(ctx, ipKey, time.Hour)

	return nil
}

// 生成订阅退订token
func GenerateUnsubscribeToken(email string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	token := ""
	for i := 0; i < 32; i++ {
		token += fmt.Sprintf("%x", r.Intn(16))
	}
	return token
}