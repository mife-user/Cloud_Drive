package repo

import (
	"drive/internal/domain"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
	rd *redis.Client
}

func NewUserRepo(db *gorm.DB, rd *redis.Client) domain.UserRepo {
	return &userRepo{db: db, rd: rd}
}
