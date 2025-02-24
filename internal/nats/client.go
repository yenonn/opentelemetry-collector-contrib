package nats

import (
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NatsClient interface {
	DialConfig(config DialConfig) error
}

type DialConfig struct {
	URL               string
	ConnectionTimeout time.Duration
	Username          string
	Password          string
	RootCA            string
}

func DefaultDialConfig() *DialConfig {
	return &DialConfig{
		URL:               "nats://localhost:4222",
		ConnectionTimeout: 5 * time.Second,
		Username:          "",
		Password:          "",
		RootCA:            "",
	}
}

type client struct {
	logger     *zap.Logger
	Connection *nats.Conn
}

func (c *client) Connect(config DialConfig) error {
	c.logger.Debug("Connecting to nats")

	nc, err := nats.Connect(config.URL, nats.UserInfo(config.Username, config.Password), nats.RootCAs(config.RootCA), nats.Timeout(config.ConnectionTimeout))
	if err != nil {
		nc.Close()
		return err
	}
	c.Connection = nc
	return nil
}

func (c *client) Closed() {
	c.logger.Debug("Disconnecting nats")
	c.Connection.Flush()
	c.Connection.Drain()
	c.Connection.Close()
}

func (c *client) IsClosed() bool {
	return c.Connection.IsClosed()
}

func (c *client) IsConnected() bool {
	return c.Connection.IsConnected()
}

func (c *client) ReconnectIfUnhealthy(config DialConfig) {
	hasConnectionError := c.IsConnected()
	if !hasConnectionError {
		c.logger.Debug("Reconnecting to nats")
		c.Connect(config)
	}
}
