// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package publisher

import (
	"context"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/nats"
	"go.uber.org/zap"
)

type Message struct {
	Subject string
	Body    []byte
}

type Publisher interface {
	Publish(ctx context.Context, message *Message) error
	PublishMessages(ctx context.Context, messages []*Message) error
	Close() error
}

type NatsPublisher struct {
	logger *zap.Logger
	client nats.NatsClient
}

func NewNatsPublisher() *NatsPublisher {
	return &NatsPublisher{
		logger: zap.NewNop(),
		client: *nats.NewNatsClient(),
	}
}

func (p *NatsPublisher) Publish(ctx context.Context, message *Message) error {
	config := nats.DefaultDialConfig()

	err := p.client.Connect(config)
	if err != nil {
		p.logger.Error("Failed to connect to nats", zap.Error(err))
	}
	if p.client.IsConnected() {
		err := p.client.Connection.Publish(message.Subject, message.Body)
		if err != nil {
			p.logger.Error("Failed to publish message", zap.Error(err))
			return err
		}
	}
	return nil
}

func (p *NatsPublisher) PublishMessages(ctx context.Context, messages []*Message) error {
	config := nats.DefaultDialConfig()

	err := p.client.Connect(config)
	if err != nil {
		p.logger.Error("Failed to connect to nats", zap.Error(err))
	}
	if p.client.IsConnected() {
		for _, message := range messages {
			err := p.client.Connection.Publish(message.Subject, message.Body)
			if err != nil {
				p.logger.Error("Failed to publish message", zap.Error(err))
				return nil
			}
		}
	}
	return nil
}

func (p *NatsPublisher) IsConnected() bool {
	return p.client.IsConnected()
}

func (p *NatsPublisher) IsClosed() bool {
	return p.client.IsClosed()
}

func (p *NatsPublisher) Close() {
	p.client.Close()
}
