package msgtest

import "github.com/rezaAmiri123/edat/msg/msgmocks"

func MockConsumer(setup func(m *msgmocks.Consumer)) *msgmocks.Consumer {
	m := &msgmocks.Consumer{}
	setup(m)
	return m
}
