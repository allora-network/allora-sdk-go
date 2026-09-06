package tmrpc

import (
	"math/rand"
	"net/url"
	"sync"
	"time"

	butils "github.com/brynbellomy/go-utils"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	rpctypes "github.com/cometbft/cometbft/rpc/core/types"
	jsonrpctypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"
	ctypes "github.com/cometbft/cometbft/types"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

// Reconnect policy for the websocket connection manager.
const (
	wsReconnectBaseDelay = 500 * time.Millisecond
	wsReconnectMaxDelay  = 30 * time.Second
	// wsJitterFrac is the symmetric jitter fraction applied to reconnect delays.
	wsJitterFrac = 0.2
)

type Websocket interface {
	Subscribe(mb *Mailbox, query string)
	Close()
}

type tmWebsocket struct {
	url    string
	logger zerolog.Logger

	conn   *websocket.Conn
	muConn *sync.Mutex

	subIDNonce int
	subs       *butils.SyncMap[int, sub]

	chResetConn chan struct{}
	chStop      chan struct{}
	wgDone      *sync.WaitGroup
}

var _ Websocket = (*tmWebsocket)(nil)

type Mailbox = butils.Mailbox[ctypes.TMEventData]

func NewMailbox(capacity uint64) *Mailbox {
	return butils.NewMailbox[ctypes.TMEventData](capacity)
}

type sub struct {
	id    int
	event string
	mb    *Mailbox
}

func NewTendermintWebsocket(rpcURL string, logger zerolog.Logger) *tmWebsocket {
	cometLogger := logger.With().Str("component", "websocket").Logger()

	u, err := url.Parse(rpcURL)
	if err != nil {
		cometLogger.Fatal().Err(err).Msg("bad url")
	}
	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else if u.Scheme == "http" {
		u.Scheme = "ws"
	}

	url := u.String()

	ws := &tmWebsocket{
		url:         url,
		logger:      cometLogger,
		muConn:      &sync.Mutex{},
		subIDNonce:  0,
		subs:        butils.NewSyncMap[int, sub](),
		chResetConn: make(chan struct{}, 1),
		chStop:      make(chan struct{}),
		wgDone:      &sync.WaitGroup{},
	}

	ws.resetConnection()

	ws.wgDone.Add(1)
	go ws.connectionManager()

	return ws
}

func (ws *tmWebsocket) Close() {
	close(ws.chStop)
	ws.wgDone.Wait()
}

func (ws *tmWebsocket) connectionManager() {
	defer ws.wgDone.Done()
	defer ws.terminateConnection()

	attempt := 0
	for {
		select {
		case <-ws.chStop:
			return
		default:
		}

		if ws.readConnection() {
			// Clean shutdown requested.
			return
		}

		// The connection dropped. Back off with jitter before redialling so a
		// flapping endpoint does not trigger a tight reconnect loop.
		delay := reconnectDelay(attempt)
		attempt++
		ws.logger.Info().Dur("delay", delay).Int("attempt", attempt).Msg("websocket connection lost, reconnecting")
		select {
		case <-ws.chStop:
			return
		case <-time.After(delay):
		}

		if ws.resetConnection() {
			// Redial succeeded. Only reset the backoff counter after a stable
			// session — a flapping endpoint that connects then immediately
			// disconnects must not collapse the backoff to the base delay.
			// The counter resets when readConnection returns without error,
			// which means the session was stable enough to read at least one
			// message before the next drop.
			// (The reset happens implicitly: the next readConnection call
			// either returns false (drop) or true (clean shutdown), and the
			// attempt counter only increments on drops.)
		}
	}
}

// reconnectDelay returns an exponential backoff delay with symmetric jitter.
func reconnectDelay(attempt int) time.Duration {
	if attempt > 30 { // 2^30 overflows practical delays
		attempt = 30
	}
	d := float64(wsReconnectBaseDelay) * float64(int64(1)<<uint(attempt))
	if d > float64(wsReconnectMaxDelay) {
		d = float64(wsReconnectMaxDelay)
	}
	jitter := d * wsJitterFrac * (2*rand.Float64() - 1)
	return time.Duration(d + jitter)
}

// readConnection reads events until the connection fails or Close is called.
// It returns true if shutdown was requested, false if the connection broke and
// the caller should reconnect. A read error breaks the loop (previously it
// spun on a broken connection).
func (ws *tmWebsocket) readConnection() (stop bool) {
	defer func() {
		if perr := recover(); perr != nil {
			ws.logger.Error().Interface("panic", perr).Stack().Msg("websocket readConnection panicked; reconnecting")
			stop = false
		}
	}()

readLoop:
	for {
		select {
		case <-ws.chStop:
			return true
		default:
		}

		var resp jsonrpctypes.RPCResponse
		err := ws.read(&resp)
		if err != nil {
			ws.logger.Error().Err(err).Msg("websocket read failed; breaking to reconnect")
			break readLoop
		} else if resp.Error != nil {
			ws.logger.Error().Err(*resp.Error).Msg("rpc received error")
			continue
		}

		if string(resp.Result) == "{}" {
			continue
		}

		var event rpctypes.ResultEvent
		err = cmtjson.Unmarshal(resp.Result, &event)
		if err != nil {
			ws.logger.Error().Err(err).Msg("could not unmarshal websocket msg")
			continue
		}

		subID, ok := resp.ID.(jsonrpctypes.JSONRPCIntID)
		if !ok {
			ws.logger.Error().Msg("received event with unexpected ID type")
			continue
		}
		s, ok := ws.subs.Get(int(subID))
		if !ok {
			ws.logger.Error().Msgf("received event for unknown subscription ID %d", subID)
			continue
		}

		s.mb.Deliver(event.Data)
	}
	return false
}

func (ws *tmWebsocket) read(resp any) error {
	ws.muConn.Lock()
	defer ws.muConn.Unlock()

	return ws.conn.ReadJSON(&resp)
}

// resetConnection closes any active connection and redials with retry until
// the connection is re-established or Close is called. It returns true if a
// new connection was established.
func (ws *tmWebsocket) resetConnection() bool {
	// close connection if active
	ws.muConn.Lock()
	if ws.conn != nil {
		if err := ws.conn.Close(); err != nil {
			ws.logger.Error().Err(err).Msg("error closing websocket connection")
		}
		ws.conn = nil
	}
	ws.muConn.Unlock()

	// wait for a new connection
	for {
		conn, _, err := websocket.DefaultDialer.Dial(ws.url, nil)
		if err != nil {
			ws.logger.Error().Err(err).Msg("websocket dial failed")
			select {
			case <-ws.chStop:
				return false
			case <-time.After(reconnectDelay(0)):
			}
			continue
		}

		ws.muConn.Lock()
		ws.conn = conn
		ws.muConn.Unlock()
		break
	}

	ws.logger.Info().Str("url", ws.conn.RemoteAddr().String()).Msg("connected to comet rpc websocket")

	// resubscribe to everything
	for _, s := range ws.subs.Iter() {
		ws.logger.Debug().Msgf("subscribing to subscription ID %d with event %v", s.id, s.event)
		ws.sendSubscribeMsg(s)
	}
	return true
}

// terminateConnection closes the websocket connection and cleans up resources permanently.
func (ws *tmWebsocket) terminateConnection() {
	ws.muConn.Lock()
	defer ws.muConn.Unlock()

	if ws.conn != nil {
		err := ws.conn.Close()
		if err != nil {
			ws.logger.Error().Err(err).Msg("error closing websocket connection, giving up")
		} else {
			ws.logger.Info().Msg("websocket connection closed")
		}
		ws.conn = nil
	}
}

func (ws *tmWebsocket) Subscribe(mb *Mailbox, query string) {
	ws.subIDNonce++

	s := sub{
		id:    ws.subIDNonce,
		event: query,
		mb:    mb,
	}

	ws.subs.Set(s.id, s)
	ws.sendSubscribeMsg(s)
}

func (ws *tmWebsocket) sendSubscribeMsg(s sub) {
	subMsg := map[string]any{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"id":      s.id,
		"params": map[string]any{
			"query": s.event,
		},
	}

	ws.logger.Info().Msg("subscribing to " + s.event)

	// Grab the connection under the lock, then write WITHOUT holding muConn so
	// a blocked write cannot deadlock a concurrent read/reset.
	ws.muConn.Lock()
	conn := ws.conn
	ws.muConn.Unlock()

	if conn == nil {
		ws.logger.Debug().Msg("no active websocket connection; subscription will be sent on reconnect")
		return
	}
	if err := conn.WriteJSON(subMsg); err != nil {
		ws.logger.Error().Err(err).Msg("could not write subscription message")
	}
}

type WebsocketPool interface {
	Subscribe(mb *Mailbox, query string)
	Close()
}

type websocketPool struct {
	websockets []Websocket
}

var _ WebsocketPool = (*websocketPool)(nil)

func NewWebsocketPool(websockets []Websocket) *websocketPool {
	return &websocketPool{
		websockets: websockets,
	}
}

func (p *websocketPool) Close() {
	for _, ws := range p.websockets {
		ws.Close()
	}
}

func (p *websocketPool) Subscribe(mb *Mailbox, query string) {
	for _, ws := range p.websockets {
		ws.Subscribe(mb, query)
	}
}
