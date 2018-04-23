package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/extensions/owner"
	r "github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

type Car struct {
	id    string
	title string
}

type Cars struct {
	router *r.Group
}

func New() *Cars {
	cc := &Cars{r.New(`my_chaincode`)}
	cc.router.Group(`car`).
		Query(`List`, cc.carsList).
		Query(`Get`, cc.carGet, p.String(`id`)).
		Invoke(`Put`, cc.carPut, p.Struct(`car`, Car{}), owner.Only)
	return cc
}

//========  Base methods ====================================

func (cc *Cars) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return owner.SetFromCreator(cc.router.Context(`init`, stub))
}

func (cc *Cars) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return cc.router.Handle(stub)
}

// ======= Chaincode methods

func (cc *Cars) carsList(c r.Context) peer.Response {

	return c.Response().Success(nil)
}

func (cc *Cars) carGet(c r.Context) peer.Response {

	return c.Response().Success(nil)
}

func (cc *Cars) carPut(c r.Context) peer.Response {

	return c.Response().Success(nil)
}
