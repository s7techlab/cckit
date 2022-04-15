// Package owner provides method for storing in chaincode state information about chaincode owner
package owner

import (
	"errors"
	"fmt"

	"github.com/s7techlab/cckit/identity"
	r "github.com/s7techlab/cckit/router"
)

var (
	// ErrOwnerNotProvided occurs when providing owner identity in init arguments
	ErrOwnerNotProvided = errors.New(`owner not provided`)

	// ErrOwnerAlreadySet owner already set
	ErrOwnerAlreadySet = errors.New(`owner already set`)
)

func IsSet(c r.Context) (bool, error) {
	return c.State().Exists(OwnerStateKey)
}

// Get returns current chaincode owner identity.Entry
// Service implementation recommended, see chaincode_owner.proto
func Get(c r.Context) (*identity.Entry, error) {
	ownerEntry, err := c.State().Get(OwnerStateKey, &identity.Entry{})
	if err != nil {
		return nil, err
	}

	o := ownerEntry.(identity.Entry)
	return &o, nil
}

// SetFromCreator sets chain code owner from stub creator
// Service implementation recommended, see chaincode_owner.proto
func SetFromCreator(c r.Context) (*identity.Entry, error) {
	if ownerSet, err := IsSet(c); err != nil {
		return nil, err
	} else if ownerSet {
		return Get(c)
	}

	creator, err := identity.FromStub(c.Stub())
	if err != nil {
		return nil, err
	}

	identityEntry, err := identity.CreateEntry(creator)
	if err != nil {
		return nil, err
	}

	return identityEntry, c.State().Insert(OwnerStateKey, identityEntry)
}

// SetFromArgs set owner from first args
func SetFromArgs(c r.Context) (*identity.Entry, error) {
	args := c.Stub().GetArgs()

	if len(args) == 2 {
		return Insert(c, string(args[0]), args[1])
	}

	if isSet, err := IsSet(c); err != nil {
		return nil, err
	} else if !isSet {
		return nil, ErrOwnerNotProvided
	}

	return Get(c)
}

// Insert information about owner to chaincode state
func Insert(c r.Context, mspID string, cert []byte) (*identity.Entry, error) {
	if ownerSet, err := IsSet(c); err != nil {
		return nil, fmt.Errorf(`check owner is set: %w`, err)
	} else if ownerSet {
		return nil, ErrOwnerAlreadySet
	}

	id, err := identity.New(mspID, cert)
	if err != nil {
		return nil, err
	}

	identityEntry, err := identity.CreateEntry(id)
	if err != nil {
		return nil, fmt.Errorf(`create owner entry: %w`, err)
	}
	return identityEntry, c.State().Insert(OwnerStateKey, identityEntry)
}

// IsInvokerOr checks tx creator and compares with owner of another identity
// Service implementation recommended, see chaincode_owner.proto
func IsInvokerOr(c r.Context, allowedTo ...identity.Identity) (bool, error) {
	if err := IsTxCreator(c); err == nil {
		return true, nil
	}
	if len(allowedTo) == 0 {
		return false, nil
	}
	invoker, err := identity.FromStub(c.Stub())
	if err != nil {
		return false, err
	}
	for _, allowed := range allowedTo {
		if allowed.GetMSPIdentifier() == invoker.GetMSPIdentifier() &&
			allowed.GetSubject() == invoker.GetSubject() {
			return true, nil
		}
	}
	return false, nil
}

// IdentityEntryFromState returns identity.Entry with chaincode owner certificate
// Service implementation recommended, see chaincode_owner.proto
func IdentityEntryFromState(c r.Context) (identity.Entry, error) {
	res, err := c.State().Get(OwnerStateKey, &identity.Entry{})
	if err != nil {
		return identity.Entry{}, err
	}

	return res.(identity.Entry), nil
}

// IsInvoker checks than tx creator is chain code owner
// Service implementation recommended, see chaincode_owner.proto
func IsInvoker(ctx r.Context) (bool, error) {
	if err := IsTxCreator(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// IsTxCreator returns error if owner identity  (msp_id + certificate) did not match tx creator identity
// Service implementation recommended, see chaincode_owner.proto
func IsTxCreator(ctx r.Context) error {
	invoker, err := identity.FromStub(ctx.Stub())
	if err != nil {
		return err
	}

	ownerEntry, err := IdentityEntryFromState(ctx)
	if err != nil {
		return err
	}

	if err = identity.Equal(invoker, ownerEntry); err != nil {
		return fmt.Errorf(`%s : %w`, err, ErrOwnerOnly)
	}

	return nil
}
