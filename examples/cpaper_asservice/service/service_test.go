package service_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/schema"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/service"
	"github.com/s7techlab/cckit/router"
	testcc "github.com/s7techlab/cckit/testing"
	"github.com/s7techlab/cckit/testing/expect"
)

func TestCommercialPaperService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commercial Paper Suite")
}

var (
	CPaperSvc = service.New()

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

	cc = testcc.NewCCService(`Commercial paper`)
)

var _ = Describe(`Commercial paper service`, func() {

	It("Allow issuer to issue new commercial paper", func() {
		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.Issue(ctx, issue)
		})).Is(cpaperInState)
	})

	It("Allow issuer to get commercial paper by composite primary key", func() {
		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.Get(ctx, id)
		})).Is(cpaperInState)
	})

	It("Allow issuer to get commercial paper by unique key", func() {
		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.GetByExternalId(ctx, &schema.ExternalId{
				Id: issue.ExternalId,
			})
		})).Is(cpaperInState)
	})

	It("Allow issuer to get a list of commercial papers", func() {
		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.List(ctx, &empty.Empty{})
		})).Is(&schema.CommercialPaperList{
			Items: []*schema.CommercialPaper{cpaperInState},
		})
	})

	It("Allow buyer to buy commercial paper", func() {
		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.Buy(ctx, buy)
		})).ProduceEvent(`BuyCommercialPaper`, buy)

		newState := proto.Clone(cpaperInState).(*schema.CommercialPaper)
		newState.Owner = buy.NewOwner
		newState.State = schema.CommercialPaper_TRADING

		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.Get(ctx, id)
		})).Is(newState)
	})

	It("Allow buyer to redeem commercial paper", func() {
		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.Redeem(ctx, redeem)
		})).ProduceEvent(`RedeemCommercialPaper`, redeem)

		newState := proto.Clone(cpaperInState).(*schema.CommercialPaper)
		newState.State = schema.CommercialPaper_REDEEMED

		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.Get(ctx, id)
		})).Is(newState)
	})

	It("Allow issuer to delete commercial paper", func() {
		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.Delete(ctx, id)
		}))

		expect.SvcResponse(cc.Exec(func(ctx router.Context) (interface{}, error) {
			return CPaperSvc.List(ctx, &empty.Empty{})
		})).Is(&schema.CommercialPaperList{
			Items: []*schema.CommercialPaper{},
		})
	})
})
