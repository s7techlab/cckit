package testdata

import (
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"path"
	"runtime"

	"github.com/s7techlab/cckit/testing"

	"github.com/s7techlab/cckit/identity"
)

type Cert struct {
	CertFilename string
	PKeyFilename string
}

var (
	Certificates = []*Cert{{
		CertFilename: `s7techlab.pem`, PKeyFilename: `s7techlab.key.pem`,
	}, {
		CertFilename: `some-person.pem`, PKeyFilename: `some-person.key.pem`,
	}, {
		CertFilename: `victor-nosov.pem`, PKeyFilename: `victor-nosov.key.pem`,
	}}
)

func ReadFile(filename string) ([]byte, error) {
	_, curFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.New(`can't load file, error accessing runtime caller'`)
	}
	return ioutil.ReadFile(path.Dir(curFile) + "/" + filename)
}

func (c *Cert) CertBytes() ([]byte, error) {
	return ReadFile(`./` + c.CertFilename)
}

func (c *Cert) PKeyBytes() ([]byte, error) {
	return ReadFile(`./` + c.PKeyFilename)
}

func (c *Cert) MustCertBytes() []byte {
	cert, err := c.CertBytes()
	if err != nil {
		panic(err)
	}
	return cert
}

func (c *Cert) MustPKeyBytes() []byte {
	cert, err := c.PKeyBytes()
	if err != nil {
		panic(err)
	}
	return cert
}

func (c *Cert) Identity(mspID string) (*identity.CertIdentity, error) {
	bb, err := c.CertBytes()
	if err != nil {
		return nil, err
	}
	return identity.New(mspID, bb)
}

// temp, todo: move signing identity from testing to identity package
func (c *Cert) SigningIdentity(mspID string) (*testing.Identity, error) {
	bb, err := c.CertBytes()
	if err != nil {
		return nil, err
	}
	return testing.IdentityFromPem(mspID, bb)
}

func (c *Cert) MustSigningIdentity(mspID string) *testing.Identity {
	bb := c.MustCertBytes()
	return testing.MustIdentityFromPem(mspID, bb)
}

func (c *Cert) MustIdentity(mspID string) *identity.CertIdentity {
	id, err := c.Identity(mspID)
	if err != nil {
		panic(err)
	}
	return id
}

func (c *Cert) Cert() (*x509.Certificate, error) {
	bb, err := c.CertBytes()
	if err != nil {
		return nil, err
	}
	return identity.Certificate(bb)
}

func (c *Cert) MustCert() *x509.Certificate {
	cert, err := c.Cert()
	if err != nil {
		panic(err)
	}
	return cert
}

func (c *Cert) Pkey() (*ecdsa.PrivateKey, error) {
	bb, err := c.PKeyBytes()
	if err != nil {
		return nil, err
	}
	return identity.PrivateKey(bb)
}

func (c *Cert) MustPKey() *ecdsa.PrivateKey {
	pkey, err := c.Pkey()
	if err != nil {
		panic(err)
	}
	return pkey
}
