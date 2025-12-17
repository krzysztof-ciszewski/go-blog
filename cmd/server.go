package main

import (
	"context"
	command "main/internal/Application/Command"
	query "main/internal/Application/Query"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	request "main/internal/UserInterface/Api/Request"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	container := dependency_injection.GetContainer()
	commandBus := container.CommandBus

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/posts", func(ctx *gin.Context) {
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
	})

	r.GET("/posts/:id", func(ctx *gin.Context) {
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
	})

	r.GET("/posts/slug/:slug", func(ctx *gin.Context) {
		slug := ctx.Param("slug")

		q := query.NewFindBySlugQuery(slug)
		post, err := container.QueryBus.Execute(context.Background(), q)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, post)
	})

	r.POST("/posts", func(ctx *gin.Context) {
		var request request.CreatePostRequest

		if err := ctx.ShouldBindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			return
		}

		command := command.NewCreatePostCommand(
			uuid.MustParse(request.Id),
			request.Slug,
			request.Title,
			request.Content,
			request.Author,
		)

		commandBus.Send(context.Background(), command)

		ctx.JSON(http.StatusAccepted, gin.H{"message": "Post created"})
	})

	r.DELETE("/posts/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		command := command.NewDeletePostCommand(uuid.MustParse(id))
		commandBus.Send(context.Background(), command)

		ctx.JSON(http.StatusAccepted, gin.H{"message": "Post deleted"})
	})

	r.Run()
}
