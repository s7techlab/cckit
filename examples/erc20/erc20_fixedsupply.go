package erc20

import (
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

const SymbolKey = `symbol`
const NameKey = `name`
const TotalSupplyKey = `totalSupply`

func NewErc20FixedSupply() *router.Chaincode {
	r := router.New(`erc20fixedSupply`).Use(p.StrictKnown).

		// Chaincode init function, initiates token smart contract with token symbol, name and totalSupply
		Init(invokeInitFixedSupply, p.String(`symbol`), p.String(`name`), p.Int(`totalSupply`)).

		// Get token symbol
		Query(`symbol`, querySymbol).

		// Get token name
		Query(`name`, queryName).

		// Get the total token supply
		Query(`totalSupply`, queryTotalSupply).

		//  get account balance
		Query(`balanceOf`, queryBalanceOf, p.String(`mspId`), p.String(`certId`)).

		//Send value amount of tokens
		Invoke(`transfer`, invokeTransfer, p.String(`toMspId`), p.String(`toCertId`), p.Int(`amount`)).

		// Allow spender to withdraw from your account, multiple times, up to the _value amount.
		// If this function is called again it overwrites the current allowance with _valu
		Invoke(`approve`, invokeApprove, p.String(`spenderMspId`), p.String(`spenderCertId`), p.Int(`amount`)).

		//    Returns the amount which _spender is still allowed to withdraw from _owner]
		Query(`allowance`, queryAllowance, p.String(`ownerMspId`), p.String(`ownerCertId`),
			p.String(`spenderMspId`), p.String(`spenderCertId`)).

		// Send amount of tokens from owner account to another
		Invoke(`transferFrom`, invokeTransferFrom, p.String(`fromMspId`), p.String(`fromCertId`),
			p.String(`toMspId`), p.String(`toCertId`), p.Int(`amount`))

	return router.NewChaincode(r)
}

func invokeInitFixedSupply(c router.Context) (interface{}, error) {
	ownerIdentity, err := owner.SetFromCreator(c)
	if err != nil {
		return nil, errors.Wrap(err, `set chaincode owner`)
	}

	// save token configuration in state
	if err := c.State().Insert(SymbolKey, c.ParamString(`symbol`)); err != nil {
		return nil, err
	}

	if err := c.State().Insert(NameKey, c.ParamString(`name`)); err != nil {
		return nil, err
	}

	if err := c.State().Insert(TotalSupplyKey, c.ParamInt(`totalSupply`)); err != nil {
		return nil, err
	}

	// set token owner initial balance
	if err := setBalance(c, ownerIdentity.GetMSPID(), ownerIdentity.GetID(), c.ParamInt(`totalSupply`)); err != nil {
		return nil, errors.Wrap(err, `set owner initial balance`)
	}

	return ownerIdentity, nil
}
