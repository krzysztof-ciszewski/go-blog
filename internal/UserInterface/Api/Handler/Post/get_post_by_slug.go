package post

import (
	"context"
	query "main/internal/Application/Query"
	query_bus "main/internal/Infrastructure/QueryBus"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPostBySlug(ctx *gin.Context, queryBus query_bus.QueryBus) {

	slug := ctx.Param("slug")

	q := query.NewFindBySlugQuery(slug)
	post, err := queryBus.Execute(context.Background(), q)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, post)
}
