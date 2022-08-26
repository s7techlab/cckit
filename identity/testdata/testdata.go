package testdata

import (
	"crypto/ecdsa"
	"crypto/x509"
	"io/ioutil"
	"path"
	"runtime"

	"github.com/hyperledger/fabric/msp"

	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/testing"
)

const DefaultMSP = `SOME_MSP`

type (
	FileReader func(filename string) ([]byte, error)

	// Cert certificate data for testing
	Cert struct {
		CertFilename string
		PKeyFilename string
		readFile     FileReader
	}

	Certs []*Cert

	IdentitySample struct {
		MspID string
		Cert  *Cert
	}
)

func (cc Certs) UseReadFile(readFile FileReader) Certs {
	for _, c := range cc {
		c.readFile = readFile
	}
	return cc
}

func (s *IdentitySample) SigningIdentity() msp.SigningIdentity {
	return s.Cert.MustIdentity(s.MspID)
}

var (
	Certificates = Certs{{
		CertFilename: `s7techlab.pem`, PKeyFilename: `s7techlab.key.pem`,
	}, {
		CertFilename: `some-person.pem`, PKeyFilename: `some-person.key.pem`,
	}, {
		CertFilename: `victor-nosov.pem`, PKeyFilename: `victor-nosov.key.pem`,
	}}.
		UseReadFile(ReadLocal())
)

func ReadLocal() func(filename string) ([]byte, error) {
	_, curFile, _, ok := runtime.Caller(1)
	curFilePath := path.Dir(curFile)
	if !ok {
		return nil
	}
	return func(filename string) ([]byte, error) {
		return ioutil.ReadFile(curFilePath + "/" + filename)
	}
}

func MustSamples(cc []*Cert, mspId string) []*IdentitySample {
	ss := make([]*IdentitySample, len(cc))
	for i, c := range Certificates {
		ss[i] = &IdentitySample{
			MspID: mspId,
			Cert:  c,
		}
	}

	return ss
}
func MustIdentities(cc []*Cert, mspId string) []*identity.CertIdentity {
	ii := make([]*identity.CertIdentity, len(cc))
	for i, c := range Certificates {
		ii[i] = c.MustIdentity(mspId)
	}

	return ii
}

func (c *Cert) MustIdentity(mspID string) *identity.CertIdentity {
	id, err := c.Identity(mspID)
	if err != nil {
		panic(err)
	}
	return id
}

func (c *Cert) CertBytes() ([]byte, error) {
	return c.readFile(`./` + c.CertFilename)
}

func (c *Cert) PKeyBytes() ([]byte, error) {
	return c.readFile(`./` + c.PKeyFilename)
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

func (c *Cert) SigningIdentity(mspID string) (*identity.CertIdentity, error) {
	return c.Identity(mspID)
}

func (c *Cert) MustSigningIdentity(mspID string) *identity.CertIdentity {
	bb := c.MustCertBytes()
	return testing.MustIdentityFromPem(mspID, bb)
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
