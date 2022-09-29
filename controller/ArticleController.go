// @Title  ArticleController
// @Description  该文件提供关于操作文章的各种方法
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

// IArticleController			定义了文章类接口
type IArticleController interface {
	RestController             // 包含增删查改功能
	PageList(ctx *gin.Context) //能够提供文章列表
}

// IArticleController			定义了文章工具类
type ArticleController struct {
	DB *gorm.DB // 含有一个数据库指针
}

// @title    Create
// @description   创建一篇文章
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Create(ctx *gin.Context) {
	var requestArticle vo.CreateArticleRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestArticle); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取登录用户
	user, _ := ctx.Get("user")

	// TODO 创建Article
	article := model.Article{
		UserId:   user.(model.User).ID,
		Title:    requestArticle.Title,
		Content:  requestArticle.Content,
		ResLong:  requestArticle.ResLong,
		ResShort: requestArticle.ResShort,
		Visible: 1,
	}

	// TODO 插入数据
	if err := a.DB.Create(&article).Error; err != nil {
		panic(err)
	}

	// TODO 成功
	response.Success(ctx, gin.H{"article": article}, "创建成功")
}

// @title    Update
// @description   更新一篇文章的内容
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Update(ctx *gin.Context) {
	var requestArticle vo.CreateArticleRequest
	// TODO 数据验证
	if err := ctx.ShouldBind(&requestArticle); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// TODO 获取path中的id
	articleId := ctx.Params.ByName("id")

	var article model.Article

	// TODO 查看文章是否存在
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// TODO 判断当前用户是否为文章的作者
	// TODO 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID

	// TODO 查看文章的所属权
	if userId != article.UserId {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	// TODO 更新文章
	if err := a.DB.Model(&article).Update(requestArticle).Error; err != nil {
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"article": article}, "更新成功")
}

// @title    Show
// @description   查看一篇文章的内容
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Show(ctx *gin.Context) {
	// TODO 获取path中的id
	articleId := ctx.Params.ByName("id")
	var article model.Article

	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// TODO 查看是否有权限
	if article.UserId != user.ID && (article.Visible == 3 || (article.Visible == 2 && !Buil.IsS(4, "Fr"+strconv.Itoa(int(user.ID)), strconv.Itoa(int(article.UserId))))) {
		response.Fail(ctx, nil, "权限不足，不能查看")
		return
	}

	response.Success(ctx, gin.H{"article": article}, "成功")
}

// @title    Delete
// @description   删除一篇文章
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) Delete(ctx *gin.Context) {
	// 获取path中的id
	articleId := ctx.Params.ByName("id")

	var article model.Article

	// TODO 查看文章是否存在
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// TODO 判断当前用户是否为文章的作者
	// TODO 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID

	// TODO 查看用户是否有操作文章的权力
	if Buil.GetH(0, "permission", strconv.Itoa(int(userId)))[0] < '4' && userId != article.UserId {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	a.DB.Delete(&article)

	// TODO 移除收藏
	for _, val := range Buil.MembersS(1, "aF"+articleId) {
		Buil.RemS(1, "Fa"+val, articleId)
	}
	Buil.Del(1, "aF"+articleId)

	// TODO 移除点赞
	Buil.Del(1, "iL"+articleId)

	// TODO 移除标签
	for _, val := range Buil.MembersS(1, "aL"+articleId) {
		Buil.RemS(1, "La"+val, articleId)
		if Buil.CardS(1, "La"+val) == 0 {
			Buil.Del(1, "La"+val)
		}
	}
	Buil.Del(1, "aL"+articleId)

	response.Success(ctx, gin.H{"article": article}, "删除成功")
}

// @title    PageList
// @description   获取多篇文章
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) PageList(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")
	usera := tuser.(model.User)

	users := Buil.MembersS(4, "Fr"+strconv.Itoa(int(usera.ID)))

	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articles []model.Article

	// TODO 查找所有分页中可见的条目
	a.DB.Where("visible = 2 and user_id in (?)", users).Or("visible = 1").Or("user_id = ?", usera.ID).Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	var total int
	a.DB.Where("visible = 2 and user_id in (?)", users).Or("visible = 1").Or("user_id = ?", usera.ID).Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"articles": articles, "total": total}, "成功")
}

// @title    NewArticleController
// @description   新建一个IArticleController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IArticleController		返回一个IArticleController用于调用各种函数
func NewArticleController() IArticleController {
	db := common.GetDB()
	db.AutoMigrate(model.Article{})
	return ArticleController{DB: db}
}
