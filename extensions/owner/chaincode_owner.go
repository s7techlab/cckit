package owner

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ErrTxInvokerIsNotOwner         = errors.New(`tx invoker is not owner`)
	ErrDeleteLastOwnerIsNotAllowed = errors.New(`delete last owner is not allowed`)
	ErrNewCertSameAsOldCert        = errors.New(`new cert same as old cert`)
)

func (x *ChaincodeOwner) GetMSPIdentifier() string {
	return x.MspId
}

func NewService() *ChaincodeOwnerService {
	return &ChaincodeOwnerService{}
}

type ChaincodeOwnerService struct {
}

// IsTxCreator - wrapper for TxCreatorIsOwner for local calls
func (c *ChaincodeOwnerService) IsTxCreator(ctx router.Context) (*ChaincodeOwner, error) {
	return c.TxCreatorIsOwner(ctx, &empty.Empty{})
}

func (c *ChaincodeOwnerService) TxCreatorIsOwner(ctx router.Context, _ *empty.Empty) (*ChaincodeOwner, error) {
	txCreator, err := identity.FromStub(ctx.Stub())
	if err != nil {
		return nil, err
	}

	owner, err := c.OwnerGet(ctx, &OwnerId{
		MspId:   txCreator.GetMSPIdentifier(),
		Subject: txCreator.GetSubject(),
	})

	if err != nil {
		return nil, fmt.Errorf(`find owner by tx creator's msp_id and cert subject: %w`, err)
	}

	if err := identity.Equal(txCreator, owner); err != nil {
		return nil, fmt.Errorf(`owner with tx creator's' msp_id and cert subject found, but: %w`, err)
	}

	return owner, nil
}

func (c *ChaincodeOwnerService) OwnersList(ctx router.Context, _ *empty.Empty) (*ChaincodeOwners, error) {
	if res, err := State(ctx).List(&ChaincodeOwner{}); err != nil {
		return nil, err
	} else {
		return res.(*ChaincodeOwners), nil
	}
}

func (c *ChaincodeOwnerService) OwnerGet(ctx router.Context, id *OwnerId) (*ChaincodeOwner, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	if res, err := State(ctx).Get(id, &ChaincodeOwner{}); err != nil {
		return nil, err
	} else {
		return res.(*ChaincodeOwner), nil
	}
}

func (c *ChaincodeOwnerService) allowToModifyBy(ctx router.Context, invoker identity.Identity) error {
	currentOwners, err := c.OwnersList(ctx, &empty.Empty{})
	if err != nil {
		return err
	}

	// no owners, allow to register
	if len(currentOwners.Items) == 0 {
		return nil
	}

	for _, owner := range currentOwners.Items {
		if err = identity.Equal(owner, invoker); err == nil {
			return nil
		}
	}

	return ErrTxInvokerIsNotOwner
}

func (c *ChaincodeOwnerService) txCreatorAllowedToModify(ctx router.Context) (identity.Identity, error) {
	txCreator, err := identity.FromStub(ctx.Stub())
	if err != nil {
		return nil, err
	}

	return txCreator, c.allowToModifyBy(ctx, txCreator)
}

func (c *ChaincodeOwnerService) OwnerRegister(ctx router.Context, registerRequest *OwnerRegisterRequest) (*ChaincodeOwner, error) {
	if err := router.ValidateRequest(registerRequest); err != nil {
		return nil, err
	}

	txCreator, err := c.txCreatorAllowedToModify(ctx)
	if err != nil {
		return nil, err
	}

	id, err := identity.New(registerRequest.MspId, registerRequest.Cert)
	if err != nil {
		return nil, fmt.Errorf(`parse certificate: %w`, err)
	}

	txTimestamp, _ := ctx.Stub().GetTxTimestamp()
	chaincodeOwner := &ChaincodeOwner{
		MspId:   id.GetMSPIdentifier(),
		Subject: id.GetSubject(),

		Issuer:         id.GetIssuer(),
		ExpiresAt:      timestamppb.New(id.ExpiresAt()),
		Cert:           registerRequest.Cert,
		UpdatedByMspId: txCreator.GetMSPIdentifier(),
		UpdatedByCert:  txCreator.GetPEM(),
		UpdatedAt:      txTimestamp,
	}

	if err = State(ctx).Insert(chaincodeOwner); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&ChaincodeOwnerRegistered{
		MspId:     chaincodeOwner.MspId,
		Subject:   chaincodeOwner.Subject,
		Issuer:    chaincodeOwner.Issuer,
		ExpiresAt: chaincodeOwner.ExpiresAt,
	}); err != nil {
		return nil, err
	}

	return chaincodeOwner, nil
}

func (c ChaincodeOwnerService) OwnerUpdate(ctx router.Context, updateRequest *OwnerUpdateRequest) (*ChaincodeOwner, error) {
	if err := router.ValidateRequest(updateRequest); err != nil {
		return nil, err
	}

	txCreator, err := c.txCreatorAllowedToModify(ctx)
	if err != nil {
		return nil, err
	}

	id, err := identity.New(updateRequest.MspId, updateRequest.Cert)
	if err != nil {
		return nil, fmt.Errorf(`parse certificate: %w`, err)
	}

	curOwner, err := c.OwnerGet(ctx, &OwnerId{
		MspId:   id.GetMSPIdentifier(),
		Subject: id.GetSubject(),
	})

	if err != nil {
		return nil, fmt.Errorf(`current owner with equal msp_id and cert_subject: %w`, err)
	}

	if bytes.Equal(curOwner.Cert, updateRequest.Cert) {
		return nil, ErrNewCertSameAsOldCert
	}

	if err = identity.Equal(curOwner, id); err != nil {
		return nil, err
	}

	txTimestamp, _ := ctx.Stub().GetTxTimestamp()
	chaincodeOwner := &ChaincodeOwner{
		MspId:   id.GetMSPIdentifier(),
		Subject: id.GetSubject(),

		Issuer:         id.GetIssuer(),
		ExpiresAt:      timestamppb.New(id.ExpiresAt()),
		Cert:           updateRequest.Cert,
		UpdatedByMspId: txCreator.GetMSPIdentifier(),
		UpdatedByCert:  txCreator.GetPEM(),
		UpdatedAt:      txTimestamp,
	}

	if err = State(ctx).Put(chaincodeOwner); err != nil {
		return nil, err
	}

	if err = Event(ctx).Set(&ChaincodeOwnerUpdated{
		MspId:     chaincodeOwner.MspId,
		Subject:   chaincodeOwner.Subject,
		ExpiresAt: chaincodeOwner.ExpiresAt,
	}); err != nil {
		return nil, err
	}

	return chaincodeOwner, nil
}

func (c ChaincodeOwnerService) OwnerDelete(ctx router.Context, id *OwnerId) (*ChaincodeOwner, error) {
	if err := router.ValidateRequest(id); err != nil {
		return nil, err
	}

	if _, err := c.txCreatorAllowedToModify(ctx); err != nil {
		return nil, err
	}

	deletedOwner, err := c.OwnerGet(ctx, id)
	if err != nil {
		return nil, err
	}

	currentOwners, err := c.OwnersList(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	if len(currentOwners.Items) == 1 {
		return nil, ErrDeleteLastOwnerIsNotAllowed
	}

	if err := State(ctx).Delete(id); err != nil {
		return nil, err
	}

	if err := Event(ctx).Set(&ChaincodeOwnerDeleted{
		MspId:   id.MspId,
		Subject: id.Subject,
	}); err != nil {
		return nil, err
	}

	return deletedOwner, nil
}
