package post

import (
	"context"
	post_query "main/internal/Application/Query/Post"
	query_bus "main/internal/Infrastructure/QueryBus"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListPosts(ctx *gin.Context, queryBus query_bus.QueryBus) {
	page := ctx.Query("page")
	pageSize := ctx.Query("pageSize")
	slug := ctx.Query("slug")
	text := ctx.Query("text")
	author := ctx.Query("author")

	var result any
	var err error

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page"})
		return
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pageSize"})
		return
	}

	q := post_query.NewFindAllByQuery(pageInt, pageSizeInt, slug, text, author)
	result, err = queryBus.Execute(context.Background(), q)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
