package coretest

import "github.com/rezaAmiri123/edat/core/coremocks"

type (
	Command             struct{ Value string }
	UnregisteredCommand struct{ Value string }
)

func (Command) CommandName() string             { return "coretest.Command" }
func (UnregisteredCommand) CommandName() string { return "coretest.UnregisteredCommand" }

func MockCommand(setup func(m *coremocks.Command)) *coremocks.Command {
	m := &coremocks.Command{}
	setup(m)
	return m
}
