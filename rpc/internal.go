package rpc

import (
	"context"
	"encoding/json"
	"github.com/NoBypass/surgo/v2/errs"
	"github.com/NoBypass/surgo/v2/rand"
	"github.com/coder/websocket"
)

func (c *WebsocketConn) listen() {
	for {
		_, msg, err := c.Read(context.TODO())
		if err != nil {
			c.logger.Error(err)
			return
		}
		go c.receive(msg)
	}
}

func (c *WebsocketConn) receive(msg []byte) {
	var res Response
	err := json.Unmarshal(msg, &res)
	if err != nil {
		c.logger.Error(err)
		return
	}

	c.mu.RLock()
	ch, ok := c.responses[res.ID]
	c.mu.RUnlock()
	if !ok {
		c.logger.Error(errs.ErrUnexpectedResponseID.Withf("id: %s", res.ID))
		return
	}

	ch <- res
}

func (c *WebsocketConn) Send(ctx context.Context, method string, params []any) (any, error) {
	ch := make(chan Response)
	id := rand.String(16)

	c.mu.Lock()
	c.responses[id] = ch
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		delete(c.responses, id)
		c.mu.Unlock()
	}()

	reqBytes, err := json.Marshal(&Request{
		ID:     id,
		Method: method,
		Params: params,
	})
	if err != nil {
		return nil, err
	}
	err = c.Write(ctx, websocket.MessageText, reqBytes)
	if err != nil {
		return nil, err
	}

	select {
	case res := <-ch:
		if res.Error != nil {
			return nil, errs.ErrDatabase.With(res.Error)
		}
		return res.Result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
