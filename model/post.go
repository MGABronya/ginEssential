package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	UserId    uint      `json:"user_id" gorm:"not null"`
	Title     string    `json:"title" gorm:"type:varchar(50);not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`
	Name      string    `json:"name" gorm:"type:varchar(50)"`
	Icon      string    `json:"icon" gorm:"type:varchar(50)"` //这里的Icon存储的是图像文件的地址
}

func (post *Post) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}