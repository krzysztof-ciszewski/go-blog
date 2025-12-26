package user_query

import (
	"context"
	"errors"
	view "main/internal/Application/View"
	entity "main/internal/Domain/Entity"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type mockUserRepository struct {
	findByProviderUserIdAndEmailFunc func(providerUserId string, userEmail string) (entity.User, error)
}

func (m *mockUserRepository) Save(user entity.User) error {
	return nil
}

func (m *mockUserRepository) FindByID(id uuid.UUID) (entity.User, error) {
	return entity.User{}, nil
}

func (m *mockUserRepository) FindByProviderUserIdAndEmail(providerUserId string, userEmail string) (entity.User, error) {
	if m.findByProviderUserIdAndEmailFunc != nil {
		return m.findByProviderUserIdAndEmailFunc(providerUserId, userEmail)
	}
	return entity.User{}, errors.New("not implemented")
}

type FindUserByQueryHandlerTestSuite struct {
	suite.Suite
	Handler        FindUserByQueryHandler
	MockRepository *mockUserRepository
}

func (s *FindUserByQueryHandlerTestSuite) SetupTest() {
	s.MockRepository = &mockUserRepository{}
	s.Handler = FindUserByQueryHandler{
		UserRepository: s.MockRepository,
	}
}

func (s *FindUserByQueryHandlerTestSuite) TestHandle_Success() {
	testUserID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	testUser := entity.User{
		ID:             testUserID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          "test@example.com",
		Provider:       "github",
		Name:           "Test User",
		FirstName:      "Test",
		LastName:       "User",
		ProviderUserId: "testprovider123",
		AvatarURL:      "https://example.com/avatar.jpg",
	}

	s.MockRepository.findByProviderUserIdAndEmailFunc = func(providerUserId string, userEmail string) (entity.User, error) {
		assert.Equal(s.T(), "testprovider123", providerUserId)
		assert.Equal(s.T(), "test@example.com", userEmail)
		return testUser, nil
	}

	ctx := context.Background()
	query := NewFindUserByQuery("testprovider123", "test@example.com")

	result, err := s.Handler.Handle(ctx, query)

	assert.NoError(s.T(), err)
	userView, ok := result.(view.UserView)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), testUserID, userView.Id)
	assert.Equal(s.T(), "test@example.com", userView.Email)
	assert.Equal(s.T(), "github", userView.Provider)
	assert.Equal(s.T(), "testprovider123", userView.ProviderUserId)
	assert.Equal(s.T(), "Test User", userView.Name)
	assert.Equal(s.T(), "Test", userView.FirstName)
	assert.Equal(s.T(), "User", userView.LastName)
	assert.Equal(s.T(), "https://example.com/avatar.jpg", userView.AvatarURL)
}

func (s *FindUserByQueryHandlerTestSuite) TestHandle_ErrorCases() {
	tests := []struct {
		name           string
		providerUserId string
		userEmail      string
	}{
		{
			name:           "UserNotFound",
			providerUserId: "nonexistent",
			userEmail:      "nonexistent@example.com",
		},
		{
			name:           "EmptyProviderUserId",
			providerUserId: "",
			userEmail:      "test@example.com",
		},
		{
			name:           "EmptyEmail",
			providerUserId: "testprovider123",
			userEmail:      "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.MockRepository.findByProviderUserIdAndEmailFunc = func(providerUserId string, userEmail string) (entity.User, error) {
				return entity.User{}, errors.New("user not found")
			}

			ctx := context.Background()
			query := NewFindUserByQuery(tt.providerUserId, tt.userEmail)

			result, err := s.Handler.Handle(ctx, query)

			assert.Error(s.T(), err)
			userView, ok := result.(view.UserView)
			assert.True(s.T(), ok)
			assert.Equal(s.T(), uuid.Nil, userView.Id)
		})
	}
}

func (s *FindUserByQueryHandlerTestSuite) TestHandle_InvalidQueryType() {
	ctx := context.Background()
	invalidQuery := "not a FindUserByQuery"

	result, err := s.Handler.Handle(ctx, invalidQuery)

	assert.NoError(s.T(), err)
	userView, ok := result.(view.UserView)
	assert.True(s.T(), ok)
	assert.Equal(s.T(), uuid.Nil, userView.Id)
}

func (s *FindUserByQueryHandlerTestSuite) TestSupports() {
	tests := []struct {
		name     string
		query    any
		expected bool
	}{
		{
			name:     "ValidQuery",
			query:    NewFindUserByQuery("testprovider123", "test@example.com"),
			expected: true,
		},
		{
			name:     "InvalidQuery",
			query:    "not a FindUserByQuery",
			expected: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			supports := s.Handler.Supports(tt.query)
			assert.Equal(s.T(), tt.expected, supports)
		})
	}
}

func TestFindUserByQueryHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FindUserByQueryHandlerTestSuite))
}
