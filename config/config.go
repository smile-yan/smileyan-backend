package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB     *gorm.DB
	Redis  *redis.Client
	Config *viper.Viper
)

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	UseSSL   bool   `mapstructure:"use_ssl"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

type UploadConfig struct {
	Path      string `mapstructure:"path"`
	AvatarPath string `mapstructure:"avatar_path"`
	MaxSize   int    `mapstructure:"max_size"`
}

type AppConfig struct {
	Server ServerConfig `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis RedisConfig `mapstructure:"redis"`
	Email EmailConfig `mapstructure:"email"`
	JWT JWTConfig `mapstructure:"jwt"`
	Upload UploadConfig `mapstructure:"upload"`
	AdminEmails []string `mapstructure:"admin_emails"`
}

// IsAdminEmail 检查邮箱是否为管理员邮箱
func (cfg *AppConfig) IsAdminEmail(email string) bool {
	for _, adminEmail := range cfg.AdminEmails {
		if email == adminEmail {
			return true
		}
	}
	return false
}

func InitConfig() *viper.Viper {
	// 加载 .env 文件（如果存在）
	_ = godotenv.Load()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// 支持环境变量覆盖
	v.AutomaticEnv()

	// 设置环境变量前缀（与 .env 文件中的变量名匹配）
	v.SetEnvPrefix("SMILEYAN_BACKEND")

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config: %v", err))
	}

	Config = v
	return v
}

func GetConfig() *AppConfig {
	v := Config
	if v == nil {
		v = InitConfig()
	}

	var appConfig AppConfig
	if err := v.Unmarshal(&appConfig); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %v", err))
	}

	// 环境变量覆盖敏感配置（使用大写以匹配 .env 文件）
	if val := os.Getenv("SMILEYAN_BACKEND_DB_HOST"); val != "" {
		appConfig.Database.Host = val
	}
	if val := os.Getenv("SMILEYAN_BACKEND_DB_USER"); val != "" {
		appConfig.Database.User = val
	}
	if val := os.Getenv("SMILEYAN_BACKEND_DB_PASSWORD"); val != "" {
		appConfig.Database.Password = val
	}
	if val := os.Getenv("SMILEYAN_BACKEND_DB_NAME"); val != "" {
		appConfig.Database.DBName = val
	}
	if val := os.Getenv("SMILEYAN_BACKEND_REDIS_HOST"); val != "" {
		appConfig.Redis.Host = val
	}
	if val := os.Getenv("SMILEYAN_BACKEND_REDIS_PASSWORD"); val != "" {
		appConfig.Redis.Password = val
	}
	if val := os.Getenv("SMILEYAN_BACKEND_REDIS_USERNAME"); val != "" {
		appConfig.Redis.Username = val
	}
	if val := os.Getenv("SMILEYAN_BACKEND_EMAIL_PASSWORD"); val != "" {
		appConfig.Email.Password = val
	}
	if val := os.Getenv("SMILEYAN_BACKEND_JWT_SECRET"); val != "" {
		appConfig.JWT.Secret = val
	}

	return &appConfig
}

func InitDatabase() *gorm.DB {
	cfg := GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get database: %v", err))
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return db
}

func InitRedis() *redis.Client {
	cfg := GetConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Username: cfg.Redis.Username,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	Redis = rdb
	return rdb
}