// @Title  PostTread
// @Description  接收前端帖子/跟帖信息时使用的结构体
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package vo

// PTRequest			接收前端帖子/跟帖信息时使用的结构体
type PTRequest struct {
	ThreadId string `json:"thread_id"`              // 帖子id
	Title    string `json:"title" binding:"max=20"` // 标题
	Content  string `json:"content"`                // 内容
	PT       bool   `json:"pt"`                     // PT为true时为帖子， 为false时为跟帖
}
