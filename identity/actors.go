package identity

type GetContent func(string) ([]byte, error)
type Actors map[string]*CertIdentity

// ActorsFromPem returns CertIdentity (MSP ID and X.509 cert) converted PEM content
func ActorsFromPem(mspID string, certPEMs map[string][]byte) (map[string]*CertIdentity, error) {
	actors := make(map[string]*CertIdentity)
	for role, certPEM := range certPEMs {
		ci, err := New(mspID, certPEM)
		if err != nil {
			return nil, err
		}
		actors[role] = ci
	}
	return actors, nil
}

// ActorsFromPemFile returns map of CertIdentity, loaded from PEM files
func ActorsFromPemFile(mspID string, files map[string]string, getContent GetContent) (Actors, error) {
	contents := make(map[string][]byte)
	for key, filename := range files {
		content, err := getContent(filename)
		if err != nil {
			return nil, err
		}
		contents[key] = content
	}
	return ActorsFromPem(mspID, contents)
}
