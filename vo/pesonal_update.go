package vo

import "github.com/jinzhu/gorm"

type PersonalChange struct {
	gorm.Model
	Name      string `json:"Newname" gorm:"type:varchar(20);"`
	Email     string `json:"Newemail" gorm:"type:varchar(50);unique"`
	Telephone string `json:"Newtelephone" gorm:"type:varchar(20)"`
	Blog      string `json:"Newblog" gorm:"type:varchar(25)"`
	QQ        string `json:"Newqq" gorm:"type:varchar(20)"`
	Sex       bool   `json:"Newsex" gorm:"type:boolean"`
	Address   string `json:"Newaddress" gorm:"type:varchar(20)"`
	Hobby     string `json:"Newhobby" gorm:"type:varchar(50)"`
}
