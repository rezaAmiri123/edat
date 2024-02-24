package msg

import "github.com/rezaAmiri123/edat/log"

// PublisherOption options for PublisherPublisher
type PublisherOption func(*Publisher)

// WithPublisherLogger is an option to set the log.Logger of the Publisher
func WithPublisherLogger(logger edatlog.Logger) PublisherOption {
	return func(publisher *Publisher) {
		publisher.logger = logger
	}
}
