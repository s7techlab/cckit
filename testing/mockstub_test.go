package testing

import (
	"testing"

	"github.com/s7techlab/cckit/examples/cars"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	"github.com/hyperledger/fabric/protos/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/identity"
)

func TestMockstub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mockstub Suite")
}

var _ = Describe(`Mockstub`, func() {

	//Create chaincode mocks
	cc := NewMockStub(`cars`, cars.New())

	// load actor certificates
	actors, err := identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
		`authority`: `s7techlab.pem`,
		`someone`:   `victor-nosov.pem`}, examplecert.Content)
	if err != nil {
		panic(err)
	}

	It("Allow to init chaincode", func() {
		//invoke chaincode method from authority actor
		expectcc.ResponseOk(cc.From(actors[`authority`]).Init()) // init chaincode from authority
	})

	It("Allow to get last event while chaincode invoke ", func() {

		expectcc.ResponseOk(cc.From(actors[`authority`]).Invoke(`carRegister`, cars.Payloads[0]))
		event := expectcc.EventPayloadIs(cc.ChaincodeEvent, &cars.Car{}).(cars.Car)

		Expect(cc.ChaincodeEvent.EventName).To(Equal(cars.CarRegisteredEvent))
		Expect(event.Id).To(Equal(cars.Payloads[0].Id))

	})

	It("Allow to get  t event while chaincode invoke using ChaincodeEventsChannel ", func() {

		var ccEvent *peer.ChaincodeEvent
		go func() {
			select {
			case ccEvent = <-cc.ChaincodeEventsChannel:
			}
		}()

		expectcc.ResponseOk(cc.From(actors[`authority`]).Invoke(`carRegister`, cars.Payloads[1]))
		Expect(ccEvent.EventName).To(Equal(cars.CarRegisteredEvent))
	})

})
