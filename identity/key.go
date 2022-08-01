package identity

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func PrivateKey(keyBytes []byte) (*ecdsa.PrivateKey, error) {
	keyPEM, _ := pem.Decode(keyBytes)
	if keyPEM == nil {
		return nil, ErrInvalidPEMStructure
	}

	key, err := x509.ParsePKCS8PrivateKey(keyPEM.Bytes)
	if err != nil {
		return nil, fmt.Errorf(`parse private key: %w`, err)
	}
	return key.(*ecdsa.PrivateKey), nil
}
