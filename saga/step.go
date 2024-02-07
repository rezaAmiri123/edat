package saga

import (
	"context"

	"github.com/rezaAmiri123/edat/core"
)

// Step interface for local, remote, ...other saga steps
type Step interface {
	hasInvocableAction(ctx context.Context, sagaData core.SagaData, compensating bool) bool
	getReplyHandler(replyName string, compensationg bool) func(ctx context.Context, data core.SagaData, reply core.Reply) error
	exexute(ctx context.Context, sagaData core.SagaData, compensation bool) func(results *stepResults)
}
