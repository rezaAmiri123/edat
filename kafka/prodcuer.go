package edatkafka

import (
	"context"

	"github.com/rezaAmiri123/edat/log"
	"github.com/rezaAmiri123/edat/msg"
	"github.com/segmentio/kafka-go"
)

// Producer implements msg.Producer
type Producer struct {
	writer     *kafka.Writer
	serializer Serializer
	logger     edatlog.Logger
}

var _ msg.Producer = (*Producer)(nil)

// NewProducer constructs a new instance of Producer
func NewProducer(brokers []string, options ...ProducerOption) *Producer {
	p := &Producer{
		writer: &kafka.Writer{
			Addr: kafka.TCP(brokers...),
		},
		serializer: DefaultSerializer,
		logger:     nil,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

func (p *Producer) Send(ctx context.Context, channel string, message msg.Message) error {
	kafkaMsg, err := p.serializer.Serialize(message)
	if err != nil {
		p.logger.Error("failed to marshal message", edatlog.Error(err))
		return err
	}

	kafkaMsg.Topic = channel

	return p.writer.WriteMessages(ctx, kafkaMsg)

}

func (p *Producer) Close(ctx context.Context) error {
	p.logger.Trace("closing message destination")
	err := p.writer.Close()
	if err!= nil{
		p.logger.Error("error closing message destination", edatlog.Error(err))
	}

	return err
}
