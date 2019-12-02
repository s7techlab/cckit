package testing_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/hyperledger/fabric/protos/peer"
	"github.com/s7techlab/cckit/examples/cars"
	idtestdata "github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
	"github.com/s7techlab/hlf-sdk-go/api"
)

func TestMockstub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mockstub Suite")
}

var (
	ids = idtestdata.MustIdentities(idtestdata.Certificates, idtestdata.DefaultMSP)

	Authority = ids[0]
	Someone   = ids[1]
)

const (
	Channel            = `my_channel`
	ChaincodeName      = `cars`
	ChaincodeProxyName = `cars_proxy`
)

var _ = Describe(`Testing`, func() {

	//Create chaincode mocks
	cc := testcc.NewMockStub(ChaincodeName, cars.New())
	ccproxy := testcc.NewMockStub(ChaincodeProxyName, cars.NewProxy(Channel, ChaincodeName))

	// ccproxy can invoke cc and vice versa
	mockedPeer := testcc.NewPeer().WithChannel(Channel, cc, ccproxy)

	Describe(`Mockstub`, func() {

		It("Allow to init chaincode", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(cc.From(Authority).Init()) // init chaincode from authority
		})

		It("Allow to get last event while chaincode invoke ", func() {

			expectcc.ResponseOk(cc.From(Authority).Invoke(`carRegister`, cars.Payloads[0]))
			event := expectcc.EventPayloadIs(cc.ChaincodeEvent, &cars.Car{}).(cars.Car)

			Expect(cc.ChaincodeEvent.EventName).To(Equal(cars.CarRegisteredEvent))
			Expect(event.Id).To(Equal(cars.Payloads[0].Id))

			Expect(len(cc.ChaincodeEventsChannel)).To(Equal(1))

		})

		It("Allow to clear events channel", func() {
			cc.ClearEvents()
			Expect(len(cc.ChaincodeEventsChannel)).To(Equal(0))

		})

		It("Allow to get events via events channel", func(done Done) {
			resp := expectcc.ResponseOk(cc.From(Authority).Invoke(`carRegister`, cars.Payloads[1]))

			Expect(<-cc.ChaincodeEventsChannel).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: cars.CarRegisteredEvent,
				Payload:   resp.Payload,
			}))

			close(done)
		}, 0.2)

		It("Allow to use multiple events subscriptions", func(done Done) {
			Expect(len(cc.ChaincodeEventsChannel)).To(Equal(0))

			sub1 := cc.EventSubscription()
			sub2 := cc.EventSubscription()

			Expect(len(sub1)).To(Equal(0))
			Expect(len(sub2)).To(Equal(0))

			resp := expectcc.ResponseOk(cc.From(Authority).Invoke(`carRegister`, cars.Payloads[2]))

			Expect(len(cc.ChaincodeEventsChannel)).To(Equal(1))
			Expect(len(sub1)).To(Equal(1))
			Expect(len(sub2)).To(Equal(1))

			Expect(<-sub1).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: cars.CarRegisteredEvent,
				Payload:   resp.Payload,
			}))

			Expect(<-sub2).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: cars.CarRegisteredEvent,
				Payload:   resp.Payload,
			}))

			Expect(<-cc.ChaincodeEventsChannel).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: cars.CarRegisteredEvent,
				Payload:   resp.Payload,
			}))

			Expect(len(cc.ChaincodeEventsChannel)).To(Equal(0))
			Expect(len(sub1)).To(Equal(0))
			Expect(len(sub2)).To(Equal(0))

			close(done)
		}, 0.2)

	})

	Describe(`Mockstub invoker`, func() {

		It("Allow to invoke mocked chaincode ", func(done Done) {
			ctx := context.Background()

			events, err := mockedPeer.Subscribe(ctx, Authority, Channel, ChaincodeName)
			Expect(err).NotTo(HaveOccurred())

			// double check interface api.Invoker
			resp, _, err := interface{}(mockedPeer).(api.Invoker).Invoke(
				ctx, Authority, Channel, ChaincodeName, `carRegister`,
				[][]byte{testcc.MustJSONMarshal(cars.Payloads[3])}, nil)
			Expect(err).NotTo(HaveOccurred())

			carFromCC := testcc.MustConvertFromBytes(resp.Payload, &cars.Car{}).(cars.Car)

			Expect(carFromCC.Id).To(Equal(cars.Payloads[3].Id))
			Expect(carFromCC.Title).To(Equal(cars.Payloads[3].Title))

			Expect(<-events.Events()).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: cars.CarRegisteredEvent,
				Payload:   resp.Payload,
			}))

			close(done)

		}, 0.3)

		It("Allow to query mocked chaincode ", func() {
			resp, err := mockedPeer.Query(
				context.Background(), Authority, Channel, ChaincodeName,
				`carGet`, [][]byte{[]byte(cars.Payloads[3].Id)}, nil)
			Expect(err).NotTo(HaveOccurred())

			carFromCC := testcc.MustConvertFromBytes(resp.Payload, &cars.Car{}).(cars.Car)

			Expect(carFromCC.Id).To(Equal(cars.Payloads[3].Id))
			Expect(carFromCC.Title).To(Equal(cars.Payloads[3].Title))
		})

		It("Allow to query mocked chaincode from chaincode", func() {
			resp, err := mockedPeer.Query(
				context.Background(), Authority, Channel, ChaincodeProxyName,
				`carGet`, [][]byte{[]byte(cars.Payloads[3].Id)}, nil)
			Expect(err).NotTo(HaveOccurred())

			carFromCC := testcc.MustConvertFromBytes(resp.Payload, &cars.Car{}).(cars.Car)

			Expect(carFromCC.Id).To(Equal(cars.Payloads[3].Id))
			Expect(carFromCC.Title).To(Equal(cars.Payloads[3].Title))
		})
	})
})
