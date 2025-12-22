package post

import (
	"context"
	query "main/internal/Application/Query"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPostById(ctx *gin.Context) {
	container := dependency_injection.GetContainer()

	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	q := query.NewGetPostQuery(parsedUUID)
	post, err := container.QueryBus.Execute(context.Background(), q)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, post)
}
