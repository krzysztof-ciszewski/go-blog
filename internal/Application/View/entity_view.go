package view

import (
	"time"

	"github.com/google/uuid"
)

type entityView struct {
	id        uuid.UUID
	createdAt time.Time
	updatedAt time.Time
}

func (e entityView) Id() uuid.UUID {
	return e.id
}

func (e entityView) CreatedAt() time.Time {
	return e.createdAt
}

func (e entityView) UpdatedAt() time.Time {
	return e.updatedAt
}

func NewEntityView(id uuid.UUID, createdAt time.Time, updatedAt time.Time) entityView {
	return entityView{id: id, createdAt: createdAt, updatedAt: updatedAt}
}
