package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	entity
	slug    string
	title   string
	content string
	author  string
}

func (p Post) Slug() string {
	return p.slug
}

func (p Post) Title() string {
	return p.title
}

func (p Post) Content() string {
	return p.content
}

func (p Post) Author() string {
	return p.author
}

func NewPost(
	id uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
	slug string,
	title string,
	content string,
	author string,
) Post {
	return Post{entity: NewEntity(id, createdAt, updatedAt), slug: slug, title: title, content: content, author: author}
}
