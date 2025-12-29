package repository

import (
	"context"
	entity "main/internal/Domain/Entity"
	repository "main/internal/Domain/Repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func (p postRepository) Save(ctx context.Context, post entity.Post) error {
	return p.db.WithContext(ctx).Create(&post).Error
}

func (p postRepository) Update(ctx context.Context, post entity.Post) error {
	return p.db.WithContext(ctx).Model(&post).Where("id = ?", post.ID).Updates(map[string]interface{}{
		"slug":       post.Slug,
		"title":      post.Title,
		"content":    post.Content,
		"updated_at": post.UpdatedAt,
	}).Error
}

func (p postRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.Post, error) {
	return gorm.G[entity.Post](p.db).Where("id = ?", id).First(ctx)
}

func (p postRepository) FindAllBy(ctx context.Context, page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
	var total int64
	tx := p.db.WithContext(ctx).Model(&entity.Post{})
	if slug != "" {
		tx = tx.Where("slug LIKE ?", "%"+slug+"%")
	}
	if text != "" {
		tx = tx.Where("content LIKE ? OR title LIKE ?", "%"+text+"%", "%"+text+"%")
	}
	if author != "" {
		tx = tx.Joins("JOIN users ON posts.author_id = users.id AND users.name LIKE ?", "%"+author+"%")
	}
	err := tx.Count(&total).Error
	if err != nil {
		return repository.PaginatedResult[entity.Post]{}, err
	}

	posts := make([]entity.Post, 0)
	err = tx.Offset((page - 1) * pageSize).Limit(pageSize).Find(&posts).Error
	if err != nil {
		return repository.PaginatedResult[entity.Post]{}, err
	}

	return repository.PaginatedResult[entity.Post]{Items: posts, Total: total, Page: page, PageSize: pageSize}, nil
}

func (p postRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return p.db.WithContext(ctx).Delete(&entity.Post{}, id).Error
}

func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return &postRepository{db: db}
}
