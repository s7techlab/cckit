package ecdh

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
)

var (
	ErrInvalidPEMStructure = errors.New(`invalid pem structure`)
)

func Marshall(pubKey *ecdsa.PublicKey) []byte {
	byteLen := (pubKey.Curve.Params().BitSize + 7) >> 3

	ret := make([]byte, 1+2*byteLen)
	ret[0] = 4 // uncompressed point

	xBytes := pubKey.X.Bytes()
	copy(ret[1+byteLen-len(xBytes):], xBytes)
	yBytes := pubKey.Y.Bytes()
	copy(ret[1+2*byteLen-len(yBytes):], yBytes)
	return ret

}

func GenerateSharedSecret(privKey *ecdsa.PrivateKey, pubKey *ecdsa.PublicKey) ([]byte, error) {
	x, _ := pubKey.Curve.ScalarMult(pubKey.X, pubKey.Y, privKey.D.Bytes())
	return x.Bytes(), nil
}

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
