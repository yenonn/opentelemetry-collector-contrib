// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package natsexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/natsexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/natsexporter/internal/publisher"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/nats"
)

type natsExporter struct {
	config         *Config
	settings       component.TelemetrySettings
	subject        string
	connectionName string
	*marshaler
	publisherFactory
	publisher publisher.Publisher
}

type publisherFactory = func(nats.NatsConfig) (publisher.Publisher, error)

func newNatsExporter(cfg *Config, set component.TelemetrySettings, publisherFactory publisherFactory, subject string, connectionName string) *natsExporter {
	exporter := &natsExporter{
		config:         cfg,
		settings:       set,
		subject:        subject,
		connectionName: connectionName,
	}
	return exporter
}

func (e *natsExporter) start(ctx context.Context, host component.Host) error {
	m, err := newMarshaler(e.config.EncodingExtensionID, host)
	if err != nil {
		return err
	}
	e.marshaler = m
	natsConfig := nats.NatsConfig{
		URL:               e.config.Connection.Endpoint,
		Username:          e.config.Connection.Auth.Plain.Username,
		Password:          e.config.Connection.Auth.Plain.Password,
		ConnectionTimeout: e.config.Connection.ConnectionTimeout,
		RootCA:            e.config.Connection.TLSConfig.CAFile,
	}
	e.settings.Logger.Info("Establishing initial connection to NATS")
	p, err := e.publisherFactory(natsConfig)
	e.publisher = p

	if err != nil {
		return err
	}
	return nil
}

func (e *natsExporter) publishTraces(context context.Context, traces ptrace.Traces) error {
	body, err := e.tracesMarshaler.MarshalTraces(traces)
	if err != nil {
		return err
	}

	message := &publisher.Message{
		Body:    body,
		Subject: e.config.Topic.Subject,
	}
	return e.publisher.Publish(context, message)
}

func (e *natsExporter) publishMetrics(context context.Context, metrics pmetric.Metrics) error {
	body, err := e.metricsMarshaler.MarshalMetrics(metrics)
	if err != nil {
		return err
	}

	message := &publisher.Message{
		Body:    body,
		Subject: e.config.Topic.Subject,
	}
	return e.publisher.Publish(context, message)
}

func (e *natsExporter) publishLogs(context context.Context, logs plog.Logs) error {
	body, err := e.logsMarshaler.MarshalLogs(logs)
	if err != nil {
		return err
	}

	message := &publisher.Message{
		Body:    body,
		Subject: e.config.Topic.Subject,
	}
	return e.publisher.Publish(context, message)
}

func (e *natsExporter) shutdown(_ context.Context) error {
	if e.publisher != nil {
		return e.publisher.Close()
	}
	return nil
}
