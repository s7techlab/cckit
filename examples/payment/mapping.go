package payment

import (
	"github.com/s7techlab/cckit/examples/payment/schema"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	// State mappings
	StateMappings = m.StateMappings{}.
			Add(&schema.Payment{}, m.PKeyAttr(`Type`, `Id`)) //key will be <'Payment',Type, Id>
	// same same
	//Add(&schema.Payment{}, m.PKeySchema(&schema.PaymentId{}))

	// Event mappings
	EventMappings = m.EventMappings{}.
			Add(&schema.PaymentEvent{})
)
