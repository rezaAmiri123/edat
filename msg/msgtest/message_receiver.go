package msgtest

import "github.com/rezaAmiri123/edat/msg/msgmocks"

func MockMessageReceiver(setup func(m *msgmocks.MessageReceiver)) *msgmocks.MessageReceiver {
	m := &msgmocks.MessageReceiver{}
	setup(m)
	return m
}
