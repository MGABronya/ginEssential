package vo

type PTRequest struct {
	ThreadId string `json:"thread_id"`
	Title    string `json:"title" binding:"max=20"`
	Content  string `json:"content"`
	PT       bool   `json:"pt"`
}
