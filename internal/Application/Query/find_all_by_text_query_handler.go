package query

import (
	"context"
	view "main/internal/Application/View"
	repository "main/internal/Domain/Repository"
)

type FindAllByTextQueryHandler struct {
	PostRepository repository.PostRepository
}

func (h FindAllByTextQueryHandler) Handle(ctx context.Context, query any) (any, error) {
	findAllByTextQuery, ok := query.(FindAllByTextQuery)
	if !ok {
		return []view.PostView{}, nil
	}

	posts, err := h.PostRepository.FindAllByText(findAllByTextQuery.Text)

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

func (h FindAllByTextQueryHandler) Supports(query any) bool {
	_, ok := query.(FindAllByTextQuery)
	return ok
}
