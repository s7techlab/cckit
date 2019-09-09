package identity

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

func PrivateKey(keyBytes []byte) (*ecdsa.PrivateKey, error) {
	keyPEM, _ := pem.Decode(keyBytes)
	if keyPEM == nil {
		return nil, ErrInvalidPEMStructure
	}

	key, err := x509.ParsePKCS8PrivateKey(keyPEM.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse private key`)
	}
	return key.(*ecdsa.PrivateKey), nil
}
