package view

import (
	"time"

	"github.com/google/uuid"
)

type PostView struct {
	entityView
	slug    string
	title   string
	content string
	author  string
}

func (p PostView) Id() uuid.UUID {
	return p.id
}

func (p PostView) Slug() string {
	return p.slug
}

func (p PostView) Title() string {
	return p.title
}

func (p PostView) Content() string {
	return p.content
}

func (p PostView) Author() string {
	return p.author
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
	return PostView{entityView: NewEntityView(id, createdAt, updatedAt), slug: slug, title: title, content: content, author: author}
}
