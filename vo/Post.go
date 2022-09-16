// @Title  Post
// @Description  接收前端帖子信息时使用的结构体
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package vo

// CreatePostRequest			接收前端帖子信息时使用的结构体
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,max=20"` // 标题
	Content string `json:"content" binding:"required"`      // 内容
}
