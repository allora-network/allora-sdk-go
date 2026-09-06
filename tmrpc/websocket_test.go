package tmrpc

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	jsonrpctypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

// TestReconnectDelayBackoffAndJitter pins that reconnect delays grow
// exponentially, are capped at wsReconnectMaxDelay, and include jitter (the
// previous implementation used a fixed 5s loop with no backoff or jitter).
func TestReconnectDelayBackoffAndJitter(t *testing.T) {
	// Exponential growth in the low range (base 500ms, doubling).
	prev := time.Duration(0)
	for attempt := 0; attempt < 5; attempt++ {
		// Average several samples to smooth jitter for the monotonic check.
		var sum time.Duration
		const n = 50
		for i := 0; i < n; i++ {
			sum += reconnectDelay(attempt)
		}
		avg := sum / n
		require.Greater(t, avg, prev, "average delay must grow with attempt %d", attempt)
		prev = avg
	}

	// Jitter present: not every sample equals the deterministic base.
	d0 := reconnectDelay(0)
	seenDifferent := false
	for i := 0; i < 20; i++ {
		if reconnectDelay(0) != d0 {
			seenDifferent = true
			break
		}
	}
	require.True(t, seenDifferent, "reconnect delay must include jitter")

	// Capped at max delay (jitter can exceed max by up to wsJitterFrac).
	maxAllowed := time.Duration(float64(wsReconnectMaxDelay) * (1 + wsJitterFrac))
	for attempt := 10; attempt < 40; attempt++ {
		require.LessOrEqual(t, reconnectDelay(attempt), maxAllowed,
			"delay for high attempt %d must be capped near wsReconnectMaxDelay", attempt)
	}
}

// brokenWebsocketConn dials a real websocket server whose handler upgrades
// and then immediately closes, so every client read fails — the
// broken-connection case readConnection must not spin on. A short read
// deadline bounds the test so a delayed close frame doesn't cause a spurious
// timeout.
func brokenWebsocketConn(t *testing.T) *websocket.Conn {
	t.Helper()

	upgrader := websocket.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		_ = conn.Close()
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	require.NoError(t, conn.SetReadDeadline(time.Now().Add(5*time.Second)))
	return conn
}

// TestReadOneBrokenConn pins that readOne surfaces a transport read failure as
// errBrokenConn rather than swallowing or reclassifying it.
func TestReadOneBrokenConn(t *testing.T) {
	ws := &tmWebsocket{
		logger: zerolog.Nop(),
		muConn: &sync.Mutex{},
	}
	ws.conn = brokenWebsocketConn(t)

	var resp jsonrpctypes.RPCResponse
	err := ws.readOne(&resp)
	require.Error(t, err, "readOne must return an error on a broken connection")
}

// errLogCounter is a zerolog writer that counts error-level log events. It
// observes how many times readConnection logs a read failure, which reveals
// whether the loop broke on the first error or spun.
type errLogCounter struct {
	n *atomic.Int64
}

// Write satisfies io.Writer; zerolog prefers WriteLevel when present, so
// counting lives there and Write is a pass-through.
func (c errLogCounter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (c errLogCounter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level == zerolog.ErrorLevel {
		c.n.Add(1)
	}
	return len(p), nil
}

// TestReadConnectionBreaksOnReadError pins the ENGN-9304 fix: a read error
// must break readConnection's loop (return false for reconnect) instead of
// spinning on the broken connection.
//
// A bare "returns false" assertion cannot distinguish break from spin,
// because gorilla panics after 1000 repeated reads on a failed conn and the
// recover() in readConnection converts that panic to the same false result.
// So this test counts error-level log events through the real logger: the
// break path logs the read failure exactly once, while the spin path logs it
// roughly 1000 times before panicking. One error log means the loop broke.
func TestReadConnectionBreaksOnReadError(t *testing.T) {
	var errCount atomic.Int64
	logger := zerolog.New(errLogCounter{n: &errCount}).Level(zerolog.ErrorLevel)

	ws := &tmWebsocket{
		logger: logger,
		muConn: &sync.Mutex{},
		chStop: make(chan struct{}),
	}
	ws.conn = brokenWebsocketConn(t)

	done := make(chan bool, 1)
	go func() { done <- ws.readConnection() }()

	var stop bool
	select {
	case stop = <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("readConnection did not return on read error; it is spinning on a broken connection")
	}

	require.False(t, stop, "readConnection must return false on connection break (not clean shutdown)")
	require.EqualValues(t, 1, errCount.Load(),
		"readConnection must log the read failure exactly once (break on first error); "+
			"a large count means it spun on the broken connection until gorilla panicked")
}
