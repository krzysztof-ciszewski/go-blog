package post

import (
	"context"
	post_command "main/internal/Application/Command/Post"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	request "main/internal/UserInterface/Api/Request"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePost(ctx *gin.Context) {
	container := dependency_injection.GetContainer()

	var req request.CreatePostRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	command := post_command.NewCreatePostCommand(
		uuid.MustParse(req.Id),
		req.Slug,
		req.Title,
		req.Content,
		req.Author,
	)

	container.CommandBus.Send(context.Background(), command)

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Post created"})
}
