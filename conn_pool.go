package broadcast

import (
	"github.com/gorilla/websocket"
	"sync"
)

// manages all ws connections and their upstream channels
type connPool struct {
	mu   sync.Mutex
	conn map[*websocket.Conn]chan *Message
}

// connect socket and issue new upstream channel
func (cp *connPool) connect(conn *websocket.Conn, errHandler func(err error, conn *websocket.Conn)) chan *Message {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	upstream := make(chan *Message, 1)
	cp.conn[conn] = upstream

	go func() {
		for msg := range upstream {
			if err := conn.WriteJSON(msg); err != nil {
				errHandler(err, conn)
			}
		}
	}()

	return upstream
}

// disconnect websocket
func (cp *connPool) disconnect(conn *websocket.Conn) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	upstream, ok := cp.conn[conn]
	if !ok {
		return
	} else {
		delete(cp.conn, conn)
	}

	close(upstream)
	conn.Close()
}

// close pool and all underlying connections
func (cp *connPool) close() {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for conn, upstream := range cp.conn {
		delete(cp.conn, conn)
		close(upstream)
		conn.Close()
	}
}
