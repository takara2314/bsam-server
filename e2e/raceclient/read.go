package raceclient

import (
	"context"
)

type MessageIterator struct {
	ctx    context.Context
	client *Client
}

func (c *Client) ReceiveStream(ctx context.Context) *MessageIterator {
	return &MessageIterator{
		ctx:    ctx,
		client: c,
	}
}

func (it *MessageIterator) Read() ([]byte, error) {
	select {
	case msg := <-it.client.receiveCh:
		return msg, nil
	case <-it.ctx.Done():
		return nil, it.ctx.Err()
	case <-it.client.closeCh:
		return nil, ErrClientClosed
	}
}

func (c *Client) readPump() {
	defer c.Close()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		select {
		case c.receiveCh <- message:
		case <-c.closeCh:
			return
		}
	}
}
