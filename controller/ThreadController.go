// @Title  ThreadController
// @Description  该文件用于提供操作跟帖的各种方法
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

// IThreadController			定义了跟帖类接口
type IThreadController interface {
	RestController             // 实现增删查改功能
	PageList(ctx *gin.Context) // 实现返回跟帖列表
}

// ThreadController			定义了跟帖工具类
type ThreadController struct {
	DB *gorm.DB //包含一个数据库指针
}

// @title    Create
// @description   创建一篇跟帖
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Create(ctx *gin.Context) {
	var requestThread vo.CreateThreadRequest

	// TODO 数据验证
	if err := ctx.ShouldBind(&requestThread); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	var post model.Post

	// TODO 查看帖子是否存在
	if t.DB.Where("id = ?", ctx.Params.ByName("id")).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	// TODO 获取登录用户
	user, _ := ctx.Get("user")

	// TODO 创建Thread
	thread := model.Thread{
		UserId:   user.(model.User).ID,
		PostId:   ctx.Params.ByName("id"),
		Content:  requestThread.Content,
		ResLong:  requestThread.ResLong,
		ResShort: requestThread.ResShort,
	}

	// TODO 插入数据
	if err := t.DB.Create(&thread).Error; err != nil {
		panic(err)
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
}

// @title    Thread
// @description   更新跟帖内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Update(ctx *gin.Context) {
	var requestThread vo.CreateThreadRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestThread); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID

	// TODO 获取path中的id
	threadId := ctx.Params.ByName("id")

	var thread model.Thread
	if t.DB.Where("id = ?", threadId).First(&thread).RecordNotFound() {
		response.Fail(ctx, nil, "跟帖不存在")
		return
	}

	// TODO 判断当前用户是否为跟帖的作者
	if userId != thread.UserId {
		response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
		return
	}

	// TODO 更新帖子
	if err := t.DB.Model(&thread).Update(requestThread).Error; err != nil {
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"thread": thread}, "更新成功")
}

// @title    Show
// @description   查看跟帖内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	threadId := ctx.Params.ByName("id")

	var thread model.Thread
	// TODO 查看跟帖是否存在
	if t.DB.Where("id = ?", threadId).First(&thread).RecordNotFound() {
		response.Fail(ctx, nil, "跟帖不存在")
		return
	}

	// TODO 返回数据
	response.Success(ctx, gin.H{"thread": thread}, "成功")
}

// @title    Delete
// @description   删除跟贴
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) Delete(ctx *gin.Context) {

	// TODO 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID

	// TODO 获取path中的id
	threadId := ctx.Params.ByName("id")

	var thread model.Thread
	if t.DB.Where("id = ?", threadId).First(&thread).RecordNotFound() {
		response.Fail(ctx, nil, "跟帖不存在")
		return
	}

	if Buil.GetH(0, "permission", strconv.Itoa(int(userId)))[0] < '4' && userId != thread.UserId {
		response.Fail(ctx, nil, "帖子不属于您，请勿非法操作")
		return
	}

	// TODO 删除帖子
	t.DB.Delete(&thread)

	// TODO 移除收藏
	for _, val := range Buil.MembersS(3, "taF"+threadId) {
		Buil.RemS(3, "tFa"+val, threadId)
	}
	Buil.Del(3, "taF"+threadId)

	// TODO 移除点赞
	Buil.Del(3, "tiL"+threadId)

	response.Success(ctx, gin.H{"thread": thread}, "删除成功")
}

// @title    PageList
// @description   返回跟帖的列表
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (t ThreadController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 获取path中的id
	postId := ctx.Params.ByName("id")

	var post model.Post

	// TODO 查看帖子是否存在
	if t.DB.Where("id = ?", postId).First(&post).RecordNotFound() {
		response.Fail(ctx, nil, "帖子不存在")
		return
	}

	// TODO 获取登录用户
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看是否有权限
	if post.UserId != user.ID && (post.Visible == 3 || (post.Visible == 2 && !Buil.IsS(4, "Fr"+strconv.Itoa(int(user.ID)), strconv.Itoa(int(post.UserId))))) {
		response.Fail(ctx, nil, "权限不足，不能查看")
		return
	}

	// TODO 分页
	var threads []model.Thread
	t.DB.Where("post_id = ?", postId).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&threads)

	// TODO 记录的总条数
	var total int
	t.DB.Model(model.Thread{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"threads": threads, "total": total}, "成功")
}

// @title    NewThreadController
// @description   新建一个跟帖的控制器
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func NewThreadController() IThreadController {
	db := common.GetDB()
	db.AutoMigrate(model.Thread{})
	return ThreadController{DB: db}
}
