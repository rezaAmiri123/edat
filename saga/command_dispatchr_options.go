package saga

import edatlog "github.com/rezaAmiri123/edat/log"

// CommandDispatcherOption options for CommandConsumers
type CommandDispatcherOption func(dispatcher *CommandDispatcher)

// WithCommandDispatcherLogger is an option to set the log.Logger of the CommandDispatcher
func WithCommandDispatcherLogger(logger edatlog.Logger) CommandDispatcherOption {
	return func(dispatcher *CommandDispatcher) {
		dispatcher.logger = logger
	}
}
