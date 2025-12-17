package view

import (
	"time"

	"github.com/google/uuid"
)

type entityView struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewEntityView(id uuid.UUID, createdAt time.Time, updatedAt time.Time) entityView {
	return entityView{Id: id, CreatedAt: createdAt, UpdatedAt: updatedAt}
}
