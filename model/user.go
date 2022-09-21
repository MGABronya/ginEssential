// @Title  user
// @Description  定义跟帖
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:46
package model

import "github.com/jinzhu/gorm"

// user			定义用户
type User struct {
	gorm.Model        // gorm的模板
	Name       string `gorm:"type:varchar(20);not null"`        // 用户名称
	Email      string `gorm:"type:varchar(50);not null;unique"` // 邮箱
	Password   string `gorm:"size:255;not null"`                // 密码
	Icon       string `gorm:"type:varchar(50)"`                 // 这里的Icon存储的是图像文件的地址
	Telephone  string `gorm:"type:varchar(20)"`                 // 电话
	Blog       string `gorm:"type:varchar(25)"`                 // 博客
	QQ         string `gorm:"type:varchar(20)"`                 // QQ
	Sex        bool   `gorm:"type:boolean"`                     // 性别
	Address    string `gorm:"type:varchar(20)"`                 // 地址
	Hobby      string `gorm:"type:varchar(50)"`                 // 爱好
	BackGround string `gorm:"type:varchar(50)"`                 // 这里的BackGround存储的是图像文件的地址
}
