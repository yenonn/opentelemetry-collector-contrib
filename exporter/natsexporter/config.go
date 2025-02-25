// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package natsexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/rabbitmqexporter"

import (
	"errors"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/exporter/exporterbatcher"
)

type Config struct {
	Connection          ConnectionConfig          `mapstructure:"connection"`
	Topic               TopicConfig               `mapstructure:"topic"`
	EncodingExtensionID *component.ID             `mapstructure:"encoding_extension"`
	RetrySettings       configretry.BackOffConfig `mapstructure:"retry_on_failure"`
	BatcherSettings     exporterbatcher.Config    `mapstructure:"batcher"`
}

type ConnectionConfig struct {
	Endpoint          string                  `mapstructure:"endpoint"`
	TLSConfig         *configtls.ClientConfig `mapstructure:"tls"`
	Auth              AuthConfig              `mapstructure:"auth"`
	ConnectionTimeout time.Duration           `mapstructure:"connection_timeout"`
	Name              string                  `mapstructure:"name"`
}

type TopicConfig struct {
	Subjects []string `mapstructure:"subjects"`
}

type AuthConfig struct {
	Plain PlainAuth `mapstructure:"plain"`
}

type PlainAuth struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if cfg.Connection.Endpoint == "" {
		return errors.New("connection.endpoint is required")
	}

	// Password-less users are possible so only validate username
	if cfg.Connection.Auth.Plain.Username == "" {
		return errors.New("connection.auth.plain.username is required")
	}

	return nil
}
