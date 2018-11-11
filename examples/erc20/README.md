# How to create ERC20 token on Hyperledger Fabric

As well as Ethereum blockchain,  Hyperledger Fabric platform (HLF) can be used for token creation, implemented as
smart contract (chaincode in HLF terminology), that holds user balances. Unlike Ethereum, HLF chaincodes can't work
with user addresses as a holder key, thus we will use combination of MSP Identifier and user certificate identifier.
Below in an simple example of how to create a token as Golang chaincode on the Hyperledger Fabric platform 
using CCKit chaincode library.

## What is ERC20

The ERC20 came about as an attempt to standardize token smart contracts in Ethereum, it describes the functions 
and events that an Ethereum token contract has to implement. Most of the major tokens on the Ethereum blockchain 
are ERC20-compliant. ERC-20 has many benefits, including unifying token wallets and ability for exchanges to list
more tokens by providing nothing more than the address of the token’s contract.

```
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

In the Hyperledger Fabric network, all actors have identity known to other participants. The default Membership Service 
Provider implementation uses X.509 certificates as identities, adopting a traditional Public Key Infrastructure (PKI) 
hierarchical model.

Using information about creator of a proposal and asset ownership the chaincode should be able implement chaincode-level 
access control mechanisms checking  is actor can initiate transactions that update the asset. The corresponding access control 
check needs to be able to store this "ownership" information associated with the asset and evaluate it with respect to the 
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

Certificate identifier introduced in the [client identity chaincode library](https://github.com/hyperledger/fabric/tree/master/core/chaincode/lib/cid) 
, it enables you to write chaincode which makes access control decisions based on the identity of the client 
(i.e. the invoker of the chaincode). 
                                     
In particular, you may make access control decisions based on either or both of the following associated with the client:
                                     
 * the client identity's MSP (Membership Service Provider) ID
 * an attribute associated with the client identity

CCkit contains [identity](https://github.com/s7techlab/cckit/tree/master/identity) package with structures and functions 
can be used for implementing access control in chaincode. 

## Getting started with example

In our example we use CCKit router for managing smart contract functions. Before you begin, be sure to get `CCkit`:

`git clone git@github.com:s7techlab/cckit.git`

and get dependencies using `dep` command:

`dep ensure -vendor-only`

ERC20 example is located in [examples/erc20](https://github.com/s7techlab/cckit/tree/master/examples/erc20) directory.

## Testing 



