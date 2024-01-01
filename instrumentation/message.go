package instrumentation

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rezaAmiri123/edat/msg"
)

func MessageWebInstrumentation() func(next msg.MessageReceiver) msg.MessageReceiver {
	resposeTime := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "message_response_time",
		Help:    "Message response time in miliseconds",
		Buckets: []float64{300, 600, 900, 1_500, 5_000, 10_000, 20_000},
	})

	return func(next msg.MessageReceiver) msg.MessageReceiver {
		return msg.ReceiveMessageFunc(func(ctx context.Context, message msg.Message) error {
			start := time.Now()
			err := next.ReceiveMessage(ctx, message)
			resposeTime.Observe(float64(time.Since(start).Microseconds()))

			return err
		})
	}
}
