# Hyperledger Fabric chaincode kit (CCKit)

## Chaincode methods router

To simplify chaincode development, we tried to compose common software development patterns, such as routing, middleware 
and invoke context:

* Routing refers to determining how an application responds to a client request to a particular endpoint. 
  Chaincode router uses rules about how to map chaincode invocation to particular handler, 
  as well as what kind of middleware need to be used during request, for example how to convert incoming argument from
  []byte to target type (string, struct etc)

* Invoke context is abstraction over ChaincodeStubInterface, represents the context of the current chaincode invocation. 
  It holds request, response and client (Identity) reference, converted parameters, as well as state and log reference. 
  As `Context` is an interface, it is easy to extend it with custom methods.

* Middleware functions are functions that have access to the invoke context, invoke result and the next middleware function 
  in the chaincode’s invoke-response cycle. The next middleware function is commonly denoted by a variable named next.


### Middleware 

Middleware functions can perform the following tasks:

* Convert input args from byte slice to desired type
* Check access control requirements
* End the request-response cycle.
* Call the next middleware function in the stack.


## Defining chaincode function and their arguments

### Delegating chaincode methods handling

[router/context.go](context.go)

```go
// Chaincode default chaincode implementation with router
type Chaincode struct {
	router *Group
}


// Init initializes chain code - sets chaincode "owner"
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.HandleInit(stub)
}

// Invoke - entry point for chain code invocations
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}

```

### Chaincode function and their arguments

Example from [ERC20](../examples/erc20) chaincode:

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
		// If this function is called again it overwrites the current allowance with _value
		Invoke(`approve`, invokeApprove, p.String(`spenderMspId`), p.String(`spenderCertId`), p.Int(`amount`)).

		//    Returns the amount which _spender is still allowed to withdraw from _owner
		Invoke(`allowance`, queryAllowance, p.String(`ownerMspId`), p.String(`ownerCertId`),
			p.String(`spenderMspId`), p.String(`spenderCertId`)).

		// Send amount of tokens from owner account to another
		Invoke(`transferFrom`, invokeTransferFrom, p.String(`fromMspId`), p.String(`fromCertId`),
			p.String(`toMspId`), p.String(`toCertId`), p.Int(`amount`))

	return router.NewChaincode(r)
}
```