package post_query

import (
	"context"
	"errors"
	view "main/internal/Application/View"
	entity "main/internal/Domain/Entity"
	repository "main/internal/Domain/Repository"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type mockPostRepository struct {
	findByIDFunc func(ctx context.Context, id uuid.UUID) (entity.Post, error)
}

func (m *mockPostRepository) Save(ctx context.Context, post entity.Post) error {
	return nil
}

func (m *mockPostRepository) Update(ctx context.Context, post entity.Post) error {
	return nil
}

func (m *mockPostRepository) FindByID(ctx context.Context, id uuid.UUID) (entity.Post, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return entity.Post{}, errors.New("not implemented")
}

func (m *mockPostRepository) FindAllBy(ctx context.Context, page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
	return repository.PaginatedResult[entity.Post]{}, nil
}

func (m *mockPostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

type GetPostQueryHandlerTestSuite struct {
	suite.Suite
	Handler        GetPostQueryHandler
	MockRepository *mockPostRepository
}

func (s *GetPostQueryHandlerTestSuite) SetupTest() {
	s.MockRepository = &mockPostRepository{}
	s.Handler = GetPostQueryHandler{
		PostRepository: s.MockRepository,
	}
}

func (s *GetPostQueryHandlerTestSuite) TestHandle() {
	testPostID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	testAuthorID := uuid.MustParse("223e4567-e89b-12d3-a456-426614174001")
	testPost := entity.Post{
		ID:        testPostID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Slug:      "test-slug",
		Title:     "Test Title",
		Content:   "Test Content",
		AuthorId:  testAuthorID,
	}

	tests := []struct {
		name             string
		query            any
		setupMock        func()
		expectedError    bool
		expectedPostID   uuid.UUID
		expectedSlug     string
		expectedTitle    string
		expectedContent  string
		expectedAuthorID uuid.UUID
	}{
		{
			name:  "Success",
			query: NewGetPostQuery(testPostID),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(ctx context.Context, id uuid.UUID) (entity.Post, error) {
					assert.Equal(s.T(), testPostID, id)
					return testPost, nil
				}
			},
			expectedError:    false,
			expectedPostID:   testPostID,
			expectedSlug:     "test-slug",
			expectedTitle:    "Test Title",
			expectedContent:  "Test Content",
			expectedAuthorID: testAuthorID,
		},
		{
			name:  "PostNotFound",
			query: NewGetPostQuery(testPostID),
			setupMock: func() {
				s.MockRepository.findByIDFunc = func(ctx context.Context, id uuid.UUID) (entity.Post, error) {
					return entity.Post{}, errors.New("post not found")
				}
			},
			expectedError:  true,
			expectedPostID: uuid.Nil,
		},
		{
			name:           "InvalidQueryType",
			query:          "not a GetPostQuery",
			setupMock:      func() {},
			expectedError:  false,
			expectedPostID: uuid.Nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			ctx := context.Background()
			result, err := s.Handler.Handle(ctx, tt.query)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			postView, ok := result.(view.PostView)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedPostID, postView.Id)
			if tt.expectedSlug != "" {
				assert.Equal(t, tt.expectedSlug, postView.Slug)
				assert.Equal(t, tt.expectedTitle, postView.Title)
				assert.Equal(t, tt.expectedContent, postView.Content)
				assert.Equal(t, tt.expectedAuthorID, postView.AuthorId)
			}
		})
	}
}

func (s *GetPostQueryHandlerTestSuite) TestSupports() {
	testPostID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name          string
		query         any
		expectedValue bool
	}{
		{
			name:          "ValidQuery",
			query:         NewGetPostQuery(testPostID),
			expectedValue: true,
		},
		{
			name:          "InvalidQuery",
			query:         "not a GetPostQuery",
			expectedValue: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			supports := s.Handler.Supports(tt.query)
			assert.Equal(t, tt.expectedValue, supports)
		})
	}
}

func TestGetPostQueryHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GetPostQueryHandlerTestSuite))
}
