package broadcast

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

// Redis based broadcast router.
type Redis struct {
	client        *redis.Client
	errHandler    func(err error, conn *websocket.Conn)
	routes        map[string][]chan *Message
	messages      chan *Message
	listen, leave chan subscriber
	stop          chan interface{}
}

// creates new redis broker
func redisBroker(cfg *RedisConfig, errHandler func(err error, conn *websocket.Conn)) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// todo: support redis cluster

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &Redis{
		client:     client,
		errHandler: errHandler,
		routes:     make(map[string][]chan *Message),
		messages:   make(chan *Message),
		listen:     make(chan subscriber),
		leave:      make(chan subscriber),
		stop:       make(chan interface{}),
	}, nil
}

// Serve serves broker.
func (r *Redis) Serve() error {
	pubsub := r.client.Subscribe()
	channel := pubsub.Channel()

	for {
		select {
		case ctx := <-r.listen:
			r.handleJoin(ctx, pubsub)

		case ctx := <-r.leave:
			r.handleLeave(ctx, pubsub)

		case msg := <-channel:
			if _, ok := r.routes[msg.Channel]; !ok {
				continue
			}

			for _, upstream := range r.routes[msg.Channel] {
				// we except that the payload is always valid json
				upstream <- &Message{Topic: msg.Channel, Payload: []byte(msg.Payload)}
			}

		case <-r.stop:
			return nil
		}
	}
}

func (r *Redis) handleLeave(sb subscriber, pubsub *redis.PubSub) {
	dropTopics := make([]string, 0)
	for _, topic := range sb.topics {
		if _, ok := r.routes[topic]; !ok {
			continue
		}

		for i, up := range r.routes[topic] {
			if up == sb.upstream {
				r.routes[topic][i] = r.routes[topic][len(r.routes[topic])-1]
				r.routes[topic][len(r.routes[topic])-1] = nil
				r.routes[topic] = r.routes[topic][:len(r.routes[topic])-1]
				break
			}
		}

		if len(r.routes[topic]) == 0 {
			// topic has no subscribers
			delete(r.routes, topic)
			dropTopics = append(dropTopics, topic)
		}
	}
	if len(dropTopics) != 0 {
		if err := pubsub.Unsubscribe(dropTopics...); err != nil {
			r.errHandler(err, nil)
		}
	}
}

func (r *Redis) handleJoin(sb subscriber, pubsub *redis.PubSub) {
	newTopics := make([]string, 0)
	for _, topic := range sb.topics {
		if _, ok := r.routes[topic]; !ok {
			r.routes[topic] = make([]chan *Message, 0)
			newTopics = append(newTopics, topic)
		}

		joined := false
		for _, up := range r.routes[topic] {
			if up == sb.upstream {
				joined = true
				break
			}
		}

		if !joined {
			r.routes[topic] = append(r.routes[topic], sb.upstream)
		}
	}
	if len(newTopics) != 0 {
		if err := pubsub.Subscribe(newTopics...); err != nil {
			r.errHandler(err, nil)
		}
	}
}

// Stop the consumption and disconnect broker.
func (r *Redis) Stop() {
	close(r.stop)
}

// Subscribe broker to one or multiple channels.
func (r *Redis) Subscribe(upstream chan *Message, topics ...string) error {
	r.listen <- subscriber{upstream: upstream, topics: topics}
	return nil
}

// Unsubscribe broker from one or multiple channels.
func (r *Redis) Unsubscribe(upstream chan *Message, topics ...string) {
	r.leave <- subscriber{upstream: upstream, topics: topics}
}

// Broadcast one or multiple messages.
func (r *Redis) Broadcast(messages ...*Message) error {
	for _, msg := range messages {
		if err := r.client.Publish(msg.Topic, []byte(msg.Payload)).Err(); err != nil {
			return err
		}
	}

	return nil
}
