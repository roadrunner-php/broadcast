package broadcast

// Broker defines the ability to operate as message passing broker.
type Broker interface {
	// Serve serves broker.
	Serve() error

	// close the consumption and disconnect broker.
	Stop()

	// Subscribe broker to one or multiple topics.
	Subscribe(upstream chan *Message, topics ...string) error

	// Unsubscribe broker from one or multiple topics.
	Unsubscribe(upstream chan *Message, topics ...string)

	// Broadcast one or multiple Messages.
	Broadcast(messages ...*Message) error
}
