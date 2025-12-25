package post_query

import (
	"context"
	view "main/internal/Application/View"
	repository "main/internal/Domain/Repository"

	"github.com/google/uuid"
)

type GetPostQueryHandler struct {
	PostRepository repository.PostRepository
}

func (h GetPostQueryHandler) Handle(ctx context.Context, query any) (any, error) {
	getPostQuery, ok := query.(GetPostQuery)
	if !ok {
		return view.PostView{}, nil
	}

	post, err := h.PostRepository.FindByID(getPostQuery.Id)

	if err != nil {
		return view.PostView{}, err
	}

	return view.NewPostView(
		uuid.MustParse(post.ID.String()),
		post.Slug,
		post.Title,
		post.Content,
		post.Author,
	), nil
}

func (h GetPostQueryHandler) Supports(query any) bool {
	_, ok := query.(GetPostQuery)
	return ok
}
