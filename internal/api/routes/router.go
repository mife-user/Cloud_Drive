package routes

import (
	"drive/internal/api/handlers"
	"drive/internal/api/middlewares"
	"drive/internal/repo"
	"drive/pkg/conf"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Router 路由结构体
type Router struct {
	authHandler *handlers.AuthHandler
	fileHandler *handlers.FileHandler
	userHandler *handlers.UserHandler
	config      *conf.Config
}

// GetRouter 获取路由实例
func GetRouter() *Router {
	return &Router{}
}

// NewRouter 初始化路由
func (r *Router) NewRouter(db *gorm.DB, config *conf.Config) bool {
	userRepo := repo.NewUserRepo(db)
	fileRepo := repo.NewFileRepo(db)
	r.authHandler = handlers.NewAuthHandler(userRepo, config)
	r.fileHandler = handlers.NewFileHandler(fileRepo, config)
	r.userHandler = handlers.NewUserHandler(userRepo, config)
	r.config = config
	return true
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(r.config.Gin.Mode)

	// 创建路由
	router := gin.Default()

	// 配置 CORS
	router.Use(middlewares.CORSMiddleware(r.config))

	// API 路由组
	api := router.Group("/api")
	{
		// 用户路由 - 公开
		user := api.Group("/user")
		{
			user.POST("/register", r.userHandler.Register)
			user.POST("/login", r.authHandler.Login)
		}

		// 文件路由 - 需要认证
		file := api.Group("/file")
		file.Use(middlewares.AuthMiddleware(r.config))
		{
			file.POST("/upload", r.fileHandler.UploadFile)
		}
	}

	return router
}
