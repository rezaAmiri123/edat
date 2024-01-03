package msgpack

import (
	"reflect"
	"sync"

	"github.com/rezaAmiri123/edat/core"
	registertypes "github.com/rezaAmiri123/edat/core/register_types"
	"github.com/shamaton/msgpack"
)

func init() {
	core.RegisterDefaultMarshaller(newMsgPackMarshaler())
	registertypes.RegisterTypes()
}

type msgPackMarshaler struct {
	items map[string]reflect.Type
	mu    sync.Mutex
}

var _ core.Marshaller = (*msgPackMarshaler)(nil)

func newMsgPackMarshaler() *msgPackMarshaler {
	return &msgPackMarshaler{
		items: map[string]reflect.Type{},
		mu:    sync.Mutex{},
	}
}

func (m *msgPackMarshaler) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m *msgPackMarshaler) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func (m *msgPackMarshaler) GetType(typeName string) reflect.Type {
	return m.items[typeName]
}

func (m *msgPackMarshaler) RegisterType(typeName string, v reflect.Type) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[typeName] = v
}
