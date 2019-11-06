package broadcast

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

func Test_Client_Consume_And_Produce(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6033",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	b, _ := c.Get(ID)
	br := b.(*Service)

	done := make(chan interface{})
	br.AddListener(func(event int, ctx interface{}) {
		if event == EventWebsocketConnect {
			close(done)
		}
	})

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6033", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer conn.Close()

	<-done

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

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`))
	assert.NoError(t, err)

	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(<-read))
}
