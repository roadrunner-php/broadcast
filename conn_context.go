package broadcast

import "github.com/gorilla/websocket"

// ConnContext represents the connection and it's state.
type ConnContext struct {
	// Upstream to push messages into.
	Upstream chan *Message

	// Conn to the client.
	Conn *websocket.Conn

	// Topics contain list of currently subscribed topics.
	Topics []string
}

func (ctx *ConnContext) addTopic(topics ...string) {
	for _, topic := range topics {
		found := false
		for _, e := range ctx.Topics {
			if e == topic {
				found = true
				break
			}
		}

		if !found {
			ctx.Topics = append(ctx.Topics, topic)
		}
	}
}

func (ctx *ConnContext) dropTopic(topics ...string) {
	for _, topic := range topics {
		for i, e := range ctx.Topics {
			if e == topic {
				ctx.Topics[i] = ctx.Topics[len(ctx.Topics)-1]
				ctx.Topics = ctx.Topics[:len(ctx.Topics)-1]
			}
		}
	}
}
