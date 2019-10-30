package testing

import (
	"github.com/s7techlab/cckit/identity"
)

type (
	Identities map[string]identity.Identity

	ReadFile func(string) ([]byte, error)
)

func MustIdentityFromPem(mspID string, certPEM []byte) *identity.CertIdentity {
	if id, err := identity.New(mspID, certPEM); err != nil {
		panic(err)
	} else {
		return id
	}
}

// IdentitiesFromPem returns CertIdentity (MSP ID and X.509 cert) converted PEM content
func IdentitiesFromPem(mspID string, certPEMs map[string][]byte) (ids Identities, err error) {
	identities := make(Identities)
	for role, certPEM := range certPEMs {
		if identities[role], err = identity.New(mspID, certPEM); err != nil {
			return
		}
	}
	return identities, nil
}

// IdentitiesFromFiles returns map of CertIdentity, loaded from PEM files
func IdentitiesFromFiles(mspID string, files map[string]string, readFile ReadFile) (Identities, error) {
	contents := make(map[string][]byte)
	for key, filename := range files {
		content, err := readFile(filename)
		if err != nil {
			return nil, err
		}
		contents[key] = content
	}
	return IdentitiesFromPem(mspID, contents)
}

// IdentityFromFile returns Identity struct containing mspId and certificate
func IdentityFromFile(mspID string, file string, readFile ReadFile) (*identity.CertIdentity, error) {
	content, err := readFile(file)
	if err != nil {
		return nil, err
	}

	return identity.New(mspID, content)
}

//  MustIdentitiesFromFiles
func MustIdentitiesFromFiles(mspID string, files map[string]string, readFile ReadFile) Identities {
	ids, err := IdentitiesFromFiles(mspID, files, readFile)
	if err != nil {
		panic(err)
	} else {
		return ids
	}
}
