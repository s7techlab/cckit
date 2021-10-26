package service_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/s7techlab/cckit/examples/cpaper_asservice"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/schema"
	cpservice "github.com/s7techlab/cckit/examples/cpaper_asservice/service"
	"github.com/s7techlab/cckit/gateway/service"
	"github.com/s7techlab/cckit/gateway/service/mock"
	idtestdata "github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mockstub Suite")
}

const (
	Channel       = `my_channel`
	ChaincodeName = `commercial_paper`
)

var (
	cPaperService *mock.ChaincodeService
	cPaperGateway *cpservice.CPaperGateway

	ctx = service.ContextWithSigner(
		context.Background(),
		idtestdata.Certificates[0].MustIdentity(idtestdata.DefaultMSP),
	)
)

var _ = Describe(`Service`, func() {

	It("Init", func() {

		ccImpl, err := cpaper_asservice.NewCC()
		Expect(err).NotTo(HaveOccurred())

		// peer imitation
		cPaperService = mock.New()
		cPaperService.Peer.WithChannel(Channel, testcc.NewMockStub(ChaincodeName, ccImpl))

		// "sdk" for deal with cpaper chaincode
		cPaperGateway = cpservice.NewCPaperGateway(cPaperService, Channel, ChaincodeName)
	})

	It("Allow to get empty commercial paper list", func() {
		pp, err := cPaperGateway.List(ctx, &empty.Empty{})
		Expect(err).NotTo(HaveOccurred())
		Expect(pp.Items).To(HaveLen(0))
	})

	It("Invoke chaincode with 'tx waiter' in context", func() {
		ctx = service.ContextWithTxWaiter(ctx, "all")
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
		ctx = service.ContextWithDefaultSigner(ctx, signer)

		_, err := cPaperGateway.Delete(ctx, &schema.CommercialPaperId{
			Issuer:      "issuer",
			PaperNumber: "1337",
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("Allow to imitate error while access to peer", func() {
		cPaperService.Invoker = mock.FailChaincode(ChaincodeName)

		_, err := cPaperGateway.List(ctx, &empty.Empty{})
		Expect(err).To(HaveOccurred())
	})
})
