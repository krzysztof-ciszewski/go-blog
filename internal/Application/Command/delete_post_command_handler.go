package command

import (
	"context"
	repository "main/internal/Domain/Repository"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type DeletePostCommandHandler struct {
	EventBus       *cqrs.EventBus
	PostRepository repository.PostRepository
}

func (h DeletePostCommandHandler) Handle(ctx context.Context, command *deletePostCommand) error {
	return h.PostRepository.Delete(command.Id)
}
