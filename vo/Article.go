// @Title  article
// @Description  文章的基本信息
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package vo

// CreateArticleRequest		通过前端发送请求接收的文章信息
type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required,max=20"` // 标题
	Content string `json:"content" binding:"required"`      // 内容
}
