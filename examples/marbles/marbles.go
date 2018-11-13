// Marbles chaincode example
package main

import (
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/identity"
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

// New chaincode marbles
func New() *router.Chaincode {
	r := router.New(`marbles`) // also initialized logger with "marbles" prefix
	r.Init(owner.InvokeSetFromCreator)
	r.Query(`owner`, owner.Query) // returns current chaincode owner

	r.Group(`marble`).
		// chain code method name is "marbleOwnerRegister"
		Invoke(`OwnerRegister`, marbleOwnerRegister, p.Struct(`identity`, &msp.SerializedIdentity{}), owner.Only).
		Query(`Init`, marbleInit). // chain code method name is "marbleInit"
		Query(`Delete`, marbleDelete, p.String(`id`)).
		Invoke(`SetOwner`, marbleSetOwner, p.String(`id`))
	return router.NewChaincode(r)
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
