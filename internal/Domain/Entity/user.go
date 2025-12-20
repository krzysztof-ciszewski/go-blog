package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	entity
	Email             string
	Password          string
	Provider          string
	Name              string
	FirstName         string
	LastName          string
	ProviderUserId    string
	AvatarURL         string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         time.Time
	IDToken           string
}

func NewUser(
	id uuid.UUID,
	createdAt time.Time,
	updatedAt time.Time,
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
) User {
	return User{
		entity:            NewEntity(id, createdAt, updatedAt),
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
