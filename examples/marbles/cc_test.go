package marbles_test

import (
	"testing"

	"github.com/s7techlab/cckit/examples/marbles"
	"github.com/s7techlab/cckit/examples/marbles/testdata"

	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/identity"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestMarbles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var _ = Describe(`Marbles`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`marbles`, marbles.New())

	// load actor certificates from github.com/s7techlab/cckit/examples/cert
	actors, err := identity.ActorsFromPemFile(
		`SOME_MSP`,
		map[string]string{`operator`: `s7techlab.pem`, `owner1`: `victor-nosov.pem`},
		examplecert.Content)
	if err != nil {
		panic(err)
	}

	BeforeSuite(func() {
		// Init chaincode from operator
		expectcc.ResponseOk(cc.From(actors[`operator`]).Init())
	})

	//Describe("Marble owner", func() {
	//
	//	It("Allow chaincode owner to register marble owner", func() {
	//		expectcc.ResponseOk(
	//			// register owner1 certificate as potential marble owner
	//			cc.From(actors[`operator`]).Invoke(`marbleOwnerRegister`, actors[`owner1`]))
	//	})
	//
	//	It("Disallow non chaincode owner to register marble owner", func() {
	//		expectcc.ResponseError(
	//			cc.From(actors[`owner1`]).Invoke(`marbleOwnerRegister`, actors[`owner1`]),
	//			owner.ErrOwnerOnly)
	//	})
	//
	//	It("Disallow chaincode owner to register duplicate marble owner", func() {
	//		expectcc.ResponseError(
	//			cc.From(actors[`operator`]).Invoke(`marbleOwnerRegister`, actors[`owner1`]),
	//			state.ErrKeyAlreadyExists)
	//	})
	//
	//	It("Disallow to pass non SerializedIdentity json", func() {
	//		expectcc.ResponseError(
	//			cc.From(actors[`owner1`]).Invoke(`marbleOwnerRegister`, `some weird string`),
	//			convert.ErrUnableToConvertValueToStruct)
	//	})
	//
	//})

	Describe("Marble init", func() {

		It("Allow to init information about marble", func() {
			expectcc.ResponseOk(
				// register owner1 certificate as potential marble owner
				cc.From(actors[`operator`]).Invoke(`marbleInit`, testdata.MarblePayloads[0]))
		})

	})

})
