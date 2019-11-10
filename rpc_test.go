package broadcast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRPC_Broadcast(t *testing.T) {
	br, rpc, c := setup(`{}`)
	defer c.Stop()

	client := br.NewClient()
	rcpClient, err := rpc.Client()
	assert.NoError(t, err)

	// must not be delivered
	ok := false
	assert.NoError(t, rcpClient.Call(
		"broadcast.Publish",
		[]*Message{newMessage("topic", `"hello1"`)},
		&ok,
	))
	assert.True(t, ok)

	assert.NoError(t, client.Subscribe("topic"))

	assert.NoError(t, rcpClient.Call(
		"broadcast.Publish",
		[]*Message{newMessage("topic", `"hello1"`)},
		&ok,
	))
	assert.True(t, ok)
	assert.Equal(t, `"hello1"`, readStr(<-client.Channel()))

	assert.NoError(t, rcpClient.Call(
		"broadcast.Publish",
		[]*Message{newMessage("topic", `"hello2"`)},
		&ok,
	))
	assert.True(t, ok)
	assert.Equal(t, `"hello2"`, readStr(<-client.Channel()))

	assert.NoError(t, client.Unsubscribe("topic"))

	assert.NoError(t, rcpClient.Call(
		"broadcast.Publish",
		[]*Message{newMessage("topic", `"hello3"`)},
		&ok,
	))
	assert.True(t, ok)

	assert.NoError(t, client.Subscribe("topic"))

	assert.NoError(t, rcpClient.Call(
		"broadcast.Publish",
		[]*Message{newMessage("topic", `"hello4"`)},
		&ok,
	))
	assert.True(t, ok)
	assert.Equal(t, `"hello4"`, readStr(<-client.Channel()))
}
