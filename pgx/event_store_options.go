package edatpgx

import "github.com/rezaAmiri123/edat/log"

type EventStoreOption func(store *EventStore)

func WithEventStoreTableName(tableName string)EventStoreOption{
	return func(store *EventStore) {
		store.tableName = tableName
	}
}

func WithEventStoreLogger(logger log.Logger)EventStoreOption{
	return func(store *EventStore) {
		store.logger = logger
	}
}
