package bootstrap

import (
	"drive/pkg/conf"
	"drive/pkg/cron"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Application struct {
	Config   *conf.Config
	Database *gorm.DB
	Redis    *redis.Client
	Engine   *gin.Engine
	Cron     *cron.Cron
}
