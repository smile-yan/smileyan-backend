package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/smileyan/backend/controllers"
	"github.com/smileyan/backend/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// 静态文件
	r.Static("/uploads", "./uploads")

	// 公共路由
	public := r.Group("/api")
	{
		// 测试接口
		public.GET("/hello", controllers.Hello)

		// 用户相关
		public.POST("/send-code", controllers.SendVerificationCode)
		public.POST("/login", controllers.Login)
		public.GET("/user", middleware.AuthMiddleware(), controllers.GetCurrentUser)

		// 文章相关
		public.GET("/posts", controllers.GetPosts)
		public.GET("/posts/:slug", controllers.GetPost)

		// 分类标签
		public.GET("/categories", controllers.GetCategories)
		public.GET("/tags", controllers.GetTags)

		// 页面
		public.GET("/pages", controllers.GetPages)
		public.GET("/pages/:slug", controllers.GetPage)

		// 评论
		public.GET("/comments", controllers.GetComments)

		// 搜索
		public.GET("/search", controllers.Search)

		// 订阅
		public.POST("/subscribe", controllers.Subscribe)
		public.GET("/unsubscribe", controllers.Unsubscribe)
	}

	// 需要登录的路由
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware())
	{
		// 用户
		auth.PUT("/user", controllers.UpdateUser)
		auth.POST("/upload-avatar", controllers.UploadAvatar)

		// 评论
		auth.POST("/comments", controllers.CreateComment)
	}

	// 管理员路由
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// 评论管理
		admin.GET("/comments", controllers.GetAdminComments)
		// 文章管理
		admin.GET("/posts", controllers.GetAdminPosts)
		admin.GET("/posts/slug/:slug", controllers.GetAdminPostBySlug)
		admin.GET("/posts/_id/:id", controllers.GetPostByID)
		admin.POST("/posts", controllers.CreatePost)
		admin.PUT("/posts/_id/:id", controllers.UpdatePost)
		admin.DELETE("/posts/_id/:id", controllers.DeletePost)
		admin.POST("/posts/_id/:id/restore", controllers.RestorePost)

		// 分类管理
		admin.POST("/categories", controllers.CreateCategory)
		admin.PUT("/categories/_id/:id", controllers.UpdateCategory)
		admin.DELETE("/categories/_id/:id", controllers.DeleteCategory)

		// 标签管理
		admin.POST("/tags", controllers.CreateTag)
		admin.DELETE("/tags/_id/:id", controllers.DeleteTag)

		// 页面管理
		admin.GET("/pages/_id/:id", controllers.GetPageByID)
		admin.POST("/pages", controllers.CreatePage)
		admin.PUT("/pages/_id/:id", controllers.UpdatePage)
		admin.DELETE("/pages/_id/:id", controllers.DeletePage)

		// 评论管理
		admin.POST("/comments/_id/:id/approve", controllers.ApproveComment)
		admin.POST("/comments/_id/:id/reject", controllers.RejectComment)
		admin.DELETE("/comments/_id/:id", controllers.DeleteComment)

		// 订阅管理
		admin.GET("/subscriptions", controllers.GetSubscriptions)
		admin.DELETE("/subscriptions/_id/:id", controllers.DeleteSubscription)
		admin.POST("/notify", controllers.NotifySubscribers)

		// 用户管理
		admin.GET("/users", controllers.GetUsers)
		admin.PUT("/users/:id/role", controllers.UpdateUserRole)
	}
}