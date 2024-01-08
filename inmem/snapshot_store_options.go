package inmem

import "github.com/rezaAmiri123/edat/es"

// SnapshotStoreOption options for SnapshotStore
type SnapshotStoreOption func(store *SnapshotStore)

// WithSnapshotStoreStrategy sets the snapshotting strategy for SnapshotStore
func WithSnapshotStoreStrategy(strategy es.SnapshotStrategy) SnapshotStoreOption {
	return func(store *SnapshotStore) {
		store.strategy = strategy
	}
}
