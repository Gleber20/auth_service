package driven

import (
	"auth_service/internal/adapter/driven/amqp"
	"context"
)

type AMQPProducer interface {
	Publish(ctx context.Context, msg amqp.Message) error
}
