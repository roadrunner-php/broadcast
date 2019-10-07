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
			"address": ":6039",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	req, err := http.NewRequest("GET", "http://localhost:6039/", nil)
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
			"address": ":6039",
			"workers":{"command": "php tests/worker-deny.php", "pool.numWorkers": 1}
		}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	req, err := http.NewRequest("GET", "http://localhost:6039/", nil)
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
	c.Register(rrhttp.ID, &rrhttp.Service{})
	c.Register(ID, &Service{})

	assert.NoError(t, c.Init(&testCfg{
		http: `{
			"address": ":6039",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	req, err := http.NewRequest("GET", "http://localhost:6039/", nil)
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
			"address": ":6039",
			"workers":{"command": "php tests/worker-ok.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6039", Path: "/ws"}

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

	out := <-read
	assert.Equal(t, `{"topic":"@join","payload":["topic"]}`, readStr(out))
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
			"address": ":6039",
			"workers":{"command": "php tests/worker-deny.php", "pool.numWorkers": 1}
		}`,
		broadcast: `{"path":"/ws"}`,
	}))

	go func() { c.Serve() }()
	time.Sleep(time.Millisecond * 100)
	defer c.Stop()

	u := url.URL{Scheme: "ws", Host: "localhost:6039", Path: "/ws"}

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
