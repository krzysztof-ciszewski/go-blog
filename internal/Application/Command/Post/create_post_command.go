package command

import "github.com/google/uuid"

type createPostCommand struct {
	Id      uuid.UUID `json:"id"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Author  uuid.UUID `json:"author"`
}

func NewCreatePostCommand(id uuid.UUID, slug string, title string, content string, author uuid.UUID) createPostCommand {
	return createPostCommand{Id: id, Slug: slug, Title: title, Content: content, Author: author}
}
