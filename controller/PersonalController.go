// @Title  PersonalController
// @Description  该文件用于提供操作个人界面的各种函数
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	Buil "Blog/util"
	"fmt"
	"ginEssential/common"
	"ginEssential/dto"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/util"
	"ginEssential/vo"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @title    PersonalPage
// @description   提供用户的个人信息
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalPage(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")

	user := tuser.(model.User)

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user)}})
}

// @title    PersonalPageArticles
// @description   提供用户的个人文章信息
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalPageArticles(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	db := common.GetDB()

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articles []model.Article
	db.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	// TODO 记录的总条数
	var total int
	db.Where("user_id = ?", user.ID).Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    PersonalPagePosts
// @description   提供用户的个人帖子信息
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalPagePosts(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	db := common.GetDB()

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var posts []model.Post
	db.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	// TODO 记录的总条数
	var total int
	db.Where("user_id = ?", user.ID).Model(model.Post{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"posts": posts, "total": total}, "成功")
}

// @title    PersonalPageThreads
// @description   提供用户的个人帖跟帖信息
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalPageThreads(ctx *gin.Context) {
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	db := common.GetDB()

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var threads []model.Thread
	db.Where("user_id = ?", user.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&threads)

	// TODO 记录的总条数
	var total int
	db.Where("user_id = ?", user.ID).Model(model.Thread{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"threads": threads, "total": total}, "成功")
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

	if user.Name != personalChange.Name {
		// TODO 重命名用户文件夹
		folder := `../../Blog/dist`
		err2 := os.Rename(folder, personalChange.Name)
		if err2 != nil {
			panic(err2)
		} else {
			println(`文件夹重命名成功`)
		}
	}

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

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user)}})
}

// @title    PersonalShowArticles
// @description   查看用户文章信息
// @auth      MGAronya（张健）       2022-9-16 12:31
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalShowArticles(ctx *gin.Context) {
	db := common.GetDB()
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	userId := ctx.Params.ByName("id")

	// TODO 查看用户是否存在
	if db.Where("id = ?", ctx.Params.ByName("id")).First(&model.User{}).RecordNotFound() {
		response.Fail(ctx, nil, "用户不存在")
		return
	}

	// TODO 分页
	var articles []model.Article
	var level int8
	if strconv.Itoa(int(user.ID)) == userId {
		level = 4
	} else if Buil.IsS(4, "Fr"+strconv.Itoa(int(user.ID)), userId) {
		level = 3
	} else {
		level = 2
	}

	// TODO 查找所有分页中可见的条目
	db.Where("user_id = ? and visible < ?", user.ID, level).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	// TODO 记录的总条数
	var total int
	db.Where("user_id = ? and visible < ?", user.ID, level).Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    PersonalShowPosts
// @description   查看用户帖子信息
// @auth      MGAronya（张健）       2022-9-16 12:31
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalShowPosts(ctx *gin.Context) {
	db := common.GetDB()
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	userId := ctx.Params.ByName("id")

	// TODO 查看用户是否存在
	if db.Where("id = ?", ctx.Params.ByName("id")).First(&model.User{}).RecordNotFound() {
		response.Fail(ctx, nil, "用户不存在")
		return
	}

	// TODO 分页
	var posts []model.Post
	var level int8
	if strconv.Itoa(int(user.ID)) == userId {
		level = 4
	} else if Buil.IsS(4, "Fr"+strconv.Itoa(int(user.ID)), userId) {
		level = 3
	} else {
		level = 2
	}

	// TODO 查找所有分页中可见的条目
	db.Where("user_id = ? and visible < ?", user.ID, level).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	// TODO 记录的总条数
	var total int
	db.Where("user_id = ? and visible < ?", user.ID, level).Model(model.Post{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"posts": posts, "total": total}, "成功")
}

// @title    PersonalShowThreads
// @description   查看用户信息
// @auth      MGAronya（张健）       2022-9-16 12:31
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalShowThreads(ctx *gin.Context) {
	db := common.GetDB()
	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	userId := ctx.Params.ByName("id")

	// TODO 查看用户是否存在
	if db.Where("id = ?", ctx.Params.ByName("id")).First(&model.User{}).RecordNotFound() {
		response.Fail(ctx, nil, "用户不存在")
		return
	}

	// TODO 分页
	var threads []model.Thread
	var level int8
	if strconv.Itoa(int(user.ID)) == userId {
		level = 4
	} else if Buil.IsS(4, "Fr"+strconv.Itoa(int(user.ID)), userId) {
		level = 3
	} else {
		level = 2
	}

	// TODO 查找所有分页中可见的条目
	db.Table("threads").Joins("join posts on threads.post_id = posts.id").Where("posts.user_id = ? and posts.visible < ?", user.ID, level).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&threads)

	// TODO 记录的总条数
	var total int
	db.Table("threads").Joins("join posts on threads.post_id = posts.id").Where("posts.user_id = ? and posts.visible < ?", user.ID, level).Model(model.Thread{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"threads": threads, "total": total}, "成功")
}

// @title    PersonalShowUsers
// @description   查看一组用户信息
// @auth      MGAronya（张健）       2022-9-16 12:31
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func PersonalShowUsers(ctx *gin.Context) {
	db := common.GetDB()

	var requestUsers vo.UsersRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestUsers); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 查询与之对应的用户列表
	users := make([]dto.UserDto, len(requestUsers.UserId))
	for i, userId := range requestUsers.UserId {
		var user model.User
		db.Where("id = ?", userId).First(&user)
		users[i] = dto.ToUserDto(user)
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"users": users}})
}
