package command

import "github.com/google/uuid"

type createPostCommand struct {
	Id      uuid.UUID `json:"id"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Author  string    `json:"author"`
}

func NewCreatePostCommand(id uuid.UUID, slug string, title string, content string, author string) createPostCommand {
	return createPostCommand{Id: id, Slug: slug, Title: title, Content: content, Author: author}
}
