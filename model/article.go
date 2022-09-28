// @Title  article
// @Description  该中间件用于“拦截”运行时恐慌的内建函数,防止程序崩溃
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Article			定义了文章的各种元素
type Article struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`    // 文章的id
	UserId    uint      `json:"user_id" gorm:"not null"`                // 文章的作者id
	Title     string    `json:"title" gorm:"type:varchar(50);not null"` // 文章的标题
	Content   string    `json:"content" gorm:"type:text;not null"`      // 文章的内容
	ResLong   string    `json:"res_long" gorm:"type:text;"`             // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text;"`            // 备用短文本
	Visible   int8      `json:"visible" gorm:"type:tinyint;default:1"`  // 可见等级
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`       // 文章的创建日期
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`       // 文章的更新日期
}

// @title    BeforeCreate
// @description   计算出一个uint
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}
