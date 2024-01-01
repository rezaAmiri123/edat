package logtest

import "github.com/rezaAmiri123/edat/log/logmocks"

func MockLogger(setup func(m *logmocks.Logger))*logmocks.Logger{
	m := &logmocks.Logger{}
	setup(m)
	
	return m
}