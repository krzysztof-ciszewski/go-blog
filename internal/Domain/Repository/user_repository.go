package repository

import (
	"context"
	entity "main/internal/Domain/Entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Save(ctx context.Context, user entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	FindByProviderUserIdAndEmail(ctx context.Context, providerUserId string, userEmail string) (entity.User, error)
}
