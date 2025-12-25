package view

import (
	"github.com/google/uuid"
)

type entityView struct {
	Id uuid.UUID `json:"id"`
}

func NewEntityView(id uuid.UUID) entityView {
	return entityView{Id: id}
}
