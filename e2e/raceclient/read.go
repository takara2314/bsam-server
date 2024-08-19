package raceclient

type MessageIterator struct {
	client *Client
}

func (c *Client) ReceiveStream() *MessageIterator {
	return &MessageIterator{
		client: c,
	}
}

func (it *MessageIterator) Read() ([]byte, error) {
	select {
	case msg := <-it.client.receiveCh:
		return msg, nil
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
