package repository

import (
	"context"
	entity "main/internal/Domain/Entity"

	"github.com/google/uuid"
)

type PostRepository interface {
	Save(ctx context.Context, post entity.Post) error
	Update(ctx context.Context, post entity.Post) error
	FindByID(ctx context.Context, id uuid.UUID) (entity.Post, error)
	FindAllBy(ctx context.Context, page int, pageSize int, slug string, text string, author string) (PaginatedResult[entity.Post], error)
	Delete(ctx context.Context, id uuid.UUID) error
}
