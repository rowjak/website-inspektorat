package models

import "github.com/goravel/framework/database/orm"

type Carousel struct {
	orm.Model
	Keterangan string `gorm:"column:keterangan"`
	ImageSm    string `gorm:"column:image_sm"`
	ImageLg    string `gorm:"column:image_lg"`
	Link       string `gorm:"column:link"`
	Status     string `gorm:"column:status"`
}

func (u *Carousel) TableName() string {
	return "carousels"
}
