package ecdh_test

import (
	"crypto/ecdsa"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/extensions/ecdh"
	"github.com/s7techlab/cckit/identity"
)

func TestDebug(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ECDH suite")
}

var _ = Describe(`ECDH`, func() {

	It("Allow to create shared key", func() {

		//load actor certificates
		actors, err := identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
			`authority`: `s7techlab.pem`,
			`someone`:   `victor-nosov.pem`}, examplecert.Content)
		if err != nil {
			panic(err)
		}

		pubKey1 := actors[`authority`].GetPublicKey().(*ecdsa.PublicKey)
		pubKey2 := actors[`someone`].GetPublicKey().(*ecdsa.PublicKey)

		privKey1Bytes, _ := examplecert.Content(`s7techlab.key.pem`)
		privKey1, err := ecdh.PrivateKey(privKey1Bytes)

		privKey2Bytes, _ := examplecert.Content(`victor-nosov.key.pem`)

		privKey2, _ := ecdh.PrivateKey(privKey2Bytes)

		secret1, err := ecdh.GenerateSharedSecret(privKey1, pubKey2)
		secret2, err := ecdh.GenerateSharedSecret(privKey2, pubKey1)
		Expect(secret1).To(Equal(secret2))
	})

})
