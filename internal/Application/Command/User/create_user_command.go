package command

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserCommand struct {
	Id                uuid.UUID `json:"id"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	Provider          string    `json:"provider"`
	Name              string    `json:"name"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	ProviderUserId    string    `json:"provider_user_id"`
	AvatarURL         string    `json:"avatar_url"`
	AccessToken       string    `json:"access_token"`
	AccessTokenSecret string    `json:"access_token_secret"`
	RefreshToken      string    `json:"refresh_token"`
	ExpiresAt         time.Time `json:"expires_at"`
	IDToken           string    `json:"id_token"`
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
	accessToken string,
	accessTokenSecret string,
	refreshToken string,
	expiresAt time.Time,
	idToken string,
) CreateUserCommand {
	return CreateUserCommand{
		Id:                id,
		Email:             email,
		Password:          password,
		Provider:          provider,
		Name:              name,
		FirstName:         firstName,
		LastName:          lastName,
		ProviderUserId:    providerUserId,
		AvatarURL:         avatarURL,
		AccessToken:       accessToken,
		AccessTokenSecret: accessTokenSecret,
		RefreshToken:      refreshToken,
		ExpiresAt:         expiresAt,
		IDToken:           idToken,
	}
}
