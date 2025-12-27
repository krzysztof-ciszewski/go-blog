package request

type UpdatePostRequest struct {
	Slug    string `binding:"required,min=3,max=255,alphanum"`
	Title   string `binding:"required,min=3,max=255"`
	Content string `binding:"required,min=10,max=10000"`
}
