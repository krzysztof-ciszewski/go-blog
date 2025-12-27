package command

import (
	"context"
	"database/sql"
	"errors"
	entity "main/internal/Domain/Entity"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wmsqlitemodernc "github.com/ThreeDotsLabs/watermill-sqlite/wmsqlitemodernc"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type mockUserRepositoryCreate struct {
	saveFunc     func(user entity.User) error
	findByIDFunc func(id uuid.UUID) (entity.User, error)
}

func (m *mockUserRepositoryCreate) Save(user entity.User) error {
	if m.saveFunc != nil {
		return m.saveFunc(user)
	}
	return nil
}

func (m *mockUserRepositoryCreate) FindByID(id uuid.UUID) (entity.User, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return entity.User{}, errors.New("not implemented")
}

func (m *mockUserRepositoryCreate) FindByProviderUserIdAndEmail(providerUserId string, userEmail string) (entity.User, error) {
	return entity.User{}, errors.New("not implemented")
}

type CreateUserCommandHandlerTestSuite struct {
	suite.Suite
	Handler        CreateUserCommandHandler
	MockRepository *mockUserRepositoryCreate
	EventBus       *cqrs.EventBus
}

func (s *CreateUserCommandHandlerTestSuite) SetupTest() {
	s.MockRepository = &mockUserRepositoryCreate{}

	db, _ := sql.Open("sqlite", ":memory:")
	db.SetMaxOpenConns(1)
	publisher, err := wmsqlitemodernc.NewPublisher(db, wmsqlitemodernc.PublisherOptions{
		InitializeSchema: true,
		Logger:           watermill.NopLogger{},
	})
	if err != nil {
		panic(err)
	}
	marshaller := cqrs.JSONMarshaler{}
	eventBus, err := cqrs.NewEventBusWithConfig(publisher, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return "events." + params.EventName, nil
		},
		Marshaler: marshaller,
		Logger:    watermill.NopLogger{},
	})
	if err != nil {
		panic(err)
	}
	s.EventBus = eventBus

	s.Handler = CreateUserCommandHandler{
		EventBus:       s.EventBus,
		UserRepository: s.MockRepository,
	}
}

func (s *CreateUserCommandHandlerTestSuite) TestHandle() {
	testUserID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	existingUser := entity.User{
		ID:             testUserID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          "existing@example.com",
		Password:       "password",
		Provider:       "github",
		Name:           "Existing User",
		FirstName:      "Existing",
		LastName:       "User",
		ProviderUserId: "provider123",
		AvatarURL:      "https://example.com/avatar.jpg",
	}

	tests := []struct {
		name          string
		command       CreateUserCommand
		setupMock     func()
		expectedError bool
		expectedSave  bool
	}{
		{
			name: "Success",
			command: NewCreateUserCommand(
				testUserID,
				"test@example.com",
				"password123",
				"github",
				"Test User",
				"Test",
				"User",
				"provider123",
				"https://example.com/avatar.jpg",
			),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(id uuid.UUID) (entity.User, error) {
					return entity.User{}, errors.New("user not found")
				}
				s.MockRepository.saveFunc = func(user entity.User) error {
					assert.Equal(s.T(), testUserID, user.ID)
					assert.Equal(s.T(), "test@example.com", user.Email)
					assert.Equal(s.T(), "password123", user.Password)
					assert.Equal(s.T(), "github", user.Provider)
					assert.Equal(s.T(), "Test User", user.Name)
					assert.Equal(s.T(), "Test", user.FirstName)
					assert.Equal(s.T(), "User", user.LastName)
					assert.Equal(s.T(), "provider123", user.ProviderUserId)
					assert.Equal(s.T(), "https://example.com/avatar.jpg", user.AvatarURL)
					return nil
				}
			},
			expectedError: false,
			expectedSave:  true,
		},
		{
			name: "UserAlreadyExists",
			command: NewCreateUserCommand(
				testUserID,
				"test@example.com",
				"password123",
				"github",
				"Test User",
				"Test",
				"User",
				"provider123",
				"https://example.com/avatar.jpg",
			),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(id uuid.UUID) (entity.User, error) {
					assert.Equal(s.T(), testUserID, id)
					return existingUser, nil
				}
				s.MockRepository.saveFunc = func(user entity.User) error {
					s.T().Error("Save should not be called when user already exists")
					return nil
				}
			},
			expectedError: false,
			expectedSave:  false,
		},
		{
			name: "SaveError",
			command: NewCreateUserCommand(
				testUserID,
				"test@example.com",
				"password123",
				"github",
				"Test User",
				"Test",
				"User",
				"provider123",
				"https://example.com/avatar.jpg",
			),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(id uuid.UUID) (entity.User, error) {
					return entity.User{}, errors.New("user not found")
				}
				s.MockRepository.saveFunc = func(user entity.User) error {
					return errors.New("database error")
				}
			},
			expectedError: true,
			expectedSave:  true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			ctx := context.Background()
			err := s.Handler.Handle(ctx, &tt.command)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateUserCommandHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CreateUserCommandHandlerTestSuite))
}
