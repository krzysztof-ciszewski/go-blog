package event

import (
	"github.com/google/uuid"
)

type UserWasCreated struct {
	ID             uuid.UUID `json:"id"`
	Email          string    `json:"email"`
	Provider       string    `json:"provider"`
	Name           string    `json:"name"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ProviderUserId string    `json:"provider_user_id"`
	AvatarURL      string    `json:"avatar_url"`
}

func NewUserWasCreated(
	ID uuid.UUID,
	Email string,
	Provider string,
	Name string,
	FirstName string,
	LastName string,
	ProviderUserId string,
	AvatarURL string,
) UserWasCreated {
	return UserWasCreated{
		ID:             ID,
		Email:          Email,
		Provider:       Provider,
		Name:           Name,
		FirstName:      FirstName,
		LastName:       LastName,
		ProviderUserId: ProviderUserId,
		AvatarURL:      AvatarURL,
	}
}
