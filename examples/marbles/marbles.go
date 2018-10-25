// Marbles chaincode example
package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/response"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
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

	r.Query(`owner`, owner.Query) // returns current chaincode owner

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
	return response.Create(owner.SetFromCreator(cc.router.Context(`init`, stub)))
}

// Invoke - entry point for chaincode invocations
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	return cc.router.Handle(stub)
}

// ======= Chain code methods

// marbleOwnerRegister  register a new owner aka end user, save user cert info (Subject, Issuer)
// into chaincode state as serialized Grant struct
func marbleOwnerRegister(c router.Context) (interface{}, error) {
	// receives mspID and certificate as arg in msp.SerializedIdentity type
	// and converts to Identity.Entry
	ownerEntry, err := identity.EntryFromSerialized(c.Arg(`identity`).(msp.SerializedIdentity))
	if err != nil {
		return nil, err
	}

	// put grant struct with owner mspID, as well as cert subject and issuer to state
	return ownerEntry, c.State().Insert(ownerEntry.GetID(), ownerEntry)
}

// marbleInit - create a new marble, store into chaincode state
func marbleInit(c router.Context) (interface{}, error) {

	return nil, nil
}

func marbleDelete(c router.Context) (interface{}, error) {

	return nil, nil
}

func marbleSetOwner(c router.Context) (interface{}, error) {

	return nil, nil
}
