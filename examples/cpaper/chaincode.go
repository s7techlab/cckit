package cpaper

import (
	"fmt"
	"time"

	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
)

type TaskState string

const (
	CommercialPaperTypeName = "CommercialPaper"

	CommercialPaperIssued   TaskState = "issued"
	CommercialPaperTrading  TaskState = "trading"
	CommercialPaperRedeemed TaskState = "redeemed"
)

type CommercialPaper struct {
	Issuer       string    `json:"issuer,omitempty"`
	PaperNumber  string    `json:"paper_number,omitempty"`
	Owner        string    `json:"owner,omitempty"`
	IssueDate    time.Time `json:"issue_date,omitempty"`
	MaturityDate time.Time `json:"maturity_date,omitempty"`
	FaceValue    int32     `json:"face_value,omitempty"`
	State        TaskState `json:"state,omitempty"`
}

// Key commercial paper has a composite key < CommercialPaper, Issuer, PaperNumber >
func (c CommercialPaper) Key() ([]string, error) {
	return []string{CommercialPaperTypeName, c.Issuer, c.PaperNumber}, nil
}

type IssueCommercialPaper struct {
	Issuer       string    `json:"issuer,omitempty"`
	PaperNumber  string    `json:"paper_number,omitempty"`
	IssueDate    time.Time `json:"issue_date,omitempty"`
	MaturityDate time.Time `json:"maturity_date,omitempty"`
	FaceValue    int32     `json:"face_value,omitempty"`
}

type BuyCommercialPaper struct {
	Issuer       string    `json:"issuer,omitempty"`
	PaperNumber  string    `json:"paper_number,omitempty"`
	CurrentOwner string    `json:"current_owner,omitempty"`
	NewOwner     string    `json:"new_owner,omitempty"`
	Price        int32     `json:"price,omitempty"`
	PurchaseDate time.Time `json:"purchase_date,omitempty"`
}

type RedeemCommercialPaper struct {
	Issuer         string    `json:"issuer,omitempty"`
	PaperNumber    string    `json:"paper_number,omitempty"`
	RedeemingOwner string    `json:"redeeming_owner,omitempty"`
	RedeemDate     time.Time `json:"redeem_date,omitempty"`
}

func NewCC() *router.Chaincode {
	r := router.New(`commercial_paper`)

	r.Init(func(context router.Context) (i interface{}, e error) {
		// No implementation required with this example
		// It could be where data migration is performed, if necessary
		return nil, nil
	})

	r.
		// Read methods
		Query(`list`, list).
		// Get method has 2 params - commercial paper primary key components
		Query(`get`, get, param.String("issuer"), param.String("paper_number")).

		// Transaction methods
		Invoke(`issue`, issue, param.Struct("issueData", &IssueCommercialPaper{})).
		Invoke(`buy`, buy, param.Struct("buyData", &BuyCommercialPaper{})).
		Invoke(`redeem`, redeem, param.Struct("redeemData", &RedeemCommercialPaper{}))

	return router.NewChaincode(r)
}

func get(c router.Context) (interface{}, error) {
	var (
		issuer      = c.ParamString("issuer")
		paperNumber = c.ParamString("paper_number")
	)

	return c.State().Get(&CommercialPaper{
		Issuer:      issuer,
		PaperNumber: paperNumber,
	})
}

func list(c router.Context) (interface{}, error) {
	// List method retrieves all entries from the ledger using GetStateByPartialCompositeKey method and passing it the
	// namespace of our contract type, in this example that's `CommercialPaper`, then it unmarshalls received bytes via
	// json.Unmarshal method and creates a JSON array of CommercialPaper entities
	return c.State().List(CommercialPaperTypeName, &CommercialPaper{})
}

func issue(c router.Context) (interface{}, error) {
	var (
		issueData       = c.Param("issueData").(IssueCommercialPaper) // Assert the chaincode parameter
		commercialPaper = &CommercialPaper{
			Issuer:       issueData.Issuer,
			PaperNumber:  issueData.PaperNumber,
			Owner:        issueData.Issuer,
			IssueDate:    issueData.IssueDate,
			MaturityDate: issueData.MaturityDate,
			FaceValue:    issueData.FaceValue,
			State:        CommercialPaperIssued, // Initial state
		}
	)

	return commercialPaper, c.State().Insert(commercialPaper)
}

func buy(c router.Context) (interface{}, error) {
	var (
		commercialPaper CommercialPaper
		// Buy transaction payload
		buyData = c.Param("buyData").(BuyCommercialPaper)

		// Get the current commercial paper state
		cp, err = c.State().Get(&CommercialPaper{
			Issuer:      buyData.Issuer,
			PaperNumber: buyData.PaperNumber,
		}, &CommercialPaper{})
	)

	if err != nil {
		return nil, fmt.Errorf("paper not found: %w", err)
	}

	commercialPaper = cp.(CommercialPaper)

	// Validate current owner
	if commercialPaper.Owner != buyData.CurrentOwner {
		return nil, fmt.Errorf(
			"paper %s %s is not owned by %s",
			commercialPaper.Issuer, commercialPaper.PaperNumber, buyData.CurrentOwner)
	}

	// First buyData moves state from ISSUED to TRADING
	if commercialPaper.State == CommercialPaperIssued {
		commercialPaper.State = CommercialPaperTrading
	}

	// Check paper is not already REDEEMED
	if commercialPaper.State == CommercialPaperTrading {
		commercialPaper.Owner = buyData.NewOwner
	} else {
		return nil, fmt.Errorf(
			"paper %s %s is not trading.current state = %s",
			commercialPaper.Issuer, commercialPaper.PaperNumber, commercialPaper.State)
	}

	return commercialPaper, c.State().Put(commercialPaper)
}

func redeem(c router.Context) (interface{}, error) {
	var (
		commercialPaper CommercialPaper

		// Buy transaction payload
		redeemData = c.Param("redeemData").(RedeemCommercialPaper)

		// Get the current commercial paper state
		cp, err = c.State().Get(&CommercialPaper{
			Issuer:      redeemData.Issuer,
			PaperNumber: redeemData.PaperNumber,
		}, &CommercialPaper{})
	)

	if err != nil {
		return nil, fmt.Errorf("paper not found: %w", err)
	}

	commercialPaper = cp.(CommercialPaper)

	// Check paper is not REDEEMED
	if commercialPaper.State == CommercialPaperRedeemed {
		return nil, fmt.Errorf(
			"paper %s %s is already redeemed",
			commercialPaper.Issuer, commercialPaper.PaperNumber)
	}

	// Verify that the redeemer owns the commercial paper before redeeming it
	if commercialPaper.Owner == redeemData.RedeemingOwner {
		commercialPaper.Owner = redeemData.Issuer
		commercialPaper.State = CommercialPaperRedeemed
	} else {
		return nil, fmt.Errorf(
			"redeeming owner does not own paper %s %s",
			commercialPaper.Issuer, commercialPaper.PaperNumber)
	}

	return commercialPaper, c.State().Put(commercialPaper)
}
