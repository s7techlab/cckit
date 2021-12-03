package identity

import (
	"errors"
)

var (
	// ErrMSPIdentifierNotEqual occurs when msp id did not match
	ErrMSPIdentifierNotEqual = errors.New(`msp identifier not equal`)

	ErrSubjectNotEqual = errors.New(`certificate subject not equal`)
	ErrIssuerNotEqual  = errors.New(`certificate issuer not equal`)
)

type (
	SubjectIssuer interface {
		GetSubject() string
		GetIssuer() string
	}

	IdentityAttrs interface {
		SubjectIssuer
		GetMSPIdentifier() string
	}
)

// Equal checks identity attributes (Msp id, cert subject and cert issuer) equality
func Equal(identity1, identity2 IdentityAttrs) error {
	if identity1.GetMSPIdentifier() != identity2.GetMSPIdentifier() {
		return ErrMSPIdentifierNotEqual
	}

	return CertEqual(identity1, identity2)
}

// CertEqual checks certificate equality
func CertEqual(cert1, cert2 SubjectIssuer) error {

	if cert1 == nil || cert2 == nil {
		return errors.New(`certificate empty`)
	}
	if cert1.GetSubject() != cert2.GetSubject() {
		return ErrSubjectNotEqual
	}

	if cert1.GetIssuer() != cert2.GetIssuer() {
		return ErrIssuerNotEqual
	}

	return nil
}
