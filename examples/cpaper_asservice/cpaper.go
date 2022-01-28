package cpaper_asservice

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"

	"github.com/s7techlab/cckit/router"
)

type CPaperService struct {
}

func NewService() *CPaperService {
	return &CPaperService{}
}

func (cc *CPaperService) List(ctx router.Context, in *empty.Empty) (*CommercialPaperList, error) {
	// List method retrieves all entries from the ledger using GetStateByPartialCompositeKey method and passing it the
	// namespace of our contract type, in this example that's "CommercialPaper", then it unmarshals received bytes via
	// proto.Ummarshal method and creates a []CommercialPaperList as defined in the
	// "StateMappings" variable at the top of the file
	if res, err := State(ctx).List(&CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaperList), nil
	}
}

func (cc *CPaperService) Get(ctx router.Context, id *CommercialPaperId) (*CommercialPaper, error) {
	if res, err := State(ctx).Get(id, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}
}

func (cc *CPaperService) GetByExternalId(ctx router.Context, id *ExternalId) (*CommercialPaper, error) {
	if res, err := State(ctx).GetByKey(
		&CommercialPaper{}, "ExternalId", []string{id.Id}, &CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*CommercialPaper), nil
	}
}

func (cc *CPaperService) Issue(ctx router.Context, issue *IssueCommercialPaper) (*CommercialPaper, error) {
	// Validate input message using the rules defined in schema
	if err := issue.Validate(); err != nil {
		return nil, errors.Wrap(err, "payload validation")
	}

	// Create state entry
	cpaper := &CommercialPaper{
		Issuer:       issue.Issuer,
		PaperNumber:  issue.PaperNumber,
		Owner:        issue.Issuer,
		IssueDate:    issue.IssueDate,
		MaturityDate: issue.MaturityDate,
		FaceValue:    issue.FaceValue,
		State:        CommercialPaper_STATE_ISSUED, // Initial state
		ExternalId:   issue.ExternalId,
	}

	if err := Event(ctx).Set(issue); err != nil {
		return nil, err
	}

	if err := State(ctx).Insert(cpaper); err != nil {
		return nil, err
	}
	return cpaper, nil
}

func (cc *CPaperService) Buy(ctx router.Context, buy *BuyCommercialPaper) (*CommercialPaper, error) {
	// Get the current commercial paper state
	cpaper, err := cc.Get(ctx, &CommercialPaperId{Issuer: buy.Issuer, PaperNumber: buy.PaperNumber})
	if err != nil {
		return nil, errors.Wrap(err, "get cpaper")
	}

	// Validate current owner
	if cpaper.Owner != buy.CurrentOwner {
		return nil, fmt.Errorf(
			"paper %s %s is not owned by %s",
			cpaper.Issuer, cpaper.PaperNumber, buy.CurrentOwner)
	}

	// First buyData moves state from ISSUED to TRADING
	if cpaper.State == CommercialPaper_STATE_ISSUED {
		cpaper.State = CommercialPaper_STATE_TRADING
	}

	// Check paper is not already REDEEMED
	if cpaper.State == CommercialPaper_STATE_TRADING {
		cpaper.Owner = buy.NewOwner
	} else {
		return nil, fmt.Errorf(
			"paper %s %s is not trading.current state = %s",
			cpaper.Issuer, cpaper.PaperNumber, cpaper.State)
	}

	if err = Event(ctx).Set(buy); err != nil {
		return nil, err
	}

	if err = State(ctx).Put(cpaper); err != nil {
		return nil, err
	}

	return cpaper, nil
}

func (cc *CPaperService) Redeem(ctx router.Context, redeem *RedeemCommercialPaper) (*CommercialPaper, error) {
	// Get the current commercial paper state
	cpaper, err := cc.Get(ctx, &CommercialPaperId{Issuer: redeem.Issuer, PaperNumber: redeem.PaperNumber})
	if err != nil {
		return nil, errors.Wrap(err, "get cpaper")
	}
	if err != nil {
		return nil, errors.Wrap(err, "paper not found")
	}

	// Check paper is not REDEEMED
	if cpaper.State == CommercialPaper_STATE_REDEEMED {
		return nil, fmt.Errorf("paper %s %s is already redeemed", cpaper.Issuer, cpaper.PaperNumber)
	}

	// Verify that the redeemer owns the commercial paper before redeeming it
	if cpaper.Owner == redeem.RedeemingOwner {
		cpaper.Owner = redeem.Issuer
		cpaper.State = CommercialPaper_STATE_REDEEMED
	} else {
		return nil, fmt.Errorf("redeeming owner does not own paper %s %s", cpaper.Issuer, cpaper.PaperNumber)
	}

	if err = Event(ctx).Set(redeem); err != nil {
		return nil, err
	}

	if err = State(ctx).Put(cpaper); err != nil {
		return nil, err
	}

	return cpaper, nil
}

func (cc *CPaperService) Delete(ctx router.Context, id *CommercialPaperId) (*CommercialPaper, error) {
	cpaper, err := cc.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get cpaper")
	}

	if err = State(ctx).Delete(id); err != nil {
		return nil, err
	}

	return cpaper, nil
}
