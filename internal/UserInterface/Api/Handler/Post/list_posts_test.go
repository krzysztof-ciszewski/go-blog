package post

import (
	"database/sql"
	test "main/internal/Infrastructure/DependencyInjection/Test"
	"net/http"
	"net/http/httptest"
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
		"/api/v1/posts",
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
		"/api/v1/posts?text=first",
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
		"/api/v1/posts?author=author1",
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
		"/api/v1/posts?text=nonexistent",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	ListPosts(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Equal(s.T(), `[]`, s.W.Body.String())
}

func (s *ListPostsTestSuite) TestListPostsByAuthorEmptyResult() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts?author=nonexistent",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	ListPosts(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Equal(s.T(), `[]`, s.W.Body.String())
}

func (s *ListPostsTestSuite) TestListPostsTextTakesPrecedenceOverAuthor() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts?text=first&author=author2",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	ListPosts(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	// When text is provided, it should take precedence over author
	assert.Contains(s.T(), s.W.Body.String(), `"slug":"slug1"`)
	assert.Contains(s.T(), s.W.Body.String(), `"title":"First Post"`)
}

func TestListPostsTestSuite(t *testing.T) {
	suite.Run(t, new(ListPostsTestSuite))
}
