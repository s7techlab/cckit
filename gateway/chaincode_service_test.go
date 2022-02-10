package gateway_test

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/convert"
	cpservice "github.com/s7techlab/cckit/examples/cpaper_asservice"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/testdata"
	"github.com/s7techlab/cckit/extensions/encryption"
	enctest "github.com/s7techlab/cckit/extensions/encryption/testing"
	"github.com/s7techlab/cckit/gateway"
	idtestdata "github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	"github.com/s7techlab/cckit/testing/gomega"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gateway suite")
}

const (
	Channel       = `my_channel`
	ChaincodeName = `commercial_paper`
)

var (
	ctx = gateway.ContextWithSigner(
		context.Background(),
		idtestdata.Certificates[0].MustIdentity(idtestdata.DefaultMSP),
	)
)

var _ = Describe(`Gateway`, func() {

	Context(`Chaincode service without options`, func() {

		var (
			ccService         *gateway.ChaincodeService
			ccInstanceService *gateway.ChaincodeInstanceService
			cPaperGateway     *cpservice.CPaperServiceGateway
			mockStub          *testcc.MockStub
		)

		It("Init", func() {
			ccImpl, err := cpservice.NewCC()
			Expect(err).NotTo(HaveOccurred())

			mockStub = testcc.NewMockStub(ChaincodeName, ccImpl)
			peer := testcc.NewPeer().WithChannel(Channel, mockStub)

			ccService = gateway.NewChaincodeService(peer)
			ccInstanceService = ccService.InstanceService(
				&gateway.ChaincodeLocator{Channel: Channel, Chaincode: ChaincodeName},
				gateway.WithEventResolver(cpservice.EventMappings))

			// "sdk" for deal with cpaper chaincode
			cPaperGateway = cpservice.NewCPaperServiceGateway(peer, Channel, ChaincodeName)
		})

		Context(`Direct calls`, func() {

			It("Require  to provide chaincode locator", func() {
				_, err := ccService.Query(ctx, &gateway.ChaincodeQueryRequest{
					Input: &gateway.ChaincodeInput{
						Args: [][]byte{[]byte(`List`), {}},
					},
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`invalid field Locator: message must exist`))
			})
		})

		Context(`Chaincode gateway`, func() {

			It("Allow to get empty commercial paper list", func() {
				pp, err := cPaperGateway.List(ctx, &empty.Empty{})
				Expect(err).NotTo(HaveOccurred())
				Expect(pp.Items).To(HaveLen(0))
			})

			It("Invoke chaincode with 'tx waiter' in context", func() {
				ctx = gateway.ContextWithTxWaiter(ctx, "all")
				_, err := cPaperGateway.Issue(ctx, testdata.Issue1)
				Expect(err).NotTo(HaveOccurred())
			})

			It("Invoke chaincode with custom identity in context", func() {
				signer := idtestdata.Certificates[1].MustIdentity(idtestdata.DefaultMSP)
				ctx = gateway.ContextWithDefaultSigner(ctx, signer)

				_, err := cPaperGateway.Delete(ctx, testdata.Id1)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context(`Events`, func() {

			It(`allow to get events as LIST  (by default from block 0 to current channel height) `, func(done Done) {
				events, err := ccInstanceService.Events(ctx, &gateway.ChaincodeInstanceEventsRequest{})

				Expect(err).NotTo(HaveOccurred())
				Expect(events.Items).To(HaveLen(1)) // 1 event on issue

				e := events.Items[0]
				var (
					eventObj  interface{}
					eventJson string
				)
				eventObj, err = convert.FromBytes(e.Event.Payload, &cpservice.IssueCommercialPaper{})
				Expect(err).NotTo(HaveOccurred())

				eventJson, err = (&jsonpb.Marshaler{EmitDefaults: true, OrigName: true}).
					MarshalToString(eventObj.(proto.Message))
				Expect(err).NotTo(HaveOccurred())
				Expect(e.Event.EventName).To(Equal(`IssueCommercialPaper`))
				Expect(e.Payload).NotTo(BeNil())
				Expect(e.Payload.Value).To(Equal([]byte(eventJson))) // check event resolving
				close(done)
			})

			It(`allow to get 0 events as LIST from chaincode with incorrect event name filter`, func(done Done) {
				events, err := ccInstanceService.Events(ctx, &gateway.ChaincodeInstanceEventsRequest{
					EventName: []string{`________IssueCommercialPaper______`},
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(events.Items).To(HaveLen(0)) // 1 event on issue

				close(done)
			})

			It(`allow to get events from block 0 to current channel height AS STREAM`, func(done Done) {
				ctxWithCancel, cancel := context.WithCancel(ctx)
				stream := gateway.NewChaincodeEventServerStream(ctxWithCancel)

				go func() {
					req := &gateway.ChaincodeEventsStreamRequest{
						Locator: &gateway.ChaincodeLocator{
							Channel:   Channel,
							Chaincode: ChaincodeName,
						},
						FromBlock: &gateway.BlockLimit{Num: 0},
						ToBlock:   &gateway.BlockLimit{Num: 0},
					}

					err := ccService.EventsStream(req, &gateway.ChaincodeEventsServer{ServerStream: stream})

					Expect(err).NotTo(HaveOccurred())
				}()

				var e *gateway.ChaincodeEvent
				err := stream.Recv(e)
				Expect(err).NotTo(HaveOccurred())
				cancel()
				close(done)
			}, 1)
		})

	})

	Context(`Chaincode service with encrypted chaincode`, func() {

		var (
			ccService  *gateway.ChaincodeService
			ccInstance *gateway.ChaincodeInstanceService

			cPaperGateway                  *cpservice.CPaperServiceGateway
			cPaperGatewayWithoutEncryption *cpservice.CPaperServiceGateway

			encMockStub *enctest.MockStub

			encryptOpts []gateway.OptFunc
		)

		It("Init", func() {
			ccImpl, err := cpservice.NewCCEncrypted()
			Expect(err).NotTo(HaveOccurred())

			encMockStub = enctest.NewMockStub(testcc.NewMockStub(ChaincodeName, ccImpl))

			encryptOpts = []gateway.OptFunc{
				gateway.WithEncryption(encMockStub.EncKey),
				// Event resolver should be AFTER encryption / decryption middleware
				gateway.WithEventResolver(cpservice.EventMappings),
			}

			// "sdk" for deal with cpaper chaincode
			peer := testcc.NewPeer().WithChannel(Channel, encMockStub.MockStub)
			ccService = gateway.NewChaincodeService(peer)

			locator := &gateway.ChaincodeLocator{
				Channel:   Channel,
				Chaincode: ChaincodeName,
			}

			ccInstance = ccService.InstanceService(locator, encryptOpts...)

			cPaperGateway = cpservice.NewCPaperServiceGateway(peer, Channel, ChaincodeName, encryptOpts...)
			cPaperGatewayWithoutEncryption = cpservice.NewCPaperServiceGateway(
				peer, Channel, ChaincodeName,
				gateway.WithEventResolver(cpservice.EventMappings))
		})

		Context(`Chaincode gateway`, func() {

			It("Disallow to query chaincode without encryption data", func() {
				_, err := cPaperGatewayWithoutEncryption.List(ctx, &empty.Empty{})
				Expect(err).To(gomega.ErrorIs(encryption.ErrKeyNotDefinedInTransientMap))

			})

			It("Allow to get empty commercial paper list", func() {
				pp, err := cPaperGateway.List(ctx, &empty.Empty{})
				Expect(err).NotTo(HaveOccurred())
				Expect(pp.Items).To(HaveLen(0))
			})

			It("Invoke chaincode", func() {
				issued, err := cPaperGateway.Issue(ctx, testdata.Issue1)
				Expect(err).NotTo(HaveOccurred())
				Expect(issued.Issuer).To(Equal(testdata.Issue1.Issuer))
			})

			It("Query chaincode", func() {
				issued, err := cPaperGateway.Get(ctx, testdata.Id1)
				Expect(err).NotTo(HaveOccurred())
				Expect(issued.Issuer).To(Equal(testdata.Issue1.Issuer))
			})

		})

		Context(`Events`, func() {

			It(`allow to get encrypted events as LIST  (by default from block 0 to current channel height) `, func(done Done) {
				events, err := ccInstance.Events(ctx, &gateway.ChaincodeInstanceEventsRequest{})

				Expect(err).NotTo(HaveOccurred())
				Expect(events.Items).To(HaveLen(1)) // 1 event on issue

				e := events.Items[0]
				var (
					eventObj  interface{}
					eventJson string
				)
				eventObj, err = convert.FromBytes(e.Event.Payload, &cpservice.IssueCommercialPaper{})
				Expect(err).NotTo(HaveOccurred())

				eventJson, err = (&jsonpb.Marshaler{EmitDefaults: true, OrigName: true}).
					MarshalToString(eventObj.(proto.Message))
				Expect(err).NotTo(HaveOccurred())
				Expect(e.Event.EventName).To(Equal(`IssueCommercialPaper`))

				Expect(e.Payload.Value).To(Equal([]byte(eventJson))) // check event resolving
				close(done)
			})

			It(`allow to get encrypted events from block 0 to current channel height AS STREAM`, func(done Done) {
				ctxWithCancel, cancel := context.WithCancel(ctx)
				stream := gateway.NewChaincodeEventServerStream(ctxWithCancel)

				go func() {
					req := &gateway.ChaincodeInstanceEventsStreamRequest{
						FromBlock: &gateway.BlockLimit{Num: 0},
						ToBlock:   &gateway.BlockLimit{Num: 0},
					}

					err := ccInstance.EventsStream(req, &gateway.ChaincodeEventsServer{ServerStream: stream})

					Expect(err).NotTo(HaveOccurred())
				}()

				var e *gateway.ChaincodeEvent
				err := stream.Recv(e)
				Expect(err).NotTo(HaveOccurred())
				cancel()
				close(done)
			}, 1)

		})

	})
	Context(`Chaincode instance service`, func() {

		var (
			ccInstanceService *gateway.ChaincodeInstanceService
		)

		It("Init", func() {
			ccImpl, err := cpservice.NewCC()
			Expect(err).NotTo(HaveOccurred())

			// peer imitation
			peer := testcc.NewPeer().WithChannel(Channel, testcc.NewMockStub(ChaincodeName, ccImpl))
			ccInstanceService = gateway.NewChaincodeInstanceService(peer, &gateway.ChaincodeLocator{
				Channel:   Channel,
				Chaincode: ChaincodeName,
			})
		})

		Context(`Direct calls`, func() {

			It("Allow to get empty commercial paper list", func() {
				resp, err := ccInstanceService.Query(ctx, &gateway.ChaincodeInstanceQueryRequest{
					Input: &gateway.ChaincodeInput{
						Args: [][]byte{[]byte(`List`), {}},
					},
				})

				Expect(err).NotTo(HaveOccurred())
				cPaperList := testcc.MustProtoUnmarshal(resp.Payload, &cpservice.CommercialPaperList{}).(*cpservice.CommercialPaperList)
				Expect(cPaperList.Items).To(HaveLen(0))
			})

			It("Invoke chaincode", func() {

				_, err := ccInstanceService.Invoke(ctx, &gateway.ChaincodeInstanceInvokeRequest{
					Input: &gateway.ChaincodeInput{
						Args: [][]byte{[]byte(`Issue`), testcc.MustProtoMarshal(testdata.Issue1)},
					}})
				Expect(err).NotTo(HaveOccurred())
			})

		})

		Context(`Chaincode gateway`, func() {

			It(`allow to get events as LIST from chaincode instance service (by default from block 0 to current channel height) `, func(done Done) {
				events, err := ccInstanceService.Events(ctx, &gateway.ChaincodeInstanceEventsRequest{
					EventName: []string{`IssueCommercialPaper`},
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(events.Items).To(HaveLen(1)) // 1 event on issue

				close(done)
			})

		})
	})

	Context(`Cross chaincode invoker`, func() {

	})

})
