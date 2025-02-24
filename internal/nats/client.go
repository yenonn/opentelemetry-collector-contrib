package nats

import (
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type INatsClient interface {
	Connect(config *DialConfig) error
	Close()
	IsCclose() bool
	IsConnected() bool
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

type NatsClient struct {
	logger     *zap.Logger
	Connection *nats.Conn
}

func NewNatsClient(logger *zap.Logger) *NatsClient {
	return &NatsClient{
		logger: logger,
	}
}

func (c *NatsClient) Connect(config DialConfig) (*nats.Conn, error) {
	c.logger.Debug("Connecting to nats")

	natsOptions := []nats.Option{
		nats.UserInfo(config.Username, config.Password),
		nats.Timeout(config.ConnectionTimeout),
		nats.RootCAs(config.RootCA),
	}

	nc, err := nats.Connect(config.URL, natsOptions...)
	if err != nil {
		nc.Close()
		return nil, err
	}
	c.Connection = nc
	return c.Connection, nil
}

func (c *NatsClient) Close() {
	c.logger.Debug("Disconnecting nats")
	c.Connection.Flush()
	c.Connection.Drain()
	c.Connection.Close()
}

func (c *NatsClient) IsClosed() bool {
	return c.Connection.IsClosed()
}

func (c *NatsClient) IsConnected() bool {
	return c.Connection.IsConnected()
}
