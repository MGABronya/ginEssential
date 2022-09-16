// @Title  routes
// @Description  程序的路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package main

import (
	"Essential/controller"
	"Essential/middleware"

	"github.com/gin-gonic/gin"
)

// @title    CollectRoute
// @description   给gin引擎挂上路由监听
// @auth      MGAronya（张健）             2022-9-16 10:52
// @param     r *gin.Engine			gin引擎
// @return    r *gin.Engine			gin引擎
func CollectRoute(r *gin.Engine) *gin.Engine {
	// TODO 添加中间件
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())

	// TODO 用户的注册路由
	r.POST("/regist", controller.Register)

	// TODO 用户的登录路由
	r.POST("/login", controller.Login)

	// TODO 用户的个人页面路由
	r.GET("/personal", middleware.AuthMiddleware(), controller.PersonalPage)

	// TODO 用户的个人页面信息更新路由
	r.PUT("/personal", middleware.AuthMiddleware(), controller.PersonalUpdate)

	// TODO 用户的头像更新路由
	r.POST("/personal", middleware.AuthMiddleware(), controller.PersonalIcon)

	// TODO 个人信息展示路由
	r.GET("/personal/:id", middleware.AuthMiddleware(), controller.PersonalShow)

	// TODO 文章路由分组
	articleRoutes := r.Group("/article")

	// TODO 添加中间件
	articleRoutes.Use(middleware.AuthMiddleware())

	// TODO 文章的controller
	articleController := controller.NewArticleController()

	// TODO 文章的创建路由
	articleRoutes.POST("", articleController.Create)

	// TODO 文章的更新路由
	articleRoutes.PUT("/:id", articleController.Update)

	// TODO 文章的展示路由
	articleRoutes.GET("/:id", articleController.Show)

	// TODO 文章的删除路由
	articleRoutes.DELETE("/:id", articleController.Delete)

	// TODO 文章的列表路由
	articleRoutes.POST("/pagelist", articleController.PageList)

	// TODO 帖子的路由分组
	postRoutes := r.Group("/post")

	// TODO 添加中间件
	postRoutes.Use(middleware.AuthMiddleware())

	// TODO 创建帖子controller
	postController := controller.NewPostController()

	// TODO 帖子的创建路由
	postRoutes.POST("", postController.Create)

	// TODO 帖子/跟帖的更新路由
	postRoutes.PUT("/:id", postController.Update)

	// TODO 帖子的查看路由
	postRoutes.GET("/:id", postController.Show)

	// TODO 帖子/跟帖的删除路由
	postRoutes.DELETE("/:id", postController.Delete)

	//TODO 帖子的列表路由
	postRoutes.POST("/pagelist", postController.PageList)

	// TODO 跟帖的创建路由
	postRoutes.POST("/:id", postController.ThreadCreate)

	return r
}
