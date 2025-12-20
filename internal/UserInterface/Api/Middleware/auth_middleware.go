package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := gothic.Store.Get(ctx.Request, os.Getenv("SESSION_NAME"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized User",
				"error":   err.Error(),
			})
			return
		}

		ctx.Next()
	}
}
