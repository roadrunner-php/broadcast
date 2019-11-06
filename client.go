package broadcast

// Client subscribes to a given topic and consumes or publish messages to it. Client will be receiving messages it
// produced.
type Client struct {
	// Messages consumed from related topic. The messages MUST be consumed, otherwise system will be blocked.
	Messages chan *Message

	// internal binding
	broadcast *Service
	topic     string
}

// NewClient creates new broadcast client.
func NewClient(topic string, broadcast *Service) (*Client, error) {
	c := &Client{
		broadcast: broadcast,
		topic:     topic,
		Messages:  make(chan *Message),
	}

	if err := c.broadcast.Subscribe(c.Messages, c.topic); err != nil {
		return nil, err
	}

	return c, nil
}

// Publish message into associated topic.
func (c *Client) Publish(payload ...interface{}) error {
	messages := make([]*Message, 0, len(payload))
	for _, p := range payload {
		messages = append(messages, NewMessage(c.topic, p))
	}

	return c.broadcast.Broadcast(messages...)
}

// Close the client and consumption.
func (c *Client) Close() {
	c.broadcast.Unsubscribe(c.Messages, c.topic)
}
