package cpaper_test

import (
	"testing"
	"time"

	"github.com/s7techlab/cckit/examples/cpaper"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	IssuerName = "SomeIssuer"
	BuyerName  = "SomeBuyer"
)

func TestCommercialPaper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commercial Paper Suite")
}

var _ = Describe(`CommercialPaper`, func() {
	paperChaincode := testcc.NewMockStub(`commercial_paper`, cpaper.NewCC())

	BeforeSuite(func() {
		expectcc.ResponseOk(paperChaincode.Init())
	})

	Describe("Commercial Paper lifecycle", func() {

		It("Allow issuer to issue new commercial paper", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(paperChaincode.Invoke(`issue`, &cpaper.IssueCommercialPaper{
				Issuer:       IssuerName,
				PaperNumber:  "0001",
				IssueDate:    time.Now(),
				MaturityDate: time.Now().Add(time.Hour * 24 * 30 * 6), // Six months later
				FaceValue:    100000,
			}))
		})

		It("Allow issuer to get commercial paper", func() {
			queryResponse := paperChaincode.Query("get", IssuerName, "0001")
			paper := expectcc.PayloadIs(queryResponse, &cpaper.CommercialPaper{}).(cpaper.CommercialPaper)

			Expect(paper.Issuer).To(Equal(IssuerName))
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(cpaper.CommercialPaperIssued))
			Expect(paper.PaperNumber).To(Equal("0001"))
			Expect(paper.FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow issuer to get a list of commercial papers", func() {
			queryResponse := paperChaincode.Query("list")
			papers := expectcc.PayloadIs(queryResponse, &[]cpaper.CommercialPaper{}).([]cpaper.CommercialPaper)

			Expect(len(papers)).To(BeNumerically("==", 1))
			Expect(papers[0].Issuer).To(Equal(IssuerName))
			Expect(papers[0].Owner).To(Equal(IssuerName))
			Expect(papers[0].State).To(Equal(cpaper.CommercialPaperIssued))
			Expect(papers[0].PaperNumber).To(Equal("0001"))
			Expect(papers[0].FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow buyer to buy commercial paper", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(paperChaincode.Invoke(`buy`, &cpaper.BuyCommercialPaper{
				Issuer:       IssuerName,
				PaperNumber:  "0001",
				CurrentOwner: IssuerName,
				NewOwner:     BuyerName,
				Price:        95000,
				PurchaseDate: time.Now(),
			}))

			queryResponse := paperChaincode.Query("get", IssuerName, "0001")
			paper := expectcc.PayloadIs(queryResponse, &cpaper.CommercialPaper{}).(cpaper.CommercialPaper)
			Expect(paper.Owner).To(Equal(BuyerName))
			Expect(paper.State).To(Equal(cpaper.CommercialPaperTrading))
		})

		It("Allow buyer to redeem commercial paper", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(paperChaincode.Invoke(`redeem`, &cpaper.RedeemCommercialPaper{
				Issuer:         IssuerName,
				PaperNumber:    "0001",
				RedeemingOwner: BuyerName,
				RedeemDate:     time.Now(),
			}))

			queryResponse := paperChaincode.Query("get", IssuerName, "0001")
			paper := expectcc.PayloadIs(queryResponse, &cpaper.CommercialPaper{}).(cpaper.CommercialPaper)
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(cpaper.CommercialPaperRedeemed))
		})
	})
})
