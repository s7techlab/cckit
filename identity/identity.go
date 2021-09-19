package identity

import (
	"github.com/hyperledger/fabric/msp"
)

// Identity interface for invoker (tx creator) and grants, stored in chain code state
type Identity interface {
	msp.Identity

	// Deprecated: GetID Identifier, based on Subject and Issuer. Use GetIdentifier instead
	GetID() string

	// Deprecated: GetMSPID msp identifier. Use GetMspIdentifier instead
	GetMSPID() string

	// GetSubject string representation of X.509 cert subject
	GetSubject() string
	// GetIssuer string representation of X.509 cert issuer
	GetIssuer() string

	// GetPublicKey *rsa.PublicKey or *dsa.PublicKey or *ecdsa.PublicKey:
	GetPublicKey() interface{}
	GetPEM() []byte
}
