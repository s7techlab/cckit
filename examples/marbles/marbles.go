package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/extensions/access"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
)

var (
	ErrMarbleOwnerAlreadyRegistered = errors.New(`marble owner already registered`)
)

// Marble represents information about marble
type Marble struct {
	ObjectType string `json:"docType"` //field for couchdb
	ID         string `json:"id"`      //the fieldtags are needed to keep case from bouncing around
	Color      string `json:"color"`
	Size       int    `json:"size"` //size in mm of marble
	Owner      string `json:"owner"`
}

// Chaincode marbles
type Chaincode struct {
	router *router.Group
}

// New chaincode marbles
func New() *Chaincode {
	r := router.New(`marbles`) // also initialized logger with "marbles" prefix

	r.Query(`owner`, owner.Get) // returns current chaincode owner

	r.Group(`marble`).
		// chain code method name is "marbleOwnerRegister"
		Invoke(`OwnerRegister`, marbleOwnerRegister, p.Struct(`identity`, &msp.SerializedIdentity{}), owner.Only).
		Query(`Init`, marbleInit). // chain code method name is "marbleInit"
		Query(`Delete`, marbleDelete, p.String(`id`)).
		Invoke(`SetOwner`, marbleSetOwner, p.String(`id`))
	return &Chaincode{r}
}

//========  Base methods ====================================
//
// Init initializes chaincode
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	// set owner of chain code with special permissions , based on tx creator certificate
	return owner.SetFromCreator(cc.router.Context(`init`, stub))
}

// Invoke - entry point for chaincode invocations
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return cc.router.Handle(stub)
}

// ======= Chain code methods

// ownerInit - register a new owner aka end user, store into chaincode state
func marbleOwnerRegister(c router.Context) peer.Response {
	//mspID and certificate
	ownerIdentity, err := identity.FromSerialized(c.Arg(`identity`).(msp.SerializedIdentity))
	if err != nil {
		return c.Response().Error(err)
	}
	ownerGrant, err := access.GrantFromIdentity(ownerIdentity)
	if err != nil {
		return c.Response().Error(err)
	}

	if exists, err := c.State().Exists(ownerGrant.GetID()); err != nil {
		return c.Response().Error(err)
	} else if exists {
		return c.Response().Error(ErrMarbleOwnerAlreadyRegistered)
	}

	// put grant struct with owner mspID, as well as cert subject and issuer to state
	return c.Response().Create(ownerGrant, c.State().Put(ownerGrant.GetID(), ownerGrant))
}

// marbleInit - create a new marble, store into chaincode state
func marbleInit(c router.Context) peer.Response {

	return c.Response().Success(nil)
}

func marbleDelete(c router.Context) peer.Response {

	return c.Response().Success(nil)
}

func marbleSetOwner(c router.Context) peer.Response {

	return c.Response().Success(nil)
}
