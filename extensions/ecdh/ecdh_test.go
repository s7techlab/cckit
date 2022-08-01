package ecdh_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/extensions/ecdh"
	"github.com/s7techlab/cckit/identity/testdata"
)

func TestDebug(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ECDH suite")
}

// load actor certificates
var (
	pubKey1 = testdata.Certificates[0].MustCert().PublicKey.(*ecdsa.PublicKey)
	pubKey2 = testdata.Certificates[1].MustCert().PublicKey.(*ecdsa.PublicKey)
	pubKey3 = testdata.Certificates[2].MustCert().PublicKey.(*ecdsa.PublicKey)

	privKey1 = testdata.Certificates[0].MustPKey()
	privKey2 = testdata.Certificates[1].MustPKey()
	privKey3 = testdata.Certificates[2].MustPKey()
)

var _ = Describe(`ECDH`, func() {

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
