package models

import (
	"github.com/goravel/framework/database/orm"
)

type Post struct {
	orm.Model
	Judul           string `gorm:"column:judul"`
	Slug            string `gorm:"column:slug;uniqueIndex"`
	Isi             string `gorm:"column:isi"`
	ThumbnailOgName string `gorm:"column:thumbnail_og_name"`
	ThumbnailName   string `gorm:"column:thumbnail_name"`

	PostImage []*PostImage `gorm:"foreignKey:PostID;references:ID" json:"post_images,omitempty"` // kenapa references:ID? karena PostID adalah foreign key yang mengacu pada ID di tabel Post, kalau berbeda kan bisa di ganti sendiri
}

func (u *Post) TableName() string {
	return "posts"
}
