package kafka

import (
	"context"
	"errors"
	"strings"

	"goproject/repository"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID, topic, username, password string, tlsEnabled bool) (*Consumer, error) {
	if len(brokers) == 0 || strings.TrimSpace(groupID) == "" || strings.TrimSpace(topic) == "" {
		return nil, errors.New("kafka brokers, group id and topic are required")
	}
	transport, err := newTransport(username, password, tlsEnabled)
	if err != nil {
		return nil, err
	}
	return &Consumer{reader: kafka.NewReader(kafka.ReaderConfig{Brokers: brokers, GroupID: groupID, Topic: topic, Dialer: newDialer(transport), MinBytes: 1, MaxBytes: 10 << 20})}, nil
}

func (consumer *Consumer) Consume(ctx context.Context, handler repository.EventHandler) error {
	if handler == nil {
		return errors.New("kafka event handler is required")
	}
	for {
		message, err := consumer.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil
			}
			return err
		}
		if err := handler(ctx, repository.Event{Key: message.Key, Value: message.Value}); err != nil {
			return err
		}
		if err := consumer.reader.CommitMessages(ctx, message); err != nil {
			return err
		}
	}
}

func (consumer *Consumer) Close() error {
	return consumer.reader.Close()
}
