package rabbitmq

import (
	"context"
	"encoding/json"
	"user-management-api/internal/utils"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rabbitMQService struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	logger  *zerolog.Logger
}

func NewRabbitMQService(amqpUrl string, logger *zerolog.Logger) (RabbitMQService, error) {
	conn, err := amqp091.Dial(amqpUrl)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to RabbitMQ")
		return nil, utils.NewWrapError("Failed connect to RabbitMQ", utils.ErrorCodeInternal, err)
	}

	ch, err := conn.Channel()

	if err != nil {
		logger.Error().Err(err).Msg("Failed to  open a channel")
		conn.Close()
		return nil, utils.NewWrapError("Failed to open a channel", utils.ErrorCodeInternal, err)
	}

	return &rabbitMQService{conn: conn, channel: ch, logger: logger}, nil
}

// send message
func (r *rabbitMQService) Publish(ctx context.Context, queue string, message any) error {
	q, err := r.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to  declare a queue for publisher")
		return err
	}

	body, err := json.Marshal(message)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to parse message to json")
		return err
	}
	err = r.channel.PublishWithContext(ctx,
		"",     // exchange default
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to publish a message")
		return err
	}
	return nil
}

func (r *rabbitMQService) Consume(ctx context.Context, queue string, handler func([]byte) error) error {
	q, err := r.channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to  declare a queue for consumer")
		return err
	}
	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to consume a message")
		return err
	}

	// listen for messages
	go func() {
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				if err := handler(msg.Body); err != nil {
					// handle failure => deny
					msg.Nack(false, false)
				} else {
					// handle success => ack
					msg.Ack(false)
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *rabbitMQService) Close() error {
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			r.logger.Error().Err(err).Msg("Failed to close a channel")
			return err
		}
	}
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			r.logger.Error().Err(err).Msg("Failed to close a connection")
			return err
		}
	}
	return nil
}
