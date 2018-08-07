package owner_test

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestOwner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Owner suite")
}

type OwnableFromCreatorChaincode struct {
	router *router.Group
}

func (cc *OwnableFromCreatorChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return owner.SetFromCreator(cc.router.Context(`init`, stub))
}

func (cc *OwnableFromCreatorChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}

//  OwnableFromArgsChaincode  - owner credentials can be passed at the time of initialization
type OwnableFromArgsChaincode struct {
	router *router.Group
}

func (cc *OwnableFromArgsChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return owner.SetFromArgs(cc.router.Context(`init`, stub))
}

func (cc *OwnableFromArgsChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}

func NewOwnableFromCreator() *OwnableFromCreatorChaincode {
	r := router.New(`ownableFromCreator`) // also initialized logger with "pingable" prefix
	r.Invoke(owner.QueryMethod, owner.Query)
	return &OwnableFromCreatorChaincode{r}
}

func NewOwnableFromArgs() *OwnableFromArgsChaincode {
	r := router.New(`ownableFromArgs`) // also initialized logger with "pingable" prefix
	r.Invoke(owner.QueryMethod, owner.Query)
	return &OwnableFromArgsChaincode{r}
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
			owner := expectcc.PayloadIs(cc1.From(actors[`owner`]).Init(), &identity.Entry{}).(identity.Entry)
			Expect(owner.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
		})

		It("Owner not changed during chaincode upgrade", func() {
			// cc upgrade
			ownerAfterSecondInit := expectcc.PayloadIs(cc1.From(actors[`someone`]).Init(), &identity.Entry{}).(identity.Entry)
			Expect(ownerAfterSecondInit.Subject).To(Equal(actors[`owner`].GetSubject()))
		})

		It("Owner can be queried", func() {
			ownerIdentity := expectcc.PayloadIs(cc1.From(actors[`someone`]).Invoke(owner.QueryMethod), &identity.Entry{}).(identity.Entry)
			Expect(ownerIdentity.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
			Expect(ownerIdentity.GetMSPID()).To(Equal(actors[`owner`].GetMSPID()))
			Expect(ownerIdentity.GetPublicKey()).To(Equal(actors[`owner`].GetPublicKey()))
		})
	})

	Describe("Owner from args", func() {

		It("Allow to set owner during chaincode init", func() {
			//invoke chaincode method from someone, but pass owner mspId and cert to init
			owner := expectcc.PayloadIs(cc2.From(actors[`someone`]).Init(actors[`owner`].MspID, actors[`owner`].GetPEM()), &identity.Entry{}).(identity.Entry)
			Expect(owner.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
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
			ownerIdentity := expectcc.PayloadIs(cc2.From(actors[`someone`]).Invoke(owner.QueryMethod), &identity.Entry{}).(identity.Entry)
			Expect(ownerIdentity.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
			Expect(ownerIdentity.GetMSPID()).To(Equal(actors[`owner`].GetMSPID()))
			Expect(ownerIdentity.GetPublicKey()).To(Equal(actors[`owner`].GetPublicKey()))
		})

	})
})
