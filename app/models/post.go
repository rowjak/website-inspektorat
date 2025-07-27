package models

import (
	"github.com/goravel/framework/database/orm"
)

type Post struct {
	orm.Model
	Judul            string `gorm:"column:judul"`
	Slug             string `gorm:"column:slug;uniqueIndex"`
	Readmore         string `gorm:"column:readmore"`
	Isi              string `gorm:"column:isi"`
	Tags             string `gorm:"column:tags"`
	ThumbnailSm      string `gorm:"column:thumbnail_sm"`
	ThumbnailLg      string `gorm:"column:thumbnail_lg"`
	AttachmentOgName string `gorm:"column:attachment_og_name"`
	AttachmentName   string `gorm:"column:attachment_name"`
	AttachmentSize   int64  `gorm:"column:attachment_size"`
	AttachmentMime   string `gorm:"column:attachment_mime"`
	Views            int    `gorm:"column:views;default:0"`
	Status           string `gorm:"column:status;default:Ditampilkan"` // publish, draft, trash
	KategoriID       uint   `gorm:"column:kategori_id"`
	UserID           uint   `gorm:"column:user_id"`
	TanggalTampil    string `gorm:"column:tanggal_tampil;type:date"`

	PostKategori *PostKategori `gorm:"foreignKey:KategoriID;references:ID" json:"post_category,omitempty"`
	PostImage    []*PostImage  `gorm:"foreignKey:PostID;references:ID" json:"post_images,omitempty"`
}

func (u *Post) TableName() string {
	return "posts"
}
