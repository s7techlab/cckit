package erc20

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/identity"
	r "github.com/s7techlab/cckit/router"
)

const (
	BalancePrefix   = `BALANCE`
	AllowancePrefix = `APPROVE`
)

var (
	ErrNotEnoughFunds                   = errors.New(`not enough funds`)
	ErrForbiddenToTransferToSameAccount = errors.New(`forbidden to transfer to same account`)
	ErrSpenderNotHaveAllowance          = errors.New(`spender not have allowance for amount`)
)

type (
	Transfer struct {
		From   identity.Id
		To     identity.Id
		Amount int
	}

	Approve struct {
		From    identity.Id
		Spender identity.Id
		Amount  int
	}
)

func querySymbol(c r.Context) (interface{}, error) {
	return c.State().Get(SymbolKey)
}

func queryName(c r.Context) (interface{}, error) {
	return c.State().Get(NameKey)
}

func queryTotalSupply(c r.Context) (interface{}, error) {
	return c.State().Get(TotalSupplyKey)
}

func queryBalanceOf(c r.Context) (interface{}, error) {
	return getBalance(c, c.ArgString(`mspId`), c.ArgString(`certId`))
}

func invokeTransfer(c r.Context) (interface{}, error) {
	// transfer target
	toMspId := c.ParamString(`toMspId`)
	toCertId := c.ParamString(`toCertId`)

	//transfer amount
	amount := c.ParamInt(`amount`)

	// get information about tx creator
	invoker, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	// Disallow to transfer token to same account
	if invoker.GetMSPID() == toMspId && invoker.GetID() == toCertId {
		return nil, ErrForbiddenToTransferToSameAccount
	}

	// get information about invoker balance from state
	invokerBalance, err := getBalance(c, invoker.GetMSPID(), invoker.GetID())
	if err != nil {
		return nil, err
	}

	// Check the funds sufficiency
	if invokerBalance-amount < 0 {
		return nil, ErrNotEnoughFunds
	}

	// Get information about recipient balance from state
	recipientBalance, err := getBalance(c, toMspId, toCertId)
	if err != nil {
		return nil, err
	}

	// Update payer and recipient balance
	if err = setBalance(c, invoker.GetMSPID(), invoker.GetID(), invokerBalance-amount); err != nil {
		return nil, err
	}

	if err = setBalance(c, toMspId, toCertId, recipientBalance+amount); err != nil {
		return nil, err
	}

	// Trigger event with name "transfer" and payload - serialized to json Transfer structure
	if err = c.SetEvent(`transfer`, &Transfer{
		From: identity.Id{
			MSP:  invoker.GetMSPID(),
			Cert: invoker.GetID(),
		},
		To: identity.Id{
			MSP:  toMspId,
			Cert: toCertId,
		},
		Amount: amount,
	}); err != nil {
		return nil, err
	}

	// return current invoker balance
	return invokerBalance - amount, nil
}

func queryAllowance(c r.Context) (interface{}, error) {
	return getAllowance(c, c.ParamString(`ownerMspId`), c.ParamString(`ownerCertId`), c.ParamString(`spenderMspId`), c.ParamString(`spenderCertId`))
}

func invokeApprove(c r.Context) (interface{}, error) {
	spenderMspId := c.ParamString(`spenderMspId`)
	spenderCertId := c.ParamString(`spenderCertId`)
	amount := c.ParamInt(`amount`)

	invoker, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	if err = setAllowance(c, invoker.GetMSPID(), invoker.GetID(), spenderMspId, spenderCertId, amount); err != nil {
		return nil, err
	}

	if err = c.SetEvent(`approve`, &Approve{
		From: identity.Id{
			MSP:  invoker.GetMSPID(),
			Cert: invoker.GetID(),
		},
		Spender: identity.Id{
			MSP:  spenderMspId,
			Cert: spenderCertId,
		},
		Amount: amount,
	}); err != nil {
		return nil, err
	}

	return true, nil
}

func invokeTransferFrom(c r.Context) (interface{}, error) {

	fromMspId := c.ParamString(`fromMspId`)
	fromCertId := c.ParamString(`fromCertId`)
	toMspId := c.ParamString(`toMspId`)
	toCertId := c.ParamString(`toCertId`)
	amount := c.ParamInt(`amount`)

	invoker, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	// check method invoker has allowances
	allowance, err := getAllowance(c, fromMspId, fromCertId, invoker.GetMSPID(), invoker.GetID())
	if err != nil {
		return nil, err
	}

	// transfer amount must be less or equal allowance
	if allowance < amount {
		return nil, ErrSpenderNotHaveAllowance
	}

	// current payer balance
	balance, err := getBalance(c, fromMspId, fromCertId)
	if err != nil {
		return nil, err
	}

	// payer balance must be greater or equal amount
	if balance-amount < 0 {
		return nil, ErrNotEnoughFunds
	}

	// current recipient balance
	recipientBalance, err := getBalance(c, toMspId, toCertId)
	if err != nil {
		return nil, err
	}

	// decrease payer balance
	if err = setBalance(c, fromMspId, fromCertId, balance-amount); err != nil {
		return nil, err
	}

	// increase recipient balance
	if err = setBalance(c, toMspId, toCertId, recipientBalance+amount); err != nil {
		return nil, err
	}

	// decrease invoker allowance
	if err = setAllowance(c, fromMspId, fromCertId, invoker.GetID(), invoker.GetID(), allowance-amount); err != nil {
		return nil, err
	}

	if err = c.Event().Set(`transfer`, &Transfer{
		From: identity.Id{
			MSP:  fromMspId,
			Cert: fromCertId,
		},
		To: identity.Id{
			MSP:  toMspId,
			Cert: toCertId,
		},
		Amount: amount,
	}); err != nil {
		return nil, err
	}

	// return current invoker balance
	return balance - amount, nil
}

// === internal functions, not "public" chaincode functions

// setBalance puts balance value to state
func balanceKey(ownerMspId, ownerCertId string) []string {
	return []string{BalancePrefix, ownerMspId, ownerCertId}
}

func allowanceKey(ownerMspId, ownerCertId, spenderMspId, spenderCertId string) []string {
	return []string{AllowancePrefix, ownerMspId, ownerCertId, spenderMspId, spenderCertId}
}

func getBalance(c r.Context, mspId, certId string) (int, error) {
	return c.State().GetInt(balanceKey(mspId, certId), 0)
}

// setBalance puts balance value to state
func setBalance(c r.Context, mspId, certId string, balance int) error {
	return c.State().Put(balanceKey(mspId, certId), balance)
}

func getAllowance(c r.Context, ownerMspId, ownerCertId, spenderMspId, spenderCertId string) (int, error) {
	return c.State().GetInt(allowanceKey(ownerMspId, ownerCertId, spenderMspId, spenderCertId), 0)
}

func setAllowance(c r.Context, ownerMspId, ownerCertId, spenderMspId, spenderCertId string, amount int) error {
	return c.State().Put(allowanceKey(ownerMspId, ownerCertId, spenderMspId, spenderCertId), amount)
}
