package edatpgx

import "github.com/rezaAmiri123/edat/log"

type MessageStoreOption func(store *MessageStore)

func WithMessageStoreTableName(tableName string) MessageStoreOption {
	return func(store *MessageStore) {
		store.tableName = tableName
	}
}

func WithMessageStoreLogger(logger edatlog.Logger) MessageStoreOption {
	return func(store *MessageStore) {
		store.logger = logger
	}
}
