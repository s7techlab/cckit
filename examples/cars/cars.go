package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/extensions/owner"
	r "github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

type Car struct {
	Id    string
	Title string
	Owner string
}

type Cars struct {
	router *r.Group
}

func New() *Cars {
	cc := &Cars{r.New(`cars`)} // also initialized logger with "cars" prefix
	cc.router.Group(`car`).
		Query(`List`, cc.carList).               // chain code method name is carList
		Query(`Get`, cc.carGet, p.String(`id`)). // chain code method name is carGet
		Invoke(`Put`, cc.carRegister, p.Struct(`car`, &Car{}), owner.Only)
	return cc
}

//========  Base methods ====================================
//
// Init initializes chain code
func (cc *Cars) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// set owner of chain code with special permissions , based on tx creator certificate
	return owner.SetFromCreator(cc.router.Context(`init`, stub))
}

// Invoke - entry point for chain code invocations
func (cc *Cars) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return cc.router.Handle(stub)
}

// ======= Chaincode methods

func (cc *Cars) carList(c r.Context) peer.Response {

	return c.Response().Success(nil)
}

func (cc *Cars) carGet(c r.Context) peer.Response {

	return c.Response().Success(nil)
}

func (cc *Cars) carRegister(c r.Context) peer.Response {

	return c.Response().Success(nil)
}

func (cc *Cars) carNewOwner(c r.Context) peer.Response {

	return c.Response().Success(nil)
}
