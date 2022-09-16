// @Title  Thread
// @Description  接收前端请求时的跟帖信息
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-9-16 0:47
package vo

// CreateThreadRequest			接收前端请求时的跟帖信息
type CreateThreadRequest struct {
	Content string `json:"content" binding:"required"` // 内容
}
