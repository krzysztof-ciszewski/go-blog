package view

import (
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
	slug string,
	title string,
	content string,
	author string,
) PostView {
	return PostView{entityView: NewEntityView(id), Slug: slug, Title: title, Content: content, Author: author}
}
