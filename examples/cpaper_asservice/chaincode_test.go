package cpaper_asservice_test

import (
	"crypto/rand"
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hyperledger/fabric/msp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/cpaper_asservice"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/schema"
	s "github.com/s7techlab/cckit/examples/cpaper_asservice/service"
	enctest "github.com/s7techlab/cckit/extensions/encryption/testing"
	"github.com/s7techlab/cckit/router"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
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

var (
	ccImpl, ccEncImpl *router.Chaincode
	err               error
	cc, ccEnc         *testcc.MockStub

	encKey       = make([]byte, 32)
	ccEncWrapped *enctest.MockStub

	issuePayload = &schema.IssueCommercialPaper{
		Issuer:       IssuerName,
		PaperNumber:  "0001",
		IssueDate:    ptypes.TimestampNow(),
		MaturityDate: testcc.MustProtoTimestamp(time.Now().AddDate(0, 2, 0)),
		FaceValue:    100000,
		ExternalId:   "EXT0001",
	}

	identity msp.SigningIdentity
)

var _ = Describe(`CommercialPaper`, func() {

	BeforeSuite(func() {

		_, err = rand.Read(encKey)
		Expect(err).NotTo(HaveOccurred())

		ccImpl, err = cpaper_asservice.NewCC()
		Expect(err).NotTo(HaveOccurred())

		ccEncImpl, err = cpaper_asservice.NewCCEncrypted()
		Expect(err).NotTo(HaveOccurred())

		cc = testcc.NewMockStub(`cpaper_as_service`, ccImpl)
		ccEnc = testcc.NewMockStub(`cpaper_as_service_encrypted`, ccEncImpl)

		// all queries/invokes arguments to cc will be encrypted
		ccEncWrapped = enctest.NewMockStub(ccEnc, encKey)

		identity, err = testcc.IdentityFromFile(MspName, `./testdata/admin.pem`, ioutil.ReadFile)
		Expect(err).NotTo(HaveOccurred())
		// Init chaincode with admin identity
		expectcc.ResponseOk(
			cc.From(identity).Init())

		//js, _ := json.Marshal(issuePayload)
		//fmt.Println(string(js))
	})

	Describe("Commercial Paper lifecycle", func() {
		It("Allow issuer to issue new commercial paper", func(done Done) {

			expectcc.ResponseOk(cc.Invoke(s.CPaperChaincode_Issue, issuePayload))

			// Validate event has been emitted with the transaction data
			expectcc.EventStringerEqual(<-cc.ChaincodeEventsChannel,
				`IssueCommercialPaper`, issuePayload)

			// Clear events channel after a test case that emits an event
			cc.ClearEvents()
			close(done)
		}, 0.1)

		It("Allow issuer to get commercial paper by composite primary key", func() {
			queryResponse := cc.Query(s.CPaperChaincode_Get, &schema.CommercialPaperId{
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
			queryResponse := cc.Query(s.CPaperChaincode_GetByExternalId, &schema.ExternalId{Id: "EXT0001"})

			paper := expectcc.PayloadIs(queryResponse, &schema.CommercialPaper{}).(*schema.CommercialPaper)

			Expect(paper.Issuer).To(Equal(IssuerName))
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(schema.CommercialPaper_ISSUED))
			Expect(paper.PaperNumber).To(Equal("0001"))
			Expect(paper.FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow issuer to get a list of commercial papers", func() {
			queryResponse := cc.Query(s.CPaperChaincode_List, &empty.Empty{})

			papers := expectcc.PayloadIs(queryResponse, &schema.CommercialPaperList{}).(*schema.CommercialPaperList)

			Expect(len(papers.Items)).To(BeNumerically("==", 1))
			Expect(papers.Items[0].Issuer).To(Equal(IssuerName))
			Expect(papers.Items[0].Owner).To(Equal(IssuerName))
			Expect(papers.Items[0].State).To(Equal(schema.CommercialPaper_ISSUED))
			Expect(papers.Items[0].PaperNumber).To(Equal("0001"))
			Expect(papers.Items[0].FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow buyer to buy commercial paper", func() {
			buyTransactionData := &schema.BuyCommercialPaper{
				Issuer:       IssuerName,
				PaperNumber:  "0001",
				CurrentOwner: IssuerName,
				NewOwner:     BuyerName,
				Price:        95000,
				PurchaseDate: ptypes.TimestampNow(),
			}

			expectcc.ResponseOk(cc.Invoke(s.CPaperChaincode_Buy, buyTransactionData))

			queryResponse := cc.Query(s.CPaperChaincode_Get, &schema.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			})

			paper := expectcc.PayloadIs(queryResponse, &schema.CommercialPaper{}).(*schema.CommercialPaper)

			Expect(paper.Owner).To(Equal(BuyerName))
			Expect(paper.State).To(Equal(schema.CommercialPaper_TRADING))

			expectcc.EventStringerEqual(<-cc.ChaincodeEventsChannel,
				`BuyCommercialPaper`, buyTransactionData)
			cc.ClearEvents()
		})

		It("Allow buyer to redeem commercial paper", func() {
			redeemTransactionData := &schema.RedeemCommercialPaper{
				Issuer:         IssuerName,
				PaperNumber:    "0001",
				RedeemingOwner: BuyerName,
				RedeemDate:     ptypes.TimestampNow(),
			}

			expectcc.ResponseOk(cc.Invoke(s.CPaperChaincode_Redeem, redeemTransactionData))

			queryResponse := cc.Query(s.CPaperChaincode_Get, &schema.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			})

			paper := expectcc.PayloadIs(queryResponse, &schema.CommercialPaper{}).(*schema.CommercialPaper)
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(schema.CommercialPaper_REDEEMED))

			expectcc.EventStringerEqual(<-cc.ChaincodeEventsChannel,
				`RedeemCommercialPaper`, redeemTransactionData)

			cc.ClearEvents()
		})

		It("Allow issuer to delete commercial paper", func() {
			expectcc.ResponseOk(cc.Invoke(s.CPaperChaincode_Delete, &schema.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			}))

			// Validate there are 0 Commercial Papers in the world state
			queryResponse := cc.Query(s.CPaperChaincode_List, &empty.Empty{})
			papers := expectcc.PayloadIs(queryResponse, &schema.CommercialPaperList{}).(*schema.CommercialPaperList)

			Expect(len(papers.Items)).To(BeNumerically("==", 0))
		})
	})

	Describe("Commercial Paper Encrypted lifecycle", func() {
		It("Allow issuer to issue new commercial paper", func(done Done) {

			expectcc.ResponseOk(ccEncWrapped.Invoke(s.CPaperChaincode_Issue, issuePayload))

			// Validate event has been emitted with the transaction data, and event name and payload is encrypted
			expectcc.EventStringerEqual(ccEncWrapped.LastEvent(),
				`IssueCommercialPaper`, issuePayload)

			// Clear events channel after a test case that emits an event
			cc.ClearEvents()
			close(done)
		}, 0.1)
	})
})
