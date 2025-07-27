package models

import "github.com/goravel/framework/database/orm"

type PostKategori struct {
	orm.Model
	Slug     string `gorm:"column:slug"`
	Kategori string `gorm:"column:kategori"`
}

func (u *PostKategori) TableName() string {
	return "post_category"
}
