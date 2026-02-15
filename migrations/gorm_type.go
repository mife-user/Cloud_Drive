package migrations

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(100);not null;default:'默认用户';comment:用户名"`
	PassWord string `gorm:"type:varchar(100);not null;comment:密码"`
	Role     string `gorm:"type:varchar(50);not null;default:'NOVIP';comment:角色"`
	Files    []File `gorm:"foreignKey:UserID"`
}

type File struct {
	gorm.Model
	FileName    string `gorm:"type:varchar(255);not null;default:'未命名文件';comment:文件名"`
	Size        int64  `gorm:"not null;default:1024;comment:文件大小"`
	Path        string `gorm:"type:varchar(255);not null;default:'/';comment:文件路径"`
	UserID      uint   `gorm:"not null;comment:用户ID"`
	Owner       string `gorm:"type:varchar(100);not null;default:'默认用户';comment:文件所有者"`
	Permissions string `gorm:"type:varchar(100);not null;default:'可以共享';comment:权限"`
}

type FileShare struct {
	gorm.Model
	FileID    uint   `gorm:"not null;index;comment:文件ID"`
	ShareID   string `gorm:"type:varchar(100);not null;uniqueIndex;comment:分享ID"`
	AccessKey string `gorm:"type:varchar(100);not null;comment:访问密钥"`
	UserID    uint   `gorm:"not null;comment:用户ID"`
	Owner     string `gorm:"type:varchar(100);not null;default:'默认用户';comment:分享所有者"`
	ExpiresAt int64  `gorm:"not null;comment:过期时间戳"`
}
