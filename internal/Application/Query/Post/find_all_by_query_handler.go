package post_query

import (
	"context"
	view "main/internal/Application/View"
	repository "main/internal/Domain/Repository"

	"github.com/google/uuid"
)

type FindAllByQueryHandler struct {
	PostRepository repository.PostRepository
}

func (h FindAllByQueryHandler) Handle(ctx context.Context, query any) (any, error) {
	findAllByQuery, ok := query.(FindAllByQuery)
	if !ok {
		return []view.PostView{}, nil
	}

	paginatedResult, err := h.PostRepository.FindAllBy(
		findAllByQuery.Filters.PaginationFilters.Page,
		findAllByQuery.Filters.PaginationFilters.PageSize,
		findAllByQuery.Filters.Slug,
		findAllByQuery.Filters.Text,
		findAllByQuery.Filters.Author,
	)

	if err != nil {
		return []view.PostView{}, err
	}

	postViews := make([]view.PostView, len(paginatedResult.Items))
	for i, post := range paginatedResult.Items {
		postViews[i] = view.NewPostView(
			uuid.MustParse(post.ID.String()),
			post.Slug,
			post.Title,
			post.Content,
			post.Author,
		)
	}

	return view.NewPaginatedView(postViews, paginatedResult.Total, paginatedResult.Page, paginatedResult.PageSize), nil
}

func (h FindAllByQueryHandler) Supports(query any) bool {
	_, ok := query.(FindAllByQuery)
	return ok
}
