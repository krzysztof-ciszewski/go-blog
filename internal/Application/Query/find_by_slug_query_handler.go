package query

import (
	"context"
	view "main/internal/Application/View"
	repository "main/internal/Domain/Repository"
)

type FindBySlugQueryHandler struct {
	PostRepository repository.PostRepository
}

func (h FindBySlugQueryHandler) Handle(ctx context.Context, query any) (any, error) {
	findBySlugQuery, ok := query.(FindBySlugQuery)
	if !ok {
		return view.PostView{}, nil
	}

	post, err := h.PostRepository.FindBySlug(findBySlugQuery.Slug())

	if err != nil {
		return view.PostView{}, err
	}

	return view.NewPostView(
		post.Id(),
		post.CreatedAt(),
		post.UpdatedAt(),
		post.Slug(),
		post.Title(),
		post.Content(),
		post.Author(),
	), nil
}

func (h FindBySlugQueryHandler) Supports(query any) bool {
	_, ok := query.(FindBySlugQuery)
	return ok
}
