package msgtest

import "github.com/rezaAmiri123/edat/core/coretest"

type Entity struct {
	coretest.Entity
}

func (Entity) DestinationChannel() string { return "entity-channel" }
