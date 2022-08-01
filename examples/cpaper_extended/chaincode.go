package cpaper_extended

import (
	"fmt"

	"github.com/s7techlab/cckit/examples/cpaper_extended/schema"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/router/param/defparam"
	m "github.com/s7techlab/cckit/state/mapping"
)

var (
	// StateMappings state mappings
	StateMappings = m.StateMappings{}.
		//  Create mapping for Commercial Paper entity
		Add(&schema.CommercialPaper{},
			m.PKeySchema(&schema.CommercialPaperId{}), // Key namespace will be <"CommercialPaper", Issuer, PaperNumber>
			m.List(&schema.CommercialPaperList{}),     // Structure of result for List method
			m.UniqKey("ExternalId"),                   // External Id is unique
		)

	// EventMappings event mappings
	EventMappings = m.EventMappings{}.
		// Event name will be "IssueCommercialPaper", payload - same as issue payload
		Add(&schema.IssueCommercialPaper{}).
		// Event name will be "BuyCommercialPaper"
		Add(&schema.BuyCommercialPaper{}).
		// Event name will be "RedeemCommercialPaper"
		Add(&schema.RedeemCommercialPaper{})
)

func NewCC() *router.Chaincode {

	r := router.New("commercial_paper")

	// Mappings for chaincode state
	r.Use(m.MapStates(StateMappings))

	// Mappings for chaincode events
	r.Use(m.MapEvents(EventMappings))

	// Store on the ledger the information about chaincode instantiator
	r.Init(owner.InvokeSetFromCreator)

	// Method for debug chaincode state
	debug.AddHandlers(r, "debug", owner.Only)

	r.
		// read methods
		Query("list", queryCPapers).

		// Get method has 2 params - commercial paper primary key components
		Query("get", queryCPaper, defparam.Proto(&schema.CommercialPaperId{})).
		Query("getByExternalId", queryCPaperGetByExternalId, param.String("externalId")).

		// txn methods
		Invoke("issue", invokeCPaperIssue, defparam.Proto(&schema.IssueCommercialPaper{})).
		Invoke("buy", invokeCPaperBuy, defparam.Proto(&schema.BuyCommercialPaper{})).
		Invoke("redeem", invokeCPaperRedeem, defparam.Proto(&schema.RedeemCommercialPaper{})).
		Invoke("delete", invokeCPaperDelete, defparam.Proto(&schema.CommercialPaperId{}))

	return router.NewChaincode(r)
}

func queryCPapers(c router.Context) (interface{}, error) {
	// List method retrieves all entries from the ledger using GetStateByPartialCompositeKey method and passing it the
	// namespace of our contract type, in this example that's "CommercialPaper", then it unmarshals received bytes via
	// proto.Ummarshal method and creates a []schema.CommercialPaperList as defined in the
	// "StateMappings" variable at the top of the file
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
		externalId = c.ParamString("externalId")
	)
	return c.State().(m.MappedState).GetByKey(&schema.CommercialPaper{}, "ExternalId", []string{externalId})
}

func invokeCPaperIssue(c router.Context) (res interface{}, err error) {
	var (
		// Input message
		issueData = c.Param().(*schema.IssueCommercialPaper) // Default parameter
	)
	// Validate input message using the rules defined in schema
	if err = issueData.Validate(); err != nil {
		return nil, fmt.Errorf("payload validation: %w", err)
	}

	// Create state entry
	cPaper := &schema.CommercialPaper{
		Issuer:       issueData.Issuer,
		PaperNumber:  issueData.PaperNumber,
		Owner:        issueData.Issuer,
		IssueDate:    issueData.IssueDate,
		MaturityDate: issueData.MaturityDate,
		FaceValue:    issueData.FaceValue,
		State:        schema.CommercialPaper_STATE_ISSUED, // Initial state
		ExternalId:   issueData.ExternalId,
	}

	if err = c.Event().Set(issueData); err != nil {
		return nil, err
	}

	return cPaper, c.State().Insert(cPaper)
}

func invokeCPaperBuy(c router.Context) (interface{}, error) {
	var (
		cPaper *schema.CommercialPaper

		// Buy transaction payload
		buyData = c.Param().(*schema.BuyCommercialPaper)

		// Get the current commercial paper state
		cp, err = c.State().Get(
			&schema.CommercialPaperId{Issuer: buyData.Issuer, PaperNumber: buyData.PaperNumber},
			&schema.CommercialPaper{})
	)

	if err != nil {
		return nil, fmt.Errorf("not found: %w", err)
	}

	cPaper = cp.(*schema.CommercialPaper)

	// Validate current owner
	if cPaper.Owner != buyData.CurrentOwner {
		return nil, fmt.Errorf(
			"paper %s %s is not owned by %s",
			cPaper.Issuer, cPaper.PaperNumber, buyData.CurrentOwner)
	}

	// First buyData moves state from ISSUED to TRADING
	if cPaper.State == schema.CommercialPaper_STATE_ISSUED {
		cPaper.State = schema.CommercialPaper_STATE_TRADING
	}

	// Check paper is not already REDEEMED
	if cPaper.State == schema.CommercialPaper_STATE_TRADING {
		cPaper.Owner = buyData.NewOwner
	} else {
		return nil, fmt.Errorf(
			"paper %s %s is not trading.current state = %s",
			cPaper.Issuer, cPaper.PaperNumber, cPaper.State)
	}

	if err = c.Event().Set(buyData); err != nil {
		return nil, err
	}

	return cPaper, c.State().Put(cPaper)
}

func invokeCPaperRedeem(c router.Context) (interface{}, error) {
	var (
		commercialPaper *schema.CommercialPaper

		// Buy transaction payload
		redeemData = c.Param().(*schema.RedeemCommercialPaper)

		// Get the current commercial paper state
		cp, err = c.State().Get(&schema.CommercialPaper{
			Issuer:      redeemData.Issuer,
			PaperNumber: redeemData.PaperNumber,
		}, &schema.CommercialPaper{})
	)

	if err != nil {
		return nil, fmt.Errorf("paper not found: %w", err)
	}

	commercialPaper = cp.(*schema.CommercialPaper)

	// Check paper is not REDEEMED
	if commercialPaper.State == schema.CommercialPaper_STATE_REDEEMED {
		return nil, fmt.Errorf(
			"paper %s %s is already redeemed",
			commercialPaper.Issuer, commercialPaper.PaperNumber)
	}

	// Verify that the redeemer owns the commercial paper before redeeming it
	if commercialPaper.Owner == redeemData.RedeemingOwner {
		commercialPaper.Owner = redeemData.Issuer
		commercialPaper.State = schema.CommercialPaper_STATE_REDEEMED
	} else {
		return nil, fmt.Errorf(
			"redeeming owner does not own paper %s %s",
			commercialPaper.Issuer, commercialPaper.PaperNumber)
	}

	if err = c.Event().Set(redeemData); err != nil {
		return nil, err
	}

	return commercialPaper, c.State().Put(commercialPaper)
}

func invokeCPaperDelete(c router.Context) (interface{}, error) {
	var (
		id = c.Param().(*schema.CommercialPaperId)
	)
	return nil, c.State().Delete(id)
}
