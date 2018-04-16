package identity

import (
	"errors"
)

var (
	// ErrPemEncodedExpected pem format error
	ErrPemEncodedExpected = errors.New("expecting a PEM-encoded X509 certificate; PEM block not found")
)

// Identity interface for invoker (tx creator) and grants, stored in chain code state
type Identity interface {
	GetId() string
	GetMSPId() string
	GetSubject() string
	GetIssuer() string
	Is(i Identity) bool
}
