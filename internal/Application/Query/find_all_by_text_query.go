package query

type FindAllByTextQuery struct {
	Text string `json:"text"`
}

func NewFindAllByTextQuery(text string) FindAllByTextQuery {
	return FindAllByTextQuery{Text: text}
}
