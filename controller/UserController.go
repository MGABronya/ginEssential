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

func Register(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	//获取参数
	email := requestUser.Email
	password := requestUser.Password
	name := requestUser.Name
	//数据验证
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

	//判断email是否存在
	if isEmailExist(DB, email) {
		response.Response(ctx, 201, 201, nil, "用户已经存在")
		return
	}
	//创建用户
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
	//发放token给前端
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()
	var requestUser = model.User{}
	ctx.Bind(&requestUser)
	//获取参数
	email := requestUser.Email
	password := requestUser.Password
	//数据验证
	if !util.VerifyEmailFormat(email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, 201, 201, nil, "密码不能少于6位")
		return
	}
	//判断邮箱是否存在
	var user model.User
	DB.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		response.Response(ctx, 201, 201, nil, "用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Fail(ctx, nil, "密码错误")
		return
	}
	//发放token给前端
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, 201, 201, nil, "系统异常")
		log.Printf("token generate error: %v", err)
		return
	}
	//返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

func isEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	return user.ID != 0
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}
