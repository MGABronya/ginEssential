package main

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"time"

	"github.com/jordan-wright/email"
)

var mod string = `
			尊敬的%s，您好！

			您于 %s 提交的邮箱验证，本次验证码为%s，为了保证账号安全，验证码有效期为5分钟。请确认为本人操作，切勿向他人泄露，感谢您的理解与使用。
			此邮箱为系统邮箱，请勿回复。

`

func SendEmailValidate(em []string) (string, error) {
	e := email.NewEmail()
	e.From = "mgAronya <2829214609@qq.com>"
	e.To = em
	// 生成6位随机验证码
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	vCode := fmt.Sprintf("%06v", rnd.Int31n(1000000))
	t := time.Now().Format("2006-01-02 15:04:05")
	//设置文件发送的内容
	content := fmt.Sprintf(mod, em[0], t, vCode)
	e.Text = []byte(content)
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:25", smtp.PlainAuth("", "2829214609@qq.com", "rmdtxokuuqyrdgii", "smtp.qq.com"))
	return vCode, err
}

/*
// GetValidateCode
// @Title GetValidateCode
// @Description  发送邮箱验证码 并存入redis（5分钟有效时间）
// @Author hyy 2022-03-05 18:18:47
// @Param c type description
func GetValidateCode(ctx *gin.Context) {
	// 获取目的邮箱
	em := []string{ctx.Param("email")}
	vCode, err := SendEmailValidate(em)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":           400,
			"msg":              "验证码发送失败",
			"ERROR-CONTROLLER": err.Error(),
		})
		return
	}
	// 验证码存入redis 并设置过期时间5分钟
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6380",
		Password: "",
		DB:       0,
	})
	_, err = client.Expire("key", 60*time.Second).Result()
	_, err = util.RedisCont.Do("set", "vCode", vCode)
	_, _ = util.RedisCont.Do("expire", "vCode", 300)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":           400,
			"msg":              "验证码存储失败",
			"ERROR-CONTROLLER": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":    "验证码发送成功",
		"status": 200,
		"vCode":  vCode,
	})
	return
}

// ValidateEmailCode
// @Title ValidateEmailCode
// @Description  验证邮箱验证码，并注册用户。
// @Author hyy 2022-03-05 18:19:18
// @Param c type description
func ValidateEmailCode(c *gin.Context) {
	vCode := c.Param("vCode")
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":           400,
			"msg":              "注册失败，json解析失败",
			"ERROR-CONTROLLER": err.Error(),
		})
		return
	}
	// 默认用户权限为2
	if user.Level == 0 {
		user.Level = 2
	}
	user.CreatAt = time.Now().Format("2006-01-02 15:04:05")
	// 设置默认头像
	if user.HandPortrait == "" {
		user.HandPortrait = "https://gin-study-1308216693.cos.ap-nanjing.myqcloud.com/handportrait.png"
	}
	// 验证结构体字段
	if err := util.Validate.Struct(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":         400,
			"msg":            "注册失败,非法字段！",
			"ERROR-VALIDATE": util.GetErr(err),
		})
		return
	}
	// 对密码进行加密
	user.Password = util.Md5Encrypt(user.Password)
	// 获取存储在redis中的验证码
	vCodeRaw, err := util.RedisCont.Do("get", "vCode")
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"status":           400,
			"msg":              "Redis获取vCode失败",
			"ERROR-CONTROLLER": err.Error(),
		})
		return
	}
	if string(vCodeRaw.([]byte)) != "" && vCode == string(vCodeRaw.([]byte)) {
		result := service.InsertUser(&user)
		if result["result"].(bool) {
			result["msg"] = "注册成功！"
			c.JSON(http.StatusOK, result)
			return
		} else {
			c.JSON(http.StatusBadRequest, result)
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "验证码已失效",
		})
		return
	}
}
*/

func main() {
	SendEmailValidate([]string{"20jzhang@stu.edu.cn"})
}
