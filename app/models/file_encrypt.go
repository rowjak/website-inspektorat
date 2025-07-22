package models

import (
	"github.com/goravel/framework/database/orm"
)

type FileEncrypt struct {
	orm.Model
	FileOgName      string `gorm:"column:file_og_name"`
	FileEncryptName string `gorm:"column:file_encrypt_name"`
	FilePath        string `gorm:"column:file_path"`
	FileMime        string `gorm:"column:file_mime"`
	FileSize        int64  `gorm:"column:file_size"`
}

func (u *FileEncrypt) TableName() string {
	return "file_encrypt"
}
