package service_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/s7techlab/cckit/testing/gomega"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/schema"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/service"
	"github.com/s7techlab/cckit/extensions/owner"
	idtestdata "github.com/s7techlab/cckit/identity/testdata"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/state"
	testcc "github.com/s7techlab/cckit/testing"
)

func TestCommercialPaperService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commercial Paper Suite")
}

var (
	CPaper = service.New()

	// service testing util
	cc, ctx = testcc.NewTxHandler(`Commercial paper`)

	ids = idtestdata.MustIdentities(idtestdata.Certificates, idtestdata.DefaultMSP)
	// actors
	Issuer = ids[0]
	Buyer  = ids[1]

	// payloads
	id = &schema.CommercialPaperId{
		Issuer:      "SomeIssuer",
		PaperNumber: "0001",
	}

	issue = &schema.IssueCommercialPaper{
		Issuer:       id.Issuer,
		PaperNumber:  id.PaperNumber,
		IssueDate:    ptypes.TimestampNow(),
		MaturityDate: testcc.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
		FaceValue:    100000,
		ExternalId:   "EXT0001",
	}

	buy = &schema.BuyCommercialPaper{
		Issuer:       id.Issuer,
		PaperNumber:  id.PaperNumber,
		CurrentOwner: id.Issuer,
		NewOwner:     "SomeBuyer",
		Price:        95000,
		PurchaseDate: ptypes.TimestampNow(),
	}

	redeem = &schema.RedeemCommercialPaper{
		Issuer:         id.Issuer,
		PaperNumber:    id.PaperNumber,
		RedeemingOwner: buy.NewOwner,
		RedeemDate:     ptypes.TimestampNow(),
	}

	cpaperInState = &schema.CommercialPaper{
		Issuer:       id.Issuer,
		Owner:        id.Issuer,
		State:        schema.CommercialPaper_ISSUED,
		PaperNumber:  id.PaperNumber,
		FaceValue:    issue.FaceValue,
		IssueDate:    issue.IssueDate,
		MaturityDate: issue.MaturityDate,
		ExternalId:   issue.ExternalId,
	}
)

var _ = Describe(`Commercial paper service`, func() {

	It("Allow to init", func() {
		cc.From(Issuer).Init(func(c router.Context) (interface{}, error) {
			return owner.SetFromCreator(c)
		}).Expect().HasError(nil)
	})

	It("Allow issuer to issue new commercial paper", func() {
		cc.From(Issuer).Tx(func() {
			cc.Expect(CPaper.Issue(ctx, issue)).Is(cpaperInState)
		})
	})

	It("Disallow issuer to issue same commercial paper", func() {
		cc.From(Issuer).Tx(func() {
			_, err := CPaper.Issue(ctx, issue)

			Expect(err).To(ErrorIs(state.ErrKeyAlreadyExists))
		})
	})

	It("Allow issuer to get commercial paper by composite primary key", func() {
		cc.Tx(func() {
			cc.Expect(CPaper.Get(ctx, id)).Is(cpaperInState)
		})
	})

	It("Allow issuer to get commercial paper by unique key", func() {
		cc.Tx(func() {
			res, err := CPaper.GetByExternalId(ctx, &schema.ExternalId{
				Id: issue.ExternalId,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(res).To(StringerEqual(cpaperInState))
		})

	})

	It("Allow issuer to get a list of commercial papers", func() {
		cc.Tx(func() {
			res, err := CPaper.List(ctx, &empty.Empty{})

			Expect(err).NotTo(HaveOccurred())
			Expect(res.Items).To(HaveLen(1))
			Expect(res.Items[0]).To(StringerEqual(cpaperInState))
		})
	})

	It("Allow buyer to buy commercial paper", func() {
		cc.From(Buyer).Tx(func() {
			cc.Expect(CPaper.Buy(ctx, buy)).
				// Produce Event - no error also
				ProduceEvent(`BuyCommercialPaper`, buy)
		})

		newState := proto.Clone(cpaperInState).(*schema.CommercialPaper)
		newState.Owner = buy.NewOwner
		newState.State = schema.CommercialPaper_TRADING

		cc.Tx(func() {
			cc.Expect(CPaper.Get(ctx, id)).Is(newState)
		})
	})

	It("Allow buyer to redeem commercial paper", func() {
		// Invoke example
		// Invoke
		cc.Invoke(func(c router.Context) (interface{}, error) {
			return CPaper.Redeem(c, redeem)
		}).Expect().ProduceEvent(`RedeemCommercialPaper`, redeem)

		newState := proto.Clone(cpaperInState).(*schema.CommercialPaper)
		newState.State = schema.CommercialPaper_REDEEMED

		cc.Invoke(func(c router.Context) (interface{}, error) {
			return CPaper.Get(c, id)
		}).Expect().Is(newState)
	})

	It("Allow issuer to delete commercial paper", func() {
		cc.From(Issuer).Tx(func() {
			_, err := CPaper.Delete(ctx, id)
			Expect(err).NotTo(HaveOccurred())
		})

		cc.Tx(func() {
			res, err := CPaper.List(ctx, &empty.Empty{})

			Expect(err).NotTo(HaveOccurred())
			Expect(res.Items).To(HaveLen(0))
		})
	})
})
