package post

import (
	"context"
	query "main/internal/Application/Query"
	query_bus "main/internal/Infrastructure/QueryBus"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPostById(ctx *gin.Context, queryBus query_bus.QueryBus) {

	id := ctx.Param("id")

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	q := query.NewGetPostQuery(parsedUUID)
	post, err := queryBus.Execute(context.Background(), q)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, post)
}
