// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package publisher

import (
	"context"
	"errors"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/nats"
	"go.uber.org/zap"
)

type Message struct {
	Subject string
	Body    []byte
}

type Publisher interface {
	Publish(ctx context.Context, message *Message) error
	Close() error
}

type NatsPublisher struct {
	logger     *zap.Logger
	client     *nats.NatsClient
	natsConfig nats.NatsConfig
}

func NewNatsPublisher(logger *zap.Logger, client *nats.NatsClient, config nats.NatsConfig) (Publisher, error) {
	p := &NatsPublisher{
		logger:     logger,
		client:     client,
		natsConfig: config,
	}
	conn, err := p.client.Connect(config)
	if err != nil {
		return nil, err
	}
	p.client.Connection = conn
	return p, nil
}

func (p *NatsPublisher) Publish(ctx context.Context, message *Message) error {
	if p.client.IsConnected() {
		err := p.client.Connection.Publish(message.Subject, message.Body)
		if err != nil {
			p.logger.Error("Failed to publish message", zap.Error(err))
			return err
		}
	} else {
		p.logger.Error("Failed to publish message, client is not connected", zap.Error(errors.New("fail to connect")))
	}
	return nil
}

func (p *NatsPublisher) IsConnected() bool {
	return p.client.IsConnected()
}

func (p *NatsPublisher) IsClosed() bool {
	return p.client.IsClosed()
}

func (p *NatsPublisher) Close() error {
	return p.client.Close()
}
