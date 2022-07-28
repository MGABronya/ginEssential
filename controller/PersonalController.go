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

func PersonalPage(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")
	db := common.GetDB()
	user := tuser.(model.User)
	var articles []model.Article
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&articles)

	var posts []model.Post
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&posts)
	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user), "articles": articles, "posts": posts}})
}

func PersonalUpdate(ctx *gin.Context) {
	db := common.GetDB()
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	var personalChange vo.PersonalChange
	// 数据验证
	if err := ctx.ShouldBind(&personalChange); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}
	if !util.VerifyEmailFormat(personalChange.Email) {
		response.Response(ctx, 201, 201, nil, "邮箱格式错误")
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

	user.Name = personalChange.Name
	user.Email = personalChange.Email
	user.Telephone = personalChange.Telephone
	user.Blog = personalChange.Blog
	user.QQ = personalChange.QQ
	user.Sex = personalChange.Sex
	user.Address = personalChange.Address
	user.Hobby = personalChange.Hobby

	db.Save(&user)
	db.Model(&model.Article{}).Where("user_id = ?", user.ID).Updates(model.Article{Name: user.Name, Email: user.Email})
	db.Model(&model.Post{}).Where("user_id = ?", user.ID).Updates(model.Post{Name: user.Name, Email: user.Email})
	db.Model(&model.Thread{}).Where("user_id = ?", user.ID).Updates(model.Thread{Name: user.Name, Email: user.Email})

	response.Success(ctx, gin.H{"user": dto.ToUserDto(user)}, "更新成功")
}

func PersonalIcon(ctx *gin.Context) {
	db := common.GetDB()
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	file, err := ctx.FormFile("file")
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
	if _, ok := allowExtMap[extName]; !ok {
		response.Fail(ctx, nil, "文件格式有误")
		return
	}

	if !util.VerifyIconFormat(user.Icon) { //非默认图片，则删除原本地文件
		if err := os.Remove("./Icon/" + user.Icon); err != nil {
			//如果删除失败则输出 file remove Error!
			fmt.Println("file remove Error!")
			//输出错误详细信息
			fmt.Printf("%s", err)
		} else {
			//如果删除成功则输出 file remove OK!
			fmt.Print("file remove OK!")
		}
	}
	file.Filename = strconv.Itoa(int(user.ID)) + extName

	ctx.SaveUploadedFile(file, "./Icon/"+file.Filename) //将文件存入本地

	db.Where("id = ?", user.ID).Take(&user)

	user.Icon = file.Filename

	db.Model(&model.Article{}).Where("user_id = ?", user.ID).Update("icon", user.Icon)
	db.Model(&model.Post{}).Where("user_id = ?", user.ID).Update("name", user.Icon)
	db.Model(&model.Thread{}).Where("user_id = ?", user.ID).Update("name", user.Icon)

	db.Save(&user)

	response.Success(ctx, gin.H{"Icon": user.Icon}, "更新成功")
}

func PersonalShow(ctx *gin.Context) {
	db := common.GetDB()
	var user model.User
	if db.Where("id = ?", ctx.Params.ByName("id")).First(&user).RecordNotFound() {
		response.Fail(ctx, nil, "用户不存在")
		return
	}
	var articles []model.Article
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&articles)

	var posts []model.Post
	db.Order("created_at desc").Where("user_id = ?", user.ID).Find(&posts)

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user), "articles": articles, "posts": posts}})
}
