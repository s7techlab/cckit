package cert

import (
	"io/ioutil"
	"path"
	"runtime"

	"github.com/s7techlab/cckit/identity"
)

func getFileContent(fixtureFile string) (content []byte, err error) {
	_, curFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("no caller information")
	}
	content, err = ioutil.ReadFile(path.Dir(curFile) + "/" + fixtureFile)
	return
}

// Plain certificate text
func Plain(filename string) ([]byte, error) {
	return getFileContent(filename)
}

// Actors returns CertIdentity loaded from certificates from filesystem
func Actors(roles map[string]string) (map[string]*identity.CertIdentity, error) {

	actors := make(map[string]*identity.CertIdentity)
	for role, filename := range roles {
		cert, err := Plain(filename)
		if err != nil {
			return nil, err
		}
		ci, err := identity.New(`SOME_MSP`, cert)
		if err != nil {
			return nil, err
		}
		actors[role] = ci
	}
	return actors, nil
}
