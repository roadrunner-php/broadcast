package broadcast

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/cobra"
	"github.com/spiral/broadcast"
	rr "github.com/spiral/roadrunner/cmd/rr/cmd"
	"github.com/spiral/roadrunner/service/metrics"
)

func init() {
	cobra.OnInitialize(func() {
		svc, _ := rr.Container.Get(metrics.ID)
		mtr, ok := svc.(*metrics.Service)
		if !ok || !mtr.Enabled() {
			return
		}

		ht, _ := rr.Container.Get(broadcast.ID)
		if bc, ok := ht.(*broadcast.Service); ok {
			collector := newCollector()

			// register metrics
			mtr.MustRegister(collector.connCounter)

			// collect events
			bc.AddListener(collector.listener)
		}
	})
}

// listener provide debug callback for system events. With colors!
type metricCollector struct {
	connCounter prometheus.Counter
}

func newCollector() *metricCollector {
	return &metricCollector{
		connCounter: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "rr_broadcast_conn_total",
				Help: "Total number of websocket connections to the broadcast service.",
			},
		),
	}
}

// listener listens to http events and generates nice looking output.
func (c *metricCollector) listener(event int, ctx interface{}) {
	switch event {
	case broadcast.EventConnect:
		c.connCounter.Inc()
	}
}
