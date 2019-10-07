package main

import (
	"github.com/gorilla/websocket"
	"net/url"
	"os"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}

	retry := 10

	var conn *websocket.Conn
	var err error
	for retry > 0 {
		conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			retry--
		} else {
			break
		}
	}

	if conn == nil {
		panic(err)
	}

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

	f, _ := os.Create("log.txt")
	defer f.Close()

	conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"join", "args":["topic"]}`))
	f.Write((<-read).([]byte)) // joined

	// read exactly 2 messages
	f.Write((<-read).([]byte))
	f.Write((<-read).([]byte))

	conn.WriteMessage(websocket.TextMessage, []byte(`{"cmd":"leave", "args":["topic"]}`))
	f.Write((<-read).([]byte)) // left
}
