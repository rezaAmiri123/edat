package inmem

import edatlog "github.com/rezaAmiri123/edat/log"

// ConsumerOption options for Consumer
type ConsumerOption func(consumer *Consumer)

// WithConsumerLogger sets the log.Logger for Consumer
func WithConsumerLogger(logger edatlog.Logger) ConsumerOption {
	return func(consumer *Consumer) {
		consumer.logger = logger
	}
}
