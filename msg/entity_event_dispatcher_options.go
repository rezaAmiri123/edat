package msg

import "github.com/rezaAmiri123/edat/log"

// EntityEventDispatcherOption options for EntityEventDispatcher
type EntityEventDispatcherOption func(consumer *EntityEventDispatcher)

// WithEntityEventDispatcherLogger is an option to set the log.Logger of the EntityEventDispatcher
func WithEntityEventDispatcherLogger(logger edatlog.Logger) EntityEventDispatcherOption {
	return func(dispatcher *EntityEventDispatcher) {
		dispatcher.logger = logger
	}
}
