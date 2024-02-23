package saga

import (
	"github.com/rezaAmiri123/edat/core"
	"github.com/rezaAmiri123/edat/msg"
)

type stepResults struct {
	commands           []msg.DomainCommand
	updatedSagaData    core.SagaData
	updatedStepContext stepContext
	local              bool
	failure            error
}
