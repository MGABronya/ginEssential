// @Title  util
// @Description  收集各种需要使用的工具函数
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package util

import (
	"Essential/common"
	"Essential/model"
	"context"
	"fmt"
	"math/rand"
	"net/smtp"
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jordan-wright/email"
)

// @title    RandomString
// @description   生成一段随机的字符串
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     n int		字符串的长度
// @return    string    一串随机的字符串
func RandomString(n int) string {
	var letters = []byte("qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	// TODO 不断用随机字母填充字符串
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// @title    VerifyEmailFormat
// @description   用于验证邮箱格式是否正确的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     email string		一串字符串，表示邮箱
// @return    bool    返回是否合法
func VerifyEmailFormat(email string) bool {
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// @title    VerifyMobileFormat
// @description   用于验证手机号格式是否正确的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     mobileNum string		一串字符串，表示手机号
// @return    bool    返回是否合法
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

// @title    VerifyQQFormat
// @description   用于验证QQ号格式是否正确的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     QQNum string		一串字符串，表示QQ
// @return    bool    返回是否合法
func VerifyQQFormat(QQNum string) bool {
	regular := "[1-9][0-9]{4,10}"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(QQNum)
}

// @title    VerifyQQFormat
// @description  用于验证Icon是否为默认图片的工具函数
// @auth      MGAronya（张健）             2022-9-16 10:29
// @param     Icon string		一串字符串，表示图像名称
// @return    bool    返回是否合法
func VerifyIconFormat(Icon string) bool {
	regular := "MGA[1-9].jpg"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(Icon)
}

// @title    isEmailExist
// @description   查看email是否在数据库中存在
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func IsEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	return user.ID != 0
}

var ctx context.Context = context.Background()

// @title    SendEmailValidate
// @description   发送验证邮件
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func SendEmailValidate(em []string) (string, error) {
	mod := `
	尊敬的%s，您好！

	您于 %s 提交的邮箱验证，本次验证码为%s，为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。
	此邮箱为系统邮箱，请勿回复。

`
	e := email.NewEmail()
	e.From = "mgAronya <2829214609@qq.com>"
	e.To = em
	// TODO 生成6位随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	t := time.Now().Format("2006-01-02 15:04:05")
	// TODO 设置文件发送的内容
	content := fmt.Sprintf(mod, em[0], t, vCode)
	e.Text = []byte(content)
	// TODO 设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2829214609@qq.com", "rmdtxokuuqyrdgii", "smtp.qq.com"))
	return vCode, err
}

// @title    IsEmailPass
// @description   验证邮箱是否通过
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    em []string       接收一个邮箱字符串
// @return   string, error     返回验证码和error值
func IsEmailPass(email string, vertify string) bool {
	client := common.GetRedisClient()
	V, _ := client.Get(ctx, email).Result()
	return V == vertify
}

// @title    SetRedis
// @description   设置验证码，并令其存活五分钟
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    email string, v string       接收一个邮箱和一个验证码
// @return   void
func SetRedis(email string, v string) {
	client := common.GetRedisClient()

	client.Set(ctx, email, v, 300*time.Second)
}
