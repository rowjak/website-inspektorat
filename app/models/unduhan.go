package models

import "github.com/goravel/framework/database/orm"

type Unduhan struct {
	orm.Model
	Jenis           string `gorm:"column:jenis"`
	FileOgName      string `gorm:"column:file_og_name"`
	FileName        string `gorm:"column:file_name"`
	FileMime        string `gorm:"column:file_mime"`
	FileSize        int64  `gorm:"column:file_size"`
	DownloadedCount uint   `gorm:"column:downloaded_count"`
	Status          string `gorm:"column:status;default:Ditampilkan"`
	Tanggal         string `gorm:"column:tanggal;type:date"`
}

func (u *Unduhan) TableName() string {
	return "unduhan"
}
