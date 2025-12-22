package post

import (
	"context"
	post_command "main/internal/Application/Command/Post"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeletePost(ctx *gin.Context) {
	container := dependency_injection.GetContainer()

	id := ctx.Param("id")

	command := post_command.NewDeletePostCommand(uuid.MustParse(id))
	container.CommandBus.Send(context.Background(), command)

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Post deleted"})
}
