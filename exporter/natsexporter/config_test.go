// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package natsexporter

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/confmap/xconfmap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/natsexporter/internal/metadata"
)

var encodingComponentID = component.NewIDWithName(component.MustNewType("otlp_encoding"), "nats")

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "test-config.yaml"))
	require.NoError(t, err)

	tests := []struct {
		id           component.ID
		expected     component.Config
		errorMessage string
	}{
		{
			id:           component.NewIDWithName(metadata.Type, "missing_endpoint"),
			errorMessage: "connection.endpoint is required",
		},
		{
			id:           component.NewIDWithName(metadata.Type, "missing_plainauth_username"),
			errorMessage: "connection.auth.plain.username is required",
		},
		{
			id: component.NewIDWithName(metadata.Type, "all_fields"),
			expected: &Config{
				Connection: ConnectionConfig{
					Endpoint: "nats://localhost:4222",
					Auth: AuthConfig{
						Plain: PlainAuth{
							Username: "user",
							Password: "pass",
						},
					},
					TLSConfig: &configtls.ClientConfig{
						Config: configtls.Config{
							CAFile: "/var/lib/myca.pem",
						},
						Insecure: true,
					},
					ConnectionTimeout: time.Millisecond,
				},
				Topic: TopicConfig{
					Subjects: ["nats.test.subject"],
				},
				EncodingExtensionID: &encodingComponentID,
				RetrySettings: configretry.BackOffConfig{
					Enabled: true,
				},
			},
		},
		{
			id: component.NewIDWithName(metadata.Type, "mandatory_fields"),
			expected: &Config{
				Connection: ConnectionConfig{
					Endpoint: "nats://localhost:4222",
					Auth: AuthConfig{
						Plain: PlainAuth{
							Username: "user",
							Password: "pass",
						},
					},
					ConnectionTimeout: defaultConnectionTimeout,
				},
				RetrySettings: configretry.BackOffConfig{
					Enabled: false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, sub.Unmarshal(cfg))

			if tt.expected == nil {
				err = errors.Join(err, xconfmap.Validate(cfg))
				assert.ErrorContains(t, err, tt.errorMessage)
				return
			}

			assert.NoError(t, xconfmap.Validate(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}
