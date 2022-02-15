package crosscc

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	StateMappings = m.StateMappings{}.
			Add(&ServiceLocator{},
			m.PKeySchema(&ServiceLocatorId{}),
			m.List(&ServiceLocators{}))

	EventMappings = m.EventMappings{}.
			Add(&ServiceLocatorSet{})
)

func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}
