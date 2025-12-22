package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func OauthLogout(ctx *gin.Context) {
	gothic.Logout(ctx.Writer, ctx.Request)
	ctx.Writer.Header().Set("Location", "/")
	ctx.Writer.WriteHeader(http.StatusTemporaryRedirect)
}
