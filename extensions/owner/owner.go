// Package owner provides method for storing in chaincode state information about chaincode owner
package owner

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/extensions/access"
	"github.com/s7techlab/cckit/identity"
	r "github.com/s7techlab/cckit/router"
)

// OwnerStateKey key used to store owner grant struct in chain code state
const OwnerStateKey = `OWNER`

var (
	// ErrToMuchArguments occurs when to much arguments passed to Init
	ErrToMuchArguments = errors.New(`too much arguments`)
)

// IdentityFromState return grant struct representing chain code owner
func IdentityFromState(stub shim.ChaincodeStubInterface) (i identity.Identity, err error) {
	owner, err := stub.GetState(OwnerStateKey)
	if err != nil {
		return nil, err
	}
	return access.FromBytes(owner)
}

// SetFromCreator sets chain code owner from stub creator
func SetFromCreator(c r.Context) peer.Response {
	var grant *access.Grant
	creator, err := identity.FromStub(c.Stub())
	if err != nil {
		return c.Response().Error(err)
	}

	grant, err = access.GrantFromIdentity(creator)
	if err != nil {
		return c.Response().Error(err)
	}
	return c.Response().Create(grant, c.State().Put(OwnerStateKey, grant))
}

// IsInvokerOr checks tx creator and compares with owner of another identity
func IsInvokerOr(stub shim.ChaincodeStubInterface, allowedTo ...identity.Identity) (bool, error) {
	if isOwner, err := IsInvoker(stub); isOwner || err != nil {
		return isOwner, err
	}
	if len(allowedTo) == 0 {
		return false, nil
	}
	creator, err := identity.FromStub(stub)
	if err != nil {
		return false, err
	}
	for _, allowed := range allowedTo {
		if allowed.Is(creator) {
			return true, nil
		}
	}
	return false, nil
}

// IsInvoker checks  than tx creator is chain code owner
func IsInvoker(stub shim.ChaincodeStubInterface) (bool, error) {
	creator, err := identity.FromStub(stub)
	if err != nil {
		return false, err
	}
	owner, err := IdentityFromState(stub)
	if err != nil {
		return false, err
	}
	return creator.Is(owner), nil
}
