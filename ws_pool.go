package broadcast

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/spiral/roadrunner/service/http/attributes"
	"net/http"
	"strings"
)

// provides an active instance of broker
type brokerProvider func() Broker

// CommandHandler handles custom commands.
type CommandHandler func(ctx *ConnContext, cmd []byte)

// manages broadcasting over web sockets
type wsPool struct {
	// path to serve websocket listener
	path string

	broker brokerProvider

	listener func(event int, ctx interface{})

	// manages wsPool upgrade protocol
	upgrade websocket.Upgrader

	// manages all open connections and broadcasting
	connPool *connPool

	// custom commands
	commands map[string]CommandHandler
}

// close the websocket pool.
func (ws *wsPool) close() {
	ws.connPool.close()
}

// middleware intercepts websocket connections.
func (ws *wsPool) middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != ws.path {
			f(w, r)
			return
		}

		broker := ws.broker()
		if broker == nil {
			f(w, r)
			return
		}

		if err := ws.assertServerAccess(f, r); err != nil {
			err.copy(w)
			return
		}

		conn, err := ws.upgrade.Upgrade(w, r, nil)
		if err != nil {
			ws.handleError(err, nil)
			return
		}

		ws.listener(EventWebsocketConnect, conn)
		upstream := ws.connPool.connect(conn, ws.handleError)

		ws.serveConn(conn, f, r, broker, upstream)
	}
}

// send and receive Messages over websocket
func (ws *wsPool) serveConn(
	conn *websocket.Conn,
	f http.HandlerFunc,
	r *http.Request,
	broker Broker,
	upstream chan *Message,
) {
	connContext := &ConnContext{
		Upstream: upstream,
		Conn:     conn,
		Topics:   make([]string, 0),
	}

	defer func() {
		ws.listener(EventWebsocketDisconnect, conn)
		broker.Unsubscribe(upstream, connContext.Topics...)
		ws.connPool.disconnect(conn)
	}()

	cmd := &Command{}
	for {
		if err := conn.ReadJSON(cmd); err != nil {
			ws.handleError(err, conn)
			return
		}

		switch cmd.Cmd {
		case "join":
			topics := make([]string, 0)
			if err := cmd.Unmarshal(&topics); err != nil {
				ws.handleError(err, conn)
				return
			}

			if err := ws.assertAccess(f, r, topics...); err != nil {
				ws.handleError(err, conn)
				return
			}

			if len(topics) == 0 {
				continue
			}

			if err := broker.Subscribe(upstream, topics...); err != nil {
				ws.handleError(err, conn)
				return
			}

			connContext.addTopic(topics...)
			upstream <- NewMessage("@join", topics)
			ws.listener(EventWebsocketJoin, &TopicEvent{Conn: conn, Topics: topics})
		case "leave":
			topics := make([]string, 0)
			if err := cmd.Unmarshal(&topics); err != nil {
				ws.handleError(err, conn)
				return
			}

			if len(topics) == 0 {
				continue
			}

			connContext.dropTopic(topics...)
			broker.Unsubscribe(upstream, topics...)
			upstream <- NewMessage("@leave", topics)
			ws.listener(EventWebsocketLeave, &TopicEvent{Conn: conn, Topics: topics})
		default:
			if handler, ok := ws.commands[cmd.Cmd]; ok {
				handler(connContext, cmd.Args)
			}
		}
	}
}

// handle connection error
func (ws *wsPool) handleError(err error, conn *websocket.Conn) {
	ws.listener(EventWebsocketError, &ErrorEvent{Conn: conn, Caused: err})
}

// assertServerAccess checks if user can join server and returns error and body if user can not. Must return nil in
// case of error
func (ws *wsPool) assertServerAccess(f http.HandlerFunc, r *http.Request) *responseWrapper {
	attributes.Set(r, "broadcast:joinServer", true)
	defer delete(attributes.All(r), "broadcast:joinServer")

	w := newResponseWrapper()
	f(w, r)

	if !w.IsOK() {
		return w
	}

	return nil
}

// assertAccess checks if user can access given channel, the application will receive all user headers and cookies.
// the decision to authorize user will be based on response code (200).
func (ws *wsPool) assertAccess(f http.HandlerFunc, r *http.Request, channels ...string) error {
	attributes.Set(r, "broadcast:joinTopics", strings.Join(channels, ","))
	defer delete(attributes.All(r), "broadcast:joinTopics")

	w := newResponseWrapper()
	f(w, r)

	if !w.IsOK() {
		return errors.New(string(w.Body()))
	}

	return nil
}
