package msg

import "context"


type Consumer interface {
	Listen(ctx context.Context, channel string, consumer ReceiveMessageFunc) error
	Close(ctx context.Context) error
}
