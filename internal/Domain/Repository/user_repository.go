package repository

import (
	entity "main/internal/Domain/Entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Save(user entity.User) error
	FindByID(id uuid.UUID) (entity.User, error)
	FindByProviderUserIdAndEmail(providerUserId string, userEmail string) (entity.User, error)
}
