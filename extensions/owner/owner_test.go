package owner

import (
	"testing"

	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/identity/testdata"
	"github.com/s7techlab/cckit/router"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	Owner   = testdata.Certificates[0].MustIdentity(`SOME_MSP`)
	Someone = testdata.Certificates[1].MustIdentity(`SOME_MSP`)
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
	Describe("Owner from creator", func() {

		It("Allow to set owner during chaincode init", func() {
			//invoke chaincode method from authority actor
			ownerEntry := expectcc.PayloadIs(cc1.From(Owner).Init(), &identity.Entry{}).(identity.Entry)
			Expect(ownerEntry.GetSubject()).To(Equal(Owner.GetSubject()))
		})

		It("Owner not changed during chaincode upgrade", func() {
			// cc upgrade
			ownerAfterSecondInit := expectcc.PayloadIs(cc1.From(Someone).Init(), &identity.Entry{}).(identity.Entry)
			Expect(ownerAfterSecondInit.Subject).To(Equal(Owner.GetSubject()))
		})

		It("Owner can be queried", func() {
			ownerIdentity := expectcc.PayloadIs(
				cc1.From(Someone).Invoke(QueryMethod), &identity.Entry{}).(identity.Entry)
			Expect(ownerIdentity.GetSubject()).To(Equal(Owner.GetSubject()))
			Expect(ownerIdentity.GetMSPIdentifier()).To(Equal(Owner.MspID))
			Expect(ownerIdentity.GetPublicKey()).To(Equal(Owner.Cert.PublicKey))
		})
	})

	Describe("Owner from args", func() {

		It("Allow to set owner during chaincode init", func() {
			//invoke chaincode method from someone, but pass owner mspId and cert to init
			ownerEntry := expectcc.PayloadIs(
				cc2.From(Someone).Init(Owner.MspID, Owner.GetPEM()), &identity.Entry{}).(identity.Entry)
			Expect(ownerEntry.GetSubject()).To(Equal(Owner.GetSubject()))
		})

		It("Owner not changed during chaincode upgrade", func() {
			// cc upgrade
			ownerAfterSecondInit := expectcc.PayloadIs(cc2.From(Someone).Init(), &identity.Entry{}).(identity.Entry)
			Expect(ownerAfterSecondInit.Subject).To(Equal(Owner.GetSubject()))
		})
		It("Disallow set owner twice", func() {
			//invoke chaincode method from someone, but pass owner mspId and cert to init
			expectcc.ResponseError(cc2.From(Someone).Init(Owner.MspID, Owner.GetPEM()))
		})

		It("Owner can be queried", func() {
			ownerIdentity := expectcc.PayloadIs(
				cc2.From(Someone).Invoke(QueryMethod), &identity.Entry{}).(identity.Entry)
			Expect(ownerIdentity.GetSubject()).To(Equal(Owner.GetSubject()))
			Expect(ownerIdentity.GetMSPIdentifier()).To(Equal(Owner.MspID))
			Expect(ownerIdentity.GetPublicKey()).To(Equal(Owner.Cert.PublicKey))
		})

	})
})
