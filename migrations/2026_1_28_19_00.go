package migrations

import "gorm.io/gorm"

var Migration_2026_1_28_19_00 = migration{
	Version:  202601281900,
	IsUpdate: 0,
	Up:       up,
	Down:     down,
}

type driveUser struct {
	gorm.Model
	UserName   string      `gorm:"type:varchar(100);not null;default:'默认用户';comment:用户名"`
	PassWord   string      `gorm:"type:varchar(100);not null;comment:密码"`
	Role       string      `gorm:"type:varchar(50);not null;default:'普通用户';comment:角色"`
	ownedFiles []driveFile `gorm:"foreignKey:FileID"`
}

type driveFile struct {
	gorm.Model
	FileName    string `gorm:"type:varchar(255);not null;default:'未命名文件';comment:文件名"`
	Size        int64  `gorm:"not null;default:1024;comment:文件大小"`
	FileID      uint   `gorm:"not null;unique;comment:文件ID"`
	Owner       driveUser
	Permissions string `gorm:"type:varchar(100);not null;default:'不可写';comment:权限"`
}

//更新
func up(db *gorm.DB) error {
	return db.Migrator().CreateTable(&driveUser{}, &driveFile{})
}

//回滚
func down(db *gorm.DB) error {
	return db.Migrator().DropTable(&driveUser{}, &driveFile{})
}
