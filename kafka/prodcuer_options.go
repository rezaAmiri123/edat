package edatkafka

import (
	"github.com/rezaAmiri123/edat/log"
	"github.com/segmentio/kafka-go"
)

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

func WithProducerTransport(transport *kafka.Transport) ProducerOption {
	return func(producer *Producer) {
		producer.writer.Transport = transport
	}
}
