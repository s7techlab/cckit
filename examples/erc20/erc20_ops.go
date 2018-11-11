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
	toMspId := c.ArgString(`toMspId`)
	toCertId := c.ArgString(`toCertId`)

	//transfer amount
	amount := c.ArgInt(`amount`)

	// get informartion about tx creator
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
	setBalance(c, invoker.GetMSPID(), invoker.GetID(), invokerBalance-amount)
	setBalance(c, toMspId, toCertId, recipientBalance+amount)

	// Trigger event with name "transfer" and payload - serialized to json Transfer structure
	c.SetEvent(`transfer`, &Transfer{
		From: identity.Id{
			MSP:  invoker.GetMSPID(),
			Cert: invoker.GetID(),
		},
		To: identity.Id{
			MSP:  toMspId,
			Cert: toCertId,
		},
		Amount: amount,
	})

	// return current invoker balance
	return invokerBalance - amount, nil
}

func queryAllowance(c r.Context) (interface{}, error) {
	return getAllowance(c, c.ArgString(`ownerMspId`), c.ArgString(`ownerCertId`), c.ArgString(`spenderMspId`), c.ArgString(`spenderCertId`))
}

func invokeApprove(c r.Context) (interface{}, error) {
	spenderMspId := c.ArgString(`spenderMspId`)
	spenderCertId := c.ArgString(`spenderCertId`)
	amount := c.ArgInt(`amount`)

	invoker, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	setAllowance(c, invoker.GetMSPID(), invoker.GetID(), spenderMspId, spenderCertId, amount)

	c.SetEvent(`approve`, &Approve{
		From: identity.Id{
			MSP:  invoker.GetMSPID(),
			Cert: invoker.GetID(),
		},
		Spender: identity.Id{
			MSP:  spenderMspId,
			Cert: spenderCertId,
		},
		Amount: amount,
	})

	return true, nil
}

func invokeTransferFrom(c r.Context) (interface{}, error) {

	fromMspId := c.ArgString(`fromMspId`)
	fromCertId := c.ArgString(`fromCertId`)
	toMspId := c.ArgString(`toMspId`)
	toCertId := c.ArgString(`toCertId`)
	amount := c.ArgInt(`amount`)

	invoker, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	// check method invoker has allowances
	allowance, err := getAllowance(c, fromMspId, fromCertId, invoker.GetMSPID(), invoker.GetID())
	if err != nil {
		return nil, err
	}

	if allowance <= amount {
		return nil, ErrSpenderNotHaveAllowance
	}

	balance, err := getBalance(c, fromMspId, fromCertId)
	if err != nil {
		return nil, err
	}

	if balance-amount < 0 {
		return nil, ErrNotEnoughFunds
	}

	recipientBalance, err := getBalance(c, toMspId, toCertId)
	if err != nil {
		return nil, err
	}

	setBalance(c, fromMspId, fromCertId, balance-amount)
	setBalance(c, toMspId, toCertId, recipientBalance+amount)
	setAllowance(c, fromMspId, fromCertId, invoker.GetID(), invoker.GetID(), allowance-amount)

	c.SetEvent(`transfer`, &Transfer{
		From: identity.Id{
			MSP:  fromMspId,
			Cert: fromCertId,
		},
		To: identity.Id{
			MSP:  toMspId,
			Cert: toCertId,
		},
		Amount: amount,
	})

	// return current invoker balance
	return balance - amount, nil
}

// === internal functions, not "public" chaincode functions

// setBalance puts balance value to state
func balanceKey(ownerMspId, ownerCertId string) []string {
	return []string{BalancePrefix, ownerMspId, ownerCertId}
}

func allowanceKey(ownerMspId, owneCertId, spenderMspId, spenderCertId string) []string {
	return []string{AllowancePrefix, ownerMspId, owneCertId, spenderMspId, spenderCertId}
}

func getBalance(c r.Context, mspId, certId string) (int, error) {
	return c.State().GetInt(balanceKey(mspId, certId), 0)
}

// setBalance puts balance value to state
func setBalance(c r.Context, mspId, certId string, balance int) error {
	return c.State().Put(balanceKey(mspId, certId), balance)
}

func getAllowance(c r.Context, ownerMspId, owneCertId, spenderMspId, spenderCertId string) (int, error) {
	return c.State().GetInt(allowanceKey(ownerMspId, owneCertId, spenderMspId, spenderCertId), 0)
}

func setAllowance(c r.Context, ownerMspId, owneCertId, spenderMspId, spenderCertId string, amount int) error {
	return c.State().Put(allowanceKey(ownerMspId, owneCertId, spenderMspId, spenderCertId), amount)
}
