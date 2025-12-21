package command

import (
	"github.com/google/uuid"
)

type CreateUserCommand struct {
	Id             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Provider       string    `json:"provider"`
	Name           string    `json:"name"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ProviderUserId string    `json:"provider_user_id"`
	AvatarURL      string    `json:"avatar_url"`
}

func NewCreateUserCommand(
	id uuid.UUID,
	email string,
	password string,
	provider string,
	name string,
	firstName string,
	lastName string,
	providerUserId string,
	avatarURL string,
) CreateUserCommand {
	return CreateUserCommand{
		Id:             id,
		Email:          email,
		Password:       password,
		Provider:       provider,
		Name:           name,
		FirstName:      firstName,
		LastName:       lastName,
		ProviderUserId: providerUserId,
		AvatarURL:      avatarURL,
	}
}
