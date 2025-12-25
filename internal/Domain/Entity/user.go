package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;column:id;default:gen_random_uuid()"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
	Email          string    `gorm:"column:email"`
	Password       string    `gorm:"column:password"`
	Provider       string    `gorm:"column:provider"`
	Name           string    `gorm:"column:name"`
	FirstName      string    `gorm:"column:first_name"`
	LastName       string    `gorm:"column:last_name"`
	ProviderUserId string    `gorm:"column:provider_user_id"`
	AvatarURL      string    `gorm:"column:avatar_url"`
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
) User {
	return User{
		ID:             id,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
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
