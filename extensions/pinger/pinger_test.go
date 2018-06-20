package pinger

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/response"
	"github.com/s7techlab/cckit/router"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestPinger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pinger suite")
}

type PingableChaincode struct {
	router *router.Group
}

func (cc *PingableChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return response.Success(nil)
}

// Invoke - entry point for chain code invocations
func (cc *PingableChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}

func New() *PingableChaincode {
	r := router.New(`pingable`) // also initialized logger with "pingable" prefix
	r.Invoke(FuncPing, Ping)
	return &PingableChaincode{r}
}

var _ = Describe(`Pinger`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`cars`, New())
	invokerIdentity, err := identity.FromFile(`SOME_MSP`, `s7techlab.pem`, examplecert.Content)
	if err != nil {
		panic(err)
	}

	Describe("Pinger", func() {

		It("Allow anynone to invoke ping method", func() {
			//invoke chaincode method from authority actor
			pingInfo := expectcc.PayloadIs(cc.From(invokerIdentity).Invoke(FuncPing), &PingInfo{}).(PingInfo)
			Expect(pingInfo.InvokerID).To(Equal(invokerIdentity.GetID()))
			Expect(pingInfo.InvokerCert).To(Equal(invokerIdentity.GetPEM()))

			//check that we have event
			pingInfoEvent := expectcc.EventPayloadIs(cc.ChaincodeEvent, &PingInfo{}).(PingInfo)
			Expect(pingInfoEvent.InvokerID).To(Equal(invokerIdentity.GetID()))
			Expect(pingInfoEvent.InvokerCert).To(Equal(invokerIdentity.GetPEM()))
		})

	})
})
