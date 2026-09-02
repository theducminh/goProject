package kafka

import (
	"context"
	"errors"
	"strings"

	"goproject/repository"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, username, password string, tlsEnabled bool) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, errors.New("kafka brokers are required")
	}
	transport, err := newTransport(username, password, tlsEnabled)
	if err != nil {
		return nil, err
	}
	return &Producer{writer: &kafka.Writer{Addr: kafka.TCP(brokers...), Balancer: &kafka.Hash{}, Transport: transport, RequiredAcks: kafka.RequireAll}}, nil
}

func (producer *Producer) Publish(ctx context.Context, topic string, event repository.Event) error {
	if strings.TrimSpace(topic) == "" {
		return errors.New("kafka topic is required")
	}
	return producer.writer.WriteMessages(ctx, kafka.Message{Topic: topic, Key: event.Key, Value: event.Value})
}

func (producer *Producer) Close() error {
	return producer.writer.Close()
}
