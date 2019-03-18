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
			m.PKeySchema(&schema.CommercialPaperId{}), //key namespace will be <`CommercialPaper`, Issuer, PaperNumber>
			m.List(&schema.CommercialPaperList{}),     // list container
			m.UniqKey(`ExternalId`),                   // external is uniq
		)

	// EventMappings
	EventMappings = m.EventMappings{}.
		// event name will be `IssueCommercialPaper`,  payload - same as issue payload
		Add(&schema.IssueCommercialPaper{}).
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
	// commercial paper key is composite key <`CommercialPaper`>, {Issuer}, {PaperNumber} >
	// where `CommercialPaper` - namespace of this type
	// list method retrieves entries from chaincode state
	// using GetStateByPartialCompositeKey method, then unmarshal received from state bytes via proto.Ummarshal method
	// and creates slice of *schema.CommercialPaper
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
		// input message
		issue = c.Param().(*schema.IssueCommercialPaper) //default parameter
	)
	// validate input message by rules, defined in schema
	if err = issue.Validate(); err != nil {
		return err, errors.Wrap(err, `payload validation`)
	}

	// create state entry
	cpaper := &schema.CommercialPaper{
		Issuer:       issue.Issuer,
		PaperNumber:  issue.PaperNumber,
		Owner:        issue.Issuer,
		IssueDate:    issue.IssueDate,
		MaturityDate: issue.MaturityDate,
		FaceValue:    issue.FaceValue,
		State:        schema.CommercialPaper_ISSUED, // initial state
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

func invokeCPaperRedeem(c router.Context) (interface{}, error) {

	return nil, nil
}

func invokeCPaperDelete(c router.Context) (interface{}, error) {
	var (
		id = c.Param().(*schema.CommercialPaperId)
	)
	return nil, c.State().Delete(id)
}
