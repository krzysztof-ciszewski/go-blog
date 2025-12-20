package entity

import (
	"time"

	"github.com/google/uuid"
)

type entity struct {
	id        uuid.UUID
	createdAt time.Time
	updatedAt time.Time
}

func (e entity) Id() uuid.UUID {
	return e.id
}

func (e entity) CreatedAt() time.Time {
	return e.createdAt
}

func (e entity) UpdatedAt() time.Time {
	return e.updatedAt
}

func NewEntity(id uuid.UUID, createdAt time.Time, updatedAt time.Time) entity {
	return entity{id: id, createdAt: createdAt, updatedAt: updatedAt}
}
