package config

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	StateMappings = m.StateMappings{}.
			Add(&TokenType{},
			m.PKeySchema(&TokenTypeId{}),
			m.List(&TokenTypes{})).
		Add(&TokenGroup{},
			m.PKeySchema(&TokenGroupId{}),
			m.List(&TokenGroups{})).
		Add(&Config{},
			m.WithConstPKey())

	EventMappings = m.EventMappings{}.
			Add(&TokenTypeCreated{}).
			Add(&TokenGroupCreated{})
)

// State with chaincode mappings
func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

// Event with chaincode mappings
func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}
