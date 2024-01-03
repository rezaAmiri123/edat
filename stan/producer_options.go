package edatstan

import "github.com/rezaAmiri123/edat/log"

type ProducerOption func(producer *Producer)

func WithProducerSerializer(serializer Serializer) ProducerOption {
	return func(producer *Producer) {
		producer.serializer = serializer
	}
}

func WithProducerLogger(logger edatlog.Logger) ProducerOption {
	return func(producer *Producer) {
		producer.logger = logger
	}
}
