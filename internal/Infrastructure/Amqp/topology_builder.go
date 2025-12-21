package amqp

import (
	"maps"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v3/pkg/amqp"
	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

type MyTopologyBuilder struct{}

func (builder *MyTopologyBuilder) ExchangeDeclare(channel *amqp091.Channel, exchangeName string, config amqp.Config) error {
	return channel.ExchangeDeclare(
		exchangeName,
		config.Exchange.Type,
		config.Exchange.Durable,
		config.Exchange.AutoDeleted,
		config.Exchange.Internal,
		config.Exchange.NoWait,
		config.Exchange.Arguments,
	)
}

func (builder *MyTopologyBuilder) BuildTopology(channel *amqp091.Channel, params amqp.BuildTopologyParams, config amqp.Config, logger watermill.LoggerAdapter) error {
	if err := channel.ExchangeDeclare(
		os.Getenv("AMQP_DLX_EXCHANGE"),
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return errors.Wrap(err, "cannot declare dlx exchange")
	}

	if _, err := channel.QueueDeclare(
		params.QueueName+"."+os.Getenv("AMQP_DLX_QUEUE_SUFFIX"),
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return errors.Wrap(err, "cannot declare dlq queue")
	}

	if err := channel.QueueBind(
		params.QueueName+"."+os.Getenv("AMQP_DLX_QUEUE_SUFFIX"),
		params.QueueName+"."+os.Getenv("AMQP_DLX_ROUTING_KEY_SUFFIX"),
		os.Getenv("AMQP_DLX_EXCHANGE"),
		false,
		nil,
	); err != nil {
		return errors.Wrap(err, "cannot bind dlq queue")
	}

	if config.Queue.Arguments == nil {
		config.Queue.Arguments = make(amqp091.Table)
	}

	queueArguments := make(amqp091.Table)
	queueArguments["x-dead-letter-exchange"] = os.Getenv("AMQP_DLX_EXCHANGE")
	queueArguments["x-dead-letter-routing-key"] = params.QueueName + "." + os.Getenv("AMQP_DLX_ROUTING_KEY_SUFFIX")

	maps.Copy(config.Queue.Arguments, queueArguments)

	if _, err := channel.QueueDeclare(
		params.QueueName,
		config.Queue.Durable,
		config.Queue.AutoDelete,
		config.Queue.Exclusive,
		config.Queue.NoWait,
		config.Queue.Arguments,
	); err != nil {
		return errors.Wrap(err, "cannot declare queue")
	}

	logger.Debug("Queue declared", nil)

	if params.ExchangeName == "" {
		logger.Debug("No exchange to declare", nil)
		return nil
	}
	if err := builder.ExchangeDeclare(channel, params.ExchangeName, config); err != nil {
		return errors.Wrap(err, "cannot declare exchange")
	}

	logger.Debug("Exchange declared", nil)


	if err := channel.QueueBind(
		params.QueueName,
		params.RoutingKey,
		params.ExchangeName,
		config.QueueBind.NoWait,
		config.QueueBind.Arguments,
	); err != nil {
		return errors.Wrap(err, "cannot bind queue")
	}
	return nil
}
