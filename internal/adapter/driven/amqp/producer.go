package amqp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type Message struct {
	Recipient string
	Subject   string
	Body      string
}

func NewProducer(url string, queueName string) (*Producer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Producer{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return p.channel.PublishWithContext(ctx,
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
}
