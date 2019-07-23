package cpaper_asservice

//go:generate make

import (
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/schema"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/service"
	"github.com/s7techlab/cckit/extensions/encryption"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	m "github.com/s7techlab/cckit/state/mapping"
)

func CCRouter(name string) (*router.Group, error) {
	r := router.New(name)
	// Store on the ledger the information about chaincode instantiation
	r.Init(owner.InvokeSetFromCreator)

	if err := service.RegisterCPaperChaincode(r, &CPaperImpl{}); err != nil {
		return nil, err
	}

	return r, nil
}

func NewCC() (*router.Chaincode, error) {
	r, err := CCRouter(`CommercialPaper`)
	if err != nil {
		return nil, err
	}

	return router.NewChaincode(r), nil
}

func NewCCEncrypted() (*router.Chaincode, error) {
	r, err := CCRouter(`CommercialPaperEncrypted`)
	if err != nil {
		return nil, err
	}

	r.
		// encryption key in transient map and encrypted args required
		Pre(encryption.ArgsDecrypt).
		// default Context replaced with EncryptedStateContext only if key is provided in transient map
		Use(encryption.EncStateContext).
		// invoke response will be encrypted cause it will be placed in blocks
		After(encryption.EncryptInvokeResponse())

	return router.NewChaincode(r), nil
}

type CPaperImpl struct {
}

func (cc *CPaperImpl) state(ctx router.Context) m.MappedState {
	return m.WrapState(ctx.State(), m.StateMappings{}.
		//  Create mapping for Commercial Paper entity
		Add(&schema.CommercialPaper{},
			m.PKeySchema(&schema.CommercialPaperId{}), // Key namespace will be <"CommercialPaper", Issuer, PaperNumber>
			m.List(&schema.CommercialPaperList{}),     // Structure of result for List method
			m.UniqKey("ExternalId"),                   // External Id is unique
		))
}

func (cc *CPaperImpl) event(ctx router.Context) state.Event {
	return m.WrapEvent(ctx.Event(), m.EventMappings{}.
		// Event name will be "IssueCommercialPaper", payload - same as issue payload
		Add(&schema.IssueCommercialPaper{}).
		// Event name will be "BuyCommercialPaper"
		Add(&schema.BuyCommercialPaper{}).
		// Event name will be "RedeemCommercialPaper"
		Add(&schema.RedeemCommercialPaper{}))
}

func (cc *CPaperImpl) List(ctx router.Context, in *empty.Empty) (*schema.CommercialPaperList, error) {
	// List method retrieves all entries from the ledger using GetStateByPartialCompositeKey method and passing it the
	// namespace of our contract type, in this example that's "CommercialPaper", then it unmarshals received bytes via
	// proto.Ummarshal method and creates a []schema.CommercialPaperList as defined in the
	// "StateMappings" variable at the top of the file
	if res, err := cc.state(ctx).List(&schema.CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*schema.CommercialPaperList), nil
	}
}

func (cc *CPaperImpl) Get(ctx router.Context, id *schema.CommercialPaperId) (*schema.CommercialPaper, error) {
	if res, err := cc.state(ctx).Get(id, &schema.CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*schema.CommercialPaper), nil
	}
}

func (cc *CPaperImpl) GetByExternalId(ctx router.Context, id *schema.ExternalId) (*schema.CommercialPaper, error) {
	if res, err := cc.state(ctx).GetByUniqKey(
		&schema.CommercialPaper{}, "ExternalId", []string{id.Id}, &schema.CommercialPaper{}); err != nil {
		return nil, err
	} else {
		return res.(*schema.CommercialPaper), nil
	}
}

func (cc *CPaperImpl) Issue(ctx router.Context, issue *schema.IssueCommercialPaper) (*schema.CommercialPaper, error) {
	// Validate input message using the rules defined in schema
	if err := issue.Validate(); err != nil {
		return nil, errors.Wrap(err, "payload validation")
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

	if err := cc.event(ctx).Set(issue); err != nil {
		return nil, err
	}

	if err := cc.state(ctx).Insert(cpaper); err != nil {
		return nil, err
	}
	return cpaper, nil
}

func (cc *CPaperImpl) Buy(ctx router.Context, buy *schema.BuyCommercialPaper) (*schema.CommercialPaper, error) {
	// Get the current commercial paper state
	cpaper, err := cc.Get(ctx, &schema.CommercialPaperId{Issuer: buy.Issuer, PaperNumber: buy.PaperNumber})
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
	if cpaper.State == schema.CommercialPaper_ISSUED {
		cpaper.State = schema.CommercialPaper_TRADING
	}

	// Check paper is not already REDEEMED
	if cpaper.State == schema.CommercialPaper_TRADING {
		cpaper.Owner = buy.NewOwner
	} else {
		return nil, fmt.Errorf(
			"paper %s %s is not trading.current state = %s",
			cpaper.Issuer, cpaper.PaperNumber, cpaper.State)
	}

	if err = cc.event(ctx).Set(buy); err != nil {
		return nil, err
	}

	if err = cc.state(ctx).Put(cpaper); err != nil {
		return nil, err
	}

	return cpaper, nil
}

func (cc *CPaperImpl) Redeem(ctx router.Context, redeem *schema.RedeemCommercialPaper) (*schema.CommercialPaper, error) {
	// Get the current commercial paper state
	cpaper, err := cc.Get(ctx, &schema.CommercialPaperId{Issuer: redeem.Issuer, PaperNumber: redeem.PaperNumber})
	if err != nil {
		return nil, errors.Wrap(err, "get cpaper")
	}
	if err != nil {
		return nil, errors.Wrap(err, "paper not found")
	}

	// Check paper is not REDEEMED
	if cpaper.State == schema.CommercialPaper_REDEEMED {
		return nil, fmt.Errorf("paper %s %s is already redeemed", cpaper.Issuer, cpaper.PaperNumber)
	}

	// Verify that the redeemer owns the commercial paper before redeeming it
	if cpaper.Owner == redeem.RedeemingOwner {
		cpaper.Owner = redeem.Issuer
		cpaper.State = schema.CommercialPaper_REDEEMED
	} else {
		return nil, fmt.Errorf("redeeming owner does not own paper %s %s", cpaper.Issuer, cpaper.PaperNumber)
	}

	if err = cc.event(ctx).Set(redeem); err != nil {
		return nil, err
	}

	if err = cc.state(ctx).Put(cpaper); err != nil {
		return nil, err
	}

	return cpaper, nil
}

func (cc *CPaperImpl) Delete(ctx router.Context, id *schema.CommercialPaperId) (*schema.CommercialPaper, error) {
	cpaper, err := cc.Get(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "get cpaper")
	}

	if err = cc.state(ctx).Delete(id); err != nil {
		return nil, err
	}

	return cpaper, nil
}
