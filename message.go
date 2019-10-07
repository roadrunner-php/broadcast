package broadcast

import "encoding/json"

// Message represent single message.
type Message struct {
	// Topic message been pushed into.
	Topic string `json:"topic"`

	// Payload to be broadcasted. Must be valid JSON.
	Payload json.RawMessage `json:"payload"`
}

// NewMessage creates new message with JSON payload.
func NewMessage(topic string, payload interface{}) *Message {
	msg := &Message{Topic: topic}
	msg.Payload, _ = json.Marshal(payload)

	return msg
}

// Command contains information send by user.
type Command struct {
	// Cmd type.
	Cmd string `json:"cmd"`

	// Args contains command specific payload.
	Args json.RawMessage `json:"args"`
}

// Unmarshal command data.
func (cmd *Command) Unmarshal(v interface{}) error {
	return json.Unmarshal(cmd.Args, v)
}
