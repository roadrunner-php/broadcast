package broadcast

import "github.com/go-redis/redis"

// Redis based broadcast Router.
type Redis struct {
	client        redis.UniversalClient
	router        *Router
	messages      chan *Message
	listen, leave chan subscriber
	stop          chan interface{}
}

// creates new redis broker
func redisBroker(cfg *RedisConfig) (*Redis, error) {
	client := cfg.redisClient()
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return &Redis{
		client:   client,
		router:   NewRouter(),
		messages: make(chan *Message),
		listen:   make(chan subscriber),
		leave:    make(chan subscriber),
		stop:     make(chan interface{}),
	}, nil
}

// Serve serves broker.
func (r *Redis) Serve() error {
	pubsub := r.client.Subscribe()
	channel := pubsub.Channel()

	for {
		select {
		case ctx := <-r.listen:
			ctx.done <- r.handleJoin(ctx, pubsub)
		case ctx := <-r.leave:
			ctx.done <- r.handleLeave(ctx, pubsub)
		case msg := <-channel:
			r.router.Dispatch(&Message{
				Topic:   msg.Channel,
				Payload: []byte(msg.Payload),
			})
		case <-r.stop:
			return nil
		}
	}
}

func (r *Redis) handleJoin(sub subscriber, pubsub *redis.PubSub) error {
	if sub.pattern != "" {
		newPatterns, err := r.router.SubscribePattern(sub.upstream, sub.pattern)
		if err != nil || len(newPatterns) == 0 {
			return err
		}

		return pubsub.PSubscribe(newPatterns...)
	}

	newTopics := r.router.Subscribe(sub.upstream, sub.topics...)
	if len(newTopics) == 0 {
		return nil
	}

	return pubsub.Subscribe(newTopics...)
}

func (r *Redis) handleLeave(sub subscriber, pubsub *redis.PubSub) error {
	if sub.pattern != "" {
		dropPatterns := r.router.UnsubscribePattern(sub.upstream, sub.pattern)
		if len(dropPatterns) == 0 {
			return nil
		}

		return pubsub.PUnsubscribe(dropPatterns...)
	}

	dropTopics := r.router.Unsubscribe(sub.upstream, sub.topics...)
	if len(dropTopics) == 0 {
		return nil
	}

	return pubsub.Unsubscribe(dropTopics...)
}

// close the consumption and disconnect broker.
func (r *Redis) Stop() {
	close(r.stop)
}

// Subscribe broker to one or multiple channels.
func (r *Redis) Subscribe(upstream chan *Message, topics ...string) error {
	ctx := subscriber{upstream: upstream, topics: topics, done: make(chan error)}

	r.listen <- ctx
	return <-ctx.done
}

func (r *Redis) SubscribePattern(upstream chan *Message, pattern string) error {
	ctx := subscriber{upstream: upstream, pattern: pattern, done: make(chan error)}

	r.listen <- ctx
	return <-ctx.done
}

// Unsubscribe broker from one or multiple channels.
func (r *Redis) Unsubscribe(upstream chan *Message, topics ...string) error {
	ctx := subscriber{upstream: upstream, topics: topics, done: make(chan error)}

	r.leave <- ctx
	return <-ctx.done
}

func (r *Redis) UnsubscribePattern(upstream chan *Message, pattern string) error {
	ctx := subscriber{upstream: upstream, pattern: pattern, done: make(chan error)}

	r.leave <- ctx
	return <-ctx.done
}

// Publish one or multiple Channel.
func (r *Redis) Publish(messages ...*Message) error {
	for _, msg := range messages {
		if err := r.client.Publish(msg.Topic, []byte(msg.Payload)).Err(); err != nil {
			return err
		}
	}

	return nil
}
