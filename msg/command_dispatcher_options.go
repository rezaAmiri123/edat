package msg

import edatlog "github.com/rezaAmiri123/edat/log"

// CommandDispatcherOption options for CommandDispatcher
type CommandDispatcherOption func(consumer *CommandDispatcher)

// WithCommandDispatcherLogger is an option to set the log.Logger of the CommandDispatcher
func WithCommandDispatcherLogger(logger edatlog.Logger) CommandDispatcherOption {
	return func(consumer *CommandDispatcher) {
		consumer.logger = logger
	}
}
