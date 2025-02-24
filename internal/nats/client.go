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

type client struct {
	logger     *zap.Logger
	Connection *nats.Conn
}

func (c *client) Connect(config DialConfig) error {
	c.logger.Debug("Connecting to nats")

	nc, err := nats.Connect(config.URL, nats.UserInfo(config.Username, config.Password), nats.RootCAs(config.RootCA))
	if err != nil {
		nc.Close()
		return err
	}
	c.Connection = nc
	return nil
}

func (c *client) Closed(config DialConfig) {
	c.logger.Debug("Disconnecting nats")
	c.Connection.Flush()
	c.Connection.Drain()
	c.Connection.Close()
}
