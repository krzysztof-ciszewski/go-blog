package post

import (
	"database/sql"
	test "main/internal/Infrastructure/DependencyInjection/Test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type DeletePostTestSuite struct {
	suite.Suite
	CommandBus *cqrs.CommandBus
	Ctx        *gin.Context
	W          *httptest.ResponseRecorder
	PubSubDb   *sql.DB
	PostUuid   uuid.UUID
}

func (s *DeletePostTestSuite) SetupTest() {
	s.CommandBus = test.GetTestContainer().CommandBus
	s.W = httptest.NewRecorder()
	s.Ctx = gin.CreateTestContextOnly(s.W, gin.Default())
	gin.SetMode(gin.TestMode)
	s.PubSubDb = test.GetPubSubDb()
	s.PubSubDb.Exec("DELETE FROM `watermill_commands.deletePostCommand`")
	test.GetTestContainer().DB.Exec("DELETE FROM users")
	userUuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
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

func (s *DeletePostTestSuite) TestDeletePost() {
	s.Ctx.Request = httptest.NewRequest(
		"DELETE",
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

	DeletePost(s.Ctx, s.CommandBus)

	assert.Equal(s.T(), http.StatusAccepted, s.W.Code)
	assert.Equal(s.T(), `{"message":"Post deleted"}`, s.W.Body.String())
	count := test.GetCommandCount("deletePostCommand")
	assert.Equal(s.T(), 1, count)
}

func TestDeletePostTestSuite(t *testing.T) {
	suite.Run(t, new(DeletePostTestSuite))
}
