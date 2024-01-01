package edatkafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type ConsumerOption func(consumer *Consumer)

func WithConsumerAckWait(ackWait time.Duration) ConsumerOption{
	return func(consumer *Consumer) {
		consumer.ackWait = ackWait
	}
}

func WithConsumerSerializer(serializer Serializer)ConsumerOption{
	return func(consumer *Consumer) {
		consumer.serializer = serializer
	}
}

func WithConsumerDialer(dialer *kafka.Dialer)ConsumerOption{
	return func(consumer *Consumer) {
		consumer.dialer = dialer
	}
}