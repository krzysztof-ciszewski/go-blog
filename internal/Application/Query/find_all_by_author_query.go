package query

type FindAllByAuthorQuery struct {
	Author string `json:"author"`
}

func NewFindAllByAuthorQuery(author string) FindAllByAuthorQuery {
	return FindAllByAuthorQuery{Author: author}
}
