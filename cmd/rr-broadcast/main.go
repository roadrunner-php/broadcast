package main

import (
	"github.com/spiral/broadcast"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"
	"github.com/spiral/roadrunner/service/rpc"
)

func main() {
	rr.Container.Register(rpc.ID, &rpc.Service{})
	rr.Container.Register(broadcast.ID, &broadcast.Service{})

	rr.Execute()
}
