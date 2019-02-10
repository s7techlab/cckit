package mapping_test

import (
	"testing"

	"github.com/hyperledger/fabric/protos/peer"

	"github.com/golang/protobuf/ptypes"

	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/examples/cpaper/testdata"
	"github.com/s7techlab/cckit/state"

	"github.com/s7techlab/cckit/examples/cpaper"

	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/identity"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var (
	actors   identity.Actors
	cPaperCC *testcc.MockStub
	err      error
)
var _ = Describe(`Mapping`, func() {

	BeforeSuite(func() {
		actors, err = identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
			`owner`: `s7techlab.pem`,
		}, examplecert.Content)

		Expect(err).To(BeNil())

		//Create commercial papers chaincode mock - protobuf based schema
		cPaperCC = testcc.NewMockStub(`cpapers`, cpaper.NewCC())
		cPaperCC.From(actors[`owner`]).Init()

	})

	Describe(`Protobuf based schema`, func() {
		It("Allow to add data to chaincode state", func(done Done) {

			events := cPaperCC.EventSubscription()
			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, &testdata.CPapers[0]))

			Expect(<-events).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: cpaper.EventIssueCommercialPaper,
				Payload:   testcc.MustProtoMarshal(&testdata.CPapers[0]),
			}))

			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, &testdata.CPapers[1]))
			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, &testdata.CPapers[2]))

			close(done)
		}, 0.2)

		It("Disallow to insert entries with same keys", func() {
			expectcc.ResponseError(cPaperCC.Invoke(`issue`, &testdata.CPapers[0]))
		})

		It("Allow to get entry list", func() {
			cpapers := expectcc.PayloadIs(cPaperCC.Query(`list`), &[]schema.CommercialPaper{}).([]schema.CommercialPaper)
			Expect(len(cpapers)).To(Equal(3))
			Expect(cpapers[0].Issuer).To(Equal(testdata.CPapers[0].Issuer))
			Expect(cpapers[0].PaperNumber).To(Equal(testdata.CPapers[0].PaperNumber))
		})

		It("Allow to get entry raw protobuf", func() {
			issue := testdata.CPapers[0]
			cpaperProtoFromCC := cPaperCC.Query(`get`, issue.Issuer, issue.PaperNumber).Payload

			stateCpaper := &schema.CommercialPaper{
				Issuer:       issue.Issuer,
				PaperNumber:  issue.PaperNumber,
				Owner:        issue.Issuer,
				IssueDate:    issue.IssueDate,
				MaturityDate: issue.MaturityDate,
				FaceValue:    issue.FaceValue,
				State:        schema.CommercialPaper_ISSUED, // initial state
			}
			cPaperProto, _ := proto.Marshal(stateCpaper)
			Expect(cpaperProtoFromCC).To(Equal(cPaperProto))
		})

		It("Allow update data in chaincode state", func() {
			cpaper := testdata.CPapers[0]
			expectcc.ResponseOk(cPaperCC.Invoke(`buy`, &schema.BuyCommercialPaper{
				Issuer:       cpaper.Issuer,
				PaperNumber:  cpaper.PaperNumber,
				CurrentOwner: cpaper.Issuer,
				NewOwner:     `some-new-owner`,
				Price:        cpaper.FaceValue - 10,
				PurchaseDate: ptypes.TimestampNow(),
			}))

			cpaperFromCC := expectcc.PayloadIs(cPaperCC.Query(`get`, cpaper.Issuer, cpaper.PaperNumber), &schema.CommercialPaper{}).(*schema.CommercialPaper)

			// state is updated
			Expect(cpaperFromCC.State).To(Equal(schema.CommercialPaper_TRADING))
			Expect(cpaperFromCC.Owner).To(Equal(`some-new-owner`))
		})

		It("Allow to delete entry", func() {
			expectcc.ResponseOk(cPaperCC.Invoke(`delete`, testdata.CPapers[0].Issuer, testdata.CPapers[0].PaperNumber))
			cpapers := expectcc.PayloadIs(cPaperCC.Invoke(`list`), &[]schema.CommercialPaper{}).([]schema.CommercialPaper)

			Expect(len(cpapers)).To(Equal(2))
			expectcc.ResponseError(cPaperCC.Invoke(`get`, testdata.CPapers[0].Issuer, testdata.CPapers[0].PaperNumber), state.ErrKeyNotFound)
		})
	})

})
