package cpaper_extended

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/examples/cpaper_extended/schema"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/router/param/defparam"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	// State mappings
	StateMappings = m.StateMappings{}.
		//  Create mapping for Commercial Paper entity
		Add(&schema.CommercialPaper{},
			m.PKeySchema(&schema.CommercialPaperId{}), // Key namespace will be <`CommercialPaper`, Issuer, PaperNumber>
			m.List(&schema.CommercialPaperList{}),     // Structure of result for List method
			m.UniqKey(`ExternalId`),                   // External Id is unique
		)

	// EventMappings
	EventMappings = m.EventMappings{}.
		// Event name will be `IssueCommercialPaper`, payload - same as issue payload
		Add(&schema.IssueCommercialPaper{}).
		// Event name will be `BuyCommercialPaper`
		Add(&schema.BuyCommercialPaper{}).
		// Event name will be `RedeemCommercialPaper`
		Add(&schema.RedeemCommercialPaper{})
)

func NewCC() *router.Chaincode {

	r := router.New(`commercial_paper`)

	// Mappings for chaincode state
	r.Use(m.MapStates(StateMappings))

	// Mappings for chaincode events
	r.Use(m.MapEvents(EventMappings))

	// Store on the ledger the information about chaincode instantiator
	r.Init(owner.InvokeSetFromCreator)

	// Method for debug chaincode state
	debug.AddHandlers(r, `debug`, owner.Only)

	r.
		// read methods
		Query(`list`, queryCPapers).

		// Get method has 2 params - commercial paper primary key components
		Query(`get`, queryCPaper, defparam.Proto(&schema.CommercialPaperId{})).
		Query(`getByExternalId`, queryCPaperGetByExternalId, param.String(`externalId`)).

		// txn methods
		Invoke(`issue`, invokeCPaperIssue, defparam.Proto(&schema.IssueCommercialPaper{})).
		Invoke(`buy`, invokeCPaperBuy, defparam.Proto(&schema.BuyCommercialPaper{})).
		Invoke(`redeem`, invokeCPaperRedeem, defparam.Proto(&schema.RedeemCommercialPaper{})).
		Invoke(`delete`, invokeCPaperDelete, defparam.Proto(&schema.CommercialPaperId{}))

	return router.NewChaincode(r)
}

func queryCPapers(c router.Context) (interface{}, error) {
	// List method retrieves all entries from the ledger using GetStateByPartialCompositeKey method and passing it the
	// namespace of our contract type, in this example that's `CommercialPaper`, then it unmarshals received bytes via
	// proto.Ummarshal method and creates a []schema.CommercialPaperList as defined in the
	// `StateMappings` variable at the top of the file
	return c.State().List(&schema.CommercialPaper{})
}

func queryCPaper(c router.Context) (interface{}, error) {
	var (
		id = c.Param().(*schema.CommercialPaperId)
	)
	return c.State().Get(id)
}

func queryCPaperGetByExternalId(c router.Context) (interface{}, error) {
	var (
		externalId = c.ParamString(`externalId`)
	)
	return c.State().(m.MappedState).GetByUniqKey(&schema.CommercialPaper{}, `ExternalId`, []string{externalId})
}

func invokeCPaperIssue(c router.Context) (res interface{}, err error) {
	var (
		// Input message
		issue = c.Param().(*schema.IssueCommercialPaper) //default parameter
	)
	// Validate input message using the rules defined in schema
	if err = issue.Validate(); err != nil {
		return err, errors.Wrap(err, `payload validation`)
	}

	// Create state entry
	cpaper := &schema.CommercialPaper{
		Issuer:       issue.Issuer,
		PaperNumber:  issue.PaperNumber,
		Owner:        issue.Issuer,
		IssueDate:    issue.IssueDate,
		MaturityDate: issue.MaturityDate,
		FaceValue:    issue.FaceValue,
		State:        schema.CommercialPaper_ISSUED, // Initial state
		ExternalId:   issue.ExternalId,
	}

	if err = c.Event().Set(issue); err != nil {
		return nil, err
	}

	return cpaper, c.State().Insert(cpaper)
}

func invokeCPaperBuy(c router.Context) (interface{}, error) {

	var (
		cpaper *schema.CommercialPaper

		// Buy transaction payload
		buy = c.Param().(*schema.BuyCommercialPaper)

		// Get the current commercial paper state
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

func invokeCPaperRedeem(c router.Context) (interface{}, error) {

	return nil, nil
}

func invokeCPaperDelete(c router.Context) (interface{}, error) {
	var (
		id = c.Param().(*schema.CommercialPaperId)
	)
	return nil, c.State().Delete(id)
}
