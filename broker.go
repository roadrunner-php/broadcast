package broadcast

import "encoding/json"

// Broker defines the ability to operate as message passing broker.
type Broker interface {
	// Serve serves broker.
	Serve() error

	// close the consumption and disconnect broker.
	Stop()

	// Subscribe broker to one or multiple topics.
	Subscribe(upstream chan *Message, topics ...string) error

	// Subscribe broker to one or multiple topics.
	SubscribePattern(upstream chan *Message, pattern string) error

	// Unsubscribe broker from one or multiple topics.
	Unsubscribe(upstream chan *Message, topics ...string) error

	// Unsubscribe broker from one or multiple topics.
	UnsubscribePattern(upstream chan *Message, pattern string) error

	// Publish one or multiple Channel.
	Publish(messages ...*Message) error
}

// Message represent single message.
type Message struct {
	// Topic message been pushed into.
	Topic string `json:"topic"`

	// Payload to be broadcasted. Must be valid json when transferred over RPC.
	Payload json.RawMessage `json:"payload"`
}
