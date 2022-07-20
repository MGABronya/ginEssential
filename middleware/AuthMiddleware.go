package middleware

import (
	"Essential/common"
	"Essential/model"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 authorization header
		tokenString := ctx.GetHeader("Authorization")

		fmt.Print("请求token", tokenString)

		//validate token formate
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(201, gin.H{
				"code": 201,
				"msg":  "格式错误，权限不足",
			})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:] //截取字符

		token, claims, err := common.ParseToken(tokenString)

		if err != nil || !token.Valid {
			ctx.JSON(201, gin.H{
				"code": 201,
				"msg":  "解析错误，权限不足",
			})
			ctx.Abort()
			return
		}

		//token通过验证, 获取claims中的UserID
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 验证用户是否存在
		if user.ID == 0 {
			ctx.JSON(201, gin.H{
				"code": 201,
				"msg":  "用户不存在，权限不足",
			})
			ctx.Abort()
			return
		}

		//用户存在 将user信息写入上下文
		ctx.Set("user", user)

		ctx.Next()
	}
}
