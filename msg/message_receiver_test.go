package msg_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rezaAmiri123/edat/msg"
)

func TestReceiveMessageFunc_ReceiveMessage(t *testing.T) {
	type args struct {
		ctx     context.Context
		message msg.Message
	}

	tests := map[string]struct {
		args    args
		f       msg.ReceiveMessageFunc
		wantErr bool
	}{
		"Success": {
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{}`)),
			},
			f: func(ctx context.Context, m msg.Message) error {
				return nil
			},
			wantErr: false,
		},
		"ReceiverError": {
			args: args{
				ctx: context.Background(),
				message: msg.NewMessage([]byte(`{}`)),
			},
			f: func(ctx context.Context, m msg.Message) error {
				return fmt.Errorf("receiver-error")
			},
			wantErr: true,
		},
		
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.f.ReceiveMessage(tt.args.ctx, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReceiveMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err!= nil{
				t.Logf("Error(): %v", err)
			}
		})
	}
}
