package query

import (
	"context"
	view "main/internal/Application/View"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
)

func GetPostQueryHandler(ctx context.Context, query *GetPostQuery) (view.PostView, error) {
	container := dependency_injection.GetContainer()

	post, err := container.PostRepository.FindByID(query.Id())

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
