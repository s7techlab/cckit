package mock_test

import (
	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/s7techlab/cckit/gateway/mock"

	"context"
	"testing"

	"github.com/s7techlab/cckit/examples/cpaper_asservice"
	cpservice "github.com/s7techlab/cckit/examples/cpaper_asservice/service"
	"github.com/s7techlab/cckit/gateway"
	idtestdata "github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gateway mock suite")
}

const (
	Channel       = `my_channel`
	ChaincodeName = `commercial_paper`
)

var (
	ccService     *mock.ChaincodeService
	cPaperGateway *cpservice.CPaperGateway

	ctx = gateway.ContextWithSigner(
		context.Background(),
		idtestdata.Certificates[0].MustIdentity(idtestdata.DefaultMSP),
	)
)

var _ = Describe(`Service`, func() {

	It("Init", func() {
		ccImpl, err := cpaper_asservice.NewCC()
		Expect(err).NotTo(HaveOccurred())

		// peer imitation
		peer := testcc.NewPeer().WithChannel(Channel, testcc.NewMockStub(ChaincodeName, ccImpl))
		ccService = mock.New(peer)

		// "sdk" for deal with cpaper chaincode
		cPaperGateway = cpservice.NewCPaperGateway(ccService, Channel, ChaincodeName)
	})

	It("Allow to imitate error while access to peer", func() {
		ccService.Invoker = mock.FailChaincode(ChaincodeName)

		_, err := cPaperGateway.List(ctx, &empty.Empty{})
		Expect(err).To(HaveOccurred())
	})
})
