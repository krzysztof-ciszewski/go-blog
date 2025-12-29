package post

import (
	post_command "main/internal/Application/Command/Post"
	user_query "main/internal/Application/Query/User"
	view "main/internal/Application/View"
	query_bus "main/internal/Infrastructure/QueryBus"
	request "main/internal/UserInterface/Api/Request"
	"net/http"
	"os"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/markbates/goth/gothic"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePost(ctx *gin.Context, commandBus *cqrs.CommandBus, queryBus query_bus.QueryBus) {
	var req request.CreatePostRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	command := post_command.NewCreatePostCommand(
		uuid.MustParse(req.Id),
		req.Slug,
		req.Title,
		req.Content,
		user.(view.UserView).Id,
	)

	commandBus.Send(ctx.Request.Context(), command)

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Post created"})
}
