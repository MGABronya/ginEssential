// @Title  BackGroundController
// @Description  该文件提供操作背景的各种方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import (
	"ginEssential/common"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/util"
	"log"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// IBackGroundController			定义了背景类接口
type IBackGroundController interface {
	Create(ctx *gin.Context) //增
	Show(ctx *gin.Context)   //查
	Update(ctx *gin.Context) //改
}

// BackGroundController			定义了背景工具类
type BackGroundController struct {
	DB *gorm.DB // 含有一个数据库指针
}

// @title    Create
// @description   创建一个背景并使用
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BackGroundController) Create(ctx *gin.Context) {
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

	file.Filename = strconv.Itoa(int(user.ID)) + extName

	// TODO 将文件存入本地
	ctx.SaveUploadedFile(file, "./BackGroundController/"+file.Filename)

	b.DB.Where("id = ?", user.ID).Take(&user)

	user.BackGround = file.Filename

	b.DB.Save(&user)

	response.Success(ctx, gin.H{"BackGround": user.BackGround}, "更新成功")
}

// @title    Update
// @description   用户选择默认背景
// @auth      MGAronya（张健）       2022-9-16 12:15
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BackGroundController) Update(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	backgroundId := ctx.Params.ByName("id")

	// TODO 非默认图片，则报错
	if !util.VerifyIconFormat(backgroundId) {
		response.Fail(ctx, nil, "该图片不存在")
		return
	}

	b.DB.Where("id = ?", user.ID).Take(&user)

	user.BackGround = backgroundId

	b.DB.Save(&user)

	response.Success(ctx, gin.H{"background": backgroundId}, "更新成功")
}

// @title    Show
// @description   查看用户的背景图片
// @auth      MGAronya（张健）       2022-9-16 12:19
// @param    ctx *gin.Context       接收一个上下文
// @return   void
func (b BackGroundController) Show(ctx *gin.Context) {
	tuser, _ := ctx.Get("user")
	user := tuser.(model.User)

	b.DB.Where("id = ?", user.ID).Take(&user)

	response.Success(ctx, gin.H{"background": user.BackGround}, "查看成功")
}

// @title    NewBackGroundController
// @description   新建一个IBackGroundController
// @auth      MGAronya（张健）       2022-9-16 12:23
// @param    void
// @return   IBackGroundController		返回一个IBackGroundController用于调用各种函数
func NewBackGroundController() IBackGroundController {
	db := common.GetDB()
	return BackGroundController{DB: db}
}
