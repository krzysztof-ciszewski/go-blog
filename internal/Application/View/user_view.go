package view

import (
	"github.com/google/uuid"
)

type UserView struct {
	entityView
	Email          string `json:"email"`
	Provider       string `json:"provider"`
	Name           string `json:"name"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProviderUserId string `json:"provider_user_id"`
	AvatarURL      string `json:"avatar_url"`
}

func NewUserView(
	id uuid.UUID,
	email string,
	provider string,
	name string,
	firstName string,
	lastName string,
	providerUserId string,
	avatarURL string,
) UserView {
	return UserView{
		entityView:     NewEntityView(id),
		Email:          email,
		Provider:       provider,
		Name:           name,
		FirstName:      firstName,
		LastName:       lastName,
		ProviderUserId: providerUserId,
		AvatarURL:      avatarURL,
	}
}
