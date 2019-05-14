package cpaper

import (
	"testing"
	"time"

	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	ISSUER_NAME = "SomeIssuer"
	BUYER_NAME  = "SomeIssuer"
)

func TestCommercialPaper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cars Suite")
}

var _ = Describe(`CommercialPaper`, func() {
	paperChaincode := testcc.NewMockStub(`commercial_paper`, NewCC())

	BeforeSuite(func() {
		expectcc.ResponseOk(paperChaincode.Init())
	})

	Describe("Commercial Paper lifecycle", func() {

		It("Allow issuer to issue new commercial paper", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(paperChaincode.Invoke(`issue`, &IssueCommercialPaper{
				Issuer:       ISSUER_NAME,
				PaperNumber:  "0001",
				IssueDate:    time.Now(),
				MaturityDate: time.Now().Add(time.Hour * 24 * 30 * 6), // Six months later
				FaceValue:    100000,
			}))
		})

		It("Allow issuer to get commercial paper", func() {
			queryResponse := paperChaincode.Query("get", ISSUER_NAME, "0001")
			paper := expectcc.PayloadIs(queryResponse, &CommercialPaper{}).(CommercialPaper)

			Expect(paper.Issuer).To(Equal(ISSUER_NAME))
			Expect(paper.Owner).To(Equal(ISSUER_NAME))
			Expect(paper.State).To(Equal(CommercialPaperIssued))
			Expect(paper.PaperNumber).To(Equal("0001"))
			Expect(paper.FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow issuer to get a list of commercial papers", func() {
			queryResponse := paperChaincode.Query("list")
			papers := expectcc.PayloadIs(queryResponse, &[]CommercialPaper{}).([]CommercialPaper)

			Expect(len(papers)).To(BeNumerically("==", 1))
			Expect(papers[0].Issuer).To(Equal(ISSUER_NAME))
			Expect(papers[0].Owner).To(Equal(ISSUER_NAME))
			Expect(papers[0].State).To(Equal(CommercialPaperIssued))
			Expect(papers[0].PaperNumber).To(Equal("0001"))
			Expect(papers[0].FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow buyer to buy commercial paper", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(paperChaincode.Invoke(`buy`, &BuyCommercialPaper{
				Issuer:       ISSUER_NAME,
				PaperNumber:  "0001",
				CurrentOwner: ISSUER_NAME,
				NewOwner:     BUYER_NAME,
				Price:        95000,
				PurchaseDate: time.Now(),
			}))

			queryResponse := paperChaincode.Query("get", ISSUER_NAME, "0001")
			paper := expectcc.PayloadIs(queryResponse, &CommercialPaper{}).(CommercialPaper)
			Expect(paper.Owner).To(Equal(BUYER_NAME))
			Expect(paper.State).To(Equal(CommercialPaperTrading))
		})

		It("Allow buyer to redeem commercial paper", func() {
			//invoke chaincode method from authority actor
			expectcc.ResponseOk(paperChaincode.Invoke(`redeem`, &RedeemCommercialPaper{
				Issuer:         ISSUER_NAME,
				PaperNumber:    "0001",
				RedeemingOwner: BUYER_NAME,
				RedeemDate:     time.Now(),
			}))

			queryResponse := paperChaincode.Query("get", ISSUER_NAME, "0001")
			paper := expectcc.PayloadIs(queryResponse, &CommercialPaper{}).(CommercialPaper)
			Expect(paper.Owner).To(Equal(ISSUER_NAME))
			Expect(paper.State).To(Equal(CommercialPaperRedeemed))
		})
	})
})
