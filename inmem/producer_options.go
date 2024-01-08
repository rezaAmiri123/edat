package inmem

import edatlog "github.com/rezaAmiri123/edat/log"

// ProducerOption options for Producer
type ProducerOption func(producer *Producer)

// WithProducerLogger sets the edatlog.Logger for Producer
func WithProducerLogger(logger edatlog.Logger)ProducerOption{
	return func(producer *Producer) {
		producer.logger = logger
	}
}