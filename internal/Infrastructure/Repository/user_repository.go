package repository

import (
	"context"
	entity "main/internal/Domain/Entity"
	repository "main/internal/Domain/Repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func (u userRepository) Save(ctx context.Context, user entity.User) error {
	return u.db.WithContext(ctx).Create(&user).Error
}

func (u userRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	var user entity.User
	err := u.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (u userRepository) FindByProviderUserIdAndEmail(ctx context.Context, providerUserId string, userEmail string) (entity.User, error) {
	var user entity.User
	err := u.db.WithContext(ctx).Where("provider_user_id = ? AND email = ?", providerUserId, userEmail).First(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}
