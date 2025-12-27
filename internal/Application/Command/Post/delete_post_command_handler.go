package command

import (
	"context"
	event "main/internal/Domain/Event"
	repository "main/internal/Domain/Repository"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type DeletePostCommandHandler struct {
	EventBus       *cqrs.EventBus
	PostRepository repository.PostRepository
}

func (h DeletePostCommandHandler) Handle(ctx context.Context, command *deletePostCommand) error {
	post, err := h.PostRepository.FindByID(command.Id)
	if err != nil {
		return err
	}

	err = h.PostRepository.Delete(command.Id)
	if err != nil {
		return err
	}

	return h.EventBus.Publish(
		context.Background(),
		event.NewPostWasDeleted(
			post.ID,
			post.CreatedAt,
			post.UpdatedAt,
			post.Slug,
			post.Title,
			post.Content,
			post.AuthorId,
		),
	)
}
