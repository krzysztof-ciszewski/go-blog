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

func (p postRepository) Save(post entity.Post) error {
	return p.db.Create(&post).Error
}

func (p postRepository) FindByID(id uuid.UUID) (entity.Post, error) {
	return gorm.G[entity.Post](p.db).Where("id = ?", id).First(context.Background())
}

func (p postRepository) FindAllBy(page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
	var total int64
	tx := p.db.Model(&entity.Post{})
	if slug != "" {
		tx = tx.Where("slug LIKE ?", "%"+slug+"%")
	}
	if text != "" {
		tx = tx.Where("(content LIKE ? OR title LIKE ?)", "%"+text+"%", "%"+text+"%")
	}
	if author != "" {
		tx = tx.Where("author LIKE ?", "%"+author+"%")
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

func (p postRepository) Delete(id uuid.UUID) error {
	return p.db.Delete(&entity.Post{}, id).Error
}

func NewPostRepository(db *gorm.DB) repository.PostRepository {
	return &postRepository{db: db}
}
