package post_query

import (
	"context"
	view "main/internal/Application/View"
	repository "main/internal/Domain/Repository"
)

type GetPostQueryHandler struct {
	PostRepository repository.PostRepository
}

func (h GetPostQueryHandler) Handle(ctx context.Context, query any) (any, error) {
	getPostQuery, ok := query.(GetPostQuery)
	if !ok {
		return view.PostView{}, nil
	}

	post, err := h.PostRepository.FindByID(ctx, getPostQuery.Id)

	if err != nil {
		return view.PostView{}, err
	}

	return view.NewPostView(
		post.ID,
		post.Slug,
		post.Title,
		post.Content,
		post.AuthorId,
	), nil
}

func (h GetPostQueryHandler) Supports(query any) bool {
	_, ok := query.(GetPostQuery)
	return ok
}
