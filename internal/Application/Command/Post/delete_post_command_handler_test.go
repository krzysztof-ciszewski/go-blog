package command

import (
	"context"
	"errors"
	entity "main/internal/Domain/Entity"
	repository "main/internal/Domain/Repository"
	"testing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type mockPostRepositoryDelete struct {
	deleteFunc func(id uuid.UUID) error
}

func (m *mockPostRepositoryDelete) Save(post entity.Post) error {
	return nil
}

func (m *mockPostRepositoryDelete) FindByID(id uuid.UUID) (entity.Post, error) {
	return entity.Post{}, errors.New("not implemented")
}

func (m *mockPostRepositoryDelete) FindAllBy(page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
	return repository.PaginatedResult[entity.Post]{}, nil
}

func (m *mockPostRepositoryDelete) Delete(id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

type DeletePostCommandHandlerTestSuite struct {
	suite.Suite
	Handler        DeletePostCommandHandler
	MockRepository *mockPostRepositoryDelete
	EventBus       *cqrs.EventBus
}

func (s *DeletePostCommandHandlerTestSuite) SetupTest() {
	s.MockRepository = &mockPostRepositoryDelete{}
	s.EventBus = nil
	s.Handler = DeletePostCommandHandler{
		EventBus:       s.EventBus,
		PostRepository: s.MockRepository,
	}
}

func (s *DeletePostCommandHandlerTestSuite) TestHandle() {
	testPostID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name          string
		command       deletePostCommand
		setupMock     func()
		expectedError bool
		expectedID    uuid.UUID
	}{
		{
			name:    "Success",
			command: NewDeletePostCommand(testPostID),
			setupMock: func() {
				s.MockRepository.deleteFunc = func(id uuid.UUID) error {
					assert.Equal(s.T(), testPostID, id)
					return nil
				}
			},
			expectedError: false,
			expectedID:    testPostID,
		},
		{
			name:    "DeleteError",
			command: NewDeletePostCommand(testPostID),
			setupMock: func() {
				s.MockRepository.deleteFunc = func(id uuid.UUID) error {
					assert.Equal(s.T(), testPostID, id)
					return errors.New("database error")
				}
			},
			expectedError: true,
			expectedID:    testPostID,
		},
		{
			name:    "PostNotFound",
			command: NewDeletePostCommand(testPostID),
			setupMock: func() {
				s.MockRepository.deleteFunc = func(id uuid.UUID) error {
					assert.Equal(s.T(), testPostID, id)
					return errors.New("post not found")
				}
			},
			expectedError: true,
			expectedID:    testPostID,
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

func TestDeletePostCommandHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DeletePostCommandHandlerTestSuite))
}
