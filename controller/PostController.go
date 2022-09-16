// @Title  PostController
// @Description  该文件用于提供操作帖子的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"Essential/common"
	"Essential/model"
	"Essential/response"
	"Essential/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// IPostController			定义了帖子类接口
type IPostController interface {
	RestController                 // 实现增删查改功能
	PageList(ctx *gin.Context)     // 实现返回帖子列表
	ThreadCreate(cts *gin.Context) // 实现创建跟帖功能
}

// PostController			定义了帖子工具类
type PostController struct {
	DB *gorm.DB //包含一个数据库指针
}

// @title    Create
// @description   创建一篇帖子
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	user, _ := ctx.Get("user")

	// TODO 创建Post
	post := model.Post{
		UserId:  user.(model.User).ID,
		Title:   requestPost.Title,
		Content: requestPost.Content,
		Icon:    user.(model.User).Icon,
		Name:    user.(model.User).Name,
		Email:   user.(model.User).Email,
	}

	// TODO 插入数据
	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    ThreadCreate
// @description   创建一篇跟帖
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) ThreadCreate(ctx *gin.Context) {
	var requestThread vo.CreateThreadRequest

	// TODO 数据验证
	if err := ctx.ShouldBind(&requestThread); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	var post model.Post

	// TODO 查看帖子是否存在
	if p.DB.Where("id = ?", ctx.Params.ByName("id")).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	// TODO 获取登录用户
	user, _ := ctx.Get("user")

	// TODO 创建Thread
	thread := model.Thread{
		UserId:  user.(model.User).ID,
		PostId:  ctx.Params.ByName("id"),
		Content: requestThread.Content,
		Icon:    user.(model.User).Icon,
		Name:    user.(model.User).Name,
		Email:   user.(model.User).Email,
	}

	// TODO 插入数据
	if err := p.DB.Create(&thread).Error; err != nil {
		panic(err)
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Update
// @description   更新帖子内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Update(ctx *gin.Context) {
	var requestPT vo.PTRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestPT); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID

	// TODO 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	// TODO 判断是更新帖子还是更新更贴
	if requestPT.PT {

		// TODO 判断当前用户是否为帖子的作者
		if userId != post.UserId {
			response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
			return
		}

		var requestPost vo.CreatePostRequest
		requestPost.Content = requestPT.Content
		requestPost.Title = requestPT.Title

		// TODO 更新帖子
		if err := p.DB.Model(&post).Update(requestPost).Error; err != nil {
			response.Fail(ctx, nil, "更新失败")
			return
		}

		response.Success(ctx, gin.H{"post": post}, "更新成功")
	} else {
		var thread model.Thread
		// TODO 判断当前用户是否为跟帖的作者
		if p.DB.Where("id = ?", requestPT.ThreadId).First(&thread).RecordNotFound() {
			response.Fail(ctx, nil, "跟帖不存在")
			return
		}
		if userId != thread.UserId {
			response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
			return
		}

		var requestThread vo.CreateThreadRequest
		requestThread.Content = requestPT.Content

		// TODO 更新帖子
		if err := p.DB.Model(&thread).Update(requestThread).Error; err != nil {
			response.Fail(ctx, nil, "更新失败")
			return
		}

		response.Success(ctx, gin.H{"thread": thread}, "更新成功")
	}

}

// @title    Show
// @description   查看帖子内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	// TODO 查看帖子是否存在
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	var threads []model.Thread
	p.DB.Order("created_at asc").Where("post_id = ?", post.ID).Find(&threads)

	// TODO 返回数据
	response.Success(ctx, gin.H{"post": post, "threads": threads}, "成功")
}

// @title    Delete
// @description   删除帖子或跟贴
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) Delete(ctx *gin.Context) {

	var requestPT vo.PTRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestPT); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID

	// TODO 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	// TODO 判断是删除帖子还是跟帖
	if requestPT.PT {

		// TODO 判断当前用户是否为帖子的作者
		if userId != post.UserId {
			response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
			return
		}

		// TODO 删除帖子
		p.DB.Where(model.Thread{PostId: post.ID.String()}).Delete(model.Thread{})
		p.DB.Delete(&post)

		response.Success(ctx, gin.H{"post": post}, "删除成功")
	} else {
		var thread model.Thread
		// TODO 判断当前用户是否为跟帖的作者
		if p.DB.Where("id = ?", requestPT.ThreadId).First(&thread).RecordNotFound() {
			response.Fail(ctx, nil, "跟帖不存在")
			return
		}
		if userId != thread.UserId {
			response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
			return
		}

		// TODO 删除帖子
		p.DB.Delete(&thread)

		response.Success(ctx, gin.H{"thread": thread}, "删除成功")
	}
}

// @title    PageList
// @description   返回帖子的列表
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (p PostController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	// TODO 记录的总条数
	var total int
	p.DB.Model(model.Post{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"data": posts, "total": total}, "成功")
}

// @title    NewPostController
// @description   新建一个帖子的控制器
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	db.AutoMigrate(model.Thread{})
	return PostController{DB: db}
}
