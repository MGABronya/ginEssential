// @Title  user_dto
// @Description  用于封装用户信息
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package dto

import (
	"Essential/model"
)

// UserDto			定义了用户的基本信息
type UserDto struct {
	Name      string `gorm:"type:varchar(20);not null"`        // 用户名称
	Email     string `gorm:"type:varchar(50);not null;unique"` // 邮箱
	Telephone string `gorm:"type:varchar(20)"`                 // 手机号
	Blog      string `gorm:"type:varchar(25)"`                 // 博客地址
	QQ        string `gorm:"type:varchar(20)"`                 // QQ账号
	Sex       bool   `gorm:"type:boolean"`                     // 性别
	Address   string `gorm:"type:varchar(20)"`                 // 地址
	Hobby     string `gorm:"type:varchar(50)"`                 // 爱好
	Icon      string `gorm:"type:varchar(50)"`                 // 这里的Icon存储的是图像文件的地址
}

// @title    ToUserDto
// @description   用户信息封装
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    user model.User       接收一个用户类
// @return   UserDto			   返回一个用户的基本信息类
func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name:      user.Name,
		Email:     user.Email,
		Telephone: user.Telephone,
		Blog:      user.Blog,
		QQ:        user.QQ,
		Sex:       user.Sex,
		Address:   user.Address,
		Hobby:     user.Hobby,
		Icon:      user.Icon,
	}
}
