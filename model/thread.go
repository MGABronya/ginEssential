// @Title  thread
// @Description  定义跟帖
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Tread			定义跟帖
type Thread struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"` // 跟帖的id
	UserId    uint      `json:"user_id" gorm:"not null"`             // 作者的id
	PostId    string    `json:"post_id" gorm:"type:char(36)"`        // 主贴的id
	Content   string    `json:"content" gorm:"type:text;not null"`   // 内容
	ResLong   string    `json:"res_long" gorm:"type:text;not null"`  // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text;not null"` // 备用短文本
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`    // 创建时间
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`    // 更新时间
}

// @title    BeforeCreate
// @description   计算出一个uint
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (thread *Thread) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.NewV4())
}
