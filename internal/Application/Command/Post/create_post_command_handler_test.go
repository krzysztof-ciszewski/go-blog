package command

import (
	"context"
	"errors"
	entity "main/internal/Domain/Entity"
	repository "main/internal/Domain/Repository"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type mockPostRepositoryCreate struct {
	saveFunc     func(post entity.Post) error
	findByIDFunc func(id uuid.UUID) (entity.Post, error)
}

func (m *mockPostRepositoryCreate) Save(post entity.Post) error {
	if m.saveFunc != nil {
		return m.saveFunc(post)
	}
	return nil
}

func (m *mockPostRepositoryCreate) FindByID(id uuid.UUID) (entity.Post, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
	return entity.Post{}, errors.New("not implemented")
}

func (m *mockPostRepositoryCreate) FindAllBy(page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
	return repository.PaginatedResult[entity.Post]{}, nil
}

func (m *mockPostRepositoryCreate) Delete(id uuid.UUID) error {
	return nil
}

type CreatePostCommandHandlerTestSuite struct {
	suite.Suite
	Handler        CreatePostCommandHandler
	MockRepository *mockPostRepositoryCreate
	EventBus       *cqrs.EventBus
}

func (s *CreatePostCommandHandlerTestSuite) SetupTest() {
	s.MockRepository = &mockPostRepositoryCreate{}
	s.EventBus = nil
	s.Handler = CreatePostCommandHandler{
		EventBus:       s.EventBus,
		PostRepository: s.MockRepository,
	}
}

func (s *CreatePostCommandHandlerTestSuite) TestHandle() {
	testPostID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	testAuthorID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174001")
	existingPost := entity.Post{
		ID:        testPostID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Slug:      "existing-slug",
		Title:     "Existing Title",
		Content:   "Existing Content",
		AuthorId:  testAuthorID,
	}

	tests := []struct {
		name           string
		command        createPostCommand
		setupMock      func()
		expectedError  bool
		expectedSave   bool
		expectedSaveID uuid.UUID
	}{
		{
			name: "Success",
			command: NewCreatePostCommand(
				testPostID,
				"test-slug",
				"Test Title",
				"Test Content",
				testAuthorID,
			),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(id uuid.UUID) (entity.Post, error) {
					return entity.Post{}, errors.New("post not found")
				}
				s.MockRepository.saveFunc = func(post entity.Post) error {
					assert.Equal(s.T(), testPostID, post.ID)
					assert.Equal(s.T(), "test-slug", post.Slug)
					assert.Equal(s.T(), "Test Title", post.Title)
					assert.Equal(s.T(), "Test Content", post.Content)
					assert.Equal(s.T(), testAuthorID, post.AuthorId)
					return nil
				}
			},
			expectedError:  false,
			expectedSave:   true,
			expectedSaveID: testPostID,
		},
		{
			name: "PostAlreadyExists",
			command: NewCreatePostCommand(
				testPostID,
				"test-slug",
				"Test Title",
				"Test Content",
				testAuthorID,
			),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(id uuid.UUID) (entity.Post, error) {
					assert.Equal(s.T(), testPostID, id)
					return existingPost, nil
				}
				s.MockRepository.saveFunc = func(post entity.Post) error {
					s.T().Error("Save should not be called when post already exists")
					return nil
				}
			},
			expectedError: false,
			expectedSave:  false,
		},
		{
			name: "SaveError",
			command: NewCreatePostCommand(
				testPostID,
				"test-slug",
				"Test Title",
				"Test Content",
				testAuthorID,
			),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(id uuid.UUID) (entity.Post, error) {
					return entity.Post{}, errors.New("post not found")
				}
				s.MockRepository.saveFunc = func(post entity.Post) error {
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

func TestCreatePostCommandHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CreatePostCommandHandlerTestSuite))
}
