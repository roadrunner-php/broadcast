package broadcast

import (
	"github.com/gorilla/websocket"
)

const (
	// EventWebsocketConnect fired when new client is connected, the context is *websocket.Conn.
	EventWebsocketConnect = iota + 2500

	// EventWebsocketDisconnect fired when websocket is disconnected, context is empty.
	EventWebsocketDisconnect

	// EventWebsocketJoin caused when topics are being consumed, context if *TopicEvent.
	EventWebsocketJoin

	// EventWebsocketLeave caused when topic consumption are stopped, context if *TopicEvent.
	EventWebsocketLeave

	// EventWebsocketError when any broadcast error occurred, the context is *ErrorEvent.
	EventWebsocketError

	// EventBrokerError the context is error.
	EventBrokerError
)

// ErrorEvent represents singular broadcast error event.
type ErrorEvent struct {
	// Conn specific to the error.
	Conn *websocket.Conn

	// Caused contains job specific error.
	Caused error
}

// TopicEvent caused when topic is joined or left.
type TopicEvent struct {
	// Conn associated with topics.
	Conn *websocket.Conn

	// Topics specific to event.
	Topics []string
}
