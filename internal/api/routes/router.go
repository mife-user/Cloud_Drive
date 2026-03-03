package routes

import (
	"drive/internal/api/handlers"
	"drive/internal/api/middlewares"
	"drive/internal/repo"
	"drive/internal/service"
	"drive/pkg/conf"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// Router 路由结构体
type Router struct {
	fileHandler *handlers.FileHandler
	userHandler *handlers.UserHandler
	config      *conf.Config
}

// GetRouter 获取路由实例
func GetRouter() *Router {
	return &Router{}
}

// NewRouter 初始化路由
func (r *Router) NewRouter(db *gorm.DB, rd *redis.Client, config *conf.Config) bool {
	userRepo := repo.NewUserRepo(db, rd)
	fileRepo := repo.NewFileRepo(db, rd)
	userServicer := service.NewUserServicer(userRepo, config)
	fileServicer := service.NewFileServicer(fileRepo, config)
	r.fileHandler = handlers.NewFileHandler(fileServicer, config)
	r.userHandler = handlers.NewUserHandler(userServicer, config)
	r.config = config
	return true
}

// Setup 设置路由
func (r *Router) Setup() *gin.Engine {
	gin.SetMode(r.config.Gin.Mode)

	router := gin.Default()

	router.Use(middlewares.CORSMiddleware(r.config))

	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/register", r.userHandler.Register)
			user.POST("/login", r.userHandler.Login)
			user.POST("/header", middlewares.AuthMiddleware(r.config), r.userHandler.UpdateHeader)
			user.GET("/header/:username", r.userHandler.GetHeader)
		}

		api.GET("/file/share/:share_id", r.fileHandler.AccessShare)

		file := api.Group("/file")
		file.Use(middlewares.AuthMiddleware(r.config))
		{
			file.GET("/view/deleted", r.fileHandler.GetDeletedFiles)
			file.POST("/upload", middlewares.TypeCheck(r.config), r.fileHandler.UploadFile)
			file.GET("/view", r.fileHandler.ViewFilesNote)
			file.GET("/view/:file_id", r.fileHandler.ViewFile)
			file.POST("/share", r.fileHandler.ShareFile)
			file.PUT("/:file_id/permissions", r.fileHandler.UpdateFilePermissions)
			file.POST("/favorite", r.fileHandler.AddFavorite)
			file.DELETE("/favorite/:file_id", r.fileHandler.RemoveFavorite)
			file.GET("/favorites", r.fileHandler.GetFavorites)
			file.DELETE("/delete/:file_id", r.fileHandler.DeleteFile)
			file.DELETE("/delete/:file_id/forever", r.fileHandler.DeleteFileForever)
		}
	}

	return router
}
