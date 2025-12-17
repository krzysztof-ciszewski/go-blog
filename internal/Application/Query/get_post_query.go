package query

import "github.com/google/uuid"

type GetPostQuery struct {
	id uuid.UUID `json:"id"`
}

func (q GetPostQuery) Id() uuid.UUID {
	return q.id
}

func NewGetPostQuery(id uuid.UUID) GetPostQuery {
	return GetPostQuery{id: id}
}
