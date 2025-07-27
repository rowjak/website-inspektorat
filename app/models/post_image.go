package models

import (
	"github.com/goravel/framework/database/orm"
)

type PostImage struct {
	orm.Model
	PostID  uint   `gorm:"column:post_id"`
	ImageSm string `gorm:"column:image_sm"`
	ImageLg string `gorm:"column:image_lg"`

	Post *Post `gorm:"foreignKey:PostID;references:ID" json:"post,omitempty"` // kenapa references:ID? karena PostID adalah foreign key yang mengacu pada ID di tabel Post
}

func (u *PostImage) TableName() string {
	return "post_images"
}
