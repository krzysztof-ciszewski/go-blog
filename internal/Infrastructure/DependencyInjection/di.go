package dependency_injection

import (
	"context"
	"log/slog"
	post_command "main/internal/Application/Command/Post"
	user_command "main/internal/Application/Command/User"
	post_query "main/internal/Application/Query/Post"
	user_query "main/internal/Application/Query/User"
	domain_repository "main/internal/Domain/Repository"
	infra_amqp "main/internal/Infrastructure/Amqp"
	config "main/internal/Infrastructure/Config"
	open_telemetry "main/internal/Infrastructure/OpenTelemetry"
	query_bus "main/internal/Infrastructure/QueryBus"
	infra_repository "main/internal/Infrastructure/Repository"
	"os"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Container struct {
	DB               *gorm.DB
	Telemetry        open_telemetry.Telemetry
	QueryBus         query_bus.QueryBus
	CommandBus       *cqrs.CommandBus
	EventBus         *cqrs.EventBus
	Router           *message.Router
	CommandProcessor *cqrs.CommandProcessor
	EventProcessor   *cqrs.EventProcessor
}

var lock = sync.Mutex{}
var container *Container

// TODO: Maybe replace with uber/dig
func GetContainer() *Container {
	if container == nil {
		lock.Lock()
		defer lock.Unlock()
		gormDb, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		telemetry, err := open_telemetry.NewTelemetry(context.Background(), *config.GetTelemetryConfig())
		if err != nil {
			panic(err)
		}

		postRepository := infra_repository.NewPostRepository(gormDb, telemetry)
		userRepository := infra_repository.NewUserRepository(gormDb, telemetry)

		queryBus := buildQueryBus(telemetry)
		registerQueryHandlers(queryBus, postRepository, userRepository, telemetry)

		logger := buildWatermillLogger()
		cqrsMarshaller := buildCqrsMarshaller()
		router := buildRouter(logger)
		amqpConfig := buildAMQPConfig(os.Getenv("AMQP_URI"))
		publisher := buildPublisher(&amqpConfig, logger)
		subscriber := buildSubscriber(&amqpConfig, logger)
		generateCommandsTopic := buildGenerateCommandsTopicFunc()
		generateEventsTopic := buildGenerateEventsTopicFunc()
		commandBus := buildCommandBus(logger, cqrsMarshaller, publisher, generateCommandsTopic)
		eventBus := buildEventBus(publisher, cqrsMarshaller, logger, generateEventsTopic)
		commandProcessor := buildCommandProcessor(router, subscriber, cqrsMarshaller, logger, generateCommandsTopic)
		registerCommandHandlers(commandProcessor, postRepository, userRepository, eventBus)
		eventProcessor := buildEventProcessor(router, subscriber, cqrsMarshaller, logger, generateEventsTopic)
		registerEventHandlers(eventProcessor, eventBus)

		container = &Container{
			DB:               gormDb,
			Telemetry:        *telemetry,
			QueryBus:         queryBus,
			CommandBus:       commandBus,
			EventBus:         eventBus,
			Router:           router,
			CommandProcessor: commandProcessor,
			EventProcessor:   eventProcessor,
		}
	}
	return container
}

func buildGenerateCommandsTopicFunc() func(commandName string) string {
	return func(commandName string) string {
		return "commands." + commandName
	}
}

func buildGenerateEventsTopicFunc() func(eventName string) string {
	return func(eventName string) string {
		return "events." + eventName
	}
}

func buildQueryBus(telemetry open_telemetry.TelemetryProvider) query_bus.QueryBus {
	return query_bus.NewQueryBus(telemetry)
}

func buildRouter(logger watermill.LoggerAdapter) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	return router
}

func buildWatermillLogger() watermill.LoggerAdapter {
	return watermill.NewSlogLoggerWithLevelMapping(nil, map[slog.Level]slog.Level{
		slog.LevelInfo: slog.LevelDebug,
	})
}

func buildCqrsMarshaller() *cqrs.JSONMarshaler {
	return &cqrs.JSONMarshaler{
		GenerateName: cqrs.StructName,
	}
}

func buildAMQPConfig(amqpURL string) amqp.Config {
	config := amqp.NewDurableQueueConfig(amqpURL)
	config.TopologyBuilder = &infra_amqp.MyTopologyBuilder{}
	config.Consume.NoRequeueOnNack = true
	return config
}

func buildPublisher(amqpConfig *amqp.Config, logger watermill.LoggerAdapter) message.Publisher {
	publisher, err := amqp.NewPublisher(*amqpConfig, logger)

	if err != nil {
		panic(err)
	}

	return publisher
}

func buildSubscriber(amqpConfig *amqp.Config, logger watermill.LoggerAdapter) message.Subscriber {
	subscriber, err := amqp.NewSubscriber(*amqpConfig, logger)
	if err != nil {
		panic(err)
	}

	return subscriber
}

func buildCommandBus(
	logger watermill.LoggerAdapter,
	cqrsMarshaller *cqrs.JSONMarshaler,
	publisher message.Publisher,
	generateCommandsTopic func(commandName string) string,
) *cqrs.CommandBus {
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

	return commandBus
}

func buildEventBus(
	publisher message.Publisher,
	cqrsMarshaller *cqrs.JSONMarshaler,
	logger watermill.LoggerAdapter,
	generateEventsTopic func(eventName string) string,
) *cqrs.EventBus {
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

	return eventBus
}

func buildCommandProcessor(
	router *message.Router,
	subscriber message.Subscriber,
	cqrsMarshaller *cqrs.JSONMarshaler,
	logger watermill.LoggerAdapter,
	generateCommandsTopic func(commandName string) string,
) *cqrs.CommandProcessor {
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

	return commandProcessor
}

func buildEventProcessor(
	router *message.Router,
	subscriber message.Subscriber,
	cqrsMarshaller *cqrs.JSONMarshaler,
	logger watermill.LoggerAdapter,
	generateEventsTopic func(eventName string) string,
) *cqrs.EventProcessor {
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

	return eventProcessor
}

func registerQueryHandlers(
	queryBus query_bus.QueryBus,
	postRepository domain_repository.PostRepository,
	userRepository domain_repository.UserRepository,
	telemetry open_telemetry.TelemetryProvider,
) {
	queryBus.RegisterHandler(post_query.GetPostQueryHandler{PostRepository: postRepository})
	queryBus.RegisterHandler(post_query.FindAllByQueryHandler{PostRepository: postRepository})
	queryBus.RegisterHandler(user_query.FindUserByQueryHandler{UserRepository: userRepository, Telemetry: telemetry})
}

func registerCommandHandlers(
	commandProcessor *cqrs.CommandProcessor,
	postRepository domain_repository.PostRepository,
	userRepository domain_repository.UserRepository,
	eventBus *cqrs.EventBus,
) {
	commandProcessor.AddHandlers(
		cqrs.NewCommandHandler("CreatePostCommandHandler", post_command.CreatePostCommandHandler{PostRepository: postRepository, EventBus: eventBus}.Handle),
		cqrs.NewCommandHandler("UpdatePostCommandHandler", post_command.UpdatePostCommandHandler{PostRepository: postRepository, EventBus: eventBus}.Handle),
		cqrs.NewCommandHandler("DeletePostCommandHandler", post_command.DeletePostCommandHandler{PostRepository: postRepository, EventBus: eventBus}.Handle),
		cqrs.NewCommandHandler("CreateUserCommandHandler", user_command.CreateUserCommandHandler{UserRepository: userRepository, EventBus: eventBus}.Handle),
	)
}

func registerEventHandlers(eventProcessor *cqrs.EventProcessor, eventBus *cqrs.EventBus) {
	eventProcessor.AddHandlers()
}
