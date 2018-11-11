# How to create ERC20 token on Hyperledger Fabric

As well as Ethereum blockchain,  Hyperledger Fabric platform (HLF) can be used for token creation, implemented as
smart contract (chaincode in HLF terminology), that holds user balances. Unlike Ethereum, HLF chaincodes can't work
with user addresses as a holder key, thus we will use combination of Membership Service Provider (MSP) Identifier 
and user certificate identifier. Below is an simple example of how to create a token as Golang chaincode on the 
Hyperledger Fabric platform using CCKit chaincode library.

## What is ERC20 token standard

The [ERC20 token standard](https://github.com/ethereum/eips/issues/20) came about as an attempt to standardize token 
smart contracts in Ethereum, it describes the functions  and events that an Ethereum token contract has to implement. 
Most of the major tokens on the Ethereum blockchain  are ERC20-compliant. ERC-20 has many benefits, including unifying 
token wallets and ability for exchanges to list more tokens by providing nothing more than the address of the token’s contract.

```solidity
// ----------------------------------------------------------------------------
// ERC Token Standard #20 Interface
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20-token-standard.md
// ----------------------------------------------------------------------------
contract ERC20Interface {
    function totalSupply() public constant returns (uint);
    function balanceOf(address tokenOwner) public constant returns (uint balance);
    function allowance(address tokenOwner, address spender) public constant returns (uint remaining);
    function transfer(address to, uint tokens) public returns (bool success);
    function approve(address spender, uint tokens) public returns (bool success);
    function transferFrom(address from, address to, uint tokens) public returns (bool success);

    event Transfer(address indexed from, address indexed to, uint tokens);
    event Approval(address indexed tokenOwner, address indexed spender, uint tokens);
}
```

## ERC20 implementation basics

Essentially, an Ethereum token contract is a smart contract that holds a map of account addresses and their balances.
The balance is a value that is defined by the contract creator - in can be fungible physical objects, another monetary value.
The unit of this balance is commonly called a token.

ERC20 functions do:

* `balanceOf` : returns the token balance of an owner identifier (account address in case of Ethereum)

* `transfer` : transfers an amount to an owner identifier of our choosing

* `approve` : sets an amount of tokens a specified owner identifier is allowed to spend on our behalf

* `allowance` : check how much an owner identifier is allowed to spend on our behalf

* `transferFrom` : specify an owner identifier to transfer from if we are allowed by that owner identifier to spend some tokens.


## Owner identifier in Hyperledger Fabric

In the Hyperledger Fabric network, all actors have an identity known to other participants. The default [Membership Service 
Provider](https://hyperledger-fabric.readthedocs.io/en/release-1.3/msp.html) implementation uses X.509 certificates as identities, adopting a traditional Public Key Infrastructure (PKI) 
hierarchical model.

Using information about creator of a proposal and asset ownership the chaincode should be able implement chaincode-level 
access control mechanisms checking  is actor can initiate transactions that update the asset. The corresponding chaincode 
logic has to be able to store this "ownership" information associated with the asset and evaluate it with respect to the 
proposal creator.

In HLF network as unique owner identifier (token balance holder) we can use combination of MSP Identifier and user 
identity identifier. Identity identifier - is concatenation of `Subject` and `Issuer` parts of X.509 certificate. 
This ID is guaranteed to be unique within the MSP.

```go
func (c *clientIdentityImpl) GetID() (string, error) {
	// The leading "x509::" distinquishes this as an X509 certificate, and
	// the subject and issuer DNs uniquely identify the X509 certificate.
	// The resulting ID will remain the same if the certificate is renewed.
	id := fmt.Sprintf("x509::%s::%s", getDN(&c.cert.Subject), getDN(&c.cert.Issuer))
	return base64.StdEncoding.EncodeToString([]byte(id)), nil
}
````
[Client identity chaincode library](https://github.com/hyperledger/fabric/tree/master/core/chaincode/lib/cid) 
allows to write chaincode which makes access control decisions based on the identity of the client 
(i.e. the invoker of the chaincode). 
                                     
In particular, you may make access control decisions based on either or both of the following associated with the client:
                                     
 * the client identity's MSP (Membership Service Provider) ID
 * an attribute associated with the client identity

CCkit contains [identity](https://github.com/s7techlab/cckit/tree/master/identity) package with structures and functions 
can that be used for implementing access control in chaincode. 

## Getting started with example

In our example we use CCKit router for managing smart contract functions. Before you begin, be sure to get `CCkit`:

`git clone git@github.com:s7techlab/cckit.git`

and get dependencies using `dep` command:

`dep ensure -vendor-only`

ERC20 example is located in [examples/erc20](https://github.com/s7techlab/cckit/tree/master/examples/erc20) directory.



## Defining token smart contract functions

First, we need to define chaincode functions. In our example we use [router](https://github.com/s7techlab/cckit/tree/master/router) 
package from CCkit, that allows us to define chaincode methods and their parameters in consistent way. 

At first we define `init` function (smart contract constructor) with arguments `symbol`, `name` and `totalSupply`. 
After that we define chaincode methods, implementing ERC20 interface, adopted to HLF owner identifiers 
(pair of MSP Id and certificate ID). All querying method are prefixed with `query`, all writing to state methods are prefixed with
`invoke`.

As a result we use [default chaincode](https://github.com/s7techlab/cckit/blob/master/router/chaincode.go) structure, 
that delegates `Init` and `Invoke` handling to router.

```go
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
        Invoke(`allowance`, queryAllowance, p.String(`ownerMspId`), p.String(`ownerCertId`),
            p.String(`spenderMspId`), p.String(`spenderCertId`)).
        // Send amount of tokens from owner account to another
        Invoke(`transferFrom`, invokeTransferFrom, p.String(`fromMspId`), p.String(`fromCertId`),
            p.String(`toMspId`), p.String(`toCertId`), p.Int(`amount`))

	return router.NewChaincode(r)
}
````

## Chaincode initialization (constructor)

Chaincode `init` function (token constructor) performs the following actions:

* puts to chaincode state information about chaincode owner, using 
  [owner](https://github.com/s7techlab/cckit/tree/master/extensions/owner) extension from CCkit
* puts to chaincode state token configuration - token symbol, name and total supply
* sets chaincode owner balance with total supply

```go
const SymbolKey = `symbol`
const NameKey = `name`
const TotalSupplyKey = `totalSupply`


func invokeInitFixedSupply(c router.Context) (interface{}, error) {
	ownerIdentity, err := owner.SetFromCreator(c)
	if err != nil {
		return nil, errors.Wrap(err, `set chaincode owner`)
	}

	// save token configuration in state
	if err := c.State().Insert(SymbolKey, c.ArgString(`symbol`)); err != nil {
		return nil, err
	}

	if err := c.State().Insert(NameKey, c.ArgString(`name`)); err != nil {
		return nil, err
	}

	if err := c.State().Insert(TotalSupplyKey, c.ArgInt(`totalSupply`)); err != nil {
		return nil, err
	}

	// set token owner initial balance
	if err := setBalance(c, ownerIdentity.GetMSPID(), ownerIdentity.GetID(), c.ArgInt(`totalSupply`)); err != nil {
		return nil, errors.Wrap(err, `set owner initial balance`)
	}

	return ownerIdentity, nil
}
```


## Defining events structure types

We use [Id](https://github.com/s7techlab/cckit/blob/master/identity/entry.go) structure from 
[identity](https://github.com/s7techlab/cckit/tree/master/identity) package:
```go
// Id structure defines short id representation
type Id struct {
	MSP  string
	Cert string
}
```

And define structures for `Transfer` and `Approve` event:

```go
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
```

## Implementing token smart contract functions

Querying function is quite simple - it's just read value from chaincode state:

```go
const SymbolKey = `symbol`

func querySymbol(c r.Context) (interface{}, error) {
	return c.State().Get(SymbolKey)
}
```

Some of changing state functions are more complicated. For example in function `invokeTransfer` we do:

* receive function invoker certificate (via tx `GetCreator()` function)
* check transfer destination 
* get current invoker (payer) balance
* check balance to transfer `amount` of tokens
* get recipient balance
* update payer and recipient balances in chaincode state

```go
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

// setBalance puts balance value to state
func setBalance(c r.Context, mspId, certId string, balance int) error {
	return c.State().Put(balanceKey(mspId, certId), balance)
}

// balanceKey creates composite key for store balance value in state
func balanceKey(ownerMspId, ownerCertId string) []string {
	return []string{BalancePrefix, ownerMspId, ownerCertId}
}
```


## Testing 

Also, we can fast test our chaincode via CCkit [MockStub](https://github.com/s7techlab/cckit/tree/master/testing).  

To start testing we init chaincode via MockStub with test parameters:

```go
var _ = Describe(`ERC-20`, func() {

	const TokenSymbol = `HLF`
	const TokenName = `HLFCoin`
	const TotalSupply = 10000
	const Decimals = 3

	//Create chaincode mock
	erc20fs := testcc.NewMockStub(`erc20`, NewErc20FixedSupply())

	// load actor certificates
	actors, err := identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
		`token_owner`:     `s7techlab.pem`,
		`account_holder1`: `victor-nosov.pem`,
		//`accoubt_holder2`: `victor-nosov.pem`
	}, examplecert.Content)
	if err != nil {
		panic(err)
	}

	BeforeSuite(func() {
		// init token haincode
		expectcc.ResponseOk(erc20fs.From(actors[`token_owner`]).Init(TokenSymbol, TokenName, TotalSupply, Decimals))
	})
```

After we can check all token operations:

```go
Describe("ERC-20 transfer", func() {
    It("Disallow to transfer token to same account", func() {
        expectcc.ResponseError(
            erc20fs.From(actors[`token_owner`]).Invoke(
                `transfer`, actors[`token_owner`].GetMSPID(), actors[`token_owner`].GetID(), 100),
            ErrForbiddenToTransferToSameAccount)
    })

    It("Disallow token holder with zero balance to transfer tokens", func() {
        expectcc.ResponseError(
            erc20fs.From(actors[`account_holder1`]).Invoke(
                `transfer`, actors[`token_owner`].GetMSPID(), actors[`token_owner`].GetID(), 100),
            ErrNotEnoughFunds)
    })

    It("Allow token holder with non zero balance to transfer tokens", func() {
        expectcc.PayloadInt(
            erc20fs.From(actors[`token_owner`]).Invoke(
                `transfer`, actors[`account_holder1`].GetMSPID(), actors[`account_holder1`].GetID(), 100),
            TotalSupply-100)

        expectcc.PayloadInt(
            erc20fs.Query(
                `balanceOf`, actors[`token_owner`].GetMSPID(), actors[`token_owner`].GetID()), TotalSupply-100)

        expectcc.PayloadInt(
            erc20fs.Query(
                `balanceOf`, actors[`account_holder1`].GetMSPID(), actors[`account_holder1`].GetID()), 100)
    })
})
```
