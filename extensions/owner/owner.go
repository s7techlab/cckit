package owner

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	r "github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/extensions/access"
	"github.com/s7techlab/cckit/identity"
)

const DefaultKey = `OWNER`

var (
	ErrToMuchArguments = errors.New(`too much arguments`)
	ErrOnlyByOwner     = errors.New(`chaincode owner required`)
)

// Get returns current owner
func Get(stub shim.ChaincodeStubInterface) ([]byte, error) {
	return stub.GetState(DefaultKey)
}

func FromState(stub shim.ChaincodeStubInterface) (i identity.Identity, err error) {
	owner, err := Get(stub)
	if err != nil {
		return nil, err
	}
	return access.FromBytes(owner)
}

// SetFromCreator chaincode owner from stub creator
func SetFromCreator(c r.Context) peer.Response {
	var grant *access.Grant
	invoker, err := access.InvokerFromStub(c.Stub())
	if err != nil {
		return c.Response().Error(err)
	}

	grant, err = access.GrantFromIdentity(invoker)
	if err != nil {
		return c.Response().Error(err)
	}
	return c.Response().Create(grant, c.State().Put(DefaultKey, grant))
}

func IsOwnerOr(stub shim.ChaincodeStubInterface, allowedTo ...identity.Identity) (bool, error) {

	if isOwner, err := InvokerIsOwner(stub); isOwner || err != nil {
		return isOwner, err
	}
	if len(allowedTo) == 0 {
		return false, nil
	}

	invoker, err := access.InvokerFromStub(stub)
	if err != nil {
		return false, err
	}
	for _, allowed := range allowedTo {
		if allowed.Is(invoker) {
			return true, nil
		}
	}

	return false, nil
}

// IsOwner checks chaincode owner
// Uses current MspID from stub creator if owner isn't presented
func InvokerIsOwner(stub shim.ChaincodeStubInterface) (bool, error) {
	invoker, err := access.InvokerFromStub(stub)
	if err != nil {
		return false, err
	}

	owner, err := FromState(stub)
	if err != nil {
		return false, err
	}

	return invoker.Is(owner), nil
}
