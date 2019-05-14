package testdata

import (
	"github.com/s7techlab/cckit/testing"
	"io/ioutil"
)

func GetTestIdentity(msp, file string) *testing.Identity {
	identity, err := testing.IdentityFromFile(msp, file, ioutil.ReadFile)
	if err != nil {
		panic(err)
	}

	return identity
}
