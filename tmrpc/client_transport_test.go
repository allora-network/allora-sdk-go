package tmrpc

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewPooledHTTPClientConfiguresConnectionPool(t *testing.T) {
	client, err := newPooledHTTPClient("http://127.0.0.1:26657", 7*time.Second)
	require.NoError(t, err)
	require.Equal(t, 7*time.Second, client.Timeout)

	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.Equal(t, maxIdleConns, transport.MaxIdleConns)
	require.Equal(t, maxIdleConnsPerHost, transport.MaxIdleConnsPerHost)
	require.Equal(t, idleConnTimeout, transport.IdleConnTimeout)
}

func TestNewPooledHTTPClientLeavesTimeoutUnsetWhenNonPositive(t *testing.T) {
	client, err := newPooledHTTPClient("http://127.0.0.1:26657", 0)
	require.NoError(t, err)
	require.Zero(t, client.Timeout)
}

func TestNewPooledHTTPClientRejectsInvalidRemote(t *testing.T) {
	_, err := newPooledHTTPClient("::not-a-url::", time.Second)
	require.Error(t, err)
}
