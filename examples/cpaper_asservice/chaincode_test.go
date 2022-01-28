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
	mockstub, ccEnc   *testcc.MockStub

	encKey      = make([]byte, 32)
	mockstubEnc *enctest.MockStub

	issuePayload = &cpaper_asservice.IssueCommercialPaper{
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

		mockstub = testcc.NewMockStub(`cpaper_as_service`, ccImpl)
		ccEnc = testcc.NewMockStub(`cpaper_as_service_encrypted`, ccEncImpl)

		// all queries/invokes arguments to cc will be encrypted
		mockstubEnc = enctest.NewMockStub(ccEnc, encKey)

		identity, err = testcc.IdentityFromFile(MspName, `./testdata/admin.pem`, ioutil.ReadFile)
		Expect(err).NotTo(HaveOccurred())
		// Init chaincode with admin identity
		expectcc.ResponseOk(
			mockstub.From(identity).Init())

		//js, _ := json.Marshal(issuePayload)
		//fmt.Println(string(js))
	})

	Describe("Commercial Paper lifecycle", func() {
		It("Allow issuer to issue new commercial paper", func(done Done) {

			expectcc.ResponseOk(mockstub.Invoke(cpaper_asservice.CPaperServiceChaincode_Issue, issuePayload))

			// Validate event has been emitted with the transaction data
			expectcc.EventStringerEqual(<-mockstub.ChaincodeEventsChannel,
				`IssueCommercialPaper`, issuePayload)

			// Clear events channel after a test case that emits an event
			mockstub.ClearEvents()
			close(done)
		}, 0.1)

		It("Allow issuer to get commercial paper by composite primary key", func() {
			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_Get, &cpaper_asservice.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			})

			paper := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaper{}).(*cpaper_asservice.CommercialPaper)

			Expect(paper.Issuer).To(Equal(IssuerName))
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(cpaper_asservice.CommercialPaper_STATE_ISSUED))
			Expect(paper.PaperNumber).To(Equal("0001"))
			Expect(paper.FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow issuer to get commercial paper by unique key", func() {
			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_GetByExternalId, &cpaper_asservice.ExternalId{Id: "EXT0001"})

			paper := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaper{}).(*cpaper_asservice.CommercialPaper)

			Expect(paper.Issuer).To(Equal(IssuerName))
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(cpaper_asservice.CommercialPaper_STATE_ISSUED))
			Expect(paper.PaperNumber).To(Equal("0001"))
			Expect(paper.FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow issuer to get a list of commercial papers", func() {
			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_List, &empty.Empty{})

			papers := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaperList{}).(*cpaper_asservice.CommercialPaperList)

			Expect(len(papers.Items)).To(BeNumerically("==", 1))
			Expect(papers.Items[0].Issuer).To(Equal(IssuerName))
			Expect(papers.Items[0].Owner).To(Equal(IssuerName))
			Expect(papers.Items[0].State).To(Equal(cpaper_asservice.CommercialPaper_STATE_ISSUED))
			Expect(papers.Items[0].PaperNumber).To(Equal("0001"))
			Expect(papers.Items[0].FaceValue).To(BeNumerically("==", 100000))
		})

		It("Allow buyer to buy commercial paper", func() {
			buyTransactionData := &cpaper_asservice.BuyCommercialPaper{
				Issuer:       IssuerName,
				PaperNumber:  "0001",
				CurrentOwner: IssuerName,
				NewOwner:     BuyerName,
				Price:        95000,
				PurchaseDate: ptypes.TimestampNow(),
			}

			expectcc.ResponseOk(mockstub.Invoke(cpaper_asservice.CPaperServiceChaincode_Buy, buyTransactionData))

			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_Get, &cpaper_asservice.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			})

			paper := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaper{}).(*cpaper_asservice.CommercialPaper)

			Expect(paper.Owner).To(Equal(BuyerName))
			Expect(paper.State).To(Equal(cpaper_asservice.CommercialPaper_STATE_TRADING))

			expectcc.EventStringerEqual(<-mockstub.ChaincodeEventsChannel,
				`BuyCommercialPaper`, buyTransactionData)
			mockstub.ClearEvents()
		})

		It("Allow buyer to redeem commercial paper", func() {
			redeemTransactionData := &cpaper_asservice.RedeemCommercialPaper{
				Issuer:         IssuerName,
				PaperNumber:    "0001",
				RedeemingOwner: BuyerName,
				RedeemDate:     ptypes.TimestampNow(),
			}

			expectcc.ResponseOk(mockstub.Invoke(cpaper_asservice.CPaperServiceChaincode_Redeem, redeemTransactionData))

			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_Get, &cpaper_asservice.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			})

			paper := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaper{}).(*cpaper_asservice.CommercialPaper)
			Expect(paper.Owner).To(Equal(IssuerName))
			Expect(paper.State).To(Equal(cpaper_asservice.CommercialPaper_STATE_REDEEMED))

			expectcc.EventStringerEqual(<-mockstub.ChaincodeEventsChannel,
				`RedeemCommercialPaper`, redeemTransactionData)

			mockstub.ClearEvents()
		})

		It("Allow issuer to delete commercial paper", func() {
			expectcc.ResponseOk(mockstub.Invoke(cpaper_asservice.CPaperServiceChaincode_Delete, &cpaper_asservice.CommercialPaperId{
				Issuer:      IssuerName,
				PaperNumber: "0001",
			}))

			// Validate there are 0 Commercial Papers in the world state
			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_List, &empty.Empty{})
			papers := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaperList{}).(*cpaper_asservice.CommercialPaperList)

			Expect(len(papers.Items)).To(BeNumerically("==", 0))
		})
	})

	Describe("Commercial Paper Encrypted lifecycle", func() {
		It("Allow issuer to issue new commercial paper", func(done Done) {

			expectcc.ResponseOk(mockstubEnc.Invoke(cpaper_asservice.CPaperServiceChaincode_Issue, issuePayload))

			// Validate event has been emitted with the transaction data, and event name and payload is encrypted
			expectcc.EventStringerEqual(mockstubEnc.LastEvent(),
				`IssueCommercialPaper`, issuePayload)

			// Clear events channel after a test case that emits an event
			mockstub.ClearEvents()
			close(done)
		}, 0.1)
	})
})
