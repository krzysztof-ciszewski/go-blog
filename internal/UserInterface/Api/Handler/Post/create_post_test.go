package post

import (
	"bytes"
	"database/sql"
	"io"
	test "main/internal/Infrastructure/DependencyInjection/Test"
	query_bus "main/internal/Infrastructure/QueryBus"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type CreatePostTestSuite struct {
	suite.Suite
	CommandBus *cqrs.CommandBus
	QueryBus   query_bus.QueryBus
	Ctx        *gin.Context
	W          *httptest.ResponseRecorder
	PubSubDb   *sql.DB
}

func (s *CreatePostTestSuite) SetupTest() {
	s.CommandBus = test.GetTestContainer().CommandBus
	s.QueryBus = test.GetTestContainer().QueryBus
	s.W = httptest.NewRecorder()
	s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
	s.Ctx.Request = httptest.NewRequest(
		"POST",
		"/api/v1/posts",
		nil,
	)
	session, err := gothic.Store.New(s.Ctx.Request, os.Getenv("SESSION_NAME"))
	if err != nil {
		panic(err)
	}
	session.Values["provider_user_id"] = "testprovideruser"
	session.Values["email"] = "test@example.com"
	if err := session.Save(s.Ctx.Request, s.Ctx.Writer); err != nil {
		panic(err)
	}
	s.Ctx.Request.Header.Set("Content-Type", "application/json")
	s.Ctx.Request.Header.Set("Cookie", s.Ctx.Writer.Header().Get("Set-Cookie"))

	gin.SetMode(gin.TestMode)
	s.PubSubDb = test.GetPubSubDb()
	s.PubSubDb.Exec("DELETE FROM `watermill_commands.createPostCommand`")
	test.GetTestContainer().DB.Exec("DELETE FROM users")
	userUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	test.GetTestContainer().DB.Exec(`
		INSERT INTO users (id, created_at, updated_at, provider, provider_user_id, email)
		VALUES (?, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'test', 'testprovideruser', 'test@example.com')
	`, userUuid.String())
}

func (s *CreatePostTestSuite) TestCreatePost() {
	s.Ctx.Request.Body = io.NopCloser(bytes.NewBufferString(`{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"slug": "testslug",
		"title": "testtitle",
		"content": "testcontent",
		"author": "testauthor"
	}`))

	CreatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusAccepted, s.W.Code)
	assert.Equal(s.T(), `{"message":"Post created"}`, s.W.Body.String())

	count := test.GetCommandCount("createPostCommand")
	assert.Equal(s.T(), 1, count)
}

func (s *CreatePostTestSuite) TestCreatePostInvalidSlugRequest() {
	s.Ctx.Request.Body = io.NopCloser(bytes.NewBufferString(`{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"slug": "test-slug",
		"title": "testtitle",
		"content": "testcontent",
		"author": "testauthor"
	}`))

	CreatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Key: 'CreatePostRequest.Slug' Error:Field validation for 'Slug' failed on the 'alphanum' tag"}`, s.W.Body.String())
	count := test.GetCommandCount("createPostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *CreatePostTestSuite) TestCreatePostInvalidTitleRequest() {
	s.Ctx.Request.Body = io.NopCloser(bytes.NewBufferString(`{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"slug": "testslug",
		"title": "t",
		"content": "testcontent",
		"author": "testauthor"
	}`))

	CreatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Key: 'CreatePostRequest.Title' Error:Field validation for 'Title' failed on the 'min' tag"}`, s.W.Body.String())
	count := test.GetCommandCount("createPostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *CreatePostTestSuite) TestCreatePostInvalidContentRequest() {
	s.Ctx.Request.Body = io.NopCloser(bytes.NewBufferString(`{
		"id": "123e4567-e89b-12d3-a456-426614174000",
		"slug": "testslug",
		"title": "testtitle",
		"content": "t",
		"author": "testauthor"
	}`))

	CreatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Key: 'CreatePostRequest.Content' Error:Field validation for 'Content' failed on the 'min' tag"}`, s.W.Body.String())
	count := test.GetCommandCount("createPostCommand")
	assert.Equal(s.T(), 0, count)
}

func TestCreatePostTestSuite(t *testing.T) {
	suite.Run(t, new(CreatePostTestSuite))
}
