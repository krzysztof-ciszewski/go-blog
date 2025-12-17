package query

type FindBySlugQuery struct {
	slug string `json:"slug"`
}

func (q FindBySlugQuery) Slug() string {
	return q.slug
}

func NewFindBySlugQuery(slug string) FindBySlugQuery {
	return FindBySlugQuery{slug: slug}
}
