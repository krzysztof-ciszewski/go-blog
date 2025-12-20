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
			avatar_url,
			access_token,
			access_token_secret,
			refresh_token,
			expires_at,
			id_token)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
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
		user.AccessToken,
		user.AccessTokenSecret,
		user.RefreshToken,
		user.ExpiresAt.Format(time.RFC3339),
		user.IDToken,
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
			avatar_url,
			access_token,
			access_token_secret,
			refresh_token,
			expires_at,
			id_token
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
	var accessToken string
	var accessTokenSecret string
	var refreshToken string
	var expiresAt time.Time
	var idToken string

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
		&accessToken,
		&accessTokenSecret,
		&refreshToken,
		&expiresAt,
		&idToken,
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
		accessToken,
		accessTokenSecret,
		refreshToken,
		expiresAt,
		idToken,
	), nil
}

func (u userRepository) FindByEmail(email string) (entity.User, error) {
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
			avatar_url,
			access_token,
			access_token_secret,
			refresh_token,
			expires_at,
			id_token
		FROM users WHERE email = $1
	`, email)

	var userId uuid.UUID
	var createdAt time.Time
	var updatedAt time.Time
	var userEmail string
	var password sql.NullString
	var provider string
	var name sql.NullString
	var firstName sql.NullString
	var lastName sql.NullString
	var providerUserId string
	var avatarURL sql.NullString
	var accessToken string
	var accessTokenSecret string
	var refreshToken string
	var expiresAt time.Time
	var idToken string

	err := row.Scan(
		&userId,
		&createdAt,
		&updatedAt,
		&userEmail,
		&password,
		&provider,
		&name,
		&firstName,
		&lastName,
		&providerUserId,
		&avatarURL,
		&accessToken,
		&accessTokenSecret,
		&refreshToken,
		&expiresAt,
		&idToken,
	)

	if err != nil {
		return entity.User{}, err
	}

	return entity.NewUser(
		userId,
		createdAt,
		updatedAt,
		userEmail,
		nullStringValue(password),
		provider,
		nullStringValue(name),
		nullStringValue(firstName),
		nullStringValue(lastName),
		providerUserId,
		nullStringValue(avatarURL),
		accessToken,
		accessTokenSecret,
		refreshToken,
		expiresAt,
		idToken,
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
