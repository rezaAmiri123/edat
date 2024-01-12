package saga

import (
	"github.com/rezaAmiri123/edat/core"
	"github.com/rezaAmiri123/edat/msg"
)

// Command is a core.Command with message header information
type Command interface {
	SagaID() string
	SagaName() string
	Command() core.Command
	Headers() msg.Headers
}

type commandMessage struct {
	sagaID   string
	sagaName string
	command  core.Command
	headers  msg.Headers
}

func (m commandMessage) SagaID() string {
	return m.sagaID
}

func (m commandMessage) SagaName() string {
	return m.sagaName
}

func (m commandMessage) Command() core.Command {
	return m.command
}

func (m commandMessage) Headers() msg.Headers {
	return m.headers
}
