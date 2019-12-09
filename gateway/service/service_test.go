package service_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/s7techlab/cckit/examples/cpaper_asservice"
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
		idtestdata.Certificates[0].MustIdentity(idtestdata.DefaultMSP))
)

var _ = Describe(`Service`, func() {

	It("Init", func() {

		ccImpl, err := cpaper_asservice.NewCC()
		Expect(err).NotTo(HaveOccurred())

		// peer imitation
		cPaperService = mock.New().WithChannel(Channel, testcc.NewMockStub(ChaincodeName, ccImpl))

		// "sdk" for deal with cpaper chaincode
		cPaperGateway = cpservice.NewCPaperGateway(cPaperService, Channel, ChaincodeName)
	})

	It("Allow to get empty commercial paper list", func() {
		pp, err := cPaperGateway.List(ctx, &empty.Empty{})
		Expect(err).NotTo(HaveOccurred())
		Expect(pp.Items).To(HaveLen(0))
	})

	It("Allow to imitate error while access to peer", func() {
		cPaperService.Invoker = mock.FailChaincode(ChaincodeName)

		_, err := cPaperGateway.List(ctx, &empty.Empty{})
		Expect(err).To(HaveOccurred())
	})
})
