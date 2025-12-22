package post

import (
	"context"
	query "main/internal/Application/Query"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPostBySlug(ctx *gin.Context) {
	container := dependency_injection.GetContainer()

	slug := ctx.Param("slug")

	q := query.NewFindBySlugQuery(slug)
	post, err := container.QueryBus.Execute(context.Background(), q)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, post)
}
