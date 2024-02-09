package edatpgx

import edatlog "github.com/rezaAmiri123/edat/log"

type SagaInstanceStoreOption func(store *SagaInstanceStore)

func WithSagaInstanceStoreTableName(tableName string) SagaInstanceStoreOption {
	return func(store *SagaInstanceStore) {
		store.tableName = tableName
	}
}

func WithSagaInstanceStoreLogger(logger edatlog.Logger) SagaInstanceStoreOption {
	return func(store *SagaInstanceStore) {
		store.logger = logger
	}
}
