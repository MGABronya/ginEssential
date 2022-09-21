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
	ResLong  string `json:"res_long"`               // 备用长文本
	ResShort string `json:"res_short"`              // 备用短文本
	PT       bool   `json:"pt"`                     // PT为true时为帖子， 为false时为跟帖
}
