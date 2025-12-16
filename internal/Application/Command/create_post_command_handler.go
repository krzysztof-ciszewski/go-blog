package command

import (
	"context"
	entity "main/internal/Domain/Entity"
	repository "main/internal/Domain/Repository"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type CreatePostCommandHandler struct {
	EventBus       *cqrs.EventBus
	PostRepository repository.PostRepository
}

func (h CreatePostCommandHandler) Handle(ctx context.Context, command *createPostCommand) error {
	post := entity.NewPost(
		command.Id,
		time.Now(),
		time.Now(),
		command.Slug,
		command.Title,
		command.Content,
		command.Author,
	)

	if _, err := h.PostRepository.FindByID(command.Id); err == nil {
		return nil
	}

	return h.PostRepository.Save(post)
}
