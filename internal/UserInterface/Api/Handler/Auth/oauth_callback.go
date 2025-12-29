package auth

import (
	command "main/internal/Application/Command/User"
	query "main/internal/Application/Query/User"
	open_telemetry "main/internal/Infrastructure/OpenTelemetry"
	query_bus "main/internal/Infrastructure/QueryBus"
	"net/http"
	"os"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
)

func OauthCallback(
	ctx *gin.Context,
	commandBus *cqrs.CommandBus,
	queryBus query_bus.QueryBus,
	telemetry open_telemetry.Telemetry,
) {
	q := ctx.Request.URL.Query()
	q.Add("provider", ctx.Param("provider"))
	ctx.Request.URL.RawQuery = q.Encode()

	gothUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error authenticating user",
			"error":   err.Error(),
		})
		return
	}

	session, err := gothic.Store.New(ctx.Request, os.Getenv("SESSION_NAME"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error stroing user session",
			"error":   err.Error(),
		})
		return
	}

	session.Values["provider_user_id"] = gothUser.UserID
	session.Values["email"] = gothUser.Email

	if err = session.Save(ctx.Request, ctx.Writer); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error saving user session",
			"error":   err.Error(),
		})
		return
	}

	user, err := queryBus.Execute(
		ctx.Request.Context(),
		query.NewFindUserByQuery(
			gothUser.UserID,
			gothUser.Email,
		),
	)

	if user != nil && err == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_URL"))
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Error generating user ID",
			"error":   err.Error(),
		})
		return
	}

	commandBus.Send(ctx.Request.Context(), command.NewCreateUserCommand(
		id,
		gothUser.Email,
		"",
		gothUser.Provider,
		gothUser.Name,
		gothUser.FirstName,
		gothUser.LastName,
		gothUser.UserID,
		gothUser.AvatarURL,
	))

	ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_URL"))
}
