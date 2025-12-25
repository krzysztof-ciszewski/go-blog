package repository

import (
	entity "main/internal/Domain/Entity"

	"github.com/google/uuid"
)

type PostRepository interface {
	Save(post entity.Post) error
	FindByID(id uuid.UUID) (entity.Post, error)
	FindAllBy(page int, pageSize int, slug string, text string, author string) (PaginatedResult[entity.Post], error)
	Delete(id uuid.UUID) error
}
