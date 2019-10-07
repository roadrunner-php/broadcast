package broadcast

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/spiral/roadrunner/service/env"
	rhttp "github.com/spiral/roadrunner/service/http"
	"github.com/spiral/roadrunner/service/http/attributes"
	"github.com/spiral/roadrunner/service/rpc"
	"net/http"
	"strings"
	"sync"
)

// ID defines public service name.
const ID = "broadcast"

// Broker defines the ability to operate as message passing broker.
type Broker interface {
	// Serve serves broker.
	Serve() error

	// Stop the consumption and disconnect broker.
	Stop()

	// Subscribe broker to one or multiple topics.
	Subscribe(upstream chan *Message, topics ...string) error

	// Unsubscribe broker from one or multiple topics.
	Unsubscribe(upstream chan *Message, topics ...string)

	// Broadcast one or multiple messages.
	Broadcast(messages ...*Message) error
}

// CommandHandler handles custom commands.
type CommandHandler func(ctx *ConnContext, cmd []byte)

// Service manages even broadcasting over websockets.
type Service struct {
	// service and broker configuration
	cfg *Config

	// manages ws upgrade protocol
	upgrade websocket.Upgrader

	// broadcast messages
	mu     sync.Mutex
	broker Broker

	// event listeners
	lsns []func(event int, ctx interface{})

	// custom commands
	commands map[string]CommandHandler

	// manages all open connections and broadcasting
	connPool *connPool
}

// AddListener attaches server event controller.
func (s *Service) AddListener(l func(event int, ctx interface{})) {
	s.lsns = append(s.lsns, l)
}

// AddCommand attached custom client command handler.
func (s *Service) AddCommand(name string, cmd CommandHandler) {
	s.commands[name] = cmd
}

// Init service.
func (s *Service) Init(cfg *Config, r *rpc.Service, h *rhttp.Service, e env.Environment) (bool, error) {
	if cfg.Path == "" || h == nil {
		return false, nil
	}

	s.cfg = cfg
	s.upgrade = websocket.Upgrader{}
	s.connPool = &connPool{conn: make(map[*websocket.Conn]chan *Message)}
	s.commands = make(map[string]CommandHandler)

	h.AddMiddleware(s.middleware)

	if e != nil {
		// ensure that underlying kernel knows what route to handle
		e.SetEnv("RR_BROADCAST_URL", cfg.Path)
	}

	return true, nil
}

// Serve broadcast broker.
func (s *Service) Serve() (err error) {
	defer s.connPool.close()

	s.mu.Lock()
	if s.cfg.Redis != nil {
		s.broker, err = redisBroker(s.cfg.Redis, s.handleError)
		if err != nil {
			return err
		}
	} else {
		s.broker = memoryBroker()
	}
	s.mu.Unlock()

	return s.broker.Serve()
}

// Stop broadcast broker.
func (s *Service) Stop() {
	broker := s.Broker()
	if broker != nil {
		broker.Stop()
	}
}

// Broker returns associated broker.
func (s *Service) Broker() Broker {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.broker
}

// middleware intercepts websocket connections.
func (s *Service) middleware(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != s.cfg.Path {
			f(w, r)
			return
		}

		broker := s.Broker()
		if broker == nil {
			f(w, r)
			return
		}

		conn, err := s.upgrade.Upgrade(w, r, nil)
		if err != nil {
			s.handleError(err, nil)
			return
		}

		s.throw(EventConnect, conn)
		upstream := s.connPool.connect(conn, s.handleError)

		s.serveConn(conn, f, r, broker, upstream)
	}
}

// send and receive messages over websocket
func (s *Service) serveConn(
	conn *websocket.Conn,
	f http.HandlerFunc,
	r *http.Request,
	broker Broker,
	upstream chan *Message,
) {
	defer s.connPool.disconnect(conn)
	defer broker.Unsubscribe(upstream)
	defer s.throw(EventDisconnect, conn)

	connContext := &ConnContext{
		Upstream: upstream,
		Conn:     conn,
		Topics:   make([]string, 0),
	}

	cmd := &Command{}
	for {
		if err := conn.ReadJSON(cmd); err != nil {
			s.handleError(err, conn)
			return
		}

		switch cmd.Cmd {
		case "join":
			topics := make([]string, 0)
			if err := cmd.Unmarshal(&topics); err != nil {
				s.handleError(err, conn)
				return
			}

			if err := s.assertAccess(f, r, topics...); err != nil {
				s.handleError(err, conn)
				return
			}

			if len(topics) == 0 {
				continue
			}

			if err := broker.Subscribe(upstream, topics...); err != nil {
				s.handleError(err, conn)
				return
			}

			connContext.addTopic(topics...)
			upstream <- NewMessage("@join", topics)
			s.throw(EventJoin, &TopicEvent{Conn: conn, Topics: topics})
		case "leave":
			topics := make([]string, 0)
			if err := cmd.Unmarshal(&topics); err != nil {
				s.handleError(err, conn)
				return
			}

			if len(topics) == 0 {
				continue
			}

			connContext.dropTopic(topics...)
			broker.Unsubscribe(upstream, topics...)
			upstream <- NewMessage("@leave", topics)
			s.throw(EventLeave, &TopicEvent{Conn: conn, Topics: topics})
		default:
			if handler, ok := s.commands[cmd.Cmd]; ok {
				handler(connContext, cmd.Args)
			}
		}
	}
}

// handle connection error
func (s *Service) handleError(err error, conn *websocket.Conn) {
	s.throw(EventError, &ErrorEvent{Conn: conn, Caused: err})
}

// throw handles service, server and pool events.
func (s *Service) throw(event int, ctx interface{}) {
	for _, l := range s.lsns {
		l(event, ctx)
	}
}

// assertAccess checks if user can access given channel, the application will receive all user headers and cookies.
// the decision to authorize user will be based on response code (200).
func (s *Service) assertAccess(f http.HandlerFunc, r *http.Request, channels ...string) error {
	w := newResponseWrapper()
	if err := attributes.Set(r, "joinTopics", strings.Join(channels, ",")); err != nil {
		return err
	}

	f(w, r)

	if !w.IsOK() {
		return errors.New(string(w.Body()))
	}

	return nil
}
