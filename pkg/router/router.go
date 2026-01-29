package router

import (
	"net/http"

	"drive/internal/domain"
	"drive/internal/repo"
	"drive/pkg/conf"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 路由结构体
type Router struct {
	userRepo domain.UserRepo
	fileRepo domain.FileRepo
	config   *conf.Config
}

// 创建路由
func NewRouter(db *gorm.DB, config *conf.Config) *Router {
	return &Router{
		userRepo: repo.NewUserRepo(db),
		fileRepo: repo.NewFileRepo(db),
		config:   config,
	}
}

// 设置路由
func (r *Router) Setup() *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(r.config.Gin.Mode)

	// 创建路由
	router := gin.Default()

	// 配置 CORS
	router.Use(func(c *gin.Context) {
		for _, origin := range r.config.Gin.Cors.AllowOrigins {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		for _, method := range r.config.Gin.Cors.AllowMethods {
			c.Header("Access-Control-Allow-Methods", method)
		}
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API 路由组
	api := router.Group("/api")
	{
		// 用户路由
		user := api.Group("/user")
		{
			user.POST("/register", r.Register)
			user.POST("/login", r.Login)
		}

		// 文件路由
		file := api.Group("/file")
		{
			file.POST("/upload", r.UploadFile)
		}
	}

	return router
}

// 用户注册
func (r *Router) Register(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := r.userRepo.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// 用户登录
func (r *Router) Login(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := r.userRepo.Logon(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登录成功", "user": user})
}

// 文件上传
func (r *Router) UploadFile(c *gin.Context) {
	var file domain.File
	if err := c.ShouldBindJSON(&file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := r.fileRepo.UploadFile(&file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "上传成功"})
}
