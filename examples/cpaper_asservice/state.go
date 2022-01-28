package cpaper_asservice

import (
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	StateMappings = m.StateMappings{}.
		//  Create mapping for Commercial Paper entity
		Add(&CommercialPaper{},
			m.PKeySchema(&CommercialPaperId{}), // Key namespace will be <"CommercialPaper", Issuer, PaperNumber>
			m.List(&CommercialPaperList{}),     // Structure of result for List method
			m.UniqKey("ExternalId"),            // External Id is unique
		)

	EventMappings = m.EventMappings{}.
		// Event name will be "IssueCommercialPaper", payload - same as issue payload
		Add(&IssueCommercialPaper{}).
		// Event name will be "BuyCommercialPaper"
		Add(&BuyCommercialPaper{}).
		// Event name will be "RedeemCommercialPaper"
		Add(&RedeemCommercialPaper{})
)

// State with chaincode mappings
func State(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), StateMappings)
}

// Event with chaincode mappings
func Event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), EventMappings)
}
