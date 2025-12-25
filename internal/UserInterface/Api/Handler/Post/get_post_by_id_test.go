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

type GetPostByIdTestSuite struct {
	suite.Suite
	QueryBus query_bus.QueryBus
	Ctx      *gin.Context
	W        *httptest.ResponseRecorder
	PubSubDb *sql.DB
	PostUuid uuid.UUID
}

func (s *GetPostByIdTestSuite) SetupTest() {
	s.QueryBus = test.GetTestContainer().QueryBus
	s.W = httptest.NewRecorder()
	s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
	gin.SetMode(gin.TestMode)
	s.PubSubDb = test.GetPubSubDb()
	test.GetTestContainer().DB.Exec("DELETE FROM users")
	userUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	test.GetTestContainer().DB.Exec(`
		INSERT INTO users (id, created_at, updated_at, provider, provider_user_id, email)
		VALUES (?, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'test', 'testprovideruser', 'test@example.com')
	`, userUuid.String())
	test.GetTestContainer().DB.Exec("DELETE FROM posts")
	postUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	s.PostUuid = postUuid
	test.GetTestContainer().DB.Exec("INSERT INTO posts (id, created_at, updated_at, slug, title, content, author_id) VALUES ($1, '2021-01-01 00:00:00', '2021-01-01 00:00:00', 'testslug', 'testtitle', 'testcontent', $2)", postUuid.String(), userUuid.String())
}

func (s *GetPostByIdTestSuite) TestGetPostById() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts/"+s.PostUuid.String(),
		nil,
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: s.PostUuid.String(),
		},
	}
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	GetPostById(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Contains(s.T(), s.W.Body.String(), `"id":"`+s.PostUuid.String()+`"`)
	assert.Contains(s.T(), s.W.Body.String(), `"slug":"testslug"`)
	assert.Contains(s.T(), s.W.Body.String(), `"title":"testtitle"`)
	assert.Contains(s.T(), s.W.Body.String(), `"content":"testcontent"`)
}

func (s *GetPostByIdTestSuite) TestGetPostByIdInvalidUUID() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts/invalid-uuid",
		nil,
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: "invalid-uuid",
		},
	}
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	GetPostById(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusBadRequest, s.W.Code)
	assert.Equal(s.T(), `{"error":"Invalid UUID"}`, s.W.Body.String())
}

func (s *GetPostByIdTestSuite) TestGetPostByIdNotFound() {
	nonExistentUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts/"+nonExistentUuid.String(),
		nil,
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "id",
			Value: nonExistentUuid.String(),
		},
	}
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	GetPostById(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusInternalServerError, s.W.Code)
	assert.Contains(s.T(), s.W.Body.String(), `"error"`)
}

func TestGetPostByIdTestSuite(t *testing.T) {
	suite.Run(t, new(GetPostByIdTestSuite))
}
