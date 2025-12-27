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

type UpdatePostTestSuite struct {
	suite.Suite
	CommandBus *cqrs.CommandBus
	QueryBus   query_bus.QueryBus
	Ctx        *gin.Context
	W          *httptest.ResponseRecorder
	PubSubDb   *sql.DB
	PostUuid   uuid.UUID
	UserUuid   uuid.UUID
}

func (s *UpdatePostTestSuite) SetupTest() {
	s.CommandBus = test.GetTestContainer().CommandBus
	s.QueryBus = test.GetTestContainer().QueryBus
	s.W = httptest.NewRecorder()
	s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
	s.Ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/v1/posts/",
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
	s.PubSubDb.Exec("DELETE FROM `watermill_commands.updatePostCommand`")
	test.GetTestContainer().DB.Exec("DELETE FROM posts")
	test.GetTestContainer().DB.Exec("DELETE FROM users")
	userUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	s.UserUuid = userUuid
	test.GetTestContainer().DB.Exec(`
		INSERT INTO users (id, created_at, updated_at, provider, provider_user_id, email)
		VALUES (?, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'test', 'testprovideruser', 'test@example.com')
	`, userUuid.String())
	postUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	s.PostUuid = postUuid
	test.GetTestContainer().DB.Exec(`INSERT INTO posts (id, created_at, updated_at, slug, title, content, author_id)
	VALUES ($1, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'testslug', 'testtitle', 'testcontent', $2)`,
		postUuid.String(),
		userUuid.String(),
	)
}

func (s *UpdatePostTestSuite) TestUpdatePost() {
	s.Ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/v1/posts/"+s.PostUuid.String(),
		io.NopCloser(bytes.NewBufferString(`{
		"slug": "updatedslug",
		"title": "updatedtitle",
		"content": "updatedcontent"
	}`)),
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: s.PostUuid.String(),
		},
	}
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

	UpdatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusAccepted, s.W.Code)
	assert.Equal(s.T(), `{"message":"Post updated"}`, s.W.Body.String())

	count := test.GetCommandCount("updatePostCommand")
	assert.Equal(s.T(), 1, count)
}

func (s *UpdatePostTestSuite) TestUpdatePostInvalidPostId() {
	s.Ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/v1/posts/invalid-uuid",
		io.NopCloser(bytes.NewBufferString(`{
		"slug": "updatedslug",
		"title": "updatedtitle",
		"content": "updatedcontent"
	}`)),
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: "invalid-uuid",
		},
	}
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

	UpdatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Invalid post ID"}`, s.W.Body.String())

	count := test.GetCommandCount("updatePostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *UpdatePostTestSuite) TestUpdatePostInvalidSlugRequest() {
	s.Ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/v1/posts/"+s.PostUuid.String(),
		io.NopCloser(bytes.NewBufferString(`{
		"slug": "test-slug",
		"title": "updatedtitle",
		"content": "updatedcontent"
	}`)),
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: s.PostUuid.String(),
		},
	}
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

	UpdatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Key: 'UpdatePostRequest.Slug' Error:Field validation for 'Slug' failed on the 'alphanum' tag"}`, s.W.Body.String())

	count := test.GetCommandCount("updatePostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *UpdatePostTestSuite) TestUpdatePostInvalidTitleRequest() {
	s.Ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/v1/posts/"+s.PostUuid.String(),
		io.NopCloser(bytes.NewBufferString(`{
		"slug": "updatedslug",
		"title": "t",
		"content": "updatedcontent"
	}`)),
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: s.PostUuid.String(),
		},
	}
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

	UpdatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Key: 'UpdatePostRequest.Title' Error:Field validation for 'Title' failed on the 'min' tag"}`, s.W.Body.String())

	count := test.GetCommandCount("updatePostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *UpdatePostTestSuite) TestUpdatePostInvalidContentRequest() {
	s.Ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/v1/posts/"+s.PostUuid.String(),
		io.NopCloser(bytes.NewBufferString(`{
		"slug": "updatedslug",
		"title": "updatedtitle",
		"content": "t"
	}`)),
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: s.PostUuid.String(),
		},
	}
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

	UpdatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Key: 'UpdatePostRequest.Content' Error:Field validation for 'Content' failed on the 'min' tag"}`, s.W.Body.String())

	count := test.GetCommandCount("updatePostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *UpdatePostTestSuite) TestUpdatePostNotFound() {
	nonExistentUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	s.Ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/v1/posts/"+nonExistentUuid.String(),
		io.NopCloser(bytes.NewBufferString(`{
		"slug": "updatedslug",
		"title": "updatedtitle",
		"content": "updatedcontent"
	}`)),
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: nonExistentUuid.String(),
		},
	}
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

	UpdatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusNotFound, s.W.Code)
	assert.Equal(s.T(), `{"error":"Post not found"}`, s.W.Body.String())

	count := test.GetCommandCount("updatePostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *UpdatePostTestSuite) TestUpdatePostUnauthorized() {
	otherUserUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	// Create a different user who is NOT the author of the post
	test.GetTestContainer().DB.Exec(`
		INSERT INTO users (id, created_at, updated_at, provider, provider_user_id, email)
		VALUES (?, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'test', 'otherprovideruser', 'other@example.com')
	`, otherUserUuid.String())

	// Create a new request with a new response recorder for this test
	w := httptest.NewRecorder()
	ctx := gin.CreateTestContextOnly(w, gin.Default())
	ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/v1/posts/"+s.PostUuid.String(),
		io.NopCloser(bytes.NewBufferString(`{
		"slug": "updatedslug",
		"title": "updatedtitle",
		"content": "updatedcontent"
	}`)),
	)
	ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: s.PostUuid.String(),
		},
	}
	session, err := gothic.Store.New(ctx.Request, os.Getenv("SESSION_NAME"))
	if err != nil {
		panic(err)
	}
	session.Values["provider_user_id"] = "otherprovideruser"
	session.Values["email"] = "other@example.com"
	if err := session.Save(ctx.Request, w); err != nil {
		panic(err)
	}
	ctx.Request.Header.Set("Content-Type", "application/json")
	ctx.Request.Header.Set("Cookie", w.Header().Get("Set-Cookie"))

	UpdatePost(ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusForbidden, w.Code)
	assert.Equal(s.T(), `{"error":"You are not authorized to update this post"}`, w.Body.String())

	count := test.GetCommandCount("updatePostCommand")
	assert.Equal(s.T(), 0, count)
}

func (s *UpdatePostTestSuite) TestUpdatePostUnauthenticated() {
	s.Ctx.Request = httptest.NewRequest(
		"PUT",
		"/api/v1/posts/"+s.PostUuid.String(),
		io.NopCloser(bytes.NewBufferString(`{
		"slug": "updatedslug",
		"title": "updatedtitle",
		"content": "updatedcontent"
	}`)),
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: s.PostUuid.String(),
		},
	}
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	UpdatePost(s.Ctx, s.CommandBus, s.QueryBus)

	assert.Equal(s.T(), http.StatusUnauthorized, s.W.Code)
	assert.Equal(s.T(), `{"error":"User not authenticated"}`, s.W.Body.String())

	count := test.GetCommandCount("updatePostCommand")
	assert.Equal(s.T(), 0, count)
}

func TestUpdatePostTestSuite(t *testing.T) {
	suite.Run(t, new(UpdatePostTestSuite))
}
