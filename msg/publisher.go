package msg

import (
	"context"
	"sync"
	"time"

	"github.com/rezaAmiri123/edat/core"
	"github.com/rezaAmiri123/edat/log"
)

// CommandMessagePublisher interface
type CommandMessagePublisher interface {
	PublisheCommand(ctx context.Context, replyChannel string, command core.Command, options ...MessageOption) error
}

// EntityEventMessagePublisher interface
type EntityEventMessagePublisher interface {
	PublishEntityEvents(ctx context.Context, entity core.Entity, options ...MessageOption) error
}

// EventMessagePublisher interface
type EventMessagePublisher interface {
	PublishEvent(ctx context.Context, event core.Event, options ...MessageOption) error
}

// ReplyMessagePublisher interface
type ReplyMessagePublisher interface {
	PublishReply(ctx context.Context, reply core.Reply, options ...MessageOption) error
}

// MessagePublisher interface
type MessagePublisher interface {
	Publish(ctx context.Context, message Message) error
}

// Publisher send domain events, commands, and replies to the publisher
type Publisher struct {
	producer Producer
	logger   edatlog.Logger
	close    sync.Once
}

var _ interface {
	CommandMessagePublisher
	EntityEventMessagePublisher
	EventMessagePublisher
	ReplyMessagePublisher
	MessagePublisher
} = (*Publisher)(nil)

// NewPublisher constructs a new Publisher
func NewPublisher(producer Producer, options ...PublisherOption)*Publisher{
	p := &Publisher{
		producer: producer,
		logger: edatlog.DefaultLogger,
	}

	for _, option := range options{
		option(p)
	}

	p.logger.Trace("msg.Publisher constructed")

	return p
}

// PublishCommand serializes a command into a message with command specific headers and publishes it to a producer
func(p *Publisher)PublisheCommand(ctx context.Context, replyChannel string, command core.Command, options ...MessageOption) error{
	msgOptions := []MessageOption{
		WithHeaders(Headers{
			MessageCommandName: command.CommandName(),
			MessageCommandReplyChannel: replyChannel,
		}),
	}

	if v, ok := command.(interface{ DestinationChannel() string });ok{
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	logger := p.logger.Sub(
		edatlog.String("CommandName", command.CommandName()),
	)

	logger.Trace("publishing command")

	payload, err := core.SerializeCommand(command)
	if err!= nil{
		logger.Error("error serializing command payload", edatlog.Error(err))
		return err
	}

	message := NewMessage(payload, msgOptions...)
	err = p.Publish(ctx, message)
	if err!= nil{
		logger.Error("error publishing command", edatlog.Error(err))
	}

	return err
}

// PublishEntityEvents serializes entity events into messages with entity specific headers and publishes it to a producer
func(p *Publisher)PublishEntityEvents(ctx context.Context, entity core.Entity, options ...MessageOption) error{
	msgOptions := []MessageOption{
		WithHeaders(Headers{
			MessageEventEntityID:entity.ID(),
			MessageEventEntityName: entity.EntityName(),
			MessageChannel: entity.EntityName(),// allow entity name and channel to overlap
		}),
	}

	if v,ok := entity.(interface{ DestinationChannel() string });ok{
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	for _, event := range entity.Events(){
		logger := p.logger.Sub(
			edatlog.String("EntityID", entity.EntityName()),
			edatlog.String("EtityName", entity.EntityName()),
		)

		err := p.PublishEvent(ctx, event, msgOptions...)
		if err!= nil{
			logger.Error("error publishing entity event", edatlog.Error(err))
			return err
		}
	}
	return nil
}

// PublishEvent serializes an event into a message with event specific headers and publishes it to a producer
func(p *Publisher)PublishEvent(ctx context.Context, event core.Event, options ...MessageOption) error{
	msgOptions := []MessageOption{
		WithHeaders(map[string]string{
			MessageEventName: event.EventName(),
		}),
	}

	if v, ok := event.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	logger := p.logger.Sub(
		edatlog.String("EventName", event.EventName()),
	)

	logger.Trace("publishing event")

	payload, err := core.SerializeEvent(event)
	if err != nil {
		logger.Error("error serializing event payload", edatlog.Error(err))
		return err
	}

	message := NewMessage(payload, msgOptions...)

	err = p.Publish(ctx, message)
	if err != nil {
		logger.Error("error publishing event", edatlog.Error(err))
	}

	return err
}

// PublishReply serializes a reply into a message with reply specific headers and publishes it to a producer
func(p *Publisher)PublishReply(ctx context.Context, reply core.Reply, options ...MessageOption) error{
	msgOptions := []MessageOption{
		WithHeaders(map[string]string{
			MessageReplyName: reply.ReplyName(),
		}),
	}

	if v, ok := reply.(interface{ DestinationChannel() string }); ok {
		msgOptions = append(msgOptions, WithDestinationChannel(v.DestinationChannel()))
	}

	msgOptions = append(msgOptions, options...)

	logger := p.logger.Sub(
		edatlog.String("ReplyName", reply.ReplyName()),
	)

	logger.Trace("publishing reply")

	payload, err := core.SerializeReply(reply)
	if err != nil {
		logger.Error("error serializing reply payload", edatlog.Error(err))
		return err
	}

	message := NewMessage(payload, msgOptions...)

	err = p.Publish(ctx, message)
	if err != nil {
		logger.Error("error publishing reply", edatlog.Error(err))
	}

	return err
}

// Publish sends a message off to a producer
func(p *Publisher)Publish(ctx context.Context, message Message) error{
	var err error
	var channel string

	channel, err = message.Headers().GetRequired(MessageChannel)
	if err!= nil{
		return err
	}

	message.Headers()[MessageDate] = time.Now().Format(time.RFC3339)

	// Published messages are request boundaries
	if	id, exists := message.Headers()[MessageCorrelationID];!exists || id ==""{
		message.Headers()[MessageCorrelationID] = core.GetCorrelationID(ctx)
	}

	if id, exists := message.Headers()[MessageCausationID]; !exists || id == "" {
		message.Headers()[MessageCausationID] = core.GetRequestID(ctx)
	}

	logger := p.logger.Sub(
		edatlog.String("MessageID", message.ID()),
		edatlog.String("CorrelationID", message.Headers()[MessageCorrelationID]),
		edatlog.String("CausationID", message.Headers()[MessageCausationID]),
		edatlog.String("Destination", channel),
		edatlog.Int("PayloadSize", len(message.Payload())),
	)

	logger.Trace("publishing message")

	err = p.producer.Send(ctx, channel, message)
	if err!= nil{
		logger.Error("error publishing message", edatlog.Error(err))
		return err
	}

	return nil
}

// Stop stops the publisher and underlying producer
func(p *Publisher)Stop(ctx context.Context)(err error){
	defer p.logger.Trace("publisher stopped")
	p.close.Do(func() {
		err = p.producer.Close(ctx)
	})
	return
}
