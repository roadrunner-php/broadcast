package main

import (
	"fmt"
	"github.com/spiral/broadcast"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"
	"github.com/spiral/roadrunner/service/rpc"
	"os"
)

type logService struct {
	broadcast *broadcast.Service
	stop      chan interface{}
}

func (l *logService) Init(service *broadcast.Service) (bool, error) {
	l.broadcast = service

	return true, nil
}

func (l *logService) Serve() error {
	l.stop = make(chan interface{})

	client := l.broadcast.NewClient()
	if err := client.SubscribePattern("tests/*"); err != nil {
		return err
	}
	defer client.Close()

	logFile, _ := os.Create("log.txt")
	defer logFile.Close()

	go func() {
		for msg := range client.Channel() {
			logFile.Write([]byte(fmt.Sprintf(
				"%s: %s\n",
				msg.Topic,
				string(msg.Payload),
			)))

			logFile.Sync()
		}
	}()

	<-l.stop
	return nil
}

func (l *logService) Stop() {
	close(l.stop)
}

func main() {
	rr.Container.Register(rpc.ID, &rpc.Service{})
	rr.Container.Register(broadcast.ID, &broadcast.Service{})
	rr.Container.Register("log", &logService{})

	rr.Execute()
}
