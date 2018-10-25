package erc20

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/identity"
	r "github.com/s7techlab/cckit/router"
)

var (
	ErrNotEnoughFunds                   = errors.New(`not enough funds`)
	ErrForbiddenToTransferToSameAccount = errors.New(`forbidden to transfer to same account`)
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

	toMspId := c.ArgString(`toMspId`)
	toCertId := c.ArgString(`toCertId`)
	amount := c.ArgInt(`amount`)

	invoker, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	if invoker.GetMSPID() == toMspId && invoker.GetID() == toCertId {
		return nil, ErrForbiddenToTransferToSameAccount
	}

	invokerBalance, err := getBalance(c, invoker.GetMSPID(), invoker.GetID())
	if err != nil {
		return nil, err
	}

	if invokerBalance-amount < 0 {
		return nil, ErrNotEnoughFunds
	}

	recipientBalance, err := getBalance(c, toMspId, toCertId)
	if err != nil {
		return nil, err
	}

	setBalance(c, invoker.GetMSPID(), invoker.GetID(), invokerBalance-amount)
	setBalance(c, toMspId, toCertId, recipientBalance+amount)

	// return current invoker balance
	return invokerBalance - amount, nil
}

// === internal functions, not "public" chaincode functions

func getBalance(c r.Context, mspId, certId string) (int, error) {
	val, err := c.State().Get([]string{mspId, certId}, convert.TypeInt, 0)
	if err != nil {
		return 0, err
	}
	return val.(int), nil
}

func setBalance(c r.Context, mspId, certId string, balance int) error {
	return c.State().Put([]string{mspId, certId}, balance)
}
