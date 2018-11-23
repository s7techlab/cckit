package testing

import (
	"testing"

	"github.com/s7techlab/cckit/examples/cars"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	"time"

	"context"

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

	})

	It("Allow to get events via events channel", func() {

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*10)

		expectcc.ResponseOk(cc.From(actors[`authority`]).Invoke(`carRegister`, cars.Payloads[1]))

		select {
		case ccEvent := <-cc.ChaincodeEventsChannel:
			Expect(ccEvent.EventName).To(Equal(cars.CarRegisteredEvent))
			event := expectcc.EventPayloadIs(ccEvent, &cars.Car{}).(cars.Car)
			Expect(event.Id).To(Equal(cars.Payloads[1].Id))
		case <-ctx.Done():
			Expect(true).To(Equal(false), `Event not received`)
		}

	})

	It("Allow to use multiple events subscriptions", func() {
		Expect(len(cc.ChaincodeEventsChannel)).To(Equal(0))

		sub1 := cc.EventSubscription()
		sub2 := cc.EventSubscription()

		Expect(len(sub1)).To(Equal(0))
		Expect(len(sub2)).To(Equal(0))

		expectcc.ResponseOk(cc.From(actors[`authority`]).Invoke(`carRegister`, cars.Payloads[2]))

		Expect(len(cc.ChaincodeEventsChannel)).To(Equal(1))
		Expect(len(sub1)).To(Equal(1))
		Expect(len(sub2)).To(Equal(1))

		event1 := <-sub1
		event2 := <-sub2
		event := <-cc.ChaincodeEventsChannel

		Expect(event1.Payload).To(Equal(event.Payload))
		Expect(event2.Payload).To(Equal(event.Payload))
		Expect(event1.Payload).To(Equal(event2.Payload))

		Expect(len(cc.ChaincodeEventsChannel)).To(Equal(0))
		Expect(len(sub1)).To(Equal(0))
		Expect(len(sub2)).To(Equal(0))
	})

})
