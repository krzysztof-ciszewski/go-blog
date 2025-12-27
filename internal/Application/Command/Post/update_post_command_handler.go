package command

import (
	"context"
	entity "main/internal/Domain/Entity"
	event "main/internal/Domain/Event"
	repository "main/internal/Domain/Repository"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type UpdatePostCommandHandler struct {
	EventBus       *cqrs.EventBus
	PostRepository repository.PostRepository
}

func (h UpdatePostCommandHandler) Handle(ctx context.Context, command *updatePostCommand) error {
	existingPost, err := h.PostRepository.FindByID(command.Id)
	if err != nil {
		return err
	}

	updatedPost := entity.NewPost(
		existingPost.ID,
		existingPost.CreatedAt,
		time.Now(),
		command.Slug,
		command.Title,
		command.Content,
		existingPost.AuthorId,
	)

	err = h.PostRepository.Update(updatedPost)
	if err != nil {
		return err
	}

	return h.EventBus.Publish(
		context.Background(),
		event.NewPostWasUpdated(
			updatedPost.ID,
			updatedPost.CreatedAt,
			updatedPost.UpdatedAt,
			updatedPost.Slug,
			updatedPost.Title,
			updatedPost.Content,
			updatedPost.AuthorId,
		),
	)
}
