package main

import (
	"context"
	"fmt"
	"log/slog"
	command "main/internal/Application/Command"
	request "main/internal/UserInterface/Api/Request"
	"net/http"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func main() {
	logger := watermill.NewSlogLoggerWithLevelMapping(nil, map[slog.Level]slog.Level{
		slog.LevelInfo: slog.LevelDebug,
	})

	cqrsMarshaller := cqrs.JSONMarshaler{
		GenerateName: cqrs.StructName,
	}

	generateCommandsTopic := func(commandName string) string {
		return "commands." + commandName
	}

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{os.Getenv("KAFKA_BROKER")},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)

	if err != nil {
		panic(err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddMiddleware(middleware.Recoverer)

	commandBus, err := cqrs.NewCommandBusWithConfig(publisher, cqrs.CommandBusConfig{
		GeneratePublishTopic: func(params cqrs.CommandBusGeneratePublishTopicParams) (string, error) {
			return generateCommandsTopic(params.CommandName), nil
		},
		OnSend: func(params cqrs.CommandBusOnSendParams) error {
			logger.Info("Sending command", watermill.LogFields{
				"command_name": params.CommandName,
			})

			params.Message.Metadata.Set("sent_at", time.Now().String())

			return nil
		},
		Marshaler: cqrsMarshaller,
		Logger:    logger,
	})

	if err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/posts/:id", func(ctx *gin.Context) {
		_ = ctx.Param("id")
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
