package broadcast

import (
	"errors"
	"github.com/spiral/roadrunner/service/rpc"
	"sync"
)

// ID defines public service name.
const ID = "broadcast"

// Service manages even broadcasting and websocket interface.
type Service struct {
	// service and broker configuration
	cfg *Config

	// broker
	mu     sync.Mutex
	broker Broker
}

// Init service.
func (s *Service) Init(cfg *Config, rpc *rpc.Service) (bool, error) {
	s.cfg = cfg

	if rpc != nil {
		if err := rpc.Register(ID, &rpcService{s: s}); err != nil {
			return false, err
		}
	}

	return true, nil
}

// Serve broadcast broker.
func (s *Service) Serve() (err error) {
	s.mu.Lock()
	if s.cfg.Redis != nil {
		if s.broker, err = redisBroker(s.cfg.Redis); err != nil {
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

// NewClient returns single connected client with ability to consume or produce into associated topic(s).
func (s *Service) NewClient() *Client {
	return &Client{upstream: make(chan *Message), broker: s.Broker()}
}

// Publish one or multiple Channel.
func (s *Service) Publish(msg ...*Message) error {
	broker := s.Broker()
	if broker == nil {
		return errors.New("no active broker")
	}

	return broker.Publish(msg...)
}
