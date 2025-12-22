package post

import (
	"context"
	query "main/internal/Application/Query"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListPosts(ctx *gin.Context) {
	container := dependency_injection.GetContainer()

	text := ctx.Query("text")
	author := ctx.Query("author")

	var result any
	var err error

	if text != "" {
		q := query.NewFindAllByTextQuery(text)
		result, err = container.QueryBus.Execute(context.Background(), q)
	} else if author != "" {
		q := query.NewFindAllByAuthorQuery(author)
		result, err = container.QueryBus.Execute(context.Background(), q)
	} else {
		q := query.NewFindAllQuery()
		result, err = container.QueryBus.Execute(context.Background(), q)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
