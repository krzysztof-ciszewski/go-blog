package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;column:id;default:gen_random_uuid()"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Slug      string    `gorm:"column:slug"`
	Title     string    `gorm:"column:title"`
	Content   string    `gorm:"column:content"`
	AuthorId  uuid.UUID `gorm:"column:author_id"`
}

func NewPost(
	id uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
	slug string,
	title string,
	content string,
	authorId uuid.UUID,
) Post {
	return Post{ID: id, CreatedAt: createdAt, UpdatedAt: updatedAt, Slug: slug, Title: title, Content: content, AuthorId: authorId}
}
