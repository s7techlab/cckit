package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/extensions/access"
	"github.com/s7techlab/cckit/identity"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestMarbles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var _ = Describe(`Marbles`, func() {

	var cc *testcc.MockStub
	actor := map[string]identity.CertIdentity{}

	BeforeSuite(func() {
		// load certificates
		for role, filename := range map[string]string{`operator`: `s7techlab.pem`, `owner1`: `victor-nosov.pem`} {
			cert, err := examplecert.Plain(filename)
			if err != nil {
				panic(err)
			}

			i, err := identity.FromCert(`SOME_MSP`, cert)
			if err != nil {
				panic(err)
			}

			actor[role] = *i.(*identity.CertIdentity)
		}

		cc = testcc.NewMockStub(`marbles`, New())
		expectcc.ResponseOk(cc.From(actor[`operator`]).Init())
	})

	Describe("Chaincode owner", func() {

		It("Allow everyone to retrieve chaincode owner", func() {
			grant := expectcc.PayloadIs(cc.Invoke(`owner`), &access.Grant{}).(*access.Grant)
			Expect(grant.GetSubject()).To(Equal(actor[`operator`].GetSubject()))
			Expect(grant.Is(actor[`operator`])).To(BeTrue())
		})

	})

	Describe("Marble owner", func() {

		It("Allow chaincode owner to register marble owner", func() {
			expectcc.ResponseOk(cc.From(actor[`operator`]).Invoke(`marbleOwnerRegister`, actor[`owner1`].ToSerialized()))
		})

	})
})
