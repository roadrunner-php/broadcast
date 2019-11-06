package broadcast

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/spiral/roadrunner/service/env"
	rhttp "github.com/spiral/roadrunner/service/http"
	"github.com/spiral/roadrunner/service/rpc"
	"sync"
)

// ID defines public service name.
const ID = "broadcast"

// Service manages even broadcasting and websocket interface.
type Service struct {
	// service and broker configuration
	cfg *Config

	// wsPool manage websockets
	wsPool *wsPool

	// broadcast Messages
	mu     sync.Mutex
	broker Broker

	// event listeners
	lsns []func(event int, ctx interface{})
}

// AddListener attaches server event controller.
func (s *Service) AddListener(l func(event int, ctx interface{})) {
	s.lsns = append(s.lsns, l)
}

// AddCommand attached custom client command handler, for websocket only.
func (s *Service) AddCommand(name string, cmd CommandHandler) {
	if s.wsPool != nil {
		s.wsPool.commands[name] = cmd
	}
}

// Init service.
func (s *Service) Init(cfg *Config, r *rpc.Service, h *rhttp.Service, e env.Environment) (bool, error) {
	s.cfg = cfg

	if s.cfg.Path != "" && h != nil {
		s.wsPool = &wsPool{
			path:     s.cfg.Path,
			broker:   s.Broker,
			listener: s.throw,
			upgrade:  websocket.Upgrader{},
			connPool: &connPool{conn: make(map[*websocket.Conn]chan *Message)},
			commands: make(map[string]CommandHandler),
		}

		h.AddMiddleware(s.wsPool.middleware)

		if e != nil {
			// ensure that underlying kernel knows what route to handle
			e.SetEnv("RR_BROADCAST_URL", cfg.Path)
		}
	}

	if r != nil {
		if err := r.Register(ID, &rpcService{s: s}); err != nil {
			return false, err
		}
	}

	return true, nil
}

// Serve broadcast broker.
func (s *Service) Serve() (err error) {
	if s.wsPool != nil {
		defer s.wsPool.close()
	}

	s.mu.Lock()
	if s.cfg.Redis != nil {
		s.broker, err = redisBroker(s.cfg.Redis, func(err error) { s.throw(EventBrokerError, err) })
		if err != nil {
			return err
		}
	} else {
		s.broker = memoryBroker()
	}
	s.mu.Unlock()

	return s.broker.Serve()
}

// close broadcast broker.
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

// NewClient returns single connected client with ability to consume or produce into topic.
func (s *Service) NewClient(upstream chan *Message) *Client {
	return &Client{
		upstream:  upstream,
		broadcast: s,
	}
}

// Subscribe broker to one or multiple topics.
func (s *Service) Subscribe(upstream chan *Message, topics ...string) error {
	return s.Broker().Subscribe(upstream, topics...)
}

// Unsubscribe broker from one or multiple topics.
func (s *Service) Unsubscribe(upstream chan *Message, topics ...string) {
	s.Broker().Unsubscribe(upstream, topics...)
}

// Broadcast one or multiple Messages.
func (s *Service) Broadcast(msg ...*Message) error {
	broker := s.Broker()
	if broker == nil {
		return errors.New("no active broker")
	}

	return broker.Broadcast(msg...)
}

// throw handles service, server and pool events.
func (s *Service) throw(event int, ctx interface{}) {
	for _, l := range s.lsns {
		l(event, ctx)
	}
}
