package cpaper

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	// State mappings
	StateMappings = m.StateMappings{}.
			Add(&schema.CommercialPaper{}, //key namespace will be []string{ CommercialPaper }
			m.StatePKeyer(func(e interface{}) ([]string, error) {
				cp := e.(*schema.CommercialPaper)
				// primary key consists of namespace, issuer and paper
				return []string{cp.Issuer, cp.PaperNumber}, nil
			}))

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
		Query(`get`, cpaperGet, p.String(`issuer`), p.String(`paperNumber`)).

		// txn methods
		Invoke(`issue`, cpaperIssue, p.Proto(p.Default, &schema.IssueCommercialPaper{})).
		Invoke(`buy`, cpaperBuy, p.Proto(p.Default, &schema.BuyCommercialPaper{})).
		Invoke(`redeem`, cpaperRedeem, p.Proto(p.Default, &schema.RedeemCommercialPaper{})).
		Invoke(`delete`, cpaperDelete, p.String(`issuer`), p.String(`paperNumber`))

	return router.NewChaincode(r)
}

func cpaperList(c router.Context) (interface{}, error) {
	return c.State().List(&schema.CommercialPaper{})
}

func cpaperIssue(c router.Context) (interface{}, error) {
	var (
		issue  = c.Param(p.Default).(*schema.IssueCommercialPaper)
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
		buy = c.Param(p.Default).(*schema.BuyCommercialPaper)

		// current commercial paper state
		cp, err = c.State().Get(&schema.CommercialPaper{
			Issuer:      buy.Issuer,
			PaperNumber: buy.PaperNumber}, &schema.CommercialPaper{})
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

	return cpaper, c.State().Put(cpaper)
}

func cpaperRedeem(c router.Context) (interface{}, error) {

	return nil, nil
}

func cpaperGet(c router.Context) (interface{}, error) {
	return c.State().Get(&schema.CommercialPaper{
		Issuer:      c.ParamString(`issuer`),
		PaperNumber: c.ParamString(`paperNumber`)})
}

func cpaperDelete(c router.Context) (interface{}, error) {
	return nil, c.State().Delete(&schema.CommercialPaper{
		Issuer:      c.ParamString(`issuer`),
		PaperNumber: c.ParamString(`paperNumber`)})
}
