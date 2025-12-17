package query

type FindAllByTextQuery struct {
	text string `json:"text"`
}

func (q FindAllByTextQuery) Text() string {
	return q.text
}

func NewFindAllByTextQuery(text string) FindAllByTextQuery {
	return FindAllByTextQuery{text: text}
}
