package identity

import (
	"errors"
)

var (
	ErrPemEncodedExpected = errors.New("expecting a PEM-encoded X509 certificate; PEM block not found")
)

type Identity interface {
	GetId() string
	GetMSPId() string
	GetSubject() string
	GetIssuer() string
	Is(i Identity) bool
}
