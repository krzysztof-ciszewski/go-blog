package repository

import (
	"context"
	entity "main/internal/Domain/Entity"
	repository "main/internal/Domain/Repository"
	open_telemetry "main/internal/Infrastructure/OpenTelemetry"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db        *gorm.DB
	telemetry open_telemetry.TelemetryProvider
}

func (u userRepository) Save(ctx context.Context, user entity.User) error {
	_, span := u.telemetry.TraceStart(ctx, "userRepository.Save")
	defer span.End()

	return u.db.Create(&user).Error
}

func (u userRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	_, span := u.telemetry.TraceStart(ctx, "userRepository.FindByID")
	defer span.End()

	var user entity.User
	err := u.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (u userRepository) FindByProviderUserIdAndEmail(ctx context.Context, providerUserId string, userEmail string) (entity.User, error) {
	_, span := u.telemetry.TraceStart(ctx, "userRepository.FindByProviderUserIdAndEmail")
	defer span.End()

	var user entity.User
	err := u.db.Where("provider_user_id = ? AND email = ?", providerUserId, userEmail).First(&user).Error
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func NewUserRepository(db *gorm.DB, telemetry open_telemetry.TelemetryProvider) repository.UserRepository {
	return &userRepository{db: db, telemetry: telemetry}
}
