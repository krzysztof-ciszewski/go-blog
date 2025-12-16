package main

import (
	"context"
	"log/slog"
	command "main/internal/Application/Command"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

func main() {
	logger := watermill.NewSlogLoggerWithLevelMapping(nil, map[slog.Level]slog.Level{
		slog.LevelInfo: slog.LevelDebug,
	})

	cqrsMarshaller := cqrs.JSONMarshaler{
		GenerateName: cqrs.StructName,
	}

	generateEventsTopic := func(eventName string) string {
		return "events." + eventName
	}

	generateCommandsTopic := func(commandName string) string {
		return "commands." + commandName
	}

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:     []string{os.Getenv("KAFKA_BROKER")},
			Unmarshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		panic(err)
	}

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{os.Getenv("KAFKA_BROKER")},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	router.AddMiddleware(middleware.Recoverer)

	commandProcessor, err := cqrs.NewCommandProcessorWithConfig(
		router,
		cqrs.CommandProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.CommandProcessorGenerateSubscribeTopicParams) (string, error) {
				return generateCommandsTopic(params.CommandName), nil
			},
			SubscriberConstructor: func(params cqrs.CommandProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return subscriber, nil
			},
			OnHandle: func(params cqrs.CommandProcessorOnHandleParams) error {
				start := time.Now()

				err := params.Handler.Handle(params.Message.Context(), params.Command)

				logger.Info("Command handled", watermill.LogFields{
					"command_name": params.CommandName,
					"duration":     time.Since(start),
					"err":          err,
				})

				return err
			},
			Marshaler: cqrsMarshaller,
			Logger:    logger,
		},
	)

	if err != nil {
		panic(err)
	}

	eventBus, err := cqrs.NewEventBusWithConfig(publisher, cqrs.EventBusConfig{
		GeneratePublishTopic: func(params cqrs.GenerateEventPublishTopicParams) (string, error) {
			return generateEventsTopic(params.EventName), nil
		},
		OnPublish: func(params cqrs.OnEventSendParams) error {
			logger.Info("Publishing event", watermill.LogFields{
				"event_name": params.EventName,
			})

			params.Message.Metadata.Set("published_at", time.Now().String())

			return nil
		},
		Marshaler: cqrsMarshaller,
		Logger:    logger,
	})
	if err != nil {
		panic(err)
	}

	eventProcessor, err := cqrs.NewEventProcessorWithConfig(
		router,
		cqrs.EventProcessorConfig{
			GenerateSubscribeTopic: func(params cqrs.EventProcessorGenerateSubscribeTopicParams) (string, error) {
				return generateEventsTopic(params.EventName), nil
			},
			SubscriberConstructor: func(params cqrs.EventProcessorSubscriberConstructorParams) (message.Subscriber, error) {
				return subscriber, nil
			},
			OnHandle: func(params cqrs.EventProcessorOnHandleParams) error {
				start := time.Now()

				err := params.Handler.Handle(params.Message.Context(), params.Event)

				logger.Info("Event handled", watermill.LogFields{
					"event_name": params.EventName,
					"duration":   time.Since(start),
					"err":        err,
				})

				return err
			},

			Marshaler: cqrsMarshaller,
			Logger:    logger,
		},
	)

	if err != nil {
		panic(err)
	}

	postRepository := dependency_injection.GetContainer().PostRepository

	err = commandProcessor.AddHandlers(
		cqrs.NewCommandHandler("CreatePostCommandHandler", command.CreatePostCommandHandler{PostRepository: postRepository, EventBus: eventBus}.Handle),
	)

	err = eventProcessor.AddHandlers()

	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}
