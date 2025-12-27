package command

import (
	"context"
	entity "main/internal/Domain/Entity"
	event "main/internal/Domain/Event"
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

	err := h.PostRepository.Save(post)
	if err != nil {
		return err
	}

	return h.EventBus.Publish(
		context.Background(),
		event.NewPostWasCreated(
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
