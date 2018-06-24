// Package owner provides method for storing in chaincode state information about chaincode owner
package owner

import (
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
	"github.com/s7techlab/cckit/identity"
	r "github.com/s7techlab/cckit/router"
)

// OwnerStateKey key used to store owner grant struct in chain code state
const OwnerStateKey = `OWNER`

var (
	// ErrToMuchArguments occurs when to much arguments passed to Init
	ErrToMuchArguments = errors.New(`too much arguments`)
)

// SetFromCreator sets chain code owner from stub creator
func SetFromCreator(c r.Context) peer.Response {
	ownerSetted, err := c.State().Exists(OwnerStateKey)
	if err != nil {
		return c.Response().Error(err)
	}

	if ownerSetted {
		return c.Response().Create(c.State().Get(OwnerStateKey, &identity.Entry{}))
	}

	creator, err := identity.FromStub(c.Stub())
	if err != nil {
		return c.Response().Error(err)
	}

	identityEntry, err := identity.CreateEntry(creator)
	if err != nil {
		return c.Response().Error(err)
	}
	return c.Response().Create(identityEntry, c.State().Insert(OwnerStateKey, identityEntry))
}

// IsInvokerOr checks tx creator and compares with owner of another identity
func IsInvokerOr(c r.Context, allowedTo ...identity.Identity) (bool, error) {
	if isOwner, err := IsInvoker(c); isOwner || err != nil {
		return isOwner, err
	}
	if len(allowedTo) == 0 {
		return false, nil
	}
	invoker, err := identity.FromStub(c.Stub())
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

// IdentityFromState
func IdentityEntryFromState(c r.Context) (interface{}, error) {
	return c.State().Get(OwnerStateKey, &identity.Entry{})
}

// IsInvoker checks  than tx creator is chain code owner
func IsInvoker(c r.Context) (bool, error) {
	invoker, err := identity.FromStub(c.Stub())
	if err != nil {
		return false, err
	}
	owner, err := IdentityEntryFromState(c)
	if err != nil {
		return false, err
	}
	return invoker.Is(owner.(identity.Entry)), nil
}
