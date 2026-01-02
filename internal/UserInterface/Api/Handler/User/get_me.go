package user

import (
	user_query "main/internal/Application/Query/User"
	view "main/internal/Application/View"
	query_bus "main/internal/Infrastructure/QueryBus"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

func GetMe(ctx *gin.Context, queryBus query_bus.QueryBus) {
	session, err := gothic.Store.Get(ctx.Request, os.Getenv("SESSION_NAME"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	providerUserId := session.Values["provider_user_id"]
	email := session.Values["email"]

	if providerUserId == nil || email == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	user, err := queryBus.Execute(
		ctx.Request.Context(),
		user_query.NewFindUserByQuery(providerUserId.(string), email.(string)),
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	userView, ok := user.(view.UserView)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	ctx.JSON(http.StatusOK, userView)
}
