package user

import (
	"database/sql"
	test "main/internal/Infrastructure/DependencyInjection/Test"
	query_bus "main/internal/Infrastructure/QueryBus"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type GetMeTestSuite struct {
	suite.Suite
	QueryBus query_bus.QueryBus
	Ctx      *gin.Context
	W        *httptest.ResponseRecorder
	PubSubDb *sql.DB
	UserUuid uuid.UUID
}

func (s *GetMeTestSuite) SetupTest() {
	if os.Getenv("SESSION_NAME") == "" {
		_ = os.Setenv("SESSION_NAME", "blog_session")
	}

	s.QueryBus = test.GetTestContainer().QueryBus
	s.W = httptest.NewRecorder()
	s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
	gin.SetMode(gin.TestMode)

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

	s.PubSubDb = test.GetPubSubDb()
}

func (s *GetMeTestSuite) TestGetMe() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/me",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	session, err := gothic.Store.New(s.Ctx.Request, os.Getenv("SESSION_NAME"))
	if err != nil {
		panic(err)
	}
	session.Values["provider_user_id"] = "testprovideruser"
	session.Values["email"] = "test@example.com"
	if err := session.Save(s.Ctx.Request, s.Ctx.Writer); err != nil {
		panic(err)
	}
	s.Ctx.Request.Header.Set("Cookie", s.Ctx.Writer.Header().Get("Set-Cookie"))

	GetMe(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusOK, s.W.Code)
	assert.Contains(s.T(), s.W.Body.String(), `"id":"`+s.UserUuid.String()+`"`)
	assert.Contains(s.T(), s.W.Body.String(), `"email":"test@example.com"`)
	assert.Contains(s.T(), s.W.Body.String(), `"provider":"test"`)
	assert.Contains(s.T(), s.W.Body.String(), `"provider_user_id":"testprovideruser"`)
}

func (s *GetMeTestSuite) TestGetMeMissingSessionValues() {
	s.Ctx.Request = httptest.NewRequest(
		"GET",
		"/api/v1/me",
		nil,
	)
	s.Ctx.Request.Header.Set("Content-Type", "application/json")

	session, err := gothic.Store.New(s.Ctx.Request, os.Getenv("SESSION_NAME"))
	if err != nil {
		panic(err)
	}
	if err := session.Save(s.Ctx.Request, s.Ctx.Writer); err != nil {
		panic(err)
	}
	s.Ctx.Request.Header.Set("Cookie", s.Ctx.Writer.Header().Get("Set-Cookie"))

	GetMe(s.Ctx, s.QueryBus)

	assert.Equal(s.T(), http.StatusInternalServerError, s.W.Code)
	assert.Equal(s.T(), `{"error":"User not found"}`, s.W.Body.String())
}

func TestGetMeTestSuite(t *testing.T) {
	suite.Run(t, new(GetMeTestSuite))
}
