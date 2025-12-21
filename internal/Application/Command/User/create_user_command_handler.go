package command

import (
	"context"
	entity "main/internal/Domain/Entity"
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

	return h.UserRepository.Save(user)
}
