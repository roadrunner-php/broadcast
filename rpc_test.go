package broadcast

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spiral/roadrunner/service"
	"github.com/spiral/roadrunner/service/env"
	rrhttp "github.com/spiral/roadrunner/service/http"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func TestRCP_Broadcast(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6056",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	b, _ := c.Get(ID)
	br := b.(*Service)

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6056", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer conn.Close()

	read := make(chan interface{})

	go func() {
		defer close(read)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			read <- message
		}
	}()

	r := &rpcService{br}

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`)))
	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(<-read))

	ok := false
	assert.NoError(t, r.Broadcast([]*Message{NewMessage("topic", "hello2")}, &ok))
	assert.True(t, ok)

	assert.Equal(t, `{"topic":"topic","payload":"hello2"}`, readStr(<-read))
}

func TestRCP_BroadcastAsync(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6055",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	b, _ := c.Get(ID)
	br := b.(*Service)

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6055", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer conn.Close()

	read := make(chan interface{})

	go func() {
		defer close(read)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			read <- message
		}
	}()

	r := &rpcService{br}

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`)))
	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(<-read))

	ok := false
	assert.NoError(t, r.BroadcastAsync([]*Message{NewMessage("topic", "hello2")}, &ok))
	assert.True(t, ok)

	assert.Equal(t, `{"topic":"topic","payload":"hello2"}`, readStr(<-read))
}
