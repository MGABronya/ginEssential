// @Title  RestController
// @Description  该文件用于封装增删查改方法
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:33
package controller

import "github.com/gin-gonic/gin"

// RestController			定义了增删查改方法
type RestController interface {
	Create(ctx *gin.Context) //增
	Update(ctx *gin.Context) //删
	Show(ctx *gin.Context)   //查
	Delete(ctx *gin.Context) //改
}
