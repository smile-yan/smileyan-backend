package main

import (
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/smileyan/backend/config"
	"github.com/smileyan/backend/middleware"
	"github.com/smileyan/backend/models"
	"github.com/smileyan/backend/routes"
	"github.com/smileyan/backend/services"
	"github.com/smileyan/backend/utils"
)

func main() {
	// 初始化配置
	config.InitConfig()

	// 初始化日志
	utils.InitLogger()
	utils.Info("Starting Smileyan Blog System...")

	// 初始化数据库
	config.InitDatabase()
	utils.Info("Database connected")

	// 初始化Redis
	config.InitRedis()
	utils.Info("Redis connected")

	// 数据库自动迁移
	models.AutoMigrate()
	utils.Info("Database migrated")

	// 初始化搜索服务
	services.InitSearchService()
	utils.Info("Search service initialized")

	// 加载已有文章到搜索索引
	var posts []models.Post
	config.DB.Find(&posts)
	for i := range posts {
		services.AddToSearchIndex(posts[i])
	}
	if len(posts) > 0 {
		utils.Info(fmt.Sprintf("Indexed %d posts", len(posts)))
	}

	// 创建上传目录
	createUploadDirs()

	// 设置Gin
	cfg := config.GetConfig()
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS配置
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	// 请求日志中间件
	r.Use(middleware.LoggerMiddleware())

	// 设置路由
	routes.SetupRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	utils.Info("Server starting on " + addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server: " + err.Error())
	}
}

func createUploadDirs() {
	// 创建上传目录
	// 在生产环境中应该使用系统调用
}