package post

import (
	"database/sql"
	test "main/internal/Infrastructure/DependencyInjection/Test"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	query_bus "main/internal/Infrastructure/QueryBus"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ListPostsTestSuite struct {
	suite.Suite
	QueryBus  query_bus.QueryBus
	Ctx       *gin.Context
	W         *httptest.ResponseRecorder
	PubSubDb  *sql.DB
	PostUuid1 uuid.UUID
	PostUuid2 uuid.UUID
	PostUuid3 uuid.UUID
}

func (s *ListPostsTestSuite) SetupTest() {
	s.QueryBus = test.GetTestContainer().QueryBus
	s.W = httptest.NewRecorder()
	s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
	gin.SetMode(gin.TestMode)
	s.PubSubDb = test.GetPubSubDb()
	test.GetTestContainer().DB.Exec("DELETE FROM posts")
	postUuid1, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	s.PostUuid1 = postUuid1
	test.GetTestContainer().DB.Exec("INSERT INTO posts (id, created_at, updated_at, slug, title, content, author) VALUES ($1, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'slug1', 'First Post', 'This is the first post content', 'author1')", postUuid1.String())

	postUuid2, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	s.PostUuid2 = postUuid2
	test.GetTestContainer().DB.Exec("INSERT INTO posts (id, created_at, updated_at, slug, title, content, author) VALUES ($1, '2021-01-02 00:00:00', '2021-01-02 00:00:00', 'slug2', 'Second Post', 'This is the second post content', 'author2')", postUuid2.String())

	postUuid3, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	s.PostUuid3 = postUuid3
	test.GetTestContainer().DB.Exec("INSERT INTO posts (id, created_at, updated_at, slug, title, content, author) VALUES ($1, '2021-01-03 00:00:00', '2021-01-03 00:00:00', 'slug3', 'Third Post', 'This is the third post content', 'author1')", postUuid3.String())
}

func (s *ListPostsTestSuite) TestListPosts() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts?page=1&pageSize=10",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	ListPosts(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Contains(s.T(), s.W.Body.String(), `"slug":"slug1"`)
	assert.Contains(s.T(), s.W.Body.String(), `"slug":"slug2"`)
	assert.Contains(s.T(), s.W.Body.String(), `"slug":"slug3"`)
	assert.Contains(s.T(), s.W.Body.String(), `"title":"First Post"`)
	assert.Contains(s.T(), s.W.Body.String(), `"title":"Second Post"`)
	assert.Contains(s.T(), s.W.Body.String(), `"title":"Third Post"`)
}

func (s *ListPostsTestSuite) TestListPostsByText() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts?page=1&pageSize=10&text=first",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	ListPosts(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Contains(s.T(), s.W.Body.String(), `"slug":"slug1"`)
	assert.Contains(s.T(), s.W.Body.String(), `"title":"First Post"`)
	assert.NotContains(s.T(), s.W.Body.String(), `"slug":"slug2"`)
	assert.NotContains(s.T(), s.W.Body.String(), `"slug":"slug3"`)
}

func (s *ListPostsTestSuite) TestListPostsByAuthor() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts?page=1&pageSize=10&author=author1",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	ListPosts(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Contains(s.T(), s.W.Body.String(), `"slug":"slug1"`)
	assert.Contains(s.T(), s.W.Body.String(), `"slug":"slug3"`)
	assert.Contains(s.T(), s.W.Body.String(), `"author":"author1"`)
	assert.NotContains(s.T(), s.W.Body.String(), `"slug":"slug2"`)
}

func (s *ListPostsTestSuite) TestListPostsByTextEmptyResult() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts?page=1&pageSize=10&text=nonexistent",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	ListPosts(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Equal(s.T(), `{"items":[],"total":0,"page":1,"page_size":10}`, s.W.Body.String())
}

func (s *ListPostsTestSuite) TestListPostsByAuthorEmptyResult() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts?page=1&pageSize=10&author=nonexistent",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	ListPosts(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Equal(s.T(), `{"items":[],"total":0,"page":1,"page_size":10}`, s.W.Body.String())
}

func (s *ListPostsTestSuite) TestListPostsPagination() {
	tests := []struct {
		name             string
		page             int
		pageSize         int
		expectedTotal    int64
		expectedCount    int
		expectedPage     int
		expectedSlugs    []string
		notExpectedSlugs []string
	}{
		{
			name:             "First page with page size 1",
			page:             1,
			pageSize:         1,
			expectedTotal:    3,
			expectedCount:    1,
			expectedPage:     1,
			expectedSlugs:    []string{"slug1"},
			notExpectedSlugs: []string{"slug2", "slug3"},
		},
		{
			name:             "Second page with page size 1",
			page:             2,
			pageSize:         1,
			expectedTotal:    3,
			expectedCount:    1,
			expectedPage:     2,
			expectedSlugs:    []string{"slug2"},
			notExpectedSlugs: []string{"slug1", "slug3"},
		},
		{
			name:             "Third page with page size 1",
			page:             3,
			pageSize:         1,
			expectedTotal:    3,
			expectedCount:    1,
			expectedPage:     3,
			expectedSlugs:    []string{"slug3"},
			notExpectedSlugs: []string{"slug1", "slug2"},
		},
		{
			name:             "First page with page size 2",
			page:             1,
			pageSize:         2,
			expectedTotal:    3,
			expectedCount:    2,
			expectedPage:     1,
			expectedSlugs:    []string{"slug1", "slug2"},
			notExpectedSlugs: []string{"slug3"},
		},
		{
			name:             "Second page with page size 2",
			page:             2,
			pageSize:         2,
			expectedTotal:    3,
			expectedCount:    1,
			expectedPage:     2,
			expectedSlugs:    []string{"slug3"},
			notExpectedSlugs: []string{"slug1", "slug2"},
		},
		{
			name:             "Page beyond available data",
			page:             10,
			pageSize:         10,
			expectedTotal:    3,
			expectedCount:    0,
			expectedPage:     10,
			expectedSlugs:    []string{},
			notExpectedSlugs: []string{"slug1", "slug2", "slug3"},
		},
		{
			name:             "Large page size covers all",
			page:             1,
			pageSize:         100,
			expectedTotal:    3,
			expectedCount:    3,
			expectedPage:     1,
			expectedSlugs:    []string{"slug1", "slug2", "slug3"},
			notExpectedSlugs: []string{},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.W = httptest.NewRecorder()
			s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
			s.Ctx.Request = httptest.NewRequest(
				"GET",
				"/api/v1/posts?page="+strconv.Itoa(tt.page)+"&pageSize="+strconv.Itoa(tt.pageSize),
				nil,
			)
			s.Ctx.Request.Header.Set("Content-Type", "application/json")

			ListPosts(s.Ctx, s.QueryBus)

			assert.Equal(s.T(), http.StatusOK, s.W.Code)
			assert.Contains(s.T(), s.W.Body.String(), `"total":`+strconv.FormatInt(tt.expectedTotal, 10))
			assert.Contains(s.T(), s.W.Body.String(), `"page":`+strconv.Itoa(tt.expectedPage))
			assert.Contains(s.T(), s.W.Body.String(), `"page_size":`+strconv.Itoa(tt.pageSize))

			for _, slug := range tt.expectedSlugs {
				assert.Contains(s.T(), s.W.Body.String(), `"slug":"`+slug+`"`)
			}

			for _, slug := range tt.notExpectedSlugs {
				assert.NotContains(s.T(), s.W.Body.String(), `"slug":"`+slug+`"`)
			}
		})
	}
}

func (s *ListPostsTestSuite) TestListPostsMultiFieldFiltering() {
	tests := []struct {
		name             string
		queryParams      map[string]string
		expectedSlugs    []string
		notExpectedSlugs []string
		expectedTotal    int64
	}{
		{
			name: "Filter by text and author",
			queryParams: map[string]string{
				"text":   "first",
				"author": "author1",
			},
			expectedSlugs:    []string{"slug1"},
			notExpectedSlugs: []string{"slug2", "slug3"},
			expectedTotal:    1,
		},
		{
			name: "Filter by text and author - no match",
			queryParams: map[string]string{
				"text":   "second",
				"author": "author1",
			},
			expectedSlugs:    []string{},
			notExpectedSlugs: []string{"slug1", "slug2", "slug3"},
			expectedTotal:    0,
		},
		{
			name: "Filter by slug and author",
			queryParams: map[string]string{
				"slug":   "slug1",
				"author": "author1",
			},
			expectedSlugs:    []string{"slug1"},
			notExpectedSlugs: []string{"slug2", "slug3"},
			expectedTotal:    1,
		},
		{
			name: "Filter by slug and author - no match",
			queryParams: map[string]string{
				"slug":   "slug2",
				"author": "author1",
			},
			expectedSlugs:    []string{},
			notExpectedSlugs: []string{"slug1", "slug2", "slug3"},
			expectedTotal:    0,
		},
		{
			name: "Filter by text and slug",
			queryParams: map[string]string{
				"text": "first",
				"slug": "slug1",
			},
			expectedSlugs:    []string{"slug1"},
			notExpectedSlugs: []string{"slug2", "slug3"},
			expectedTotal:    1,
		},
		{
			name: "Filter by text, slug and author",
			queryParams: map[string]string{
				"text":   "first",
				"slug":   "slug1",
				"author": "author1",
			},
			expectedSlugs:    []string{"slug1"},
			notExpectedSlugs: []string{"slug2", "slug3"},
			expectedTotal:    1,
		},
		{
			name: "Filter by text, slug and author - no match",
			queryParams: map[string]string{
				"text":   "second",
				"slug":   "slug1",
				"author": "author1",
			},
			expectedSlugs:    []string{},
			notExpectedSlugs: []string{"slug1", "slug2", "slug3"},
			expectedTotal:    0,
		},
		{
			name: "Filter by text matching multiple posts with author filter",
			queryParams: map[string]string{
				"text":   "post",
				"author": "author1",
			},
			expectedSlugs:    []string{"slug1", "slug3"},
			notExpectedSlugs: []string{"slug2"},
			expectedTotal:    2,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.W = httptest.NewRecorder()
			s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())

			queryString := "page=1&pageSize=10"
			for key, value := range tt.queryParams {
				queryString += "&" + key + "=" + value
			}

			s.Ctx.Request = httptest.NewRequest(
				"GET",
				"/api/v1/posts?"+queryString,
				nil,
			)
			s.Ctx.Request.Header.Set("Content-Type", "application/json")

			ListPosts(s.Ctx, s.QueryBus)

			assert.Equal(s.T(), http.StatusOK, s.W.Code)
			assert.Contains(s.T(), s.W.Body.String(), `"total":`+strconv.FormatInt(tt.expectedTotal, 10))

			for _, slug := range tt.expectedSlugs {
				assert.Contains(s.T(), s.W.Body.String(), `"slug":"`+slug+`"`)
			}

			for _, slug := range tt.notExpectedSlugs {
				assert.NotContains(s.T(), s.W.Body.String(), `"slug":"`+slug+`"`)
			}
		})
	}
}

func (s *ListPostsTestSuite) TestListPostsPaginationWithFilters() {
	tests := []struct {
		name             string
		page             int
		pageSize         int
		text             string
		author           string
		slug             string
		expectedTotal    int64
		expectedCount    int
		expectedSlugs    []string
		notExpectedSlugs []string
	}{
		{
			name:             "First page with filter by author",
			page:             1,
			pageSize:         1,
			author:           "author1",
			expectedTotal:    2,
			expectedCount:    1,
			expectedSlugs:    []string{"slug1"},
			notExpectedSlugs: []string{"slug2", "slug3"},
		},
		{
			name:             "Second page with filter by author",
			page:             2,
			pageSize:         1,
			author:           "author1",
			expectedTotal:    2,
			expectedCount:    1,
			expectedSlugs:    []string{"slug3"},
			notExpectedSlugs: []string{"slug1", "slug2"},
		},
		{
			name:             "First page with filter by text",
			page:             1,
			pageSize:         1,
			text:             "post",
			expectedTotal:    3,
			expectedCount:    1,
			expectedSlugs:    []string{"slug1"},
			notExpectedSlugs: []string{"slug2", "slug3"},
		},
		{
			name:             "Pagination with text and author filter",
			page:             1,
			pageSize:         1,
			text:             "post",
			author:           "author1",
			expectedTotal:    2,
			expectedCount:    1,
			expectedSlugs:    []string{"slug1"},
			notExpectedSlugs: []string{"slug2", "slug3"},
		},
		{
			name:             "Pagination with text and author filter - second page",
			page:             2,
			pageSize:         1,
			text:             "post",
			author:           "author1",
			expectedTotal:    2,
			expectedCount:    1,
			expectedSlugs:    []string{"slug3"},
			notExpectedSlugs: []string{"slug1", "slug2"},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.W = httptest.NewRecorder()
			s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())

			queryString := "page=" + strconv.Itoa(tt.page) + "&pageSize=" + strconv.Itoa(tt.pageSize)
			if tt.text != "" {
				queryString += "&text=" + tt.text
			}
			if tt.author != "" {
				queryString += "&author=" + tt.author
			}
			if tt.slug != "" {
				queryString += "&slug=" + tt.slug
			}

			s.Ctx.Request = httptest.NewRequest(
				"GET",
				"/api/v1/posts?"+queryString,
				nil,
			)
			s.Ctx.Request.Header.Set("Content-Type", "application/json")

			ListPosts(s.Ctx, s.QueryBus)

			assert.Equal(s.T(), http.StatusOK, s.W.Code)
			assert.Contains(s.T(), s.W.Body.String(), `"total":`+strconv.FormatInt(tt.expectedTotal, 10))
			assert.Contains(s.T(), s.W.Body.String(), `"page":`+strconv.Itoa(tt.page))
			assert.Contains(s.T(), s.W.Body.String(), `"page_size":`+strconv.Itoa(tt.pageSize))

			for _, slug := range tt.expectedSlugs {
				assert.Contains(s.T(), s.W.Body.String(), `"slug":"`+slug+`"`)
			}

			for _, slug := range tt.notExpectedSlugs {
				assert.NotContains(s.T(), s.W.Body.String(), `"slug":"`+slug+`"`)
			}
		})
	}
}

func TestListPostsTestSuite(t *testing.T) {
	suite.Run(t, new(ListPostsTestSuite))
}
