package amqp

import (
	"fmt"
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
	dlxRoutingKeySuffix := os.Getenv("AMQP_DLX_ROUTING_KEY_SUFFIX")
	dlxExchangeName := os.Getenv("AMQP_DLX_EXCHANGE")
	dlxQueueSuffix := os.Getenv("AMQP_DLX_QUEUE_SUFFIX")
	if err := channel.ExchangeDeclare(
		dlxExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot declare dlx exchange %s", dlxExchangeName))
	}

	dlqQueueName := params.QueueName + "." + dlxQueueSuffix
	dlqRoutingKey := params.QueueName + "." + dlxRoutingKeySuffix
	if _, err := channel.QueueDeclare(
		dlqQueueName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot declare dlq queue %s", dlqQueueName))
	}

	if err := channel.QueueBind(
		dlqQueueName,
		dlqRoutingKey,
		dlxExchangeName,
		false,
		nil,
	); err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("cannot bind dlq queue %s to exchange %s with routing key %s",
				dlqQueueName,
				dlxExchangeName,
				dlqRoutingKey,
			),
		)
	}

	if config.Queue.Arguments == nil {
		config.Queue.Arguments = make(amqp091.Table)
	}

	queueArguments := make(amqp091.Table)
	queueArguments["x-dead-letter-exchange"] = dlxExchangeName
	queueArguments["x-dead-letter-routing-key"] = dlqRoutingKey

	maps.Copy(config.Queue.Arguments, queueArguments)

	if _, err := channel.QueueDeclare(
		params.QueueName,
		config.Queue.Durable,
		config.Queue.AutoDelete,
		config.Queue.Exclusive,
		config.Queue.NoWait,
		config.Queue.Arguments,
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot declare queue %s", params.QueueName))
	}

	logger.Debug("Queue declared", nil)

	if params.ExchangeName == "" {
		logger.Debug("No exchange to declare", nil)
		return nil
	}
	if err := builder.ExchangeDeclare(channel, params.ExchangeName, config); err != nil {
		return errors.Wrap(err, fmt.Sprintf("cannot declare exchange %s", params.ExchangeName))
	}

	logger.Debug("Exchange declared", nil)

	if err := channel.QueueBind(
		params.QueueName,
		params.RoutingKey,
		params.ExchangeName,
		config.QueueBind.NoWait,
		config.QueueBind.Arguments,
	); err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("cannot bind queue %s to exchange %s with routing key %s",
				params.QueueName,
				params.ExchangeName,
				params.RoutingKey,
			),
		)
	}
	return nil
}
