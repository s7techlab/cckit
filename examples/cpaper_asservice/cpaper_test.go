package cpaper_asservice_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/cpaper_asservice"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/testdata"
	"github.com/s7techlab/cckit/extensions/owner"
	idtestdata "github.com/s7techlab/cckit/identity/testdata"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	testcc "github.com/s7techlab/cckit/testing"
	. "github.com/s7techlab/cckit/testing/gomega"
)

func TestCommercialPaperService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commercial Paper Suite")
}

var (
	CPaper = cpaper_asservice.NewService()

	// service testing util
	cc, ctx = testcc.NewTxHandler(`Commercial paper`)

	ids = idtestdata.MustSamples(idtestdata.Certificates, idtestdata.DefaultMSP)
	// actors
	Issuer = ids[0]
	Buyer  = ids[1]
)

var _ = Describe(`Commercial paper service`, func() {

	It("Allow to init", func() {
		cc.From(Issuer).Init(func(c router.Context) (interface{}, error) {
			return owner.SetFromCreator(c)
		}).Expect().HasError(nil)
	})

	It("Allow issuer to issue new commercial paper", func() {
		cc.From(Issuer).Tx(func() {
			cc.Expect(CPaper.Issue(ctx, testdata.Issue1)).Is(testdata.CPaperInState1)
		})
	})

	// Use TxFunc helper - return closure func() {} with some assertions
	It("Disallow issuer to issue same commercial paper", cc.From(Issuer).TxFunc(func() {
		// Expect helper return TxRes
		// we don't check result payload, error only
		cc.Expect(CPaper.Issue(ctx, testdata.Issue1)).HasError(state.ErrKeyAlreadyExists)
	}))

	It("Allow issuer to get commercial paper by composite primary key", func() {
		cc.Tx(func() {
			// Expect helper, check error is empty and result
			cc.Expect(CPaper.Get(ctx, testdata.Id1)).Is(testdata.CPaperInState1)
		})
	})

	It("Allow issuer to get commercial paper by unique key", func() {
		cc.Tx(func() {
			// without helpers
			res, err := CPaper.GetByExternalId(ctx, &cpaper_asservice.ExternalId{
				Id: testdata.Issue1.ExternalId,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(StringerEqual(testdata.CPaperInState1))
		})

	})

	It("Allow issuer to get a list of commercial papers", func() {
		cc.Tx(func() {
			res, err := CPaper.List(ctx, &empty.Empty{})

			Expect(err).NotTo(HaveOccurred())
			Expect(res.Items).To(HaveLen(1))
			Expect(res.Items[0]).To(StringerEqual(testdata.CPaperInState1))
		})
	})

	It("Allow buyer to buy commercial paper", func() {
		cc.From(Buyer).Tx(func() {
			cc.Expect(CPaper.Buy(ctx, testdata.Buy1)).
				// Produce Event - no error and event name and payload check
				ProduceEvent(`BuyCommercialPaper`, testdata.Buy1)
		})

		newState := proto.Clone(testdata.CPaperInState1).(*cpaper_asservice.CommercialPaper)
		newState.Owner = testdata.Buy1.NewOwner
		newState.State = cpaper_asservice.CommercialPaper_STATE_TRADING

		cc.Tx(func() {
			cc.Expect(CPaper.Get(ctx, testdata.Id1)).Is(newState)
		})
	})

	It("Allow buyer to redeem commercial paper", func() {
		// Invoke example
		cc.Invoke(func(c router.Context) (interface{}, error) {
			return CPaper.Redeem(c, testdata.Redeem1)
		}).Expect().ProduceEvent(`RedeemCommercialPaper`, testdata.Redeem1)

		newState := proto.Clone(testdata.CPaperInState1).(*cpaper_asservice.CommercialPaper)
		newState.State = cpaper_asservice.CommercialPaper_STATE_REDEEMED

		cc.Invoke(func(c router.Context) (interface{}, error) {
			return CPaper.Get(c, testdata.Id1)
		}).Expect().Is(newState)
	})

	It("Allow issuer to delete commercial paper", func() {
		cc.From(Issuer).Tx(func() {
			cc.Expect(CPaper.Delete(ctx, testdata.Id1)).HasNoError()
		})

		cc.Tx(func() {
			res, err := CPaper.List(ctx, &empty.Empty{})

			Expect(err).NotTo(HaveOccurred())
			Expect(res.Items).To(HaveLen(0))
		})
	})
})
