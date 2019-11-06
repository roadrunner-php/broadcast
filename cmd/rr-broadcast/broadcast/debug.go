package broadcast

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spiral/broadcast"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"
	"github.com/spiral/roadrunner/cmd/util"
	"strings"
)

func init() {
	cobra.OnInitialize(func() {
		if rr.Debug {
			svc, _ := rr.Container.Get(broadcast.ID)
			if svc, ok := svc.(*broadcast.Service); ok {
				svc.AddListener((&debugger{logger: rr.Logger}).listener)
			}
		}
	})
}

// listener provide debug callback for system events. With colors!
type debugger struct{ logger *logrus.Logger }

// listener listens to http events and generates nice looking output.
func (s *debugger) listener(event int, ctx interface{}) {

	switch event {
	case broadcast.EventWebsocketConnect:
		conn := ctx.(*websocket.Conn)
		s.logger.Debug(util.Sprintf(
			"[broadcast] <green+hb>%s</reset> connected",
			conn.RemoteAddr(),
		))

	case broadcast.EventWebsocketDisconnect:
		conn := ctx.(*websocket.Conn)
		s.logger.Debug(util.Sprintf(
			"[broadcast] <yellow+hb>%s</reset> disconnected",
			conn.RemoteAddr(),
		))

	case broadcast.EventWebsocketJoin:
		e := ctx.(*broadcast.TopicEvent)
		s.logger.Debug(util.Sprintf(
			"[broadcast] <white+hb>%s</reset> join <cyan+hb>[%s]</reset>",
			e.Conn.RemoteAddr(),
			strings.Join(e.Topics, ", "),
		))

	case broadcast.EventWebsocketLeave:
		e := ctx.(*broadcast.TopicEvent)
		s.logger.Debug(util.Sprintf(
			"[broadcast] <white+hb>%s</reset> leave <magenta+hb>[%s]</reset>",
			e.Conn.RemoteAddr(),
			strings.Join(e.Topics, ", "),
		))

	case broadcast.EventWebsocketError:
		e := ctx.(*broadcast.ErrorEvent)
		if e.Conn != nil {
			s.logger.Debug(util.Sprintf(
				"[broadcast] <grey+hb>%s</reset> <yellow>%s</reset>",
				e.Conn.RemoteAddr(),
				e.Caused.Error(),
			))
		} else {
			s.logger.Error(util.Sprintf(
				"[broadcast]: <red+hb>%s</reset>",
				e.Caused.Error(),
			))
		}
	}
}
