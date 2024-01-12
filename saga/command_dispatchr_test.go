package saga_test

import (
	"context"
	"testing"

	"github.com/rezaAmiri123/edat/core"
	"github.com/rezaAmiri123/edat/core/coretest"
	edatlog "github.com/rezaAmiri123/edat/log"
	"github.com/rezaAmiri123/edat/log/logmocks"
	"github.com/rezaAmiri123/edat/log/logtest"
	"github.com/rezaAmiri123/edat/msg"
	"github.com/rezaAmiri123/edat/msg/msgmocks"
	"github.com/rezaAmiri123/edat/msg/msgtest"
	"github.com/rezaAmiri123/edat/saga"
	"github.com/stretchr/testify/mock"
)

type (
	sagaCommand        struct{ Value string }
	unregistredCommand struct{ Value string }
)

func (sagaCommand) CommandName() string        { return "saga_test.sagaCommand" }
func (unregistredCommand) CommandName() string { return "saga_test.unregistredCommand" }

func TestCommandDispatcher_ReceiveMessage(t *testing.T) {
	type handler struct {
		cmd core.Command
		fn  saga.CommandHandlerFunc
	}
	type fields struct {
		publisher msg.ReplyMessagePublisher
		handlers  []handler
		logger    edatlog.Logger
	}
	type args struct{
		ctx context.Context
		message msg.Message
	}

	core.RegisterDefaultMarshaller(coretest.NewTestMarshaller())
	core.RegisterCommands(sagaCommand{})

	tests := map[string]struct{
		fields fields
		args args
		wantErr bool
	}{
		"Success": {
			fields: fields{
				publisher: msgtest.MockReplyMessagePublisher(func(m *msgmocks.ReplyMessagePublisher) {
					m.On("PublishReply", mock.Anything, mock.AnythingOfType("msg.Success"), mock.Anything, mock.Anything, mock.Anything).Return(nil)
				}),
				handlers: []handler{{
					cmd: sagaCommand{},
					fn: func(ctx context.Context, c saga.Command) ([]msg.Reply, error) {
						return []msg.Reply{msg.WithSuccess()},nil
					},
				}},
				logger: logtest.MockLogger(func(m *logmocks.Logger) {
					m.On("Sub", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(m)
					m.On("Trace", mock.AnythingOfType("string"), mock.Anything)
					m.On("Debug", mock.AnythingOfType("string"), mock.Anything)
				}),
			},
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{"Value":""}`), msg.WithHeaders(msg.Headers{
					msg.MessageCommandName:         sagaCommand{}.CommandName(),
					saga.MessageCommandSagaID:      "test-id",
					saga.MessageCommandSagaName:    "test",
					msg.MessageCommandReplyChannel: "reply-channel",
				})),
			},
			 wantErr: false,
		},
	}

	for name, tt := range tests{
		t.Run(name, func(t *testing.T) {
			d := saga.NewCommandDispatcher(tt.fields.publisher, saga.WithCommandDispatcherLogger(tt.fields.logger))
			for _, handler := range tt.fields.handlers{
				d.Handle(handler.cmd, handler.fn)
			}
			err := d.ReceiveMessage(tt.args.ctx,tt.args.message)
			if err != nil != tt.wantErr{
				t.Errorf("ReceiveMessage() error = %v, watErr %v", err, tt.wantErr)
			}
			mock.AssertExpectationsForObjects(t, tt.fields.publisher, tt.fields.logger)
		})
	}
}
