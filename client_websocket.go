package allora

import (
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

type CometRPCWebsocket struct {
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

type sub struct {
	id    int
	event string
	mb    *butils.Mailbox[ctypes.TMEventData]
}

func NewCometRPCWebsocket(rpcURL string, logger zerolog.Logger) *CometRPCWebsocket {
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

	ws := &CometRPCWebsocket{
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

func (ws *CometRPCWebsocket) Close() {
	close(ws.chStop)
	ws.wgDone.Wait()
}

func (ws *CometRPCWebsocket) connectionManager() {
	defer ws.wgDone.Done()
	defer ws.terminateConnection()

	for {
		ws.readConnection()

		select {
		case <-ws.chStop:
			return
		default:
		}

		ws.resetConnection()
	}
}

func (ws *CometRPCWebsocket) readConnection() {
	defer func() {
		if perr := recover(); perr != nil {
			ws.logger.Error().Msgf("recovered from panic: %v", perr)
		}
	}()

	for {
		select {
		case <-ws.chStop:
			return
		default:
		}

		var resp jsonrpctypes.RPCResponse
		err := ws.read(&resp)
		if err != nil {
			ws.logger.Error().Err(err).Msg("could not read websocket msg")
			continue
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

		subID := int(resp.ID.(jsonrpctypes.JSONRPCIntID))
		sub, ok := ws.subs.Get(subID)
		if !ok {
			ws.logger.Error().Msgf("received event for unknown subscription ID %d", subID)
			continue
		}

		sub.mb.Deliver(event.Data)
	}
}

func (ws *CometRPCWebsocket) read(resp any) error {
	ws.muConn.Lock()
	defer ws.muConn.Unlock()

	return ws.conn.ReadJSON(&resp)
}

func (ws *CometRPCWebsocket) resetConnection() {
	ws.muConn.Lock()
	defer ws.muConn.Unlock()

	// close connection if active
	if ws.conn != nil {
		ws.logger.Info().Msg("websocket connection closed, reconnecting...")
		err := ws.conn.Close()
		if err != nil {
			ws.logger.Error().Err(err).Msg("error closing websocket connection")
		} else {
			ws.logger.Info().Msg("websocket connection closed, reconnecting...")
		}
		ws.conn = nil
	}

	// wait for a new connection
	for {
		conn, _, err := websocket.DefaultDialer.Dial(ws.url, nil)
		if err != nil {
			ws.logger.Error().Err(err).Msg("websocket dial failed")
			select {
			case <-ws.chStop:
				return
			case <-time.After(5 * time.Second):
			}
			continue
		}

		ws.conn = conn
		break
	}

	ws.logger.Info().Str("url", ws.conn.RemoteAddr().String()).Msg("connected to comet rpc websocket")

	// resubscribe to everything
	for _, sub := range ws.subs.Iter() {
		ws.logger.Debug().Msgf("subscribing to subscription ID %d with event %v", sub.id, sub.event)
		ws.sendSubscribeMsg(sub)
	}
}

// terminateConnection closes the websocket connection and cleans up resources permanently.
func (ws *CometRPCWebsocket) terminateConnection() {
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

func (ws *CometRPCWebsocket) Subscribe(mb *butils.Mailbox[ctypes.TMEventData], query string) {
	ws.subIDNonce++

	sub := sub{
		id:    ws.subIDNonce,
		event: query,
		mb:    mb,
	}

	ws.subs.Set(sub.id, sub)
	ws.sendSubscribeMsg(sub)
}

func (ws *CometRPCWebsocket) sendSubscribeMsg(sub sub) {
	subMsg := map[string]any{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"id":      sub.id,
		"params": map[string]any{
			"query": sub.event,
		},
	}

	ws.logger.Info().Msg("subscribing to " + sub.event)

	err := ws.conn.WriteJSON(subMsg)
	if err != nil {
		ws.logger.Error().Err(err).Msg("could not write subscription message")
	}
}
