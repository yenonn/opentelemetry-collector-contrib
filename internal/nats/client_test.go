package nats

import (
	"testing"
	"time"

	"gotest.tools/assert"
)

func TestDefaultDailConfig(t *testing.T) {
	config := DefaultDialConfig()
	assert.Equal(t, "nats://localhost:4222", config.URL)
	assert.Equal(t, 5*time.Second, config.ConnectionTimeout)
	assert.Equal(t, "", config.Username)
	assert.Equal(t, "", config.Password)
	assert.Equal(t, "", config.RootCA)
}
