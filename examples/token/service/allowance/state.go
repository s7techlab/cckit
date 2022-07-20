package allowance

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	StateMappings = m.StateMappings{}.
		//  Create mapping for Allowance entity
		// key `Allowance`,`{OwnerAddress}`,`{SpenderAddress}`,`{Path[0]}`..., `{Path[n]`
		Add(&Allowance{},
			m.PKeySchema(&AllowanceId{}),
			m.List(&Allowances{}), // Structure of result for List method
		)

	EventMappings = m.EventMappings{}.
			Add(&Approved{}).
			Add(&TransferredFrom{})
)

// State with chaincode mappings
func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

// Event with chaincode mappings
func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}
