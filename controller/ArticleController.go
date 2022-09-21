// @Title  ArticleController
// @Description  该文件提供关于操作文章的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"Essential/common"
	"Essential/dto"
	"Essential/model"
	"Essential/response"
	"Essential/vo"
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
		Icon:     user.(model.User).Icon,
		Name:     user.(model.User).Name,
		Email:    user.(model.User).Email,
	}

	// TODO 插入数据
	if err := a.DB.Create(&article).Error; err != nil {
		panic(err)
	}

	// TODO 成功
	response.Success(ctx, nil, "创建成功")
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
	user, _ := ctx.Get("user")
	var article model.Article

	// TODO 查看文章是否在数据库中存在
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"article": article, "user": dto.ToUserDto(user.(model.User))}, "成功")
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
	if userId != article.UserId {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	a.DB.Delete(&article)

	response.Success(ctx, gin.H{"article": article}, "成功")
}

// @title    PageList
// @description   获取多篇文章
// @auth      MGAronya（张健）       2022-9-16 12:20
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (a ArticleController) PageList(ctx *gin.Context) {
	// TODO 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// TODO 分页
	var articles []model.Article
	a.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	// TODO 记录的总条数
	var total int
	a.DB.Model(model.Article{}).Count(&total)

	// TODO 返回数据
	response.Success(ctx, gin.H{"data": articles, "total": total}, "成功")
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
