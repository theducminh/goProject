package repository

import "context"

type Event struct {
	Key   []byte
	Value []byte
}

type EventHandler func(ctx context.Context, event Event) error

type EventPublisher interface {
	Publish(ctx context.Context, topic string, event Event) error
	Close() error
}

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type EventConsumer interface {
	Consume(ctx context.Context, handler EventHandler) error
	Close() error
}
