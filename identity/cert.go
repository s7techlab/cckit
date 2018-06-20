package identity

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

// Certificate decode and parse .pem []byte x509 certificate structure
func Certificate(c []byte) (cert *x509.Certificate, err error) {
	block, _ := pem.Decode(c)
	if block == nil {
		return nil, ErrPemEncodedExpected
	}
	return x509.ParseCertificate(block.Bytes)
}

// ID returns identifier from .509  certificate and base64 encode
func ID(subject, issuer string) string {
	return base64.StdEncoding.EncodeToString([]byte(IDRaw(subject, issuer)))
}

// IDRaw generates string identifier from .509  certificate
func IDRaw(subject, issuer string) string {
	return fmt.Sprintf("x509::%s::%s", subject, issuer)
}

// IDByCert returns id by certificate subject and issuer
func IDByCert(cert *x509.Certificate) string {
	return ID(GetDN(&cert.Subject), GetDN(&cert.Issuer))
}

// GetDN (distinguished name) associated with a pkix.Name.
// NOTE: This code is almost a direct copy of the String() function in
// https://go-review.googlesource.com/c/go/+/67270/1/src/crypto/x509/pkix/pkix.go#26
// which returns a DN as defined by RFC 2253.
func GetDN(name *pkix.Name) string {
	r := name.ToRDNSequence()
	s := ""
	for i := 0; i < len(r); i++ {
		rdn := r[len(r)-1-i]
		if i > 0 {
			s += ","
		}
		for j, tv := range rdn {
			if j > 0 {
				s += "+"
			}
			typeString := tv.Type.String()
			typeName, ok := attributeTypeNames[typeString]
			if !ok {
				derBytes, err := asn1.Marshal(tv.Value)
				if err == nil {
					s += typeString + "=#" + hex.EncodeToString(derBytes)
					continue // No value escaping necessary.
				}
				typeName = typeString
			}
			s += typeName + "=" + getEscaped(tv.Value)
		}
	}
	return s
}

func getEscaped(val interface{}) string {

	valueString := fmt.Sprint(val)
	escaped := ""
	begin := 0
	for idx, c := range valueString {
		if (idx == 0 && (c == ' ' || c == '#')) ||
			(idx == len(valueString)-1 && c == ' ') {
			escaped += valueString[begin:idx]
			escaped += "\\" + string(c)
			begin = idx + 1
			continue
		}
		switch c {
		case ',', '+', '"', '\\', '<', '>', ';':
			escaped += valueString[begin:idx]
			escaped += "\\" + string(c)
			begin = idx + 1
		}
	}
	escaped += valueString[begin:]
	return escaped
}

var attributeTypeNames = map[string]string{
	"2.5.4.6":  "C",
	"2.5.4.10": "O",
	"2.5.4.11": "OU",
	"2.5.4.3":  "CN",
	"2.5.4.5":  "SERIALNUMBER",
	"2.5.4.7":  "L",
	"2.5.4.8":  "ST",
	"2.5.4.9":  "STREET",
	"2.5.4.17": "POSTALCODE",
}
