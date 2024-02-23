package saga

import "github.com/rezaAmiri123/edat/msg"

// WithSagaInfo is an option to set additional Saga specific headers
func WithSagaInfo(instance *Instance) msg.MessageOption {
	return msg.WithHeaders(map[string]string{
		MessageCommandSagaID:   instance.sagaID,
		MessageCommandSagaName: instance.sagaName,
	})
}
