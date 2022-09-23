// @Title  PersonalController
// @Description  该文件用于提供操作个人界面的各种函数
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"Essential/common"
	"Essential/dto"
	"Essential/model"
	"Essential/response"
	"Essential/util"
	"Essential/vo"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @title    PersonalPage
// @description   提供用户的个人页面信息
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalPage(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")
	db := common.GetDB()

	user := tuser.(model.User)

	var articles []model.Article

	// TODO 查看用户的所有文章
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&articles)

	var posts []model.Post

	// TODO 查看用户的所有帖子
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&posts)

	var threads []model.Thread

	// TODO 查看用户的所有跟帖
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&threads)

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user), "articles": articles, "posts": posts, "threads": threads}})
}

// @title    PersonalUpdate
// @description   个人页面更新
// @auth      MGAronya（张健）       2022-9-16 12:31
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalUpdate(ctx *gin.Context) {
	db := common.GetDB()
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	var personalChange vo.PersonalChange
	// TODO 数据验证
	if err := ctx.ShouldBind(&personalChange); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	if !util.VerifyEmailFormat(personalChange.Email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
		return
	}
	if personalChange.Email != user.Email && util.IsEmailExist(db, personalChange.Email) {
		response.Response(ctx, 201, 201, nil, "邮箱重复")
		return
	}
	if personalChange.Name != user.Name && util.IsNameExist(db, personalChange.Name) {
		response.Response(ctx, 201, 201, nil, "用户名已被注册")
		return
	}
	if personalChange.Telephone != "" && !util.VerifyMobileFormat(personalChange.Telephone) {
		response.Response(ctx, 201, 201, nil, "手机号格式错误")
		return
	}
	if personalChange.QQ != "" && !util.VerifyQQFormat(personalChange.QQ) {
		response.Response(ctx, 201, 201, nil, "QQ格式错误")
		return
	}

	db.Where("id = ?", user.ID).Take(&user)

	// TODO 更新用户信息
	user.Name = personalChange.Name
	user.Email = personalChange.Email
	user.Telephone = personalChange.Telephone
	user.Blog = personalChange.Blog
	user.QQ = personalChange.QQ
	user.Sex = personalChange.Sex
	user.Address = personalChange.Address
	user.Hobby = personalChange.Hobby

	db.Save(&user)

	response.Success(ctx, gin.H{"user": dto.ToUserDto(user)}, "更新成功")
}

// @title    PersonalIcon
// @description   个人头像图片上传
// @auth      MGAronya（张健）       2022-9-16 12:31
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalIcon(ctx *gin.Context) {
	db := common.GetDB()
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	file, err := ctx.FormFile("file")

	//TODO 数据验证
	if err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}

	// TODO 格式验证
	if _, ok := allowExtMap[extName]; !ok {
		response.Fail(ctx, nil, "文件格式有误")
		return
	}

	// TODO 非默认图片，则删除原本地文件
	if !util.VerifyIconFormat(user.Icon) {
		if err := os.Remove("./Icon/" + user.Icon); err != nil {
			// TODO 如果删除失败则输出 file remove Error!
			fmt.Println("file remove Error!")
			// TODO 输出错误详细信息
			fmt.Printf("%s", err)
		} else {
			// TODO 如果删除成功则输出 file remove OK!
			fmt.Print("file remove OK!")
		}
	}
	file.Filename = strconv.Itoa(int(user.ID)) + extName

	// TODO 将文件存入本地
	ctx.SaveUploadedFile(file, "./Icon/"+file.Filename)

	db.Where("id = ?", user.ID).Take(&user)

	user.Icon = file.Filename

	// TODO 更新用户简易信息
	db.Model(&model.Article{}).Where("user_id = ?", user.ID).Update("icon", user.Icon)
	db.Model(&model.Post{}).Where("user_id = ?", user.ID).Update("name", user.Icon)
	db.Model(&model.Thread{}).Where("user_id = ?", user.ID).Update("name", user.Icon)

	db.Save(&user)

	response.Success(ctx, gin.H{"Icon": user.Icon}, "更新成功")
}

// @title    PersonalShow
// @description   查看用户信息
// @auth      MGAronya（张健）       2022-9-16 12:31
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalShow(ctx *gin.Context) {
	db := common.GetDB()
	var user model.User

	// TODO 查看用户是否存在
	if db.Where("id = ?", ctx.Params.ByName("id")).First(&user).RecordNotFound() {
		response.Fail(ctx, nil, "用户不存在")
		return
	}

	// TODO 查询用户的文章、帖子、跟帖
	var articles []model.Article
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&articles)

	var posts []model.Post
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&posts)

	var threads []model.Thread
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&threads)

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user), "articles": articles, "posts": posts, "threads": threads}})
}
