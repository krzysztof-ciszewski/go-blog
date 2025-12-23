package post

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	post_command "main/internal/Application/Command/Post"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeletePost(ctx *gin.Context, commandBus *cqrs.CommandBus) {
	id := ctx.Param("id")

	command := post_command.NewDeletePostCommand(uuid.MustParse(id))
	commandBus.Send(context.Background(), command)

	ctx.JSON(http.StatusAccepted, gin.H{"message": "Post deleted"})
}
