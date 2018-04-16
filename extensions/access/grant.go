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
func (g Grant) GetId() (string) {
	return identity.Id(g.Subject, g.Issuer)
}

func (g Grant) GetMSPId() (string) {
	return g.MSPId
}

func (g Grant) GetSubject() (string) {
	return g.Subject
}

func (g Grant) GetIssuer() (string) {
	return g.Issuer
}



func (g Grant) Is(id identity.Identity) bool {
	return g.MSPId == id.GetMSPId() && g.Subject == id.GetSubject()
}

func (g Grant) ToBytes() (marshalled []byte) {
	marshalled, _ = json.Marshal(g)
	return
}

func FromBytes(marshalled []byte) (grant *Grant, err error) {
	g := new(Grant)
	err = json.Unmarshal(marshalled, g)
	return g, err
}

func GrantFromIdentity(i identity.Identity) (g *Grant, err error) {
	return &Grant{
		MSPId:   i.GetMSPId(),
		Subject: i.GetSubject(),
		Issuer:  i.GetIssuer(),
	}, nil
}