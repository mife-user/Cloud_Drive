package repo

import (
	"drive/internal/domain"
	"sync"

	"golang.org/x/sync/singleflight"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type fileRepo struct {
	db    *gorm.DB
	rd    *redis.Client
	rds   *singleflight.Group
	locks sync.Map
}

// NewFileRepo 创建文件仓库
func NewFileRepo(db *gorm.DB, rd *redis.Client) domain.FileRepo {
	return &fileRepo{db: db, rd: rd, rds: &singleflight.Group{}}
}

// LockByID 按ID锁
func (r *fileRepo) LockByID(lockKey string) func() {
	lock, _ := r.locks.LoadOrStore(lockKey, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)
	mutex.Lock()
	return mutex.Unlock
}
