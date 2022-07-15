package vo

type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required,max=20"`
	Content string `json:"content" binding:"required"`
}
