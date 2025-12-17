package query

import (
	"context"
	view "main/internal/Application/View"
	repository "main/internal/Domain/Repository"
)

type FindAllByAuthorQueryHandler struct {
	PostRepository repository.PostRepository
}

func (h FindAllByAuthorQueryHandler) Handle(ctx context.Context, query any) (any, error) {
	findAllByAuthorQuery, ok := query.(FindAllByAuthorQuery)
	if !ok {
		return []view.PostView{}, nil
	}

	posts, err := h.PostRepository.FindAllByAuthor(findAllByAuthorQuery.Author())

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

func (h FindAllByAuthorQueryHandler) Supports(query any) bool {
	_, ok := query.(FindAllByAuthorQuery)
	return ok
}
