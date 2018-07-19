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

type OwnableChaincode struct {
	router *router.Group
}

func (cc *OwnableChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return owner.SetFromCreator(cc.router.Context(`init`, stub))
}

func (cc *OwnableChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}

func New() *OwnableChaincode {
	r := router.New(`ownable`) // also initialized logger with "pingable" prefix
	r.Invoke(owner.QueryMethod, owner.Query)
	return &OwnableChaincode{r}
}

var _ = Describe(`Ownable`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`ownable`, New())
	actors, err := identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
		`owner`:   `s7techlab.pem`,
		`someone`: `victor-nosov.pem`}, examplecert.Content)
	if err != nil {
		panic(err)
	}

	Describe("Owner", func() {

		It("Allow set owner during chaincode init", func() {
			//invoke chaincode method from authority actor
			owner := expectcc.PayloadIs(cc.From(actors[`owner`]).Init(), &identity.Entry{}).(identity.Entry)
			Expect(owner.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
		})

		It("Owner not changed during chaincode upgrade", func() {
			// cc upgrade
			ownerAfterSecondInit := expectcc.PayloadIs(cc.From(actors[`someone`]).Init(), &identity.Entry{}).(identity.Entry)
			Expect(ownerAfterSecondInit.Subject).To(Equal(actors[`owner`].GetSubject()))
		})

		It("Can be queried", func() {
			ownerIdentity := expectcc.PayloadIs(cc.From(actors[`someone`]).Invoke(owner.QueryMethod), &identity.Entry{}).(identity.Entry)
			Expect(ownerIdentity.GetSubject()).To(Equal(actors[`owner`].GetSubject()))
			Expect(ownerIdentity.GetMSPID()).To(Equal(actors[`owner`].GetMSPID()))
			Expect(ownerIdentity.GetPublicKey()).To(Equal(actors[`owner`].GetPublicKey()))
		})
	})
})
