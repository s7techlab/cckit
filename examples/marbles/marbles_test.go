package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/extensions/access"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestMarbles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var _ = Describe(`Marbles`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`marbles`, New())

	// load actor certificates
	actors, err := examplecert.Actors(map[string]string{`operator`: `s7techlab.pem`, `owner1`: `victor-nosov.pem`})
	if err != nil {
		panic(err)
	}

	BeforeSuite(func() {
		// Init chaincode from operator
		expectcc.ResponseOk(cc.From(actors[`operator`]).Init())
	})

	Describe("Chaincode owner", func() {
		It("Allow everyone to retrieve chaincode owner", func() {
			grant := expectcc.PayloadIs(cc.Invoke(`owner`), &access.Grant{}).(*access.Grant)
			Expect(grant.GetSubject()).To(Equal(actors[`operator`].GetSubject()))
			Expect(grant.Is(actors[`operator`])).To(BeTrue())
		})
	})

	Describe("Marble owner", func() {

		It("Allow chaincode owner to register marble owner", func() {
			expectcc.ResponseOk(cc.From(actors[`operator`]).Invoke(`marbleOwnerRegister`, actors[`owner1`].ToSerialized()))
		})

	})
})
