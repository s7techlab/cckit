package identity

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
)

var (
	// ErrPemEncodedExpected pem format error
	ErrPemEncodedExpected = errors.New("expecting a PEM-encoded X509 certificate; PEM block not found")
)

// Identity interface for invoker (tx creator) and grants, stored in chain code state
type Identity interface {
	GetID() string
	GetMSPID() string
	GetSubject() string
	GetIssuer() string
	Is(i Identity) bool
}

// FromCert creates CertIdentity struct from an mspID and certificate
func FromCert(mspID string, c []byte) (ci *CertIdentity, err error) {
	cert, err := Certificate(c)
	if err != nil {
		return nil, err
	}
	return &CertIdentity{mspID, cert}, nil
}

// InvokerFromStub creates Identity interface  from tx creator mspID and certificate (stub.GetCreator)
func FromStub(stub shim.ChaincodeStubInterface) (i Identity, err error) {
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
	return &CertIdentity{mspID, cert}, nil
}

// FromSerialized converts  msp.SerializedIdentity struct  to Identity interface{}
func FromSerialized(s msp.SerializedIdentity) (i Identity, err error) {
	return FromCert(s.Mspid, s.IdBytes)
}

// CertIdentity  structs holds data of an tx creator
type CertIdentity struct {
	MspID string
	Cert  *x509.Certificate
}

// GetID get id based in certificate subject and issuer
func (ci CertIdentity) GetID() string {
	return IDByCert(ci.Cert)
}

// GetMSPID returns invoker's membership service provider id
func (ci CertIdentity) GetMSPID() string {
	return ci.MspID
}

// GetSubject returns invoker's certificate subject
func (ci CertIdentity) GetSubject() string {
	return GetDN(&ci.Cert.Subject)
}

// GetIssuer returns invoker's certificate issuer
func (ci CertIdentity) GetIssuer() string {
	return GetDN(&ci.Cert.Issuer)
}

// Is checks invoker is equal an identity
func (ci CertIdentity) Is(id Identity) bool {
	return ci.MspID == id.GetMSPID() && ci.GetSubject() == id.GetSubject()
}

func (ci CertIdentity) PemEncode() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  `CERTIFICATE`,
		Bytes: ci.Cert.Raw,
	})
}

// ToSerialized converts CertIdentity to *msp.SerializedIdentity
func (ci CertIdentity) ToSerialized() *msp.SerializedIdentity {
	return &msp.SerializedIdentity{
		Mspid:   ci.MspID,
		IdBytes: ci.PemEncode(),
	}
}

// ToBytes converts to serializedIdentity and then to json
func (ci CertIdentity) ToBytes() ([]byte, error) {
	return json.Marshal(ci.ToSerialized())
}
