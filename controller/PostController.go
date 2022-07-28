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

type IPostController interface {
	RestController
	PageList(ctx *gin.Context)
	ThreadCreate(cts *gin.Context)
}

type PostController struct {
	DB *gorm.DB
}

func (p PostController) Create(ctx *gin.Context) {
	var requestPost vo.CreatePostRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")

	// 创建Post
	post := model.Post{
		UserId:  user.(model.User).ID,
		Title:   requestPost.Title,
		Content: requestPost.Content,
		Icon:    user.(model.User).Icon,
		Name:    user.(model.User).Name,
		Email:   user.(model.User).Email,
	}

	// 插入数据
	if err := p.DB.Create(&post).Error; err != nil {
		panic(err)
	}

	// 成功
	response.Success(ctx, nil, "创建成功")
}

func (p PostController) ThreadCreate(ctx *gin.Context) {
	var requestThread vo.CreateThreadRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestThread); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	var post model.Post
	if p.DB.Where("id = ?", ctx.Params.ByName("id")).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")

	// 创建Thread
	thread := model.Thread{
		UserId:  user.(model.User).ID,
		PostId:  ctx.Params.ByName("id"),
		Content: requestThread.Content,
		Icon:    user.(model.User).Icon,
		Name:    user.(model.User).Name,
		Email:   user.(model.User).Email,
	}

	// 插入数据
	if err := p.DB.Create(&thread).Error; err != nil {
		panic(err)
	}

	// 成功
	response.Success(ctx, nil, "创建成功")
}

func (p PostController) Update(ctx *gin.Context) {
	var requestPT vo.PTRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestPT); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID

	// 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	if requestPT.PT {

		// 判断当前用户是否为帖子的作者
		if userId != post.UserId {
			response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
			return
		}

		var requestPost vo.CreatePostRequest
		requestPost.Content = requestPT.Content
		requestPost.Title = requestPT.Title

		// 更新帖子
		if err := p.DB.Model(&post).Update(requestPost).Error; err != nil {
			response.Fail(ctx, nil, "更新失败")
			return
		}

		response.Success(ctx, gin.H{"post": post}, "更新成功")
	} else {
		var thread model.Thread
		// 判断当前用户是否为跟帖的作者
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

		// 更新帖子
		if err := p.DB.Model(&thread).Update(requestThread).Error; err != nil {
			response.Fail(ctx, nil, "更新失败")
			return
		}

		response.Success(ctx, gin.H{"thread": thread}, "更新成功")
	}

}

func (p PostController) Show(ctx *gin.Context) {
	// 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	var threads []model.Thread
	p.DB.Where("post_id = ?", post.ID).Find(&threads)

	// 返回数据
	response.Success(ctx, gin.H{"post": post, "threads": threads}, "成功")
}

func (p PostController) Delete(ctx *gin.Context) {

	var requestPT vo.PTRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestPT); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID

	// 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post
	if p.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	if requestPT.PT {

		// 判断当前用户是否为帖子的作者
		if userId != post.UserId {
			response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
			return
		}

		// 删除帖子
		p.DB.Where(model.Thread{PostId: post.ID.String()}).Delete(model.Thread{})
		p.DB.Delete(&post)

		response.Success(ctx, gin.H{"post": post}, "删除成功")
	} else {
		var thread model.Thread
		// 判断当前用户是否为跟帖的作者
		if p.DB.Where("id = ?", requestPT.ThreadId).First(&thread).RecordNotFound() {
			response.Fail(ctx, nil, "跟帖不存在")
			return
		}
		if userId != thread.UserId {
			response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
			return
		}

		// 删除帖子
		p.DB.Delete(&thread)

		response.Success(ctx, gin.H{"thread": thread}, "删除成功")
	}
}

func (p PostController) PageList(ctx *gin.Context) {
	// 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 分页
	var posts []model.Post
	p.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&posts)

	//var user []model.User

	//a.DB.Where("id = ?", article.UserId).First(&user)

	// 记录的总条数
	var total int
	p.DB.Model(model.Post{}).Count(&total)

	//"user": dto.ToUserDto(user.(model.User))
	// 返回数据
	response.Success(ctx, gin.H{"data": posts, "total": total}, "成功")
}

func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	db.AutoMigrate(model.Thread{})
	return PostController{DB: db}
}
