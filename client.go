package broadcast

import "sync"

// NewClient subscribes to a given topic and consumes or publish messages to it. NewClient will be receiving messages it
// produced.
type Client struct {
	upstream  chan *Message
	broadcast *Service
	mu        sync.Mutex
	topics    []string
}

// NewClient client to specific topics.
func (c *Client) Subscribe(topics ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	newTopics := make([]string, 0)
	for _, topic := range topics {
		found := false
		for _, e := range c.topics {
			if e == topic {
				found = true
				break
			}
		}

		if !found {
			newTopics = append(newTopics, topic)
		}
	}

	c.topics = append(c.topics, newTopics...)
	if len(newTopics) == 0 {
		return nil
	}

	return c.broadcast.Subscribe(c.upstream, newTopics...)
}

// Unsubscribe client from specific topics
func (c *Client) Unsubscribe(topics ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	dropTopics := make([]string, 0)
	for _, topic := range topics {
		for i, e := range c.topics {
			if e == topic {
				c.topics[i] = c.topics[len(c.topics)-1]
				c.topics = c.topics[:len(c.topics)-1]
				dropTopics = append(dropTopics, topic)
			}
		}
	}

	if len(dropTopics) == 0 {
		return
	}

	c.broadcast.Unsubscribe(c.upstream, dropTopics...)
}

// Publish message into associated topic or topics.
func (c *Client) Publish(msg ...*Message) error {
	return c.broadcast.Broadcast(msg...)
}

// Topics return all the topics client subscribed to.
func (c *Client) Topics() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.topics
}

// Close the client and consumption.
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.topics) != 0 {
		c.broadcast.Unsubscribe(c.upstream, c.topics...)
	}

	close(c.upstream)
}
