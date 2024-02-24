package msgtest

import "github.com/rezaAmiri123/edat/msg/msgmocks"

func MockProducer(setup func(m *msgmocks.Producer)) *msgmocks.Producer {
	m := &msgmocks.Producer{}
	setup(m)
	return m
}
