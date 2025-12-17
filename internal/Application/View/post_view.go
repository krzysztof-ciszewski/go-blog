package view

import (
	"time"

	"github.com/google/uuid"
)

type PostView struct {
	entityView
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func NewPostView(
	id uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
	slug string,
	title string,
	content string,
	author string,
) PostView {
	return PostView{entityView: NewEntityView(id, createdAt, updatedAt), Slug: slug, Title: title, Content: content, Author: author}
}
