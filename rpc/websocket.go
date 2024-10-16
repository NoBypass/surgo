package rpc

import (
	"context"
	"github.com/coder/websocket"
	"sync"
)

type WebsocketConn struct {
	*websocket.Conn
	mu        sync.RWMutex
	responses map[string]chan Response
	logger    Logger
}

type Logger interface {
	Error(err error)
}

func NewWebsocketConn(url string, logger Logger) (*WebsocketConn, error) {
	c, _, err := websocket.Dial(context.Background(), url, nil)
	if err != nil {
		return nil, err
	}

	conn := &WebsocketConn{
		Conn:      c,
		responses: make(map[string]chan Response),
		logger:    logger,
	}
	go conn.listen()
	return conn, nil
}

func (c *WebsocketConn) Close() error {
	return c.Conn.Close(websocket.StatusNormalClosure, "")
}
