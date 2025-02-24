// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package natsexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/natsexporter"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/natsexporter/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/natsexporter/internal/publisher"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/nats"
)

const (
	defaultConnectionTimeout          = time.Second * 10
	defaultConnectionHeartbeat        = time.Second * 5
	defaultPublishConfirmationTimeout = time.Second * 5

	deadletterMetricsSubject = "deadletter.otlp_metrics"
	deadletterSpansSubject   = "deadletter.otlp_spans"
	deadletterLogsSubject    = "deadletter.otlp_logs"

	defaultSpansConnectionName   = "otel-collector-spans"
	defaultMetricsConnectionName = "otel-collector-metrics"
	defaultLogsConnectionName    = "otel-collector-logs"
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.LogsStability),
		exporter.WithMetrics(createMetricsExporter, metadata.TracesStability),
		exporter.WithTraces(createTracesExporter, metadata.LogsStability),
	)
}

func createDefaultConfig() component.Config {
	retrySettings := configretry.BackOffConfig{
		Enabled: false,
	}
	return &Config{
		RetrySettings: retrySettings,
		Connection: ConnectionConfig{
			ConnectionTimeout: defaultConnectionTimeout,
		},
	}
}

func createTracesExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Traces, error) {
	config := cfg.(*Config)

	spansSubject := getTopicSubjectOrDefault(config, deadletterSpansSubject)
	connectionName := defaultSpansConnectionName
	if config.Connection.Name != "" {
		connectionName = config.Connection.Name
	}
	r := newNatsExporter(config, set.TelemetrySettings, newPublisherFactory(set), spansSubject, connectionName)

	return exporterhelper.NewTraces(
		ctx,
		set,
		cfg,
		r.publishTraces,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(r.start),
		exporterhelper.WithShutdown(r.shutdown),
		exporterhelper.WithRetry(config.RetrySettings),
	)
}

func createMetricsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Metrics, error) {
	config := (cfg.(*Config))

	metricsSubject := getTopicSubjectOrDefault(config, deadletterMetricsSubject)

	connectionName := defaultMetricsConnectionName
	if config.Connection.Name != "" {
		connectionName = config.Connection.Name
	}
	r := newNatsExporter(config, set.TelemetrySettings, newPublisherFactory(set), metricsSubject, connectionName)

	return exporterhelper.NewMetrics(
		ctx,
		set,
		cfg,
		r.publishMetrics,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(r.start),
		exporterhelper.WithShutdown(r.shutdown),
		exporterhelper.WithRetry(config.RetrySettings),
	)
}

func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	config := (cfg.(*Config))

	logsSubject := getTopicSubjectOrDefault(config, deadletterLogsSubject)
	connectionName := defaultLogsConnectionName
	if config.Connection.Name != "" {
		connectionName = config.Connection.Name
	}
	r := newNatsExporter(config, set.TelemetrySettings, newPublisherFactory(set), logsSubject, connectionName)

	return exporterhelper.NewLogs(
		ctx,
		set,
		cfg,
		r.publishLogs,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(r.start),
		exporterhelper.WithShutdown(r.shutdown),
		exporterhelper.WithRetry(config.RetrySettings),
	)
}

func getTopicSubjectOrDefault(config *Config, fallback string) string {
	subject := fallback
	if config.Topic.Subject != "" {
		subject = config.Topic.Subject
	}
	return subject
}

func newPublisherFactory(set exporter.Settings) publisherFactory {
	return func(natsConfig nats.NatsConfig) (publisher.Publisher, error) {
		return publisher.NewNatsPublisher(set.Logger, nats.NewNatsClient(set.Logger), natsConfig)
	}
}
