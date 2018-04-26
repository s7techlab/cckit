package access

import (
	"encoding/json"

	"github.com/s7techlab/cckit/identity"
)

// Grant structure for storing grants for certificate
type Grant struct {
	MSPId   string `json:"MspId"`
	Subject string `json:"Subject"`
	Issuer  string `json:"Issuer"`
}

// ========  Identity interface ===================

// GetID identifier by certificate subject and issuer
func (g Grant) GetID() string {
	return identity.ID(g.Subject, g.Issuer)
}

// GetMSPID membership service provider identifier
func (g Grant) GetMSPID() string {
	return g.MSPId
}

// GetSubject certificate subject
func (g Grant) GetSubject() string {
	return g.Subject
}

// GetIssuer certificate issuer
func (g Grant) GetIssuer() string {
	return g.Issuer
}

// Is checks grant equal an identity
func (g Grant) Is(id identity.Identity) bool {
	return g.MSPId == id.GetMSPID() && g.Subject == id.GetSubject()
}

// FromBytes unmarshal from json bytes
func FromBytes(marshalled []byte) (grant *Grant, err error) {
	g := new(Grant)
	err = json.Unmarshal(marshalled, g)
	return g, err
}

// GrantFromIdentity creates grant structure from an identity interface
func GrantFromIdentity(i identity.Identity) (g *Grant, err error) {
	return &Grant{
		MSPId:   i.GetMSPID(),
		Subject: i.GetSubject(),
		Issuer:  i.GetIssuer(),
	}, nil
}
