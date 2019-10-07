package broadcast

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spiral/roadrunner/service"
	"github.com/spiral/roadrunner/service/env"
	rrhttp "github.com/spiral/roadrunner/service/http"
	"github.com/spiral/roadrunner/service/rpc"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

type testCfg struct {
	http      string
	rpc       string
	broadcast string
	target    string
}

func (cfg *testCfg) Get(name string) service.Config {
	if name == rrhttp.ID {
		return &testCfg{target: cfg.http}
	}

	if name == ID {
		return &testCfg{target: cfg.broadcast}
	}

	if name == rpc.ID {
		return &testCfg{target: cfg.rpc}
	}

	return nil
}
func (cfg *testCfg) Unmarshal(out interface{}) error {
	return json.Unmarshal([]byte(cfg.target), out)
}

func readStr(m interface{}) string {
	return strings.TrimRight(string(m.([]byte)), "\n")
}

func Test_HttpService_Echo(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(rrhttp.ID, &rrhttp.Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6041",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	req, err := http.NewRequest("GET", "http://localhost:6041/", nil)
	assert.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer r.Body.Close()

	b, _ := ioutil.ReadAll(r.Body)

	assert.NoError(t, err)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, []byte(""), b)
}

func Test_HttpService_Echo400(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(rrhttp.ID, &rrhttp.Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6040",
			"workers":{"command": "php tests/worker-deny.php", "pool.numWorkers": 1}
		}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	req, err := http.NewRequest("GET", "http://localhost:6040/", nil)
	assert.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer r.Body.Close()

	assert.NoError(t, err)
	assert.Equal(t, 401, r.StatusCode)
}

func Test_Service_EnvPath(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rpc.ID, &rpc.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6029",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		rpc:       `{"listen":"tcp://127.0.0.1:6002"}`,
		broadcast: `{"path":"/ws"}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	req, err := http.NewRequest("GET", "http://localhost:6029/", nil)
	assert.NoError(t, err)

	r, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer r.Body.Close()

	b, _ := ioutil.ReadAll(r.Body)

	assert.NoError(t, err)
	assert.Equal(t, 200, r.StatusCode)
	assert.Equal(t, []byte("/ws"), b)
}

func Test_Service_JoinTopic(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6038",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6038", Path: "/ws"}

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

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`))
	assert.NoError(t, err)

	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(<-read))
}

func Test_Service_DenyJoin(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6037",
			"workers":{"command": "php tests/worker-deny.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6037", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer conn.Close()

	read := make(chan interface{})

	go func() {
		defer close(read)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				read <- err
				continue
			}
			read <- message
		}
	}()

	err = conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`))
	assert.NoError(t, err)

	out := <-read
	assert.Error(t, out.(error))
}

func Test_Service_EmptyTopics(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6036",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6036", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer conn.Close()

	read := make(chan interface{})

	go func() {
		defer close(read)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				read <- err
				continue
			}
			read <- message
		}
	}()

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":[]}`)))

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["a"]}`)))
	assert.Equal(t, `{"topic":"@join","payload":["a"]}`, readStr(<-read))

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"leave", "args":[]}`)))

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"leave", "args":["a"]}`)))
	assert.Equal(t, `{"topic":"@leave","payload":["a"]}`, readStr(<-read))
}

func Test_Service_BadTopics(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6035",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6035", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer conn.Close()

	read := make(chan interface{})

	go func() {
		defer close(read)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				read <- err
				continue
			}
			read <- message
		}
	}()

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":"hello"}`)))
	assert.Error(t, (<-read).(error))
}

func Test_Service_BadTopicsLeave(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6034",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6034", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer conn.Close()

	read := make(chan interface{})

	go func() {
		defer close(read)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				read <- err
				continue
			}
			read <- message
		}
	}()

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"leave", "args":"hello"}`)))
	assert.Error(t, (<-read).(error))
}

func Test_Service_Events(t *testing.T) {
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
		if event == EventConnect {
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

func Test_Service_Command(t *testing.T) {
	logger, _ := test.NewNullLogger()
	logger.SetLevel(logrus.DebugLevel)

	c := service.NewContainer(logger)
	c.Register(env.ID, &env.Service{})
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6032",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	b, _ := c.Get(ID)
	br := b.(*Service)

	br.AddCommand("send", func(ctx *ConnContext, cmd []byte) {
		assert.Equal(t, []byte(`"send-message"`), cmd)
		assert.Equal(t, []string{"topic"}, ctx.Topics)
		assert.NoError(t, br.Broker().Broadcast(NewMessage("topic", "custom-message")))
	})

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6032", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	defer conn.Close()

	read := make(chan interface{})

	go func() {
		defer close(read)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				read <- err
				continue
			}
			read <- message
		}
	}()

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`)))
	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(<-read))

	assert.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"send", "args":"send-message"}`)))
	assert.Equal(t, `{"topic":"topic","payload":"custom-message"}`, readStr(<-read))
}
