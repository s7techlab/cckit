package pinger

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/identity/testdata"
	"github.com/s7techlab/cckit/router"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestPinger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pinger suite")
}

var (
	Someone = testdata.Certificates[1].MustIdentity(`SOME_MSP`)
)

func New() *router.Chaincode {
	r := router.New(`pingable`).
		Init(router.EmptyContextHandler).
		Invoke(FuncPing, Ping)
	return router.NewChaincode(r)
}

var _ = Describe(`Pinger`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`cars`, New())

	Describe("Pinger", func() {

		It("Allow anynone to invoke ping method", func() {
			//invoke chaincode method from authority actor
			pingInfo := expectcc.PayloadIs(cc.From(Someone).Invoke(FuncPing), &PingInfo{}).(PingInfo)
			Expect(pingInfo.InvokerID).To(Equal(Someone.GetID()))
			Expect(pingInfo.InvokerCert).To(Equal(Someone.GetPEM()))

			//check that we have event
			pingInfoEvent := expectcc.EventPayloadIs(cc.ChaincodeEvent, &PingInfo{}).(PingInfo)
			Expect(pingInfoEvent.InvokerID).To(Equal(Someone.GetID()))
			Expect(pingInfoEvent.InvokerCert).To(Equal(Someone.GetPEM()))
		})

	})
})
