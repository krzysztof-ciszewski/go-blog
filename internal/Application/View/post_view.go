package view

import (
	"github.com/google/uuid"
)

type PostView struct {
	entityView
	Slug     string    `json:"slug"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	AuthorId uuid.UUID `json:"author_id"`
}

func NewPostView(
	id uuid.UUID,
	slug string,
	title string,
	content string,
	authorId uuid.UUID,
) PostView {
	return PostView{entityView: NewEntityView(id), Slug: slug, Title: title, Content: content, AuthorId: authorId}
}
