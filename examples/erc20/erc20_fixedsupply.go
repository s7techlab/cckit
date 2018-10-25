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
		Init(invokeInitFixedSupply, p.String(`symbol`), p.String(`name`), p.Int(`totalSupply`)).
		Query(`symbol`, querySymbol).                                                                  // Get token symbol
		Query(`name`, queryName).                                                                      // Get token name
		Query(`totalSupply`, queryTotalSupply).                                                        // Get the total token supply
		Query(`balanceOf`, queryBalanceOf, p.String(`mspId`), p.String(`certId`)).                     //  get account balance
		Invoke(`transfer`, invokeTransfer, p.String(`toMspId`), p.String(`toCertId`), p.Int(`amount`)) //Send value amount of tokens tos

	// non implemented yet
	//transferFrom(address _from, address _to, uint256 _value) public returns (bool success)[Send _value amount of tokens from address _from to address _to]
	//approve(address _spender, uint256 _value) public returns (bool success) [Allow _spender to withdraw from your account, multiple times, up to the _value amount. If this function is called again it overwrites the current allowance with _value]
	//allowance(address _owner, address _spender) public view returns (uint256 remaining) [Returns the amount which _spender is still allowed to withdraw from _owner]

	return router.NewChaincode(r)
}

func invokeInitFixedSupply(c router.Context) (interface{}, error) {
	ownerIdentity, err := owner.SetFromCreator(c)
	if err != nil {
		return nil, errors.Wrap(err, `set chaincode owner`)
	}

	if err := c.State().Insert(SymbolKey, c.ArgString(`symbol`)); err != nil {
		return nil, err
	}

	if err := c.State().Insert(NameKey, c.ArgString(`name`)); err != nil {
		return nil, err
	}

	if err := c.State().Insert(TotalSupplyKey, c.ArgInt(`totalSupply`)); err != nil {
		return nil, err
	}
	if err := setBalance(c, ownerIdentity.GetMSPID(), ownerIdentity.GetID(), c.ArgInt(`totalSupply`)); err != nil {
		return nil, errors.Wrap(err, `set owner initial balance`)
	}

	return ownerIdentity, nil
}
