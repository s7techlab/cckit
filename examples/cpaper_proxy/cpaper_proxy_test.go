package cpaper_proxy_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cpservice "github.com/s7techlab/cckit/examples/cpaper_asservice"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/testdata"
	"github.com/s7techlab/cckit/examples/cpaper_proxy"
	"github.com/s7techlab/cckit/extensions/crosscc"
	"github.com/s7techlab/cckit/gateway"
	idtestdata "github.com/s7techlab/cckit/identity/testdata"
	testcc "github.com/s7techlab/cckit/testing"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gateway suite")
}

const (
	Channel1             = `channel1`
	Channel2             = `channel2`
	ChaincodeCPaperProxy = `cpaper_proxy`
	ChaincodeCPaper      = `cpaper`
)

var (
	ctx = gateway.ContextWithSigner(
		context.Background(),
		idtestdata.Certificates[0].MustIdentity(idtestdata.DefaultMSP),
	)

	Id = &cpaper_proxy.Id{
		Issuer:      testdata.Id1.Issuer,
		PaperNumber: testdata.Id1.PaperNumber,
	}
)

var _ = Describe(`Chaincode service resolving`, func() {

	Context(`Services in one chaincode`, func() {

		var (
			cPaperGateway      *cpservice.CPaperServiceGateway
			cPaperProxyGateway *cpaper_proxy.CPaperProxyServiceGateway
		)

		It("Init", func() {
			cPaperProxyCC, err := cpaper_proxy.NewCCWithLocalCPaper()
			Expect(err).NotTo(HaveOccurred())

			mockStub := testcc.NewMockStub(ChaincodeCPaperProxy, cPaperProxyCC)
			peer := testcc.NewPeer().WithChannel(Channel1, mockStub)

			// both gw are looking to same channel / chaincode
			cPaperProxyGateway = cpaper_proxy.NewCPaperProxyServiceGateway(peer, Channel1, ChaincodeCPaperProxy)
			cPaperGateway = cpservice.NewCPaperServiceGateway(peer, Channel1, ChaincodeCPaperProxy)
		})

		It("Initial empty result", func() {

			_, err := cPaperProxyGateway.GetFromCPaper(ctx, Id)
			Expect(err).To(HaveOccurred())
		})

		It("Allow to add cpaper", func() {
			_, err := cPaperGateway.Issue(ctx, testdata.Issue1)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Allow to get data from cpaper service via proxy service", func() {
			_, err := cPaperProxyGateway.GetFromCPaper(ctx, Id)
			Expect(err).NotTo(HaveOccurred())
		})

	})

	Context(`Services in separate chaincodes`, func() {

		var (
			cPaperGateway      *cpservice.CPaperServiceGateway
			cPaperProxyGateway *cpaper_proxy.CPaperProxyServiceGateway

			cPaperProxyCCSettingGateway *crosscc.SettingServiceGateway
		)

		It("Init", func() {
			cPaperProxyCC, err := cpaper_proxy.NewCCWithRemoteCPaper()
			Expect(err).NotTo(HaveOccurred())

			cPaperCC, err := cpservice.NewCC()
			Expect(err).NotTo(HaveOccurred())

			mockStub1 := testcc.NewMockStub(ChaincodeCPaperProxy, cPaperProxyCC)
			mockStub2 := testcc.NewMockStub(ChaincodeCPaper, cPaperCC)
			peer := testcc.NewPeer().
				WithChannel(Channel1, mockStub1).
				WithChannel(Channel2, mockStub2)

			// gw are looking to separate channels / chaincodes
			cPaperProxyGateway = cpaper_proxy.NewCPaperProxyServiceGateway(peer, Channel1, ChaincodeCPaperProxy)
			cPaperGateway = cpservice.NewCPaperServiceGateway(peer, Channel2, ChaincodeCPaper)
			cPaperProxyCCSettingGateway = crosscc.NewSettingServiceGateway(peer, Channel1, ChaincodeCPaperProxy)
		})

		It("Initial empty result", func() {
			_, err := cPaperProxyGateway.GetFromCPaper(ctx, Id)
			Expect(err).To(HaveOccurred())

		})

		It("Allow to add cpaper", func() {
			_, err := cPaperGateway.Issue(ctx, testdata.Issue1)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Disallow to get data from cpaper service via proxy service while chaincode locator setting is empty", func() {
			_, err := cPaperProxyGateway.GetFromCPaper(ctx, Id)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`chaincode locator not found`))
		})

		It("Allow to set setting for chaincode", func() {
			//_, err := cPaperGateway.Issue(ctx, testdata.Issue1)
			//Expect(err).NotTo(HaveOccurred())

			_, err := cPaperProxyCCSettingGateway.ServiceLocatorSet(ctx, &crosscc.ServiceLocatorSetRequest{
				Service:   cPaperGateway.ServiceDef().Name(),
				Channel:   Channel2,
				Chaincode: ChaincodeCPaper,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("Allow to get data from cpaper service via proxy service", func() {
			_, err := cPaperProxyGateway.GetFromCPaper(ctx, Id)
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
