package models

import "github.com/goravel/framework/database/orm"

type Tag struct {
	orm.Model
	Slug string `gorm:"column:slug"`
	Nama string `gorm:"column:nama"`
}

func (u *Tag) TableName() string {
	return "tags"
}
