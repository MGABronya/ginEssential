package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Email     string `gorm:"type:varchar(50);not null;unique"`
	Password  string `gorm:"size:255;not null"`
	Icon      string `gorm:"type:varchar(50)"` //这里的Icon存储的是图像文件的地址
	Telephone string `gorm:"type:varchar(20)"`
	Blog      string `gorm:"type:varchar(25)"`
	QQ        string `gorm:"type:varchar(20)"`
	Sex       bool   `gorm:"type:boolean"`
	Address   string `gorm:"type:varchar(20)"`
	Hobby     string `gorm:"type:varchar(50)"`
}
