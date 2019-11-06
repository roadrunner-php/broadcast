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

func TestMemory_Broadcast(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6050",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	b, _ := c.Get(ID)
	br := b.(*Service)

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6050", Path: "/ws"}

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

	assert.NoError(t, br.getBroker().Broadcast(NewMessage("topic", "hello1"))) // must not be delivered

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`)))
	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(<-read))

	// double join is OK
	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`)))
	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(<-read))

	assert.NoError(t, br.getBroker().Broadcast(NewMessage("topic", "hello2")))
	assert.Equal(t, `{"topic":"topic","payload":"hello2"}`, readStr(<-read))

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"leave", "args":["topic"]}`)))
	assert.Equal(t, `{"topic":"@leave","payload":["topic"]}`, readStr(<-read))

	assert.NoError(t, br.getBroker().Broadcast(NewMessage("topic", "hello2")))

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`)))
	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(<-read))
}

func TestMemory_Broadcast_Error(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6051",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	b, _ := c.Get(ID)
	br := b.(*Service)

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6051", Path: "/ws"}

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

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`)))
	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(<-read))

	assert.NoError(t, br.getBroker().Broadcast(&Message{Topic: "topic", Payload: []byte("broken")}))
	assert.Equal(t, ``, readStr(<-read))

	assert.NoError(t, br.getBroker().Broadcast(NewMessage("topic", "hello2")))
	assert.Equal(t, `{"topic":"topic","payload":"hello2"}`, readStr(<-read))
}
