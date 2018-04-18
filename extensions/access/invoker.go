package access

import (
	"crypto/x509"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/s7techlab/cckit/identity"
)

// Invoker structs holds data of an tx creator
type Invoker struct {
	mspID string
	cert  *x509.Certificate
}

// GetID get id based in certificate subject and issuer
func (i Invoker) GetID() string {
	return identity.IDByCert(i.cert)
}

// GetMSPID returns invoker's membership service provider id
func (i Invoker) GetMSPID() string {
	return i.mspID
}

// GetSubject returns invoker's certificate subject
func (i Invoker) GetSubject() string {
	return identity.GetDN(&i.cert.Subject)
}

// GetIssuer returns invoker's certificate issuer
func (i Invoker) GetIssuer() string {
	return identity.GetDN(&i.cert.Issuer)
}

// Is checks invoker is equal an identity
func (i Invoker) Is(id identity.Identity) bool {
	return i.mspID == id.GetMSPID() && i.GetSubject() == id.GetSubject()
}

// InvokerFromCert creates invoker struct from an mspID and certificate
func InvokerFromCert(mspID string, c []byte) (i identity.Identity, err error) {
	cert, err := identity.Certificate(c)
	if err != nil {
		return nil, err
	}
	return &Invoker{mspID, cert}, nil
}

// InvokerFromStub creates invoker struct from tx creator mspID and certificate
func InvokerFromStub(stub shim.ChaincodeStubInterface) (i identity.Identity, err error) {
	clientIdentity, err := cid.New(stub)
	if err != nil {
		return
	}
	mspID, err := clientIdentity.GetMSPID()
	if err != nil {
		return
	}
	cert, err := clientIdentity.GetX509Certificate()
	if err != nil {
		return
	}
	return &Invoker{mspID, cert}, nil
}
