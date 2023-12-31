package edatkafka

import (
	"context"
	"io"
	"time"

	edatlog "github.com/rezaAmiri123/edat/log"
	"github.com/rezaAmiri123/edat/msg"
	"github.com/segmentio/kafka-go"
)

// DefaultAckWait is a time.Duration representing the maximum amount of time for a consumer to finish
var DefaultAckWait = time.Second * 30

// Consumer implements msg.Consumer
type Consumer struct {
	brokers    []string
	groupID    string
	ackWait    time.Duration
	serializer Serializer
	logger     edatlog.Logger
	dialer     *kafka.Dialer
}

var _ msg.Consumer = (*Consumer)(nil)

// NewConsumer constructs a new instance of Consumer
func NewConsumer(brokers []string, groupID string, options ...ConsumerOption) *Consumer {
	c := &Consumer{
		brokers:    brokers,
		groupID:    groupID,
		ackWait:    DefaultAckWait,
		serializer: DefaultSerializer,
		logger:     edatlog.DefaultLogger,
		dialer:     kafka.DefaultDialer,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func (c *Consumer) Listen(ctx context.Context, channel string, consumer msg.ReceiveMessageFunc) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: c.brokers,
		GroupID: c.groupID,
		Topic:   channel,
		Dialer:  c.dialer,
	})

	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			c.logger.Error("error closing kafka-go reader", edatlog.Error(err))
		}
	}(reader)

	for {
		err := c.receiveMessage(ctx, reader, consumer)
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

func (c *Consumer) Close(ctx context.Context) error {
	c.logger.Trace("closing message source")
	return nil
}

func (c *Consumer) receiveMessage(ctx context.Context, reader *kafka.Reader, consumer msg.ReceiveMessageFunc) error {
	m, err := reader.FetchMessage(ctx)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	var message msg.Message
	message, err = c.serializer.Deserialize(m)
	if err != nil {
		return err
	}

	wCtx, cancel := context.WithTimeout(ctx, c.ackWait)
	defer cancel()

	errc := make(chan error)
	go func() {
		errc <- consumer(wCtx, message)
	}()

	select {
	case err = <-errc:
		if err == nil {
			if ackErr := reader.CommitMessages(ctx, m); ackErr != nil {
				c.logger.Error("error acknowledging message", edatlog.Error(err))
			}
		}
	case <-ctx.Done():
		c.logger.Trace("listener has closed; in-progress message processing is termnated")
	case <-wCtx.Done():
		c.logger.Warn("timed out waiting for message to finish")
	}
	return nil
}
