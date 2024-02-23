package saga

import edatlog "github.com/rezaAmiri123/edat/log"

// OrchestratorOption options for Orchestrator
type OrchestratorOption func(o *Orchestrator)

// WithOrchestratorLogger is an option to set the log.Logger of the Orchestrator
func WithOrchestratorLogger(logger edatlog.Logger) OrchestratorOption {
	return func(o *Orchestrator) {
		o.logger = logger
	}
}
