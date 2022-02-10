package cpaper_asservice_test

import (
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hyperledger/fabric/msp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/s7techlab/cckit/examples/cpaper_asservice"
	"github.com/s7techlab/cckit/examples/cpaper_asservice/testdata"
	enctest "github.com/s7techlab/cckit/extensions/encryption/testing"
	"github.com/s7techlab/cckit/router"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

const (
	MspName = "msp"
)

func TestCommercialPaper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Commercial Paper Suite")
}

var (
	ccImpl, ccEncImpl *router.Chaincode
	err               error
	mockstub          *testcc.MockStub
	mockstubEnc       *enctest.MockStub

	identity msp.SigningIdentity
)

var _ = Describe(`CommercialPaper`, func() {

	BeforeSuite(func() {

		ccImpl, err = cpaper_asservice.NewCC()
		Expect(err).NotTo(HaveOccurred())

		ccEncImpl, err = cpaper_asservice.NewCCEncrypted()
		Expect(err).NotTo(HaveOccurred())

		mockstub = testcc.NewMockStub(`cpaper_as_service`, ccImpl)
		// all queries/invokes arguments to cc will be encrypted
		mockstubEnc = enctest.NewMockStub(testcc.NewMockStub(`cpaper_as_service_encrypted`, ccEncImpl))

		identity, err = testcc.IdentityFromFile(MspName, `./testdata/admin.pem`, ioutil.ReadFile)
		Expect(err).NotTo(HaveOccurred())
		// Init chaincode with admin identity
		expectcc.ResponseOk(mockstub.From(identity).Init())
	})

	Describe("Commercial Paper lifecycle", func() {
		It("Allow issuer to issue new commercial paper", func() {
			expectcc.ResponseOk(mockstub.Invoke(cpaper_asservice.CPaperServiceChaincode_Issue, testdata.Issue1))

			// Validate event has been emitted with the transaction data
			expectcc.EventStringerEqual(mockstub.ChaincodeEvent, `IssueCommercialPaper`, testdata.Issue1)
		}, 0.1)

		It("Allow issuer to get commercial paper by composite primary key", func() {
			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_Get, testdata.Id1)

			paper := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaper{}).(*cpaper_asservice.CommercialPaper)

			Expect(paper.Issuer).To(Equal(testdata.Issue1.Issuer))
			Expect(paper.Owner).To(Equal(testdata.Issue1.Issuer))
			Expect(paper.State).To(Equal(cpaper_asservice.CommercialPaper_STATE_ISSUED))
			Expect(paper.PaperNumber).To(Equal(testdata.Issue1.PaperNumber))
			Expect(paper.FaceValue).To(BeNumerically("==", testdata.Issue1.FaceValue))
		})

		It("Allow issuer to get commercial paper by unique key", func() {
			queryResponse := mockstub.Query(
				cpaper_asservice.CPaperServiceChaincode_GetByExternalId, testdata.ExternalId1)

			paper := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaper{}).(*cpaper_asservice.CommercialPaper)

			Expect(paper.Issuer).To(Equal(testdata.Issue1.Issuer))
			Expect(paper.Owner).To(Equal(testdata.Issue1.Issuer))
			Expect(paper.State).To(Equal(cpaper_asservice.CommercialPaper_STATE_ISSUED))
			Expect(paper.PaperNumber).To(Equal(testdata.Issue1.PaperNumber))
			Expect(paper.FaceValue).To(BeNumerically("==", testdata.Issue1.FaceValue))
		})

		It("Allow issuer to get a list of commercial papers", func() {
			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_List, &empty.Empty{})

			papers := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaperList{}).(*cpaper_asservice.CommercialPaperList)

			Expect(len(papers.Items)).To(BeNumerically("==", 1))
			Expect(papers.Items[0].Issuer).To(Equal(testdata.Issue1.Issuer))
			Expect(papers.Items[0].Owner).To(Equal(testdata.Issue1.Issuer))
			Expect(papers.Items[0].State).To(Equal(cpaper_asservice.CommercialPaper_STATE_ISSUED))
			Expect(papers.Items[0].PaperNumber).To(Equal(testdata.Issue1.PaperNumber))
			Expect(papers.Items[0].FaceValue).To(BeNumerically("==", testdata.Issue1.FaceValue))
		})

		It("Allow buyer to buy commercial paper", func() {
			expectcc.ResponseOk(mockstub.Invoke(cpaper_asservice.CPaperServiceChaincode_Buy, testdata.Buy1))
			expectcc.EventStringerEqual(mockstub.ChaincodeEvent, `BuyCommercialPaper`, testdata.Buy1)

			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_Get, testdata.Id1)

			paper := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaper{}).(*cpaper_asservice.CommercialPaper)

			Expect(paper.Owner).To(Equal(testdata.Buy1.NewOwner))
			Expect(paper.State).To(Equal(cpaper_asservice.CommercialPaper_STATE_TRADING))

		})

		It("Allow buyer to redeem commercial paper", func() {
			expectcc.ResponseOk(mockstub.Invoke(cpaper_asservice.CPaperServiceChaincode_Redeem, testdata.Redeem1))
			expectcc.EventStringerEqual(mockstub.ChaincodeEvent, `RedeemCommercialPaper`, testdata.Redeem1)

			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_Get, testdata.Id1)

			paper := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaper{}).(*cpaper_asservice.CommercialPaper)
			Expect(paper.Owner).To(Equal(testdata.Issue1.Issuer))
			Expect(paper.State).To(Equal(cpaper_asservice.CommercialPaper_STATE_REDEEMED))
		})

		It("Allow issuer to delete commercial paper", func() {
			expectcc.ResponseOk(mockstub.Invoke(cpaper_asservice.CPaperServiceChaincode_Delete, testdata.Id1))

			// Validate there are 0 Commercial Papers in the world state
			queryResponse := mockstub.Query(cpaper_asservice.CPaperServiceChaincode_List, &empty.Empty{})
			papers := expectcc.PayloadIs(queryResponse, &cpaper_asservice.CommercialPaperList{}).(*cpaper_asservice.CommercialPaperList)

			Expect(len(papers.Items)).To(BeNumerically("==", 0))
		})
	})

	Describe("Commercial Paper Encrypted lifecycle", func() {
		It("Allow issuer to issue new commercial paper", func() {

			expectcc.ResponseOk(mockstubEnc.Invoke(cpaper_asservice.CPaperServiceChaincode_Issue, testdata.Issue1))

			// Validate event has been emitted with the transaction data, and event name and payload is encrypted
			expectcc.EventStringerEqual(mockstubEnc.LastEvent(),
				`IssueCommercialPaper`, testdata.Issue1)
		})
	})
})
