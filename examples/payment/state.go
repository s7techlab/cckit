package payment

import (
	"github.com/s7techlab/cckit/examples/payment/schema"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	// State mappings
	StateMappings = m.StateMappings{}.
		Add(&schema.Payment{}, //key namespace will be []string{ CommercialPaper }
			m.UseStatePKeyer(func(e interface{}) ([]string, error) {
				cp := e.(*schema.Payment)
				// primary key consists of namespace, issuer and paper
				return []string{cp.Type, cp.Id}, nil
			}))
)
