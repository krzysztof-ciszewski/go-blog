package request

type CreatePostRequest struct {
	Id      string `binding:"required,uuid"`
	Slug    string `binding:"required,min=3,max=255,alphanum"`
	Title   string `binding:"required,min=3,max=255"`
	Content string `binding:"required,min=10,max=10000"`
	Author  string `binding:"required,min=3,max=255"`
}
