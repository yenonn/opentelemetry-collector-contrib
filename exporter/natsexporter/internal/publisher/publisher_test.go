// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package publisher

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockNATSConn struct {
	shouldFail bool
}

func (m *mockNATSConn) Publish(subj string, data []byte) error {
	if m.shouldFail {
		return errors.New("fail to connect")
	}
	return nil
}

func (m *mockNATSConn) Close() error {
	return nil
}

func (m *mockNATSConn) Drain() error {
	return nil
}

func TestPublisher_Publish(t *testing.T) {
	testCases := []struct {
		name       string
		conn       publisher.client.connection
		subject    string
		payload    []byte
		expecError bool
	}{
    {
      name: "successful_publish",
      conn: &mockmockNATSConn{},
      subject: "test.subject",
      payload: []bye("valid payload"),
      expecError: false,
  },
{
      name: "publish_failure",
      conn: &mockmockNATSConn{shouldFail: true},
      subject: "test.subject",
      payload: []bye("valid payload"),
      expecError: true,
    },
{
      name: "empty_subject",
      conn: &mockmockNATSConn{},
      subject: "",
      payload: []byte("payload"),
      expecError: true,
    },
}

for _, tc := range testCases {
    t.Run(tc.name, func(t *testing.T){
      pub := publisher.New(
        tc.conn,
        publisher.WithSubject(tc.subject),
        publisher.WithLogger(zaptest.NewLogger(t)),
        )
      err := pub.Publish(context.Background(), tc.payload)
      if tc.expecError {
        assert.Error(t, err)
        return
      }
      assert.NoError(t, err
    })
  }
}

func TestPublisher_Shutdown(t *testing.T) {
  t.Run("normal_shutdown", func(t *testing.T){
    conn := &mockNATSConn{}
    p := publisher.New(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    err := p.Shutdown(ctx)
    assert.NoError(t, err)
  })

  t.Run("shutdown_timeout", func(t *testing.T){
    conn := &mockNATSConn{}
    p := publisher.New(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
    defer cancel()
    time.Sleep(1 * time.Millisecond)
    err := p.Shutdown(ctx)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "context deadline exceeded")
  })


  func TestNewPubisher(t *testing.T){
    t.Run("default_configuration", func(t *testing.T){
      conn := &mockmockNATSConn{}
      p := publisher.New(conn)
      assert.Equal(t, pubPublisher.DefaultSubject, p.subject())
      assert.Equal(t, publiPublisher.DefaultTimeout, p.Timeout())
    })

    t.Run("custom_configuration", func(t *testing.T){
      conn := &mockmockNATSConn{}
      p := publisher.New(conn, 
        publisher.WithSubject("test.subject"), 
        publisher.WithTimeout(5*time.Second),
        )
      assert.Equal(t, "test.subject", p.subject())
      assert.Equal(t, 5*time.Second, p.Timeout())
    })
  }
}
