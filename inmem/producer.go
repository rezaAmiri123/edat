package inmem

import (
	"context"

	edatlog "github.com/rezaAmiri123/edat/log"
	"github.com/rezaAmiri123/edat/msg"
)

// Producer implements msg.Producer
type Producer struct{
	logger edatlog.Logger
}

var _ msg.Producer = (*Producer)(nil)

// NewProducer constructs a new Producer
func NewProducer(options ...ProducerOption)*Producer{
	p:= &Producer{
		logger: edatlog.DefaultLogger,
	}

	for _, option := range options{
		option(p)
	}

	return p
}

// Send implements msg.Producer.Send
func(p *Producer)Send(ctx context.Context, channel string, message msg.Message)error{
	if result, exists := channels.Load(channel);exists{
		destination := result.(chan msg.Message)

		destination <- message

		p.logger.Trace("message sent to inmem channel", edatlog.String("Channel", channel))
	}

	return nil
}

// Close implements msg.Producer.Close
func(p *Producer)Close(ctx context.Context)error{
	p.logger.Trace("closing message destination")

	return nil
}
