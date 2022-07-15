package vo

type CreateThreadRequest struct {
	Content string `json:"content" binding:"required"`
}
