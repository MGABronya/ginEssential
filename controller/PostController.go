// @Title  PostController
// @Description  该文件用于提供操作帖子的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	Buil "Blog/util"
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/vo"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// IPostController			定义了帖子类接口
type IPostController interface {
	RestController             // 实现增删查改功能
	PageList(ctx *gin.Context) // 实现返回帖子列表
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
		UserId:   user.(model.User).ID,
		Title:    requestPost.Title,
		Content:  requestPost.Content,
		ResLong:  requestPost.ResLong,
		ResShort: requestPost.ResShort,
	}

	// TODO 插入数据
	if err := p.DB.Create(&post).Error; err != nil {
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
	var requestPost vo.CreatePostRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestPost); err != nil {
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

	// TODO 判断当前用户是否为帖子的作者
	if userId != post.UserId {
		response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
		return
	}

	// TODO 更新帖子
	if err := p.DB.Model(&post).Update(requestPost).Error; err != nil {
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"post": post}, "更新成功")
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

	// TODO 判断当前用户是否为帖子的作者
	if Buil.GetH(0, "permission", strconv.Itoa(int(userId)))[0] < '4' && userId != post.UserId {
		response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
		return
	}

	// TODO 移除收藏
	for _, val := range Buil.MembersS(3, "paF"+postId) {
		Buil.RemS(3, "pFa"+val, postId)
	}
	Buil.Del(3, "paF"+postId)

	// TODO 移除点赞
	Buil.Del(3, "piL"+postId)

	// TODO 移除标签
	for _, val := range Buil.MembersS(3, "paL"+postId) {
		Buil.RemS(3, "pLa"+val, postId)
	}
	Buil.Del(3, "paL"+postId)

	// TODO 移除跟帖的相关属性
	var threads []model.Thread
	p.DB.Where("post_id = ?", postId).Find(&threads)

	for _, thread := range threads {
		// TODO 移除收藏
		for _, val := range Buil.MembersS(3, "taF"+thread.ID.String()) {
			Buil.RemS(3, "tFa"+val, thread.ID.String())
		}
		Buil.Del(3, "taF"+thread.ID.String())

		// TODO 移除点赞
		Buil.Del(3, "tiL"+thread.ID.String())
	}

	// TODO 删除帖子
	p.DB.Where(model.Thread{PostId: post.ID.String()}).Delete(model.Thread{})
	p.DB.Delete(&post)

	response.Success(ctx, gin.H{"post": post}, "删除成功")
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
	response.Success(ctx, gin.H{"posts": posts, "total": total}, "成功")
}

// @title    NewPostController
// @description   新建一个帖子的控制器
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func NewPostController() IPostController {
	db := common.GetDB()
	db.AutoMigrate(model.Post{})
	return PostController{DB: db}
}
