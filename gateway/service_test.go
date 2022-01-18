package gateway_test

import (
	"context"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/s7techlab/cckit/convert"
	"github.com/s7techlab/cckit/examples/cpaper_asservice"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/schema"
	cpservice "github.com/s7techlab/cckit/examples/cpaper_asservice/service"
	"github.com/s7techlab/cckit/gateway"
	idtestdata "github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
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

	Context(`Chaincode service`, func() {

		var (
			ccService     *gateway.ChaincodeService
			cPaperGateway *cpservice.CPaperGateway
			mockStub      *testcc.MockStub
		)

		It("Init", func() {
			ccImpl, err := cpaper_asservice.NewCC()
			Expect(err).NotTo(HaveOccurred())

			mockStub = testcc.NewMockStub(ChaincodeName, ccImpl)

			ccService = gateway.NewChaincodeService(
				testcc.NewPeer().WithChannel(Channel, mockStub),
				gateway.WithEventResolver(cpservice.EventMappings),
			)

			// "sdk" for deal with cpaper chaincode
			cPaperGateway = cpservice.NewCPaperGateway(ccService, Channel, ChaincodeName)
		})

		Context(`Direct calls`, func() {

			It("Require  to provide chaincode locator", func() {
				_, err := ccService.Query(ctx, &gateway.ChaincodeInput{
					Args: [][]byte{[]byte(`List`), []byte{}},
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring(`invalid field Chaincode: message must exist`))
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
				_, err := cPaperGateway.Issue(ctx, &schema.IssueCommercialPaper{
					Issuer:       "issuer",
					PaperNumber:  "1337",
					ExternalId:   "228",
					IssueDate:    timestamppb.Now(),
					MaturityDate: timestamppb.Now(),
					FaceValue:    2,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("Invoke chaincode with custom identity in context", func() {
				signer := idtestdata.Certificates[1].MustIdentity(idtestdata.DefaultMSP)
				ctx = gateway.ContextWithDefaultSigner(ctx, signer)

				_, err := cPaperGateway.Delete(ctx, &schema.CommercialPaperId{
					Issuer:      "issuer",
					PaperNumber: "1337",
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It(`allow to get events as LIST  (by default from block 0 to current channel height) `, func(done Done) {
				events, err := ccService.Events(ctx, &gateway.ChaincodeEventsRequest{
					Chaincode: &gateway.ChaincodeLocator{
						Channel:   Channel,
						Chaincode: ChaincodeName,
					},
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(events.Items).To(HaveLen(1)) // 1 event on issue

				e := events.Items[0]
				var (
					eventObj  interface{}
					eventJson string
				)
				eventObj, err = convert.FromBytes(e.Event.Payload, &schema.IssueCommercialPaper{})
				Expect(err).NotTo(HaveOccurred())

				eventJson, err = (&jsonpb.Marshaler{EmitDefaults: true, OrigName: true}).
					MarshalToString(eventObj.(proto.Message))
				Expect(e.Event.EventName).To(Equal(`IssueCommercialPaper`))
				Expect(e.Payload).To(Equal([]byte(eventJson)))
				close(done)
			})

			It(`allow to get events as LIST from chaincode instance service (by default from block 0 to current channel height) `, func(done Done) {
				events, err := ccService.InstanceService(Channel, ChaincodeName).Events(ctx, &gateway.ChaincodeInstanceEventsRequest{})

				Expect(err).NotTo(HaveOccurred())
				Expect(events.Items).To(HaveLen(1)) // 1 event on issue

				close(done)
			})

			It(`allow to get events from block 0 to current channel height AS STREAM`, func(done Done) {
				ctxWithCancel, cancel := context.WithCancel(ctx)
				stream := gateway.NewChaincodeEventServerStream(ctxWithCancel)

				go func() {
					req := &gateway.ChaincodeEventsStreamRequest{
						Chaincode: &gateway.ChaincodeLocator{
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

	Context(`Chaincode instance service`, func() {

		var (
			ccInstanceService *gateway.ChaincodeInstanceService
		)

		It("Init", func() {
			ccImpl, err := cpaper_asservice.NewCC()
			Expect(err).NotTo(HaveOccurred())

			// peer imitation
			peer := testcc.NewPeer().WithChannel(Channel, testcc.NewMockStub(ChaincodeName, ccImpl))
			ccInstanceService = gateway.NewChaincodeInstanceService(peer, Channel, ChaincodeName)
		})

		Context(`Direct calls`, func() {

			It("Allow to get empty commercial paper list", func() {
				resp, err := ccInstanceService.Query(ctx, &gateway.ChaincodeInstanceInput{
					Args: [][]byte{[]byte(`List`), []byte{}},
				})

				Expect(err).NotTo(HaveOccurred())
				cPaperList := testcc.MustProtoUnmarshal(resp.Payload, &schema.CommercialPaperList{}).(*schema.CommercialPaperList)
				Expect(cPaperList.Items).To(HaveLen(0))
			})

		})
	})

})
