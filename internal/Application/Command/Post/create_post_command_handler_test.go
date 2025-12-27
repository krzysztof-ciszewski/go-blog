package command

import (
	"context"
	"database/sql"
	"errors"
	entity "main/internal/Domain/Entity"
	event "main/internal/Domain/Event"
	repository "main/internal/Domain/Repository"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wmsqlitemodernc "github.com/ThreeDotsLabs/watermill-sqlite/wmsqlitemodernc"
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

func (m *mockPostRepositoryCreate) Update(post entity.Post) error {
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
	Handler         CreatePostCommandHandler
	MockRepository  *mockPostRepositoryCreate
	EventBus        *cqrs.EventBus
	PublishedEvents []interface{}
}

func (s *CreatePostCommandHandlerTestSuite) SetupTest() {
	s.MockRepository = &mockPostRepositoryCreate{}
	s.PublishedEvents = make([]interface{}, 0)

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
		OnPublish: func(params cqrs.OnEventSendParams) error {
			s.PublishedEvents = append(s.PublishedEvents, params.Event)
			return nil
		},
		Marshaler: marshaller,
		Logger:    watermill.NopLogger{},
	})
	if err != nil {
		panic(err)
	}
	s.EventBus = eventBus

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
		name            string
		command         createPostCommand
		setupMock       func()
		expectedError   bool
		expectedSave    bool
		expectedSaveID  uuid.UUID
		expectedPublish bool
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
			expectedError:   false,
			expectedSave:    true,
			expectedSaveID:  testPostID,
			expectedPublish: true,
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
			expectedError:   false,
			expectedSave:    false,
			expectedPublish: false,
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
			expectedError:   true,
			expectedSave:    true,
			expectedPublish: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			s.PublishedEvents = make([]interface{}, 0)
			tt.setupMock()

			ctx := context.Background()
			err := s.Handler.Handle(ctx, &tt.command)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedPublish {
				assert.Greater(t, len(s.PublishedEvents), 0)
				if len(s.PublishedEvents) > 0 {
					publishedEvent, ok := s.PublishedEvents[0].(event.PostWasCreated)
					assert.True(t, ok)
					assert.Equal(t, testPostID, publishedEvent.ID)
					assert.Equal(t, "test-slug", publishedEvent.Slug)
					assert.Equal(t, "Test Title", publishedEvent.Title)
					assert.Equal(t, "Test Content", publishedEvent.Content)
					assert.Equal(t, testAuthorID, publishedEvent.AuthorId)
					assert.False(t, publishedEvent.CreatedAt.IsZero())
					assert.False(t, publishedEvent.UpdatedAt.IsZero())
				}
			} else {
				assert.Equal(t, 0, len(s.PublishedEvents))
			}
		})
	}
}

func TestCreatePostCommandHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(CreatePostCommandHandlerTestSuite))
}
