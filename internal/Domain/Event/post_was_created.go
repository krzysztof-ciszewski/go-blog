package event

import (
	"github.com/google/uuid"
	"time"
)

type PostWasCreated struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Slug      string    `json:"slug"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorId  uuid.UUID `json:"author_id"`
}

func NewPostWasCreated(
	ID uuid.UUID,
	CreatedAt time.Time,
	UpdatedAt time.Time,
	Slug string,
	Title string,
	Content string,
	AuthorId uuid.UUID,
) PostWasCreated {
	return PostWasCreated{
		ID:        ID,
		CreatedAt: CreatedAt,
		UpdatedAt: UpdatedAt,
		Slug:      Slug,
		Title:     Title,
		Content:   Content,
		AuthorId:  AuthorId,
	}
}
