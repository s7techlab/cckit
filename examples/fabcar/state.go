package fabcar

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	StateMappings = m.StateMappings{}.
			Add(&Maker{},
			m.PKeySchema(&MakerName{}),
			m.List(&Makers{})).
		Add(&Car{},
			m.PKeySchema(&CarId{}),
			m.List(&Cars{})).
		Add(&CarOwner{},
			m.PKeySchema(&CarOwnerId{}),
			m.List(&CarOwners{})).
		Add(&CarDetail{},
			m.PKeySchema(&CarDetailId{}),
			m.List(&CarDetails{}))

	EventMappings = m.EventMappings{}.
			Add(&MakerCreated{}).
			Add(&MakerDeleted{}).
			Add(&CarCreated{}).
			Add(&CarDeleted{}).
			Add(&CarUpdated{}).
			Add(&CarOwnerDeleted{}).
			Add(&CarOwnersUpdated{}).
			Add(&CarDetailDeleted{}).
			Add(&CarDetailsUpdated{})
)

func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}
