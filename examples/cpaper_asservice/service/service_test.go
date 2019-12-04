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
	testcc "github.com/s7techlab/cckit/testing"
)

func TestCommercialPaperService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commercial Paper Suite")
}

var (
	CPaper = service.New()

	// service testing util
	hdl, ctx = testcc.NewTxHandler(`Commercial paper`)

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
		hdl.From(Issuer).Init(func() (interface{}, error) {
			return owner.SetFromCreator(ctx)
		}).Expect().HasError(nil)
	})

	It("Allow issuer to issue new commercial paper", func() {
		hdl.From(Issuer).Invoke(func() (interface{}, error) {
			return CPaper.Issue(ctx, issue)
		}).Expect().Is(cpaperInState)
	})

	It("Allow issuer to get commercial paper by composite primary key", func() {
		hdl.Invoke(func() (interface{}, error) {
			return CPaper.Get(ctx, id)
		}).Expect().Is(cpaperInState)
	})

	It("Allow issuer to get commercial paper by unique key", func() {
		res := hdl.Invoke(func() (interface{}, error) {
			return CPaper.GetByExternalId(ctx, &schema.ExternalId{
				Id: issue.ExternalId,
			})
		})

		Expect(res.Result).To(StringerEqual(cpaperInState))
	})

	It("Allow issuer to get a list of commercial papers", func() {
		hdl.Invoke(func() (interface{}, error) {
			return CPaper.List(ctx, &empty.Empty{})
		}).Expect().Is(&schema.CommercialPaperList{
			Items: []*schema.CommercialPaper{cpaperInState},
		})
	})

	It("Allow buyer to buy commercial paper", func() {
		hdl.From(Buyer).Invoke(func() (interface{}, error) {
			return CPaper.Buy(ctx, buy)
		}).Expect().ProduceEvent(`BuyCommercialPaper`, buy)

		newState := proto.Clone(cpaperInState).(*schema.CommercialPaper)
		newState.Owner = buy.NewOwner
		newState.State = schema.CommercialPaper_TRADING

		hdl.Invoke(func() (interface{}, error) {
			return CPaper.Get(ctx, id)
		}).Expect().Is(newState)
	})

	It("Allow buyer to redeem commercial paper", func() {
		hdl.Invoke(func() (interface{}, error) {
			return CPaper.Redeem(ctx, redeem)
		}).Expect().ProduceEvent(`RedeemCommercialPaper`, redeem)

		newState := proto.Clone(cpaperInState).(*schema.CommercialPaper)
		newState.State = schema.CommercialPaper_REDEEMED

		hdl.Invoke(func() (interface{}, error) {
			return CPaper.Get(ctx, id)
		}).Expect().Is(newState)
	})

	It("Allow issuer to delete commercial paper", func() {
		hdl.From(Issuer).Invoke(func() (interface{}, error) {
			return CPaper.Delete(ctx, id)
		}).Expect().HasError(nil)

		hdl.Invoke(func() (interface{}, error) {
			return CPaper.List(ctx, &empty.Empty{})
		}).Expect().Is(&schema.CommercialPaperList{
			Items: []*schema.CommercialPaper{},
		})
	})
})
