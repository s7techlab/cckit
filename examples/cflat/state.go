package cflat

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	StateMappings = m.StateMappings{}.
			Add(&Flat{},
			m.PKeySchema(&FlatId{}),
			m.List(&Flats{})).
		Add(&FlatResident{},
			m.PKeySchema(&FlatResidentId{}),
			m.List(&FlatResidents{})).
		Add(&FlatRoom{},
			m.PKeySchema(&FlatRoomId{}),
			m.List(&FlatRooms{}))

	EventMappings = m.EventMappings{}.
			Add(&FlatCreated{}).
			Add(&FlatDeleted{}).
			Add(&FlatUpdated{}).
			Add(&FlatResidentsUpdated{}).
			Add(&FlatRoomsUpdated{})
)

func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}
