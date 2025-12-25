package post_query

import "github.com/google/uuid"

type GetPostQuery struct {
	Id uuid.UUID `json:"id"`
}

func NewGetPostQuery(id uuid.UUID) GetPostQuery {
	return GetPostQuery{Id: id}
}
