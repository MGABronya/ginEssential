package dto

import (
	"Essential/model"
)

type UserDto struct {
	Name      string `gorm:"type:varchar(20);not null"`
	Email     string `gorm:"type:varchar(50);not null;unique"`
	Telephone string `gorm:"type:varchar(20)"`
	Blog      string `gorm:"type:varchar(25)"`
	QQ        string `gorm:"type:varchar(20)"`
	Sex       bool   `gorm:"type:boolean"`
	Address   string `gorm:"type:varchar(20)"`
	Hobby     string `gorm:"type:varchar(50)"`
	Icon      string `gorm:"type:varchar(50)"` //这里的Icon存储的是图像文件的地址
}

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

/*
func test(){
	r := gin.Default()
	r.POST("/test/load", func(c *gin.Context){
		file, _ := c.FormFile("file")
		name := c.PostForm("name")                            //其它需要的传入信息
		c.SaveUploadedFile(file, "./Icon/" + file.Filename)   //将文件存入本地
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
		c.File("./Icon/" + file.Filename)             //以上为返回文件
	})
}*/
