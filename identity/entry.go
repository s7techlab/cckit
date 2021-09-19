// Package access contains structs for storing chaincode access control information
package identity

import (
	"crypto/x509"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	protomsp "github.com/hyperledger/fabric-protos-go/msp"
)

// Entry structure for storing identity information
// string representation certificate Subject and Issuer can be used for reach query searching
type Entry struct {
	MSPId   string
	Subject string
	Issuer  string
	PEM     []byte
	Cert    *x509.Certificate `json:"-"` // temporary cert
}

// Id structure defines short id representation
type Id struct {
	MSP  string
	Cert string
}

// IdentityEntry interface
type IdentityEntry interface {
	GetIdentityEntry() Entry
}

// ========  Identity interface ===================

// GetID identifier by certificate subject and issuer
func (e Entry) GetID() string {
	return ID(e.Subject, e.Issuer)
}

// GetMSPID membership service provider identifier
func (e Entry) GetMSPID() string {
	return e.MSPId
}

// GetSubject certificate subject
func (e Entry) GetSubject() string {
	return e.Subject
}

// GetIssuer certificate issuer
func (e Entry) GetIssuer() string {
	return e.Issuer
}

// GetPK certificate issuer
func (e Entry) GetPEM() []byte {
	return e.PEM
}

func (e Entry) GetPublicKey() interface{} {

	if e.Cert == nil {
		cert, err := Certificate(e.PEM)

		if err != nil {
			return err
		}

		e.Cert = cert
	}

	return e.Cert.PublicKey
}

// Is checks IdentityEntry is equal to an other Identity
func (e Entry) Is(id Identity) bool {
	return e.MSPId == id.GetMSPID() && e.Subject == id.GetSubject()
}

//func (e Entry) FromBytes(bb []byte) (interface{}, error) {
//	entry := new(Entry)
//	err := json.Unmarshal(bb, entry)
//	return entry, err
//}

func (e Entry) GetIdentityEntry() Entry {
	return e
}

// CreateEntry creates IdentityEntry structure from an identity interface
func CreateEntry(i Identity) (g *Entry, err error) {
	return &Entry{
		MSPId:   i.GetMSPIdentifier(),
		Subject: i.GetSubject(),
		Issuer:  i.GetIssuer(),
		PEM:     i.GetPEM(),
	}, nil
}

func EntryFromStub(stub shim.ChaincodeStubInterface) (g *Entry, err error) {
	id, err := FromStub(stub)
	if err != nil {
		return nil, err
	}
	return CreateEntry(id)
}

// EntryFromSerialized creates Entry from SerializedEntry
func EntryFromSerialized(s protomsp.SerializedIdentity) (g *Entry, err error) {
	id, err := FromSerialized(s)
	if err != nil {
		return nil, err
	}
	return CreateEntry(id)
}
