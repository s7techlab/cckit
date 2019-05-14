package cpaper_extended

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/s7techlab/cckit/examples/cpaper_extended/schema"
	"github.com/s7techlab/cckit/examples/cpaper_extended/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
	"path"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	MspName = "msp"

	IssuerName = "SomeIssuer"
	BuyerName  = "SomeBuyer"
)

func TestCommercialPaper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commercial Paper Suite")
}

var _ = Describe(`CommercialPaper`, func() {
	paperChaincode := testcc.NewMockStub(`commercial_paper`, NewCC())

	BeforeSuite(func() {
		// Init chaincode with admin identity
		expectcc.ResponseOk(
			paperChaincode.
				From(testdata.GetTestIdentity(MspName, path.Join("testdata", "admin", "admin.pem"))).
				Init())
	})

	Describe("Commercial Paper lifecycle", func() {
		It("Allow issuer to issue new commercial paper", func() {
			expectcc.ResponseOk(paperChaincode.Invoke(`issue`, &schema.IssueCommercialPaper{
				Issuer:       IssuerName,
				PaperNumber:  "0001",
				IssueDate:    ptypes.TimestampNow(),
				MaturityDate: testcc.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
				FaceValue:    100000,
				ExternalId:   "EXT0001",
			}))
		})

		It("Allow issuer to get commercial paper by composite primary key", func() {
			queryResponse := paperChaincode.Query("get", &schema.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			})
			paper := expectcc.PayloadIs(queryResponse, &schema.CommercialPaper{}).(*schema.CommercialPaper)

			Expect(paper.Issuer).To(Equal(IssuerName))
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(schema.CommercialPaper_ISSUED))
			Expect(paper.PaperNumber).To(Equal("0001"))
			Expect(paper.FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow issuer to get commercial paper by unique key", func() {
			queryResponse := paperChaincode.Query("getByExternalId", "EXT0001")
			paper := expectcc.PayloadIs(queryResponse, &schema.CommercialPaper{}).(*schema.CommercialPaper)

			Expect(paper.Issuer).To(Equal(IssuerName))
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(schema.CommercialPaper_ISSUED))
			Expect(paper.PaperNumber).To(Equal("0001"))
			Expect(paper.FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow issuer to get a list of commercial papers", func() {
			queryResponse := paperChaincode.Query("list")
			papers := expectcc.PayloadIs(queryResponse, &schema.CommercialPaperList{}).(*schema.CommercialPaperList)

			Expect(len(papers.Items)).To(BeNumerically("==", 1))
			Expect(papers.Items[0].Issuer).To(Equal(IssuerName))
			Expect(papers.Items[0].Owner).To(Equal(IssuerName))
			Expect(papers.Items[0].State).To(Equal(schema.CommercialPaper_ISSUED))
			Expect(papers.Items[0].PaperNumber).To(Equal("0001"))
			Expect(papers.Items[0].FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow buyer to buy commercial paper", func() {
			expectcc.ResponseOk(paperChaincode.Invoke(`buy`, &schema.BuyCommercialPaper{
				Issuer:       IssuerName,
				PaperNumber:  "0001",
				CurrentOwner: IssuerName,
				NewOwner:     BuyerName,
				Price:        95000,
				PurchaseDate: ptypes.TimestampNow(),
			}))

			queryResponse := paperChaincode.Query("get", &schema.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			})
			paper := expectcc.PayloadIs(queryResponse, &schema.CommercialPaper{}).(*schema.CommercialPaper)
			Expect(paper.Owner).To(Equal(BuyerName))
			Expect(paper.State).To(Equal(schema.CommercialPaper_TRADING))
		})

		It("Allow buyer to redeem commercial paper", func() {
			expectcc.ResponseOk(paperChaincode.Invoke(`redeem`, &schema.RedeemCommercialPaper{
				Issuer:         IssuerName,
				PaperNumber:    "0001",
				RedeemingOwner: BuyerName,
				RedeemDate:     ptypes.TimestampNow(),
			}))

			queryResponse := paperChaincode.Query("get", &schema.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			})
			paper := expectcc.PayloadIs(queryResponse, &schema.CommercialPaper{}).(*schema.CommercialPaper)
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(schema.CommercialPaper_REDEEMED))
		})

		It("Allow issuer to redeem delete commercial paper", func() {
			expectcc.ResponseOk(paperChaincode.Invoke(`delete`, &schema.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			}))

			// Validate there are 0 Commercial Papers in the world state
			queryResponse := paperChaincode.Query("list")
			papers := expectcc.PayloadIs(queryResponse, &schema.CommercialPaperList{}).(*schema.CommercialPaperList)

			Expect(len(papers.Items)).To(BeNumerically("==", 0))
		})
	})
})
