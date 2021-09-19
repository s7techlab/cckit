// Package owner provides method for storing in chaincode state information about chaincode owner
package owner

import (
	"github.com/pkg/errors"

	"github.com/s7techlab/cckit/identity"
	r "github.com/s7techlab/cckit/router"
)

// OwnerStateKey key used to store owner grant struct in chain code state
const OwnerStateKey = `OWNER`

var (
	// ErrOwnerNotProvided occurs when providing owner identity in init arguments
	ErrOwnerNotProvided = errors.New(`owner not provided`)

	// ErrOwnerAlreadySetted owner already setted
	ErrOwnerAlreadySetted = errors.New(`owner already setted`)

	// ErrMSPIdentifierNotEqual occurs when tx creator and cc owner certificate did not match
	ErrMSPIdentifierNotEqual = errors.New(`msp identifier not equal`)
)

func IsSetted(c r.Context) (bool, error) {
	return c.State().Exists(OwnerStateKey)
}

func Get(c r.Context) (*identity.Entry, error) {
	ownerEntry, err := c.State().Get(OwnerStateKey, &identity.Entry{})
	if err != nil {
		return nil, err
	}

	o := ownerEntry.(identity.Entry)
	return &o, nil
}

// SetFromCreator sets chain code owner from stub creator
func SetFromCreator(c r.Context) (*identity.Entry, error) {
	if ownerSetted, err := IsSetted(c); err != nil {
		return nil, err
	} else if ownerSetted {
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

// SetFromArgs set owner fron first args
func SetFromArgs(c r.Context) (*identity.Entry, error) {
	args := c.Stub().GetArgs()

	if len(args) == 2 {
		return Insert(c, string(args[0]), args[1])
	}

	if isSetted, err := IsSetted(c); err != nil {
		return nil, err
	} else if !isSetted {
		return nil, ErrOwnerNotProvided
	}

	return Get(c)
}

// Insert information about owner to chaincode state
func Insert(c r.Context, mspID string, cert []byte) (*identity.Entry, error) {

	if ownerSetted, err := IsSetted(c); err != nil {
		return nil, errors.Wrap(err, `check owner is set`)
	} else if ownerSetted {
		return nil, ErrOwnerAlreadySetted
	}

	id, err := identity.New(mspID, cert)
	if err != nil {
		return nil, err
	}

	identityEntry, err := identity.CreateEntry(id)
	if err != nil {
		return nil, errors.Wrap(err, `create owner entry`)
	}
	return identityEntry, c.State().Insert(OwnerStateKey, identityEntry)
}

// IsInvokerOr checks tx creator and compares with owner of another identity
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
func IdentityEntryFromState(c r.Context) (identity.Entry, error) {
	res, err := c.State().Get(OwnerStateKey, &identity.Entry{})
	if err != nil {
		return identity.Entry{}, err
	}

	return res.(identity.Entry), nil
}

// Deprecated: IsInvoker checks  than tx creator is chain code owner
// use IsTxCreator
func IsInvoker(ctx r.Context) (bool, error) {
	if err := IsTxCreator(ctx); err != nil {
		return false, err
	}

	return true, nil
}

// IsTxCreator returns error if owner identity  (msp_id + certificate) did not match tx creator identity
func IsTxCreator(ctx r.Context) error {
	invoker, err := identity.FromStub(ctx.Stub())
	if err != nil {
		return err
	}

	ownerEntry, err := IdentityEntryFromState(ctx)
	if err != nil {
		return err
	}

	if ownerEntry.GetMSPID() != invoker.GetMSPIdentifier() {
		return ErrMSPIdentifierNotEqual
	}

	return identity.CertEqual(invoker, ownerEntry)
}
