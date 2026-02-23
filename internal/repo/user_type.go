package repo

import (
	"drive/internal/domain"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type userRepo struct {
	db    *gorm.DB
	rd    *redis.Client
	locks sync.Map
}

func NewUserRepo(db *gorm.DB, rd *redis.Client) domain.UserRepo {
	return &userRepo{db: db, rd: rd}
}

// LockByID 按ID锁
func (r *userRepo) LockByID(userID uint) func() {
	lock, _ := r.locks.LoadOrStore(userID, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)
	mutex.Lock()
	return mutex.Unlock
}
