package query

type FindAllByAuthorQuery struct {
	author string `json:"author"`
}

func (q FindAllByAuthorQuery) Author() string {
	return q.author
}

func NewFindAllByAuthorQuery(author string) FindAllByAuthorQuery {
	return FindAllByAuthorQuery{author: author}
}
