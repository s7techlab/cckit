package router

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// Chaincode default chaincode implementation with router
type Chaincode struct {
	router *Group
}

// NewChaincode new default chaincode implementation
func NewChaincode(r *Group) *Chaincode {
	return &Chaincode{r}
}

// Init initializes chain code - sets chaincode "owner"
//========  Base methods ====================================
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.HandleInit(stub)
}

// Invoke - entry point for chain code invocations
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}
