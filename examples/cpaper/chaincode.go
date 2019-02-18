package cpaper

import (
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param/defparam"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	// State mappings
	StateMappings = m.StateMappings{}.
			Add(&schema.CommercialPaper{}, m.PKeySchema(&schema.CommercialPaperId{})) //key namespace will be <`CommercialPaper`, Issuer, PaperNumber>

	// EventMappings
	EventMappings = m.EventMappings{}.
			Add(&schema.IssueCommercialPaper{}) // event name will be `IssueCommercialPaper`,  payload - same as issue payload

)

func NewCC() *router.Chaincode {

	r := router.New(`commercial_paper`)

	// Mappings for chaincode state
	r.Use(m.MapStates(StateMappings))

	// Mappings for chaincode events
	r.Use(m.MapEvents(EventMappings))

	// store in chaincode state information about chaincode first instantiator
	r.Init(owner.InvokeSetFromCreator)

	// method for debug chaincode state
	debug.AddHandlers(r, `debug`, owner.Only)

	r.
		// read methods
		Query(`list`, cpaperList).

		// Get method has 2 params - commercial paper primary key components
		Query(`get`, cpaperGet, defparam.Proto(&schema.CommercialPaperId{})).

		// txn methods
		Invoke(`issue`, cpaperIssue, defparam.Proto(&schema.IssueCommercialPaper{})).
		Invoke(`buy`, cpaperBuy, defparam.Proto(&schema.BuyCommercialPaper{})).
		Invoke(`redeem`, cpaperRedeem, defparam.Proto(&schema.RedeemCommercialPaper{})).
		Invoke(`delete`, cpaperDelete, defparam.Proto(&schema.CommercialPaperId{}))

	return router.NewChaincode(r)
}

func NewEncryptedCC() *router.Chaincode {

	r := router.New(`commercial_paper_encrypted`)

	// Mappings for chaincode state
	r.Use(m.MapStates(StateMappings))

	// Mappings for chaincode events
	r.Use(m.MapEvents(EventMappings))

	r.Pre(encryption.ArgsDecrypt)
	r.Use(encryption.EncStateContext) // use encrypted state by default

	// store in chaincode state information about chaincode first instantiator
	r.Init(owner.InvokeSetFromCreator)

	// method for debug chaincode state
	debug.AddHandlers(r, `debug`, owner.Only)

	r.
		// read methods
		Query(`list`, cpaperList).

		// Get method has 2 params - commercial paper primary key components
		Query(`get`, cpaperGet, defparam.Proto(&schema.CommercialPaperId{})).

		// txn methods
		Invoke(`issue`, cpaperIssue, defparam.Proto(&schema.IssueCommercialPaper{})).
		Invoke(`buy`, cpaperBuy, defparam.Proto(&schema.BuyCommercialPaper{})).
		Invoke(`redeem`, cpaperRedeem, defparam.Proto(&schema.RedeemCommercialPaper{})).
		Invoke(`delete`, cpaperDelete, defparam.Proto(&schema.CommercialPaperId{}))

	return router.NewChaincode(r)
}
