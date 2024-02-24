package msg

import (
	"context"
	"strings"

	"github.com/rezaAmiri123/edat/core"
	edatlog "github.com/rezaAmiri123/edat/log"
)

// CommandHandlerFunc function handlers for msg.Command
type CommandHandlerFunc func(context.Context, Command) ([]Reply, error)

// CommandDispatcher is a MessageReceiver for Commands
type CommandDispatcher struct {
	publisher ReplyMessagePublisher
	handlers  map[string]CommandHandlerFunc
	logger    edatlog.Logger
}

var _ MessageReceiver = (*CommandDispatcher)(nil)

// NewCommandDispatcher constructs a new CommandDispatcher
func NewCommandDispatcher(publisher ReplyMessagePublisher, options ...CommandDispatcherOption) *CommandDispatcher {
	c := &CommandDispatcher{
		publisher: publisher,
		handlers:  map[string]CommandHandlerFunc{},
		logger:    edatlog.DefaultLogger,
	}

	for _, option := range options {
		option(c)
	}

	c.logger.Trace("msg.CommandDispatcher constructed")

	return c
}

// Handle adds a new Command that will be handled by handler
func (d *CommandDispatcher) Handle(cmd core.Command, handler CommandHandlerFunc) *CommandDispatcher {
	d.logger.Trace("command handler added", edatlog.String("CommandName", cmd.CommandName()))
	d.handlers[cmd.CommandName()] = handler
	return d
}

// ReceiveMessage implements MessageReceiver.ReceiveMessage
func (d *CommandDispatcher) ReceiveMessage(ctx context.Context, message Message) error {
	commandName, err := message.Headers().GetRequired(MessageCommandName)
	if err != nil {
		d.logger.Error("error reading command name", edatlog.Error(err))
		return nil
	}

	logger := d.logger.Sub(
		edatlog.String("CommandName", commandName),
		edatlog.String("MessageID", message.ID()),
	)

	logger.Debug("received command message")

	// check first for a handler of the command; It is possible commands might be published into channels
	// that haven't been registered in our application
	handler, exists := d.handlers[commandName]
	if !exists {
		return nil
	}

	logger.Trace("command handler found")

	command, err := core.DeserializeCommand(commandName, message.Payload())
	if err != nil {
		logger.Error("error decoding command message payload", edatlog.Error(err))
		return nil
	}

	replyChannel, err := message.Headers().GetRequired(MessageCommandReplyChannel)
	if err != nil {
		logger.Error("error reading reply channel", edatlog.Error(err))
		return nil
	}

	correlationHeaders := d.correlationHeaders(message.Headers())

	cmdMsg := commandMessage{command, correlationHeaders}

	replies, err := handler(ctx, cmdMsg)
	if err != nil {
		logger.Error("command handler returned an error", edatlog.Error(err))
		rerr := d.sendReplies(ctx, replyChannel, []Reply{WithFailure()}, correlationHeaders)
		if rerr != nil {
			logger.Error("error sending replies", edatlog.Error(rerr))
			return nil
		}
		return nil
	}

	err = d.sendReplies(ctx, replyChannel, replies, correlationHeaders)
	if err != nil {
		logger.Error("error sending replies", edatlog.Error(err))
		return nil
	}

	return nil
}

func (d *CommandDispatcher) sendReplies(ctx context.Context, replyChannel string, replies []Reply, correlationHeaders Headers) error {
	for _, reply := range replies {
		err := d.publisher.PublishReply(ctx, reply.Reply(),
			WithHeaders(reply.Headers()),
			WithHeaders(correlationHeaders),
			WithDestinationChannel(replyChannel),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *CommandDispatcher) correlationHeaders(headers Headers) Headers {
	replyHeaders := make(map[string]string)
	for key, value := range headers {
		if key == MessageCommandName {
			continue
		}

		if strings.HasPrefix(key, MessageCommandPrefix) {
			replyHeader := MessageReplyPrefix + key[len(MessageCommandPrefix):]
			replyHeaders[replyHeader] = value
		}
	}

	return replyHeaders
}
