package repo

import (
	"drive/internal/domain"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type fileRepo struct {
	db    *gorm.DB
	rd    *redis.Client
	locks sync.Map
}

// NewFileRepo 创建文件仓库
func NewFileRepo(db *gorm.DB, rd *redis.Client) domain.FileRepo {
	return &fileRepo{db: db, rd: rd}
}

// LockByID 按ID锁
func (r *fileRepo) LockByID(fileID uint) func() {
	lock, _ := r.locks.LoadOrStore(fileID, &sync.Mutex{})
	mutex := lock.(*sync.Mutex)
	mutex.Lock()
	return mutex.Unlock
}
