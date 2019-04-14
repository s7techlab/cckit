package cpaper

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param/defparam"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	// State mappings
	StateMappings = m.StateMappings{}.Add(
		&schema.CommercialPaper{},                 // define mapping for this structure
		m.PKeySchema(&schema.CommercialPaperId{}), // key  will be <`CommercialPaper`, Issuer, PaperNumber>
		m.List(&schema.CommercialPaperList{}))     // structure-result of list method

	// EventMappings
	EventMappings = m.EventMappings{}.
			Add(&schema.IssueCommercialPaper{}). // event name will be `IssueCommercialPaper`,  payload - same as issue payload
			Add(&schema.BuyCommercialPaper{}).
			Add(&schema.RedeemCommercialPaper{})
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

func cpaperList(c router.Context) (interface{}, error) {
	// commercial paper key is composite key <`CommercialPaper`>, {Issuer}, {PaperNumber} >
	// where `CommercialPaper` - namespace of this type
	// list method retrieves entries from chaincode state
	// using GetStateByPartialCompositeKey method, then unmarshal received from state bytes via proto.Ummarshal method
	// and creates slice of *schema.CommercialPaper
	return c.State().List(&schema.CommercialPaper{})
}

func cpaperIssue(c router.Context) (interface{}, error) {
	var (
		issue  = c.Param().(*schema.IssueCommercialPaper) //default parameter
		cpaper = &schema.CommercialPaper{
			Issuer:       issue.Issuer,
			PaperNumber:  issue.PaperNumber,
			Owner:        issue.Issuer,
			IssueDate:    issue.IssueDate,
			MaturityDate: issue.MaturityDate,
			FaceValue:    issue.FaceValue,
			State:        schema.CommercialPaper_ISSUED, // initial state
		}
		err error
	)

	if err = c.Event().Set(issue); err != nil {
		return nil, err
	}

	return cpaper, c.State().Insert(cpaper)
}

func cpaperBuy(c router.Context) (interface{}, error) {

	var (
		cpaper *schema.CommercialPaper

		// but tx payload
		buy = c.Param().(*schema.BuyCommercialPaper)

		// current commercial paper state
		cp, err = c.State().Get(
			&schema.CommercialPaperId{Issuer: buy.Issuer, PaperNumber: buy.PaperNumber},
			&schema.CommercialPaper{})
	)

	if err != nil {
		return nil, errors.Wrap(err, `not found`)
	}
	cpaper = cp.(*schema.CommercialPaper)

	// Validate current owner
	if cpaper.Owner != buy.CurrentOwner {
		return nil, fmt.Errorf(`paper %s %s is not owned by %s`, cpaper.Issuer, cpaper.PaperNumber, buy.CurrentOwner)
	}

	// First buy moves state from ISSUED to TRADING
	if cpaper.State == schema.CommercialPaper_ISSUED {
		cpaper.State = schema.CommercialPaper_TRADING
	}

	// Check paper is not already REDEEMED
	if cpaper.State == schema.CommercialPaper_TRADING {
		cpaper.Owner = buy.NewOwner
	} else {
		return nil, fmt.Errorf(`paper %s %s is not trading.current state = %s`, cpaper.Issuer, cpaper.PaperNumber, cpaper.State)
	}

	if err = c.Event().Set(buy); err != nil {
		return nil, err
	}

	return cpaper, c.State().Put(cpaper)
}

func cpaperRedeem(c router.Context) (interface{}, error) {
	// implement me
	return nil, nil
}

func cpaperGet(c router.Context) (interface{}, error) {
	return c.State().Get(c.Param().(*schema.CommercialPaperId))
}

func cpaperDelete(c router.Context) (interface{}, error) {
	return nil, c.State().Delete(c.Param().(*schema.CommercialPaperId))
}
