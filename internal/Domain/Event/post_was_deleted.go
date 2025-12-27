package event

import (
	"github.com/google/uuid"
	"time"
)

type PostWasDeleted struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorId  uuid.UUID `json:"author_id"`
}

func NewPostWasDeleted(
	ID uuid.UUID,
	CreatedAt time.Time,
	UpdatedAt time.Time,
	Slug string,
	Title string,
	Content string,
	AuthorId uuid.UUID,
) PostWasDeleted {
	return PostWasDeleted{
		ID:        ID,
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
		Slug:      Slug,
		Title:     Title,
		Content:   Content,
		AuthorId:  AuthorId,
	}
}
