package edatpgx

import (
	"github.com/rezaAmiri123/edat/es"
	edatlog "github.com/rezaAmiri123/edat/log"
)

type SnapshotStoreOption func(store *SnapshotStore)

func WithSnapshotStoreTableName(tableName string) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.tableName = tableName
	}
}

func WithSnapshotStoreStrategy(strategy es.SnapshotStrategy) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.strategy = strategy
	}
}

func WithSnapshotStoreLogger(logger edatlog.Logger) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.logger = logger
	}
}
