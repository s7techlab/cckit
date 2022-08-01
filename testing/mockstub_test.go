package testing_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"context"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/s7techlab/cckit/examples/cars"
	idtestdata "github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
	"github.com/s7techlab/cckit/testing/testdata"
)

func TestMockStub(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MockStub Suite")
}

var (
	ids = idtestdata.MustIdentities(idtestdata.Certificates, idtestdata.DefaultMSP)

	Authority = ids[0]
)

const (
	Channel              = `my_channel`
	CarsChaincode        = `cars`
	CarsProxyChaincode   = `cars_proxy`
	TxIsolationChaincode = `tx_isolation`
)

var _ = Describe(`Testing`, func() {

	//Create chaincode mocks
	cc := testcc.NewMockStub(CarsChaincode, cars.New())
	ccproxy := testcc.NewMockStub(CarsProxyChaincode, cars.NewProxy(Channel, CarsChaincode))

	txIsolationCC := testcc.NewMockStub(TxIsolationChaincode, testdata.NewTxIsolationCC())

	// ccproxy can invoke cc and vice versa
	mockedPeer := testcc.NewPeer().WithChannel(Channel, cc, ccproxy)

	Describe(`MockStub`, func() {

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
				TxId:        cc.LastTxID,
				ChaincodeId: cc.Name,
				EventName:   cars.CarRegisteredEvent,
				Payload:     resp.Payload,
			}))

			close(done)
		}, 0.2)

		It("Allow to use multiple events subscriptions", func(done Done) {
			Expect(len(cc.ChaincodeEventsChannel)).To(Equal(0))

			sub1, closer1 := cc.EventSubscription()
			sub2, closer2 := cc.EventSubscription()

			Expect(len(sub1)).To(Equal(0))
			Expect(len(sub2)).To(Equal(0))

			resp := expectcc.ResponseOk(cc.From(Authority).Invoke(`carRegister`, cars.Payloads[2]))

			Expect(len(cc.ChaincodeEventsChannel)).To(Equal(1))
			Expect(len(sub1)).To(Equal(1))
			Expect(len(sub2)).To(Equal(1))

			Expect(<-sub1).To(BeEquivalentTo(&peer.ChaincodeEvent{
				ChaincodeId: cc.Name,
				TxId:        cc.LastTxID,
				EventName:   cars.CarRegisteredEvent,
				Payload:     resp.Payload,
			}))

			Expect(<-sub2).To(BeEquivalentTo(&peer.ChaincodeEvent{
				ChaincodeId: cc.Name,
				TxId:        cc.LastTxID,
				EventName:   cars.CarRegisteredEvent,
				Payload:     resp.Payload,
			}))

			Expect(<-cc.ChaincodeEventsChannel).To(BeEquivalentTo(&peer.ChaincodeEvent{
				ChaincodeId: cc.Name,
				TxId:        cc.LastTxID,
				EventName:   cars.CarRegisteredEvent,
				Payload:     resp.Payload,
			}))

			Expect(len(cc.ChaincodeEventsChannel)).To(Equal(0))
			Expect(len(sub1)).To(Equal(0))
			Expect(len(sub2)).To(Equal(0))

			_ = closer1()
			_ = closer2()

			close(done)
		}, 0.2)

	})

	Describe(`MockStub invoker`, func() {

		It("Allow to invoke mocked chaincode ", func(done Done) {
			ctx := context.Background()

			events, closer, err := mockedPeer.Events(ctx, Channel, CarsChaincode, Authority)
			Expect(err).NotTo(HaveOccurred())

			// double check interface Peer
			resp, _, err := mockedPeer.Invoke(
				ctx, Channel, CarsChaincode,
				[][]byte{[]byte(`carRegister`), testcc.MustJSONMarshal(cars.Payloads[3])},
				Authority, nil, ``)
			Expect(err).NotTo(HaveOccurred())

			carFromCC := testcc.MustConvertFromBytes(resp.Payload, &cars.Car{}).(cars.Car)

			Expect(carFromCC.Id).To(Equal(cars.Payloads[3].Id))
			Expect(carFromCC.Title).To(Equal(cars.Payloads[3].Title))

			event := <-events
			Expect(event.Event()).To(BeEquivalentTo(&peer.ChaincodeEvent{
				ChaincodeId: cc.Name,
				TxId:        cc.LastTxID,
				EventName:   cars.CarRegisteredEvent,
				Payload:     resp.Payload,
			}))

			_ = closer()
			close(done)

		}, 0.3)

		It("Allow to query mocked chaincode ", func() {
			resp, err := mockedPeer.Query(
				context.Background(),
				Channel, CarsChaincode,
				[][]byte{[]byte(`carGet`), []byte(cars.Payloads[3].Id)},
				Authority, nil)
			Expect(err).NotTo(HaveOccurred())

			carFromCC := testcc.MustConvertFromBytes(resp.Payload, &cars.Car{}).(cars.Car)

			Expect(carFromCC.Id).To(Equal(cars.Payloads[3].Id))
			Expect(carFromCC.Title).To(Equal(cars.Payloads[3].Title))
		})

		It("Allow to query mocked chaincode from chaincode", func() {
			resp, err := mockedPeer.Query(
				context.Background(),
				Channel, CarsProxyChaincode,
				[][]byte{[]byte(`carGet`), []byte(cars.Payloads[3].Id)},
				Authority, nil)
			Expect(err).NotTo(HaveOccurred())

			carFromCC := testcc.MustConvertFromBytes(resp.Payload, &cars.Car{}).(cars.Car)

			Expect(carFromCC.Id).To(Equal(cars.Payloads[3].Id))
			Expect(carFromCC.Title).To(Equal(cars.Payloads[3].Title))
		})

		It("Should return error when unknown channel provided", func() {
			_, err := mockedPeer.Query(
				context.Background(),
				"unknown-channel", CarsProxyChaincode,
				[][]byte{[]byte(`carGet`), []byte(cars.Payloads[3].Id)},
				Authority, nil)
			Expect(err).To(HaveOccurred())

		})

		It("Should return error when unknown carID provided", func() {
			_, err := mockedPeer.Query(
				context.Background(),
				Channel, CarsProxyChaincode,
				[][]byte{[]byte(`carGet`), []byte("unknown_car_id")},
				Authority, nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe(`Tx isolation`, func() {
		It("Read after write returns empty", func() {
			res := txIsolationCC.Invoke(testdata.TxIsolationReadAfterWrite)
			Expect(int(res.Status)).To(Equal(shim.OK))
			Expect(res.Payload).To(Equal([]byte{}))
		})

		It("Read after delete returns state entry", func() {
			res := txIsolationCC.Invoke(testdata.TxIsolationReadAfterDelete)
			Expect(int(res.Status)).To(Equal(shim.OK))
			Expect(res.Payload).To(Equal(testdata.Value1))
		})

	})
})
