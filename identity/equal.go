package identity

import (
	"errors"
)

var (
	ErrSubjectNotEqual = errors.New(`certificate subject not equal`)
	ErrIssuerNotEqual  = errors.New(`certificate issuer not equal`)
)

type (
	SubjectIssuer interface {
		GetSubject() string
		GetIssuer() string
	}
)

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
