package auth

import (
	"context"
	query "main/internal/Application/Query"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func OauthInitial(ctx *gin.Context) {
	container := dependency_injection.GetContainer()

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

	user, err := container.QueryBus.Execute(
		context.Background(),
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
