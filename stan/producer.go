package edatstan

import (
	"context"
	"fmt"

	"github.com/nats-io/stan.go"
	"github.com/rezaAmiri123/edat/log"
	"github.com/rezaAmiri123/edat/msg"
)

type Producer struct{
	conn stan.Conn
	serializer Serializer
	logger edatlog.Logger
}

var _ msg.Producer = (*Producer)(nil)

func NewProducer(conn stan.Conn, options ...ProducerOption)*Producer{
	p := &Producer{
		conn: conn,
		serializer: DefaultSerializer,
		logger: edatlog.DefaultLogger,
	}

	for _, option := range options{
		option(p)
	}

	return p
}

func(p *Producer)Send(ctx context.Context, channel string, message msg.Message)error{
	logger := p.logger.Sub(
		edatlog.String("channel", channel),
	)

	data, err := p.serializer.Serialize(message)
	if err!= nil{
		logger.Error("failed to marshal message", edatlog.Error(err))
		return fmt.Errorf("message could not be marshalled")
	}

	if err = p.conn.Publish(channel, data);err!= nil{
		return err
	}

	return nil
}

func(p *Producer)Close(ctx context.Context)error{
	p.logger.Trace("closing message destination")
	err := p.conn.Close()
	if err!= nil{
		p.logger.Error("error closing message destination", edatlog.Error(err))
	}

	return err
}
