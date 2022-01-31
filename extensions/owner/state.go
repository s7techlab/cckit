package owner

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

// OwnerStateKey key used to store owner grant struct in chain code state
// "handler" part of owner extension supports only one owner
// "service" part of owner extension supports multiple owners
const OwnerStateKey = `OWNER`

var (
	StateMappings = m.StateMappings{}.
			Add(&ChaincodeOwner{},
			m.PKeySchema(&OwnerId{}),
			m.List(&ChaincodeOwners{}))

	EventMappings = m.EventMappings{}.
			Add(&ChaincodeOwnerCreated{}).
			Add(&ChaincodeOwnerUpdated{}).
			Add(&ChaincodeOwnerDeleted{})
)

func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}
