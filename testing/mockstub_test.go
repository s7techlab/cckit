package testing

import (
	"testing"

	"github.com/s7techlab/cckit/examples/cars"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	"time"

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

		Expect(len(cc.ChaincodeEventsChannel)).To(Equal(1))

	})

	It("Allow to clear events channel", func() {

		cc.ClearEvents()
		Expect(len(cc.ChaincodeEventsChannel)).To(Equal(0))
		timeout := make(chan bool, 1)

		go func() {
			time.Sleep(time.Millisecond * 10)
			timeout <- true
		}()

		expectcc.ResponseOk(cc.From(actors[`authority`]).Invoke(`carRegister`, cars.Payloads[1]))

		select {
		case ccEvent := <-cc.ChaincodeEventsChannel:
			Expect(ccEvent.EventName).To(Equal(cars.CarRegisteredEvent))
			event := expectcc.EventPayloadIs(ccEvent, &cars.Car{}).(cars.Car)
			Expect(event.Id).To(Equal(cars.Payloads[1].Id))
		case <-timeout:
			Expect(true).To(Equal(false), `Event not received`)
		}

	})

})
