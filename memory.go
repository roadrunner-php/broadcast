package broadcast

type subscriber struct {
	upstream chan *Message
	done     chan interface{}
	topics   []string
}

// Memory manages broadcasting in memory.
type Memory struct {
	routes        map[string][]chan *Message
	messages      chan *Message
	listen, leave chan subscriber
	stop          chan interface{}
}

// memoryBroker creates new memory based message broker.
func memoryBroker() *Memory {
	return &Memory{
		routes:   make(map[string][]chan *Message),
		messages: make(chan *Message),
		listen:   make(chan subscriber),
		leave:    make(chan subscriber),
		stop:     make(chan interface{}),
	}
}

// Serve serves broker.
func (m *Memory) Serve() error {
	for {
		select {
		case ctx := <-m.listen:
			m.handleJoin(ctx)
			close(ctx.done)
		case ctx := <-m.leave:
			m.handleLeave(ctx)
			close(ctx.done)
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

func (m *Memory) handleLeave(sb subscriber) {
	for _, topic := range sb.topics {
		if _, ok := m.routes[topic]; !ok {
			continue
		}

		for i, up := range m.routes[topic] {
			if up == sb.upstream {
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
}

func (m *Memory) handleJoin(sb subscriber) {
	for _, topic := range sb.topics {
		if _, ok := m.routes[topic]; !ok {
			m.routes[topic] = make([]chan *Message, 0)
		}

		joined := false
		for _, up := range m.routes[topic] {
			if up == sb.upstream {
				joined = true
				break
			}
		}

		if !joined {
			m.routes[topic] = append(m.routes[topic], sb.upstream)
		}
	}
}

// close the consumption and disconnect broker.
func (m *Memory) Stop() {
	close(m.stop)
}

// Subscribe broker to one or multiple channels.
func (m *Memory) Subscribe(upstream chan *Message, topics ...string) error {
	ctx := subscriber{upstream: upstream, topics: topics, done: make(chan interface{})}

	m.listen <- ctx
	<-ctx.done

	return nil
}

// Unsubscribe broker from one or multiple channels.
func (m *Memory) Unsubscribe(upstream chan *Message, topics ...string) {
	ctx := subscriber{upstream: upstream, topics: topics, done: make(chan interface{})}

	m.leave <- ctx
	<-ctx.done
}

// Broadcast one or multiple messages.
func (m *Memory) Broadcast(messages ...*Message) error {
	for _, msg := range messages {
		m.messages <- msg
	}

	return nil
}
