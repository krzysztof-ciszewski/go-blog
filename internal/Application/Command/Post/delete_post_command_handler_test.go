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

type mockPostRepositoryDelete struct {
	deleteFunc   func(id uuid.UUID) error
	findByIDFunc func(id uuid.UUID) (entity.Post, error)
}

func (m *mockPostRepositoryDelete) Save(post entity.Post) error {
	return nil
}

func (m *mockPostRepositoryDelete) FindByID(id uuid.UUID) (entity.Post, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(id)
	}
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
	Handler         DeletePostCommandHandler
	MockRepository  *mockPostRepositoryDelete
	EventBus        *cqrs.EventBus
	PublishedEvents []interface{}
}

func (s *DeletePostCommandHandlerTestSuite) SetupTest() {
	s.MockRepository = &mockPostRepositoryDelete{}
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

	s.Handler = DeletePostCommandHandler{
		EventBus:       s.EventBus,
		PostRepository: s.MockRepository,
	}
}

func (s *DeletePostCommandHandlerTestSuite) TestHandle() {
	testPostID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	testAuthorID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174001")
	existingPost := entity.Post{
		ID:        testPostID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Slug:      "test-slug",
		Title:     "Test Title",
		Content:   "Test Content",
		AuthorId:  testAuthorID,
	}

	tests := []struct {
		name            string
		command         deletePostCommand
		setupMock       func()
		expectedError   bool
		expectedID      uuid.UUID
		expectedPublish bool
	}{
		{
			name:    "Success",
			command: NewDeletePostCommand(testPostID),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(id uuid.UUID) (entity.Post, error) {
					assert.Equal(s.T(), testPostID, id)
					return existingPost, nil
				}
				s.MockRepository.deleteFunc = func(id uuid.UUID) error {
					assert.Equal(s.T(), testPostID, id)
					return nil
				}
			},
			expectedError:   false,
			expectedID:      testPostID,
			expectedPublish: true,
		},
		{
			name:    "DeleteError",
			command: NewDeletePostCommand(testPostID),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(id uuid.UUID) (entity.Post, error) {
					assert.Equal(s.T(), testPostID, id)
					return existingPost, nil
				}
				s.MockRepository.deleteFunc = func(id uuid.UUID) error {
					assert.Equal(s.T(), testPostID, id)
					return errors.New("database error")
				}
			},
			expectedError:   true,
			expectedID:      testPostID,
			expectedPublish: false,
		},
		{
			name:    "PostNotFound",
			command: NewDeletePostCommand(testPostID),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(id uuid.UUID) (entity.Post, error) {
					assert.Equal(s.T(), testPostID, id)
					return entity.Post{}, errors.New("post not found")
				}
				s.MockRepository.deleteFunc = func(id uuid.UUID) error {
					s.T().Error("Delete should not be called when post is not found")
					return nil
				}
			},
			expectedError:   true,
			expectedID:      testPostID,
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
					publishedEvent, ok := s.PublishedEvents[0].(event.PostWasDeleted)
					assert.True(t, ok)
					assert.Equal(t, testPostID, publishedEvent.ID)
					assert.Equal(t, existingPost.CreatedAt, publishedEvent.CreatedAt)
					assert.Equal(t, existingPost.UpdatedAt, publishedEvent.UpdatedAt)
					assert.Equal(t, existingPost.Slug, publishedEvent.Slug)
					assert.Equal(t, existingPost.Title, publishedEvent.Title)
					assert.Equal(t, existingPost.Content, publishedEvent.Content)
					assert.Equal(t, existingPost.AuthorId, publishedEvent.AuthorId)
				}
			} else {
				assert.Equal(t, 0, len(s.PublishedEvents))
			}
		})
	}
}

func TestDeletePostCommandHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(DeletePostCommandHandlerTestSuite))
}
