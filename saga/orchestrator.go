package saga

import (
	"context"

	"github.com/google/uuid"
	"github.com/rezaAmiri123/edat/core"
	edatlog "github.com/rezaAmiri123/edat/log"
	"github.com/rezaAmiri123/edat/msg"
)

// Orchestrator orchestrates local and distributed processes
type Orchestrator struct {
	definition    Definition
	instanceStore InstanceStore
	publisher     msg.CommandMessagePublisher
	logger        edatlog.Logger
}

const sagaNotStated = -1

var _ msg.MessageReceiver = (*Orchestrator)(nil)

// NewOrchestrator constructs a new Orchestrator
func NewOrchestrator(definition Definition, store InstanceStore, publisher msg.CommandMessagePublisher, options ...OrchestratorOption) *Orchestrator {
	o := &Orchestrator{
		definition:    definition,
		instanceStore: store,
		publisher:     publisher,
		logger:        edatlog.DefaultLogger,
	}

	for _, option := range options{
		option(o)
	}

	o.logger.Trace("saga.Orchestrator constructed", edatlog.String("SagaName", definition.SagaName()))

	return o
}

// Start creates a new instance of the saga and begins execution
func(o *Orchestrator)Start(ctx context.Context, sagaData core.SagaData)(*Instance, error){
	instance := &Instance{
		sagaID: uuid.New().String(),
		sagaName: o.definition.SagaName(),
		sagaData: sagaData,
	}

	err := o.instanceStore.Save(ctx, instance)
	if err != nil{
		return nil,err
	}

	logger := o.logger.Sub(
		edatlog.String("SagaName", o.definition.SagaName()),
		edatlog.String("SagaID", instance.sagaID)
	)

	logger.Trace("executing saga starting hook")
	o.definition.OnHook(SagaStarting, instance)

	// results := o.executeNextStep(ctx, stepContext{step: sagaNotStated}, sagaData)
}
