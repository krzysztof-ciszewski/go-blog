package command

import "github.com/google/uuid"

type deletePostCommand struct {
	Id uuid.UUID `json:"id"`
}

func NewDeletePostCommand(id uuid.UUID) deletePostCommand {
	return deletePostCommand{Id: id}
}
