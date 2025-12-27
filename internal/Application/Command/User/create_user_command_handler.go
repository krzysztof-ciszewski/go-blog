package command

import (
	"context"
	entity "main/internal/Domain/Entity"
	event "main/internal/Domain/Event"
	repository "main/internal/Domain/Repository"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type CreateUserCommandHandler struct {
	EventBus       *cqrs.EventBus
	UserRepository repository.UserRepository
}

func (h CreateUserCommandHandler) Handle(ctx context.Context, command *CreateUserCommand) error {
	user := entity.NewUser(
		command.Id,
		time.Now(),
		time.Now(),
		command.Email,
		command.Password,
		command.Provider,
		command.Name,
		command.FirstName,
		command.LastName,
		command.ProviderUserId,
		command.AvatarURL,
	)

	if _, err := h.UserRepository.FindByID(command.Id); err == nil {
		return nil
	}

	err := h.UserRepository.Save(user)
	if err != nil {
		return err
	}

	return h.EventBus.Publish(
		context.Background(),
		event.NewUserWasCreated(
			user.ID,
			user.Email,
			user.Provider,
			user.Name,
			user.FirstName,
			user.LastName,
			user.ProviderUserId,
			user.AvatarURL,
		),
	)
}
