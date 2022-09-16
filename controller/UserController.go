// @Title  UserController
// @Description  该文件用于提供操作用户的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"Essential/common"
	"Essential/dto"
	"Essential/model"
	"Essential/response"
	"Essential/util"
	"log"

	"math/rand"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// @title    Register
// @description   用户注册
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Register(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	// TODO 获取参数
	email := requestUser.Email
	password := requestUser.Password
	name := requestUser.Name
	// TODO 数据验证
	if !util.VerifyEmailFormat(email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, 201, 201, nil, "密码不能少于6位")
		return
	}

	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, email, password)

	// TODO 判断email是否存在
	if isEmailExist(DB, email) {
		response.Response(ctx, 201, 201, nil, "用户已经存在")
		return
	}

	// TODO 创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:     name,
		Email:    email,
		Password: string(hasedPassword),
		Icon:     "MGA" + strconv.Itoa(rand.Intn(9)+1) + ".jpg",
	}
	DB.Create(&newUser)
	// TODO 发放token给前端
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	// TODO 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// @title    Login
// @description   用户登录
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Login(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	// TODO 获取参数
	email := requestUser.Email
	password := requestUser.Password
	// TODO 数据验证
	if !util.VerifyEmailFormat(email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, 201, 201, nil, "密码不能少于6位")
		return
	}
	// TODO 判断邮箱是否存在
	var user model.User
	DB.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		response.Response(ctx, 201, 201, nil, "用户不存在")
		return
	}
	// TODO 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(ctx, nil, "密码错误")
		return
	}
	// TODO 发放token给前端
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	// TODO 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// @title    isEmailExist
// @description   查看email是否在数据库中存在
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func isEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	return user.ID != 0
}

// @title    Info
// @description   解析上下文中的token并返回user
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}
