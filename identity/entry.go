// Package access contains structs for storing chaincode access control information
package identity

import "github.com/hyperledger/fabric/core/chaincode/shim"

// Entry structure for storing identity information
type Entry struct {
	MSPId   string
	Subject string
	Issuer  string
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
		MSPId:   i.GetMSPID(),
		Subject: i.GetSubject(),
		Issuer:  i.GetIssuer(),
	}, nil
}

func EntryFromStub(stub shim.ChaincodeStubInterface) (g *Entry, err error) {
	id, err := FromStub(stub)
	if err != nil {
		return nil, err
	}
	return CreateEntry(id)
}
