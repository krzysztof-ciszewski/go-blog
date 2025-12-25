package post_query

import query "main/internal/Application/Query"

type Filters struct {
	PaginationFilters query.PaginationFilters
	Slug              string
	Text              string
	Author            string
}
