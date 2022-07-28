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

type IArticleController interface {
	RestController
	PageList(ctx *gin.Context)
}

type ArticleController struct {
	DB *gorm.DB
}

func (a ArticleController) Create(ctx *gin.Context) {
	var requestArticle vo.CreateArticleRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestArticle); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取登录用户
	user, _ := ctx.Get("user")

	// 创建Article
	article := model.Article{
		UserId:  user.(model.User).ID,
		Title:   requestArticle.Title,
		Content: requestArticle.Content,
		Icon:    user.(model.User).Icon,
		Name:    user.(model.User).Name,
		Email:   user.(model.User).Email,
	}

	// 插入数据
	if err := a.DB.Create(&article).Error; err != nil {
		panic(err)
	}

	// 成功
	response.Success(ctx, nil, "创建成功")
}

func (a ArticleController) Update(ctx *gin.Context) {
	var requestArticle vo.CreateArticleRequest
	// 数据验证
	if err := ctx.ShouldBind(&requestArticle); err != nil {
		log.Print(err.Error())
		response.Fail(ctx, nil, "数据验证错误")
		return
	}

	// 获取path中的id
	articleId := ctx.Params.ByName("id")

	var article model.Article
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != article.UserId {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	// 更新文章
	if err := a.DB.Model(&article).Update(requestArticle).Error; err != nil {
		response.Fail(ctx, nil, "更新失败")
		return
	}

	response.Success(ctx, gin.H{"article": article}, "更新成功")
}

func (a ArticleController) Show(ctx *gin.Context) {
	// 获取path中的id
	articleId := ctx.Params.ByName("id")
	user, _ := ctx.Get("user")
	var article model.Article
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	response.Success(ctx, gin.H{"article": article, "user": dto.ToUserDto(user.(model.User))}, "成功")
}

func (a ArticleController) Delete(ctx *gin.Context) {
	// 获取path中的id
	articleId := ctx.Params.ByName("id")

	var article model.Article
	if a.DB.Where("id = ?", articleId).First(&article).RecordNotFound() {
		response.Fail(ctx, nil, "文章不存在")
		return
	}

	// 判断当前用户是否为文章的作者
	// 获取登录用户
	user, _ := ctx.Get("user")
	userId := user.(model.User).ID
	if userId != article.UserId {
		response.Fail(ctx, nil, "文章不属于您，请勿非法操作")
		return
	}

	a.DB.Delete(&article)

	response.Success(ctx, gin.H{"article": article}, "成功")
}

func (a ArticleController) PageList(ctx *gin.Context) {
	// 获取分页参数
	pageNum, _ := strconv.Atoi(ctx.DefaultQuery("pageNum", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "20"))

	// 分页
	var articles []model.Article
	a.DB.Order("created_at desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&articles)

	//var user []model.User

	//a.DB.Where("id = ?", article.UserId).First(&user)

	// 记录的总条数
	var total int
	a.DB.Model(model.Article{}).Count(&total)

	//"user": dto.ToUserDto(user.(model.User))
	// 返回数据
	response.Success(ctx, gin.H{"data": articles, "total": total}, "成功")
}

func NewArticleController() IArticleController {
	db := common.GetDB()
	db.AutoMigrate(model.Article{})
	return ArticleController{DB: db}
}
