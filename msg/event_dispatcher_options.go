package msg

import edatlog "github.com/rezaAmiri123/edat/log"

// EventDispatcherOption options for EventDispatcher
type EventDispatcherOption func(consumer *EventDispatcher)

// WithEventDispatcherLogger is an option to set the log.Logger of the EventDispatcher
func WithEventDispatcherLogger(logger edatlog.Logger) EventDispatcherOption {
	return func(dispatcher *EventDispatcher) {
		dispatcher.logger = logger
	}
}
