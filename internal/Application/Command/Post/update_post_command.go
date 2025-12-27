package command

import "github.com/google/uuid"

type updatePostCommand struct {
	Id      uuid.UUID `json:"id"`
	Slug    string    `json:"slug"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

func NewUpdatePostCommand(id uuid.UUID, slug string, title string, content string) updatePostCommand {
	return updatePostCommand{Id: id, Slug: slug, Title: title, Content: content}
}
