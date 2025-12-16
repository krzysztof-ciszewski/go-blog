package repository

import (
	entity "main/internal/Domain/Entity"

	"github.com/google/uuid"
)

type PostRepository interface {
	Save(post entity.Post) error
	FindByID(id uuid.UUID) (entity.Post, error)
	FindBySlug(slug string) (entity.Post, error)
	FindAll() ([]entity.Post, error)
	FindAllByAuthor(author string) ([]entity.Post, error)
	FindAllByText(text string) ([]entity.Post, error)
	Delete(id uuid.UUID) error
}
