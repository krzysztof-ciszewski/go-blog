package auth

import (
	query "main/internal/Application/Query/User"
	open_telemetry "main/internal/Infrastructure/OpenTelemetry"
	query_bus "main/internal/Infrastructure/QueryBus"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func OauthInitial(ctx *gin.Context, queryBus query_bus.QueryBus, telemetry open_telemetry.Telemetry) {
	spanCtx, span := telemetry.TraceStart(ctx.Request.Context(), "OauthInitial")
	defer span.End()

	q := ctx.Request.URL.Query()
	q.Add("provider", ctx.Param("provider"))
	ctx.Request.URL.RawQuery = q.Encode()

	session, err := gothic.Store.Get(ctx.Request, os.Getenv("SESSION_NAME"))
	if err != nil {
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		return
	}

	providerUserId := session.Values["provider_user_id"]
	email := session.Values["email"]

	if providerUserId == nil || email == nil {
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		return
	}

	user, err := queryBus.Execute(
		spanCtx,
		query.NewFindUserByQuery(
			providerUserId.(string),
			email.(string),
		),
	)

	if user != nil && err == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_URL"))
		return
	}

	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}
