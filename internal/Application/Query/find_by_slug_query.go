package query

type FindBySlugQuery struct {
	Slug string `json:"slug"`
}

func NewFindBySlugQuery(slug string) FindBySlugQuery {
	return FindBySlugQuery{Slug: slug}
}
