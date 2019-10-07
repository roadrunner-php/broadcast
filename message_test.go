package broadcast

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMessage(t *testing.T) {
	d, _ := json.Marshal(NewMessage("topic", "message"))

	assert.Equal(t, []byte(`{"topic":"topic","payload":"message"}`), d)
}

func TestCommandPayload(t *testing.T) {
	c := &Command{
		Cmd:  "test",
		Args: []byte(`{"topic":"topic","payload":"message"}`),
	}

	m := new(Message)
	assert.NoError(t, c.Unmarshal(m))

	assert.Equal(t, NewMessage("topic", "message"), m)
}
