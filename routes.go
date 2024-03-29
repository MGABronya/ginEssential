// @Title  routes
// @Description  程序的路由均集中在这个文件里
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:50
package main

import (
	"ginEssential/controller"
	"ginEssential/middleware"

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

	// TODO 用户的邮箱验证
	r.GET("/verify/:id", controller.VerifyEmail)

	// TODO 用户找回密码
	r.PUT("/security", controller.Security)

	// TODO 用户更改密码
	r.PUT("/updatepass", middleware.AuthMiddleware(), controller.UpdatePass)

	// TODO 用户的登录路由
	r.POST("/login", controller.Login)

	// TODO 用户的个人页面路由
	r.GET("/personal", middleware.AuthMiddleware(), controller.PersonalPage)

	// TODO 用户的个人文章路由
	r.GET("/personal/articles", middleware.AuthMiddleware(), controller.PersonalPageArticles)

	// TODO 用户的个人帖子路由
	r.GET("/personal/posts", middleware.AuthMiddleware(), controller.PersonalPagePosts)

	// TODO 用户的个人跟帖路由
	r.GET("/personal/threads", middleware.AuthMiddleware(), controller.PersonalPageThreads)

	// TODO 用户的个人页面信息更新路由
	r.PUT("/personal", middleware.AuthMiddleware(), controller.PersonalUpdate)

	// TODO 用户的头像更新路由
	r.PUT("/personal/icon", middleware.AuthMiddleware(), controller.PersonalIcon)

	// TODO 个人信息展示路由
	r.GET("/personal/:id", middleware.AuthMiddleware(), controller.PersonalShow)

	// TODO 个人文章展示路由
	r.GET("/personal/articles/:id", middleware.AuthMiddleware(), controller.PersonalShowArticles)

	// TODO 个人帖子展示路由
	r.GET("/personal/posts/:id", middleware.AuthMiddleware(), controller.PersonalShowPosts)

	// TODO 个人跟帖展示路由
	r.GET("/personal/threads/:id", middleware.AuthMiddleware(), controller.PersonalShowThreads)

	// TODO 查找一组用户的信息
	r.POST("/personal/users", middleware.AuthMiddleware(), controller.PersonalShowUsers)

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
	articleRoutes.GET("/pagelist", articleController.PageList)

	// TODO 帖子的路由分组
	postRoutes := r.Group("/post")

	// TODO 添加中间件
	postRoutes.Use(middleware.AuthMiddleware())

	// TODO 创建帖子controller
	postController := controller.NewPostController()

	// TODO 帖子的创建路由
	postRoutes.POST("", postController.Create)

	// TODO 帖子的更新路由
	postRoutes.PUT("/:id", postController.Update)

	// TODO 帖子的查看路由
	postRoutes.GET("/:id", postController.Show)

	// TODO 帖子的删除路由
	postRoutes.DELETE("/:id", postController.Delete)

	//TODO 帖子的列表路由
	postRoutes.GET("/pagelist", postController.PageList)

	// TODO 跟帖的路由分组
	threadRoutes := r.Group("/thread")

	// TODO 添加中间件
	threadRoutes.Use(middleware.AuthMiddleware())

	// TODO 创建跟帖controller
	threadController := controller.NewThreadController()

	// TODO 跟帖的创建路由
	threadRoutes.POST("/:id", threadController.Create)

	// TODO 跟帖的更新路由
	threadRoutes.PUT("/:id", threadController.Update)

	// TODO 跟帖的查看路由
	threadRoutes.GET("/:id", threadController.Show)

	// TODO 跟帖的删除路由
	threadRoutes.DELETE("/:id", threadController.Delete)

	//TODO 跟帖的列表路由
	threadRoutes.GET("/pagelist/:id", threadController.PageList)

	// TODO 背景图片的路由分组
	backgroundRoutes := r.Group("/background")

	// TODO 添加中间件
	backgroundRoutes.Use(middleware.AuthMiddleware())

	// TODO 背景图片的controller
	backgroundController := controller.NewBackGroundController()

	// TODO 查看用户使用的背景图片
	backgroundRoutes.GET("/show", backgroundController.Show)

	// TODO 用户选择新的背景图片
	backgroundRoutes.PUT("/update/:id", backgroundController.Update)

	// TODO 用户上传自己的用户图片
	backgroundRoutes.PUT("/create", backgroundController.Create)

	return r
}
