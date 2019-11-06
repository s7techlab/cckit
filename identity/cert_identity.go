package identity

import (
	"crypto/x509"
	"encoding/pem"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/msp"
	protomsp "github.com/hyperledger/fabric/protos/msp"
	"github.com/pkg/errors"
)

// New creates CertIdentity struct from an mspID and certificate
func New(mspID string, certPEM []byte) (ci *CertIdentity, err error) {
	cert, err := Certificate(certPEM)
	if err != nil {
		return nil, err
	}
	return &CertIdentity{mspID, cert}, nil
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
func FromSerialized(s protomsp.SerializedIdentity) (ci *CertIdentity, err error) {
	return New(s.Mspid, s.IdBytes)
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

// Deprecated: GetMSPID returns invoker's membership service provider id
func (ci CertIdentity) GetMSPID() string {
	return ci.MspID
}

func (ci CertIdentity) ExpiresAt() time.Time {
	return ci.Cert.NotAfter
}

func (ci CertIdentity) GetMSPIdentifier() string {
	return ci.MspID
}

func (ci CertIdentity) GetIdentifier() *msp.IdentityIdentifier {
	return &msp.IdentityIdentifier{
		Mspid: ci.MspID,
		Id:    ci.GetID(),
	}
}

func (ci CertIdentity) Validate() error {
	return nil
}

func (ci CertIdentity) Verify(msg []byte, sig []byte) error {
	return nil
}

func (ci CertIdentity) Anonymous() bool {
	return false
}

func (ci CertIdentity) GetOrganizationalUnits() []*msp.OUIdentifier {
	return nil
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

// GetPEM certificate encoded to PEM
func (ci CertIdentity) GetPEM() []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  `CERTIFICATE`,
		Bytes: ci.Cert.Raw,
	})
}

// ToSerialized converts CertIdentity to *msp.SerializedIdentity
func (ci CertIdentity) ToSerialized() *protomsp.SerializedIdentity {
	return &protomsp.SerializedIdentity{
		Mspid:   ci.MspID,
		IdBytes: ci.GetPEM(),
	}
}

func (ci CertIdentity) Serialize() ([]byte, error) {
	return ci.ToBytes()
}

// ToBytes converts to serializedIdentity and then to json
func (ci CertIdentity) ToBytes() ([]byte, error) {
	return proto.Marshal(ci.ToSerialized())
}

func (ci CertIdentity) SatisfiesPrincipal(principal *protomsp.MSPPrincipal) error {
	return nil
}

func (ci CertIdentity) Sign(msg []byte) ([]byte, error) {
	return nil, nil
}

func (ci CertIdentity) GetPublicVersion() msp.Identity {
	return nil
}
