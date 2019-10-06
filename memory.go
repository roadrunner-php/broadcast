package broadcast

type streamContext struct {
	upstream chan *Message
	topics   []string
}

// Memory manages broadcasting in memory.
type Memory struct {
	routes        map[string][]chan *Message
	messages      chan *Message
	listen, leave chan streamContext
	stop          chan interface{}
}

// memoryBroker creates new memory based message broker.
func memoryBroker() *Memory {
	return &Memory{
		routes:   make(map[string][]chan *Message),
		messages: make(chan *Message),
		listen:   make(chan streamContext),
		leave:    make(chan streamContext),
		stop:     make(chan interface{}),
	}
}

// Serve serves broker.
func (m *Memory) Serve() error {
	for {
		select {
		case ctx := <-m.listen:
			for _, topic := range ctx.topics {
				if _, ok := m.routes[topic]; !ok {
					m.routes[topic] = make([]chan *Message, 0)
				}

				joined := false
				for _, up := range m.routes[topic] {
					if up == ctx.upstream {
						joined = true
						break
					}
				}

				if !joined {
					m.routes[topic] = append(m.routes[topic], ctx.upstream)
				}
			}
		case ctx := <-m.leave:
			for _, topic := range ctx.topics {
				if _, ok := m.routes[topic]; !ok {
					continue
				}

				for i, up := range m.routes[topic] {
					if up == ctx.upstream {
						m.routes[topic][i] = m.routes[topic][len(m.routes[topic])-1]
						m.routes[topic][len(m.routes[topic])-1] = nil
						m.routes[topic] = m.routes[topic][:len(m.routes[topic])-1]
						break
					}
				}

				if len(m.routes[topic]) == 0 {
					// topic has no subscribers
					delete(m.routes, topic)
				}
			}
		case msg := <-m.messages:
			if _, ok := m.routes[msg.Topic]; !ok {
				continue
			}

			for _, upstream := range m.routes[msg.Topic] {
				upstream <- msg
			}

		case <-m.stop:
			return nil
		}
	}
}

// Stop the consumption and disconnect broker.
func (m *Memory) Stop() {
	close(m.stop)
}

// Subscribe broker to one or multiple channels.
func (m *Memory) Subscribe(upstream chan *Message, topics ...string) error {
	m.listen <- streamContext{upstream: upstream, topics: topics}
	return nil
}

// Unsubscribe broker from one or multiple channels.
func (m *Memory) Unsubscribe(upstream chan *Message, topics ...string) {
	m.leave <- streamContext{upstream: upstream, topics: topics}
}

// Broadcast one or multiple messages.
func (m *Memory) Broadcast(messages ...*Message) error {
	for _, msg := range messages {
		m.messages <- msg
	}

	return nil
}
