package query

import (
	"context"
	view "main/internal/Application/View"
	repository "main/internal/Domain/Repository"
)

type FindAllQueryHandler struct {
	PostRepository repository.PostRepository
}

func (h FindAllQueryHandler) Handle(ctx context.Context, query any) (any, error) {
	_, ok := query.(FindAllQuery)
	if !ok {
		return []view.PostView{}, nil
	}

	posts, err := h.PostRepository.FindAll()

	if err != nil {
		return []view.PostView{}, err
	}

	postViews := make([]view.PostView, len(posts))
	for i, post := range posts {
		postViews[i] = view.NewPostView(
			post.Id(),
			post.CreatedAt(),
			post.UpdatedAt(),
			post.Slug(),
			post.Title(),
			post.Content(),
			post.Author(),
		)
	}

	return postViews, nil
}

func (h FindAllQueryHandler) Supports(query any) bool {
	_, ok := query.(FindAllQuery)
	return ok
}
