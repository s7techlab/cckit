package identity

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"

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
	// GetID Identifier, based on Subject and Issuer
	GetID() string

	// GetMSPID msp identifier
	GetMSPID() string

	//  GetSubject string representation of X.509 cert subject
	GetSubject() string
	//  GetIssuer string representation of X.509 cert issuer
	GetIssuer() string

	// GetPublicKey *rsa.PublicKey or *dsa.PublicKey or *ecdsa.PublicKey:
	GetPublicKey() interface{}
	GetPEM() []byte
	Is(i Identity) bool
}

type GetContent func(string) ([]byte, error)

// New creates CertIdentity struct from an mspID and certificate
func New(mspID string, certPEM []byte) (ci *CertIdentity, err error) {
	cert, err := Certificate(certPEM)
	if err != nil {
		return nil, err
	}
	return &CertIdentity{mspID, cert}, nil
}

// FromFile creates certIdentity from file determined by filename
func FromFile(mspID string, filename string, getContent GetContent) (ci *CertIdentity, err error) {
	PEM, err := getContent(filename)
	if err != nil {
		return nil, err
	}
	return New(mspID, PEM)
}

// FromStub creates Identity interface  from tx creator mspID and certificate (stub.GetCreator)
func FromStub(stub shim.ChaincodeStubInterface) (ci *CertIdentity, err error) {
	clientIdentity, err := cid.New(stub)
	if err != nil {
		return nil, errors.Wrap(err, `client identity from stub`)
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
func FromSerialized(s msp.SerializedIdentity) (ci *CertIdentity, err error) {
	return New(s.Mspid, s.IdBytes)
}

// EntryFromSerialized creates Entry from SerializedEntry
func EntryFromSerialized(s msp.SerializedIdentity) (g *Entry, err error) {
	id, err := FromSerialized(s)
	if err != nil {
		return nil, err
	}
	return CreateEntry(id)
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

func (ci CertIdentity) GetPublicKey() interface{} {
	return ci.Cert.PublicKey
}

// Is checks invoker is equal an identity
func (ci CertIdentity) Is(id Identity) bool {
	return ci.MspID == id.GetMSPID() && ci.GetSubject() == id.GetSubject()
}

// GetPEM certificate encoded to PEM
func (ci CertIdentity) GetPEM() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  `CERTIFICATE`,
		Bytes: ci.Cert.Raw,
	})
}

// ToSerialized converts CertIdentity to *msp.SerializedIdentity
func (ci CertIdentity) ToSerialized() *msp.SerializedIdentity {
	return &msp.SerializedIdentity{
		Mspid:   ci.MspID,
		IdBytes: ci.GetPEM(),
	}
}

// ToBytes converts to serializedIdentity and then to json
func (ci CertIdentity) ToBytes() ([]byte, error) {
	return json.Marshal(ci.ToSerialized())
}
