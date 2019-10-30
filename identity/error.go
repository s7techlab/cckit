package identity

import (
	"github.com/pkg/errors"
)

var (
	ErrInvalidPEMStructure = errors.New(`invalid pem structure`)

	// ErrPemEncodedExpected pem format error
	ErrPemEncodedExpected = errors.New("expecting a PEM-encoded X509 certificate; PEM block not found")
)
