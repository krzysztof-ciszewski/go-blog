package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := gothic.Store.Get(ctx.Request, os.Getenv("SESSION_NAME"))
		if err != nil || session.Values["provider_user_id"] == nil || session.Values["email"] == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized User",
			})
			return
		}

		ctx.Next()
	}
}
