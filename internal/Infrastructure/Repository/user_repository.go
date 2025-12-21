package repository

import (
	"database/sql"
	entity "main/internal/Domain/Entity"
	repository "main/internal/Domain/Repository"
	"time"

	"github.com/google/uuid"
)

type userRepository struct {
	db *sql.DB
}

func (u userRepository) Save(user entity.User) error {
	_, err := u.db.Exec(`
		INSERT INTO users (
			id,
			created_at,
			updated_at,
			email,
			password,
			provider,
			name,
			first_name,
			last_name,
			provider_user_id,
			avatar_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`,
		user.Id().String(),
		user.CreatedAt().Format(time.RFC3339),
		user.UpdatedAt().Format(time.RFC3339),
		user.Email,
		sqlNullString(user.Password),
		user.Provider,
		sqlNullString(user.Name),
		sqlNullString(user.FirstName),
		sqlNullString(user.LastName),
		user.ProviderUserId,
		sqlNullString(user.AvatarURL),
	)

	return err
}

func (u userRepository) FindByID(id uuid.UUID) (entity.User, error) {
	row := u.db.QueryRow(`
		SELECT
			id,
			created_at,
			updated_at,
			email,
			password,
			provider,
			name,
			first_name,
			last_name,
			provider_user_id,
			avatar_url
		FROM users WHERE id = $1
	`, id)

	var userId uuid.UUID
	var createdAt time.Time
	var updatedAt time.Time
	var email string
	var password sql.NullString
	var provider string
	var name sql.NullString
	var firstName sql.NullString
	var lastName sql.NullString
	var providerUserId string
	var avatarURL sql.NullString

	err := row.Scan(
		&userId,
		&createdAt,
		&updatedAt,
		&email,
		&password,
		&provider,
		&name,
		&firstName,
		&lastName,
		&providerUserId,
		&avatarURL,
	)

	if err != nil {
		return entity.User{}, err
	}

	return entity.NewUser(
		userId,
		createdAt,
		updatedAt,
		email,
		nullStringValue(password),
		provider,
		nullStringValue(name),
		nullStringValue(firstName),
		nullStringValue(lastName),
		providerUserId,
		nullStringValue(avatarURL),
	), nil
}

func (u userRepository) FindByProviderUserIdAndEmail(filteredProviderUserId string, filteredUserEmail string) (entity.User, error) {
	row := u.db.QueryRow(`
		SELECT
			id,
			created_at,
			updated_at,
			email,
			password,
			provider,
			name,
			first_name,
			last_name,
			provider_user_id,
			avatar_url
		FROM users WHERE email = $1 AND provider_user_id = $2
	`, filteredUserEmail, filteredProviderUserId)

	var userId uuid.UUID
	var createdAt time.Time
	var updatedAt time.Time
	var email string
	var password sql.NullString
	var provider string
	var name sql.NullString
	var firstName sql.NullString
	var lastName sql.NullString
	var providerUserId string
	var avatarURL sql.NullString

	err := row.Scan(
		&userId,
		&createdAt,
		&updatedAt,
		&email,
		&password,
		&provider,
		&name,
		&firstName,
		&lastName,
		&providerUserId,
		&avatarURL,
	)

	if err != nil {
		return entity.User{}, err
	}

	return entity.NewUser(
		userId,
		createdAt,
		updatedAt,
		email,
		nullStringValue(password),
		provider,
		nullStringValue(name),
		nullStringValue(firstName),
		nullStringValue(lastName),
		providerUserId,
		nullStringValue(avatarURL),
	), nil
}

func sqlNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullStringValue(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}
