package ecdh_test

import (
	"crypto/ecdsa"
	"testing"

	"crypto/elliptic"

	"fmt"

	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/extensions/ecdh"
	"github.com/s7techlab/cckit/identity"
)

func TestDebug(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ECDH suite")
}

var _ = Describe(`ECDH`, func() {

	//load actor certificates
	actors, err := identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
		`authority`: `s7techlab.pem`,
		`someone1`:  `victor-nosov.pem`,
		`someone2`:  `some-person.pem`,
	}, examplecert.Content)
	if err != nil {
		panic(err)
	}

	pubKey1 := actors[`authority`].GetPublicKey().(*ecdsa.PublicKey)
	pubKey2 := actors[`someone1`].GetPublicKey().(*ecdsa.PublicKey)
	pubKey3 := actors[`someone2`].GetPublicKey().(*ecdsa.PublicKey)

	privKey1Bytes, _ := examplecert.Content(`s7techlab.key.pem`)
	privKey1, err := ecdh.PrivateKey(privKey1Bytes)

	if err != nil {
		panic(err)
	}

	privKey2Bytes, _ := examplecert.Content(`victor-nosov.key.pem`)
	privKey2, err := ecdh.PrivateKey(privKey2Bytes)
	if err != nil {
		panic(err)
	}

	privKey3Bytes, _ := examplecert.Content(`some-person.key.pem`)
	privKey3, err := ecdh.PrivateKey(privKey3Bytes)
	if err != nil {
		panic(err)
	}

	It("Allow to create shared key for 2 parties", func() {

		secret12, err := ecdh.GenerateSharedSecret(privKey1, pubKey2)
		Expect(err).To(BeNil())
		secret21, err := ecdh.GenerateSharedSecret(privKey2, pubKey1)
		Expect(err).To(BeNil())

		Expect(secret12).To(Equal(secret21))

		secret23, err := ecdh.GenerateSharedSecret(privKey2, pubKey3)
		Expect(err).To(BeNil())
		secret32, err := ecdh.GenerateSharedSecret(privKey3, pubKey2)
		Expect(err).To(BeNil())

		Expect(secret23).To(Equal(secret32))

	})

	It("Allow to create shared key for 3 parties", func() {

		secret12, err := ecdh.GenerateSharedSecret(privKey1, pubKey2)
		Expect(err).To(BeNil())

		x, y := elliptic.Unmarshal(pubKey1.Curve, secret12)

		fmt.Println(x, y)
		//if x == nil || y == nil {
		//	return key, false
		//}
		//key = &ellipticPublicKey{
		//	Curve: e.curve,
		//	X:     x,
		//	Y:     y,
		//}

	})

})
