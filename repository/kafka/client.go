package kafka

import (
	"crypto/tls"
	"errors"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

func newTransport(username, password string, tlsEnabled bool) (*kafka.Transport, error) {
	transport := &kafka.Transport{TLS: &tls.Config{MinVersion: tls.VersionTLS12}}
	if !tlsEnabled {
		transport.TLS = nil
	}
	if strings.TrimSpace(username) == "" && password == "" {
		return transport, nil
	}
	if strings.TrimSpace(username) == "" || password == "" {
		return nil, errors.New("kafka username and password must be provided together")
	}
	mechanism, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		return nil, err
	}
	transport.SASL = mechanism
	return transport, nil
}

func newDialer(transport *kafka.Transport) *kafka.Dialer {
	return &kafka.Dialer{Timeout: 10 * time.Second, DualStack: true, TLS: transport.TLS, SASLMechanism: transport.SASL}
}
