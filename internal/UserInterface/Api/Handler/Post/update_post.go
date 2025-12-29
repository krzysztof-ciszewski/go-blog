package post

import (
	post_command "main/internal/Application/Command/Post"
	post_query "main/internal/Application/Query/Post"
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

func UpdatePost(ctx *gin.Context, commandBus *cqrs.CommandBus, queryBus query_bus.QueryBus) {
	id := ctx.Param("id")
	postId, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var req request.UpdatePostRequest
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	post, err := queryBus.Execute(
		ctx.Request.Context(),
		post_query.NewGetPostQuery(postId),
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	postView, ok := post.(view.PostView)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid post data"})
		return
	}

	userView, ok := user.(view.UserView)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	if postView.AuthorId != userView.Id {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this post"})
		return
	}

	command := post_command.NewUpdatePostCommand(
		postId,
		req.Slug,
		req.Title,
		req.Content,
	)

	commandBus.Send(ctx.Request.Context(), command)

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Post updated"})
}
