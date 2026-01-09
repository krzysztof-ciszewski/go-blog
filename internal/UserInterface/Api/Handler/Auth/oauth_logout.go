package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func OauthLogout(ctx *gin.Context) {
	session, err := gothic.Store.Get(ctx.Request, os.Getenv("SESSION_NAME"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	session.Options.MaxAge = -1
	session.Values = make(map[any]any)

	err = session.Save(ctx.Request, ctx.Writer)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	gothic.Store.Save(ctx.Request, ctx.Writer, session)

	ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_URL"))
}
