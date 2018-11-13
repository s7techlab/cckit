package owner

import (
	"testing"

	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOwner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Owner suite")
}

func NewOwnableFromCreator() *router.Chaincode {
	return router.NewChaincode(router.
		New(`ownableFromCreator`).
		Init(InvokeSetFromCreator).
		Invoke(QueryMethod, Query))
}

// NewOwnableFromArgs - owner credentials can be passed at the time of initialization
func NewOwnableFromArgs() *router.Chaincode {
	return router.NewChaincode(router.
		New(`ownableFromArgs`).
		Init(InvokeSetFromArgs).
		Invoke(QueryMethod, Query))
}

var _ = Describe(`Ownable`, func() {

	//Create chaincode mock
	cc1 := testcc.NewMockStub(`ownableFromCreator`, NewOwnableFromCreator())
	cc2 := testcc.NewMockStub(`ownableFromArgs`, NewOwnableFromArgs())
	actors, err := identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
		`owner`:   `s7techlab.pem`,
		`someone`: `victor-nosov.pem`}, examplecert.Content)
	if err != nil {
		panic(err)
	}

	Describe("Owner from creator", func() {

		It("Allow to set owner during chaincode init", func() {
			//invoke chaincode method from authority actor
			ownerEntry := expectcc.PayloadIs(cc1.From(actors[`owner`]).Init(), &identity.Entry{}).(identity.Entry)
			Expect(ownerEntry.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
		})

		It("Owner not changed during chaincode upgrade", func() {
			// cc upgrade
			ownerAfterSecondInit := expectcc.PayloadIs(cc1.From(actors[`someone`]).Init(), &identity.Entry{}).(identity.Entry)
			Expect(ownerAfterSecondInit.Subject).To(Equal(actors[`owner`].GetSubject()))
		})

		It("Owner can be queried", func() {
			ownerIdentity := expectcc.PayloadIs(cc1.From(actors[`someone`]).Invoke(QueryMethod), &identity.Entry{}).(identity.Entry)
			Expect(ownerIdentity.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
			Expect(ownerIdentity.GetMSPID()).To(Equal(actors[`owner`].GetMSPID()))
			Expect(ownerIdentity.GetPublicKey()).To(Equal(actors[`owner`].GetPublicKey()))
		})
	})

	Describe("Owner from args", func() {

		It("Allow to set owner during chaincode init", func() {
			//invoke chaincode method from someone, but pass owner mspId and cert to init
			ownerEntry := expectcc.PayloadIs(cc2.From(actors[`someone`]).Init(actors[`owner`].MspID, actors[`owner`].GetPEM()), &identity.Entry{}).(identity.Entry)
			Expect(ownerEntry.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
		})

		It("Owner not changed during chaincode upgrade", func() {
			// cc upgrade
			ownerAfterSecondInit := expectcc.PayloadIs(cc2.From(actors[`someone`]).Init(), &identity.Entry{}).(identity.Entry)
			Expect(ownerAfterSecondInit.Subject).To(Equal(actors[`owner`].GetSubject()))
		})
		It("Disallow set owner twice", func() {
			//invoke chaincode method from someone, but pass owner mspId and cert to init
			expectcc.ResponseError(cc2.From(actors[`someone`]).Init(actors[`owner`].MspID, actors[`owner`].GetPEM()))
		})

		It("Owner can be queried", func() {
			ownerIdentity := expectcc.PayloadIs(cc2.From(actors[`someone`]).Invoke(QueryMethod), &identity.Entry{}).(identity.Entry)
			Expect(ownerIdentity.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
			Expect(ownerIdentity.GetMSPID()).To(Equal(actors[`owner`].GetMSPID()))
			Expect(ownerIdentity.GetPublicKey()).To(Equal(actors[`owner`].GetPublicKey()))
		})

	})
})
