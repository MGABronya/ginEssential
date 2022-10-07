// @Title  post
// @Description  定义帖子
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Post			定义了帖子
type Post struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`                                 // 帖子的id
	UserId    uint      `json:"user_id" gorm:"index:idx_userId;not null"`                            // 作者的id
	Title     string    `json:"title" gorm:"type:char(50);not null;index:search_idx,class:FULLTEXT"` // 停止的标题
	Content   string    `json:"content" gorm:"type:text;not null;index:search_idx,class:FULLTEXT"`   // 停止的内容
	ResLong   string    `json:"res_long" gorm:"type:text;index:search_idx,class:FULLTEXT"`           // 备用长文本
	ResShort  string    `json:"res_short" gorm:"type:text;index:search_idx,class:FULLTEXT"`          // 备用短文本
	Visible   int8      `json:"visible" gorm:"type:tinyint;default:1"`                               // 可见等级
	CreatedAt Time      `json:"created_at" gorm:"type:timestamp"`                                    // 帖子的创建时间
	UpdatedAt Time      `json:"updated_at" gorm:"type:timestamp"`                                    // 帖子的更新时间
}

// @title    BeforeCreate
// @description   计算出一个uint
// @auth      MGAronya（张健）             2022-9-16 10:19
// @param     scope *gorm.Scope
// @return    error
func (post *Post) BeforeCreate(scope *gorm.DB) error {
	post.ID = uuid.NewV4()
	return nil
}
