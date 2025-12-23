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

type GetPostBySlugTestSuite struct {
	suite.Suite
	QueryBus query_bus.QueryBus
	Ctx      *gin.Context
	W        *httptest.ResponseRecorder
	PubSubDb *sql.DB
	PostUuid uuid.UUID
	PostSlug string
}

func (s *GetPostBySlugTestSuite) SetupTest() {
	s.QueryBus = test.GetTestContainer().QueryBus
	s.W = httptest.NewRecorder()
	s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
	gin.SetMode(gin.TestMode)
	s.PubSubDb = test.GetPubSubDb()
	postUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	s.PostUuid = postUuid
	s.PostSlug = "testslug"
	test.GetTestContainer().DB.Exec("DELETE FROM posts")
	test.GetTestContainer().DB.Exec("INSERT INTO posts (id, created_at, updated_at, slug, title, content, author) VALUES ($1, '2021-01-01 00:00:00', '2021-01-01 00:00:00', $2, 'testtitle', 'testcontent', 'testauthor')", postUuid.String(), s.PostSlug)
}

func (s *GetPostBySlugTestSuite) TestGetPostBySlug() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts/slug/"+s.PostSlug,
		nil,
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "slug",
			Value: s.PostSlug,
		},
	}
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	GetPostBySlug(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Contains(s.T(), s.W.Body.String(), `"id":"`+s.PostUuid.String()+`"`)
	assert.Contains(s.T(), s.W.Body.String(), `"slug":"`+s.PostSlug+`"`)
	assert.Contains(s.T(), s.W.Body.String(), `"title":"testtitle"`)
	assert.Contains(s.T(), s.W.Body.String(), `"content":"testcontent"`)
	assert.Contains(s.T(), s.W.Body.String(), `"author":"testauthor"`)
}

func (s *GetPostBySlugTestSuite) TestGetPostBySlugNotFound() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/posts/slug/nonexistentslug",
		nil,
	)
	s.Ctx.Params = gin.Params{
		gin.Param{
			Key:   "slug",
			Value: "nonexistentslug",
		},
	}
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	GetPostBySlug(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusInternalServerError, s.W.Code)
	assert.Contains(s.T(), s.W.Body.String(), `"error"`)
}

func TestGetPostBySlugTestSuite(t *testing.T) {
	suite.Run(t, new(GetPostBySlugTestSuite))
}
