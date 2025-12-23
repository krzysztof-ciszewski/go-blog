package post

import (
	"bytes"
	"database/sql"
	test "main/internal/Infrastructure/DependencyInjection/Test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type CreatePostTestSuite struct {
	suite.Suite
	CommandBus *cqrs.CommandBus
	Ctx        *gin.Context
	W          *httptest.ResponseRecorder
	PubSubDb   *sql.DB
}

func (s *CreatePostTestSuite) SetupTest() {
	s.CommandBus = test.GetTestContainer().CommandBus
	s.W = httptest.NewRecorder()
	s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
	gin.SetMode(gin.TestMode)
	s.PubSubDb = test.GetPubSubDb()
	s.PubSubDb.Exec("DELETE FROM `watermill_commands.createPostCommand`")
}

func (s *CreatePostTestSuite) TestCreatePost() {
	s.Ctx.Request = httptest.NewRequest(
		"POST",
		"/api/v1/posts",
		bytes.NewBufferString(`{
			"id": "123e4567-e89b-12d3-a456-426614174000",
			"slug": "testslug",
			"title": "testtitle",
			"content": "testcontent",
			"author": "testauthor"
		}`))
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	CreatePost(s.Ctx, s.CommandBus)

	assert.Equal(s.T(), http.StatusAccepted, s.W.Code)
	assert.Equal(s.T(), `{"message":"Post created"}`, s.W.Body.String())

	count := test.GetCommandCount("createPostCommand")
	assert.Equal(s.T(), 1, count)
}

func (s *CreatePostTestSuite) TestCreatePostInvalidSlugRequest() {
	s.Ctx.Request = httptest.NewRequest(
		"POST",
		"/api/v1/posts",
		bytes.NewBufferString(`{
			"id": "123e4567-e89b-12d3-a456-426614174000",
			"slug": "test-slug",
			"title": "testtitle",
			"content": "testcontent",
			"author": "testauthor"
		}`))
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	CreatePost(s.Ctx, s.CommandBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Key: 'CreatePostRequest.Slug' Error:Field validation for 'Slug' failed on the 'alphanum' tag"}`, s.W.Body.String())
	count := test.GetCommandCount("createPostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *CreatePostTestSuite) TestCreatePostInvalidTitleRequest() {
	s.Ctx.Request = httptest.NewRequest(
		"POST",
		"/api/v1/posts",
		bytes.NewBufferString(`{
			"id": "123e4567-e89b-12d3-a456-426614174000",
			"slug": "testslug",
			"title": "t",
			"content": "testcontent",
			"author": "testauthor"
		}`))
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	CreatePost(s.Ctx, s.CommandBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Key: 'CreatePostRequest.Title' Error:Field validation for 'Title' failed on the 'min' tag"}`, s.W.Body.String())
	count := test.GetCommandCount("createPostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *CreatePostTestSuite) TestCreatePostInvalidContentRequest() {
	s.Ctx.Request = httptest.NewRequest(
		"POST",
		"/api/v1/posts",
		bytes.NewBufferString(`{
			"id": "123e4567-e89b-12d3-a456-426614174000",
			"slug": "testslug",
			"title": "testtitle",
			"content": "t",
			"author": "testauthor"
		}`))
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	CreatePost(s.Ctx, s.CommandBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Key: 'CreatePostRequest.Content' Error:Field validation for 'Content' failed on the 'min' tag"}`, s.W.Body.String())
	count := test.GetCommandCount("createPostCommand")
	assert.Equal(s.T(), 0, count)
}

func TestCreatePostTestSuite(t *testing.T) {
	suite.Run(t, new(CreatePostTestSuite))
}
