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

type mockPostRepositoryForFindAll struct {
	findAllByFunc func(page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error)
}

func (m *mockPostRepositoryForFindAll) Save(post entity.Post) error {
	return nil
}

func (m *mockPostRepositoryForFindAll) FindByID(id uuid.UUID) (entity.Post, error) {
	return entity.Post{}, nil
}

func (m *mockPostRepositoryForFindAll) FindAllBy(page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
	if m.findAllByFunc != nil {
		return m.findAllByFunc(page, pageSize, slug, text, author)
	}
	return repository.PaginatedResult[entity.Post]{}, errors.New("not implemented")
}

func (m *mockPostRepositoryForFindAll) Delete(id uuid.UUID) error {
	return nil
}

type FindAllByQueryHandlerTestSuite struct {
	suite.Suite
	Handler        FindAllByQueryHandler
	MockRepository *mockPostRepositoryForFindAll
}

func (s *FindAllByQueryHandlerTestSuite) SetupTest() {
	s.MockRepository = &mockPostRepositoryForFindAll{}
	s.Handler = FindAllByQueryHandler{
		PostRepository: s.MockRepository,
	}
}

func (s *FindAllByQueryHandlerTestSuite) TestHandle() {
	testPostID1 := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	testPostID2 := uuid.MustParse("223e4567-e89b-12d3-a456-426614174001")
	testPostID3 := uuid.MustParse("323e4567-e89b-12d3-a456-426614174002")
	testAuthorID := uuid.MustParse("423e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name               string
		query              any
		setupMock          func()
		expectedError      bool
		expectedTotal      int64
		expectedPage       int
		expectedPageSize   int
		expectedItemsLen   int
		expectedItems      []entity.Post
		expectedResultType string
	}{
		{
			name:  "Success",
			query: NewFindAllByQuery(1, 10, "", "", ""),
			setupMock: func() {
				testPosts := []entity.Post{
					{
						ID:        testPostID1,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Slug:      "test-slug-1",
						Title:     "Test Title 1",
						Content:   "Test Content 1",
						AuthorId:  testAuthorID,
					},
					{
						ID:        testPostID2,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Slug:      "test-slug-2",
						Title:     "Test Title 2",
						Content:   "Test Content 2",
						AuthorId:  testAuthorID,
					},
				}
				s.MockRepository.findAllByFunc = func(page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
					assert.Equal(s.T(), 1, page)
					assert.Equal(s.T(), 10, pageSize)
					assert.Equal(s.T(), "", slug)
					assert.Equal(s.T(), "", text)
					assert.Equal(s.T(), "", author)
					return repository.PaginatedResult[entity.Post]{
						Items:    testPosts,
						Total:    2,
						Page:     1,
						PageSize: 10,
					}, nil
				}
			},
			expectedError:    false,
			expectedTotal:    2,
			expectedPage:     1,
			expectedPageSize: 10,
			expectedItemsLen: 2,
			expectedItems: []entity.Post{
				{ID: testPostID1, Slug: "test-slug-1", Title: "Test Title 1", Content: "Test Content 1", AuthorId: testAuthorID},
				{ID: testPostID2, Slug: "test-slug-2", Title: "Test Title 2", Content: "Test Content 2", AuthorId: testAuthorID},
			},
			expectedResultType: "PaginatedView",
		},
		{
			name:  "WithFilters",
			query: NewFindAllByQuery(2, 20, "test-slug", "search text", "author-name"),
			setupMock: func() {
				testPosts := []entity.Post{
					{
						ID:        testPostID3,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
						Slug:      "filtered-slug",
						Title:     "Filtered Title",
						Content:   "Filtered Content",
						AuthorId:  testAuthorID,
					},
				}
				s.MockRepository.findAllByFunc = func(page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
					assert.Equal(s.T(), 2, page)
					assert.Equal(s.T(), 20, pageSize)
					assert.Equal(s.T(), "test-slug", slug)
					assert.Equal(s.T(), "search text", text)
					assert.Equal(s.T(), "author-name", author)
					return repository.PaginatedResult[entity.Post]{
						Items:    testPosts,
						Total:    1,
						Page:     2,
						PageSize: 20,
					}, nil
				}
			},
			expectedError:    false,
			expectedTotal:    1,
			expectedPage:     2,
			expectedPageSize: 20,
			expectedItemsLen: 1,
			expectedItems: []entity.Post{
				{ID: testPostID3, Slug: "filtered-slug", Title: "Filtered Title", Content: "Filtered Content", AuthorId: testAuthorID},
			},
			expectedResultType: "PaginatedView",
		},
		{
			name:  "EmptyResult",
			query: NewFindAllByQuery(1, 10, "", "", ""),
			setupMock: func() {
				s.MockRepository.findAllByFunc = func(page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
					return repository.PaginatedResult[entity.Post]{
						Items:    []entity.Post{},
						Total:    0,
						Page:     1,
						PageSize: 10,
					}, nil
				}
			},
			expectedError:      false,
			expectedTotal:      0,
			expectedPage:       1,
			expectedPageSize:   10,
			expectedItemsLen:   0,
			expectedItems:      []entity.Post{},
			expectedResultType: "PaginatedView",
		},
		{
			name:  "RepositoryError",
			query: NewFindAllByQuery(1, 10, "", "", ""),
			setupMock: func() {
				s.MockRepository.findAllByFunc = func(page int, pageSize int, slug string, text string, author string) (repository.PaginatedResult[entity.Post], error) {
					return repository.PaginatedResult[entity.Post]{}, errors.New("database error")
				}
			},
			expectedError:      true,
			expectedResultType: "Slice",
		},
		{
			name:               "InvalidQueryType",
			query:              "not a FindAllByQuery",
			setupMock:          func() {},
			expectedError:      false,
			expectedResultType: "Slice",
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

			switch tt.expectedResultType {
			case "PaginatedView":
				paginatedView, ok := result.(view.PaginatedView[view.PostView])
				assert.True(t, ok)
				assert.Equal(t, tt.expectedTotal, paginatedView.Total)
				assert.Equal(t, tt.expectedPage, paginatedView.Page)
				assert.Equal(t, tt.expectedPageSize, paginatedView.PageSize)
				assert.Len(t, paginatedView.Items, tt.expectedItemsLen)
				for i, expectedPost := range tt.expectedItems {
					if i < len(paginatedView.Items) {
						assert.Equal(t, expectedPost.ID, paginatedView.Items[i].Id)
						assert.Equal(t, expectedPost.Slug, paginatedView.Items[i].Slug)
						assert.Equal(t, expectedPost.Title, paginatedView.Items[i].Title)
						assert.Equal(t, expectedPost.Content, paginatedView.Items[i].Content)
						assert.Equal(t, expectedPost.AuthorId, paginatedView.Items[i].AuthorId)
					}
				}
			case "Slice":
				postViews, ok := result.([]view.PostView)
				assert.True(t, ok)
				assert.Len(t, postViews, 0)
			}
		})
	}
}

func (s *FindAllByQueryHandlerTestSuite) TestSupports() {
	tests := []struct {
		name          string
		query         any
		expectedValue bool
	}{
		{
			name:          "ValidQuery",
			query:         NewFindAllByQuery(1, 10, "", "", ""),
			expectedValue: true,
		},
		{
			name:          "InvalidQuery",
			query:         "not a FindAllByQuery",
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

func TestFindAllByQueryHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(FindAllByQueryHandlerTestSuite))
}
