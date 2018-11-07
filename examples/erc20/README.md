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

Essentially, an Ethereum token contract is a smart contract that holds a map of account addresses and their balances
The balance is a value that is defined by the contract creator - in can be fungible physical objects, another monetary value.
The unit of this balance is commonly called a token.


