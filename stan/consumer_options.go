package edatstan

import (
	"time"

	"github.com/nats-io/stan.go"
	"github.com/rezaAmiri123/edat/log"
)

type ConsumerOption func(cosumer *Consumer)

func WithConsumerActWait(actWait time.Duration) ConsumerOption {
	return func(cosumer *Consumer) {
		cosumer.ackWait = actWait
		cosumer.subOptions = append(cosumer.subOptions, stan.AckWait(actWait))
	}
}

func WithConsumerSubscriptionOptions(optios ...stan.SubscriptionOption) ConsumerOption {
	return func(cosumer *Consumer) {
		cosumer.subOptions = append(cosumer.subOptions, optios...)
	}
}

func WithConsumerSerializer(serializer Serializer) ConsumerOption {
	return func(cosumer *Consumer) {
		cosumer.serializer = serializer
	}
}

func WithConsumerLogger(logger log.Logger)ConsumerOption{
	return func(cosumer *Consumer) {
		cosumer.logger = logger
	}
}
