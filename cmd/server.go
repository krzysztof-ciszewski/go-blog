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

	r.GET("/posts/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		q  := query.NewGetPostQuery(uuid.MustParse(id))
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

	r.Run()
}
