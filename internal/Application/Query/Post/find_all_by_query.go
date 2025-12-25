package post_query

import query "main/internal/Application/Query"

type FindAllByQuery struct {
	Filters Filters
}

func NewFindAllByQuery(page int, pageSize int, slug string, text string, author string) FindAllByQuery {
	return FindAllByQuery{Filters: Filters{
		PaginationFilters: query.PaginationFilters{
			Page:     page,
			PageSize: pageSize,
		},
		Slug:   slug,
		Text:   text,
		Author: author,
	}}
}
