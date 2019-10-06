package broadcast

import "encoding/json"

// Payload represent single message.
type Message struct {
	// Topic message been pushed into.
	Topic string `json:"topic"`

	// Payload to be broadcasted. Must be valid JSON.
	Payload json.RawMessage `json:"message"`
}

// NewMessage creates new message with JSON payload.
func NewMessage(topic string, payload interface{}) *Message {
	msg := &Message{Topic: topic}
	msg.Payload, _ = json.Marshal(payload)

	return msg
}

// Command contains information send by user.
type Command struct {
	// Command type.
	Command string `json:"command"`

	// Data contains command specific payload.
	Data json.RawMessage `json:"data"`
}

// Unmarshal command data.
func (cmd *Command) Unmarshal(v interface{}) error {
	return json.Unmarshal(cmd.Data, v)
}
