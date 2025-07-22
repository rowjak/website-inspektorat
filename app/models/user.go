package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	NamaLengkap string `gorm:"column:nama_lengkap"`
	Email       string `gorm:"column:email"`
	Password    string `gorm:"column:password" json:"-"` //json - adalah untuk hidden password
	Role        string `gorm:"column:role"`
}

func (u *User) TableName() string {
	return "users"
}
