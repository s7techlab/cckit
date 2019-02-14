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
				EventName: `IssueCommercialPaper`,
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
			cp := testdata.CPapers[0]
			cpaperProtoFromCC := cPaperCC.Query(`get`, &schema.CommercialPaperId{Issuer: cp.Issuer, PaperNumber: cp.PaperNumber}).Payload

			stateCpaper := &schema.CommercialPaper{
				Issuer:       cp.Issuer,
				PaperNumber:  cp.PaperNumber,
				Owner:        cp.Issuer,
				IssueDate:    cp.IssueDate,
				MaturityDate: cp.MaturityDate,
				FaceValue:    cp.FaceValue,
				State:        schema.CommercialPaper_ISSUED, // initial state
			}
			cPaperProto, _ := proto.Marshal(stateCpaper)
			Expect(cpaperProtoFromCC).To(Equal(cPaperProto))
		})

		It("Allow update data in chaincode state", func() {
			cp := testdata.CPapers[0]
			expectcc.ResponseOk(cPaperCC.Invoke(`buy`, &schema.BuyCommercialPaper{
				Issuer:       cp.Issuer,
				PaperNumber:  cp.PaperNumber,
				CurrentOwner: cp.Issuer,
				NewOwner:     `some-new-owner`,
				Price:        cp.FaceValue - 10,
				PurchaseDate: ptypes.TimestampNow(),
			}))

			cpaperFromCC := expectcc.PayloadIs(
				cPaperCC.Query(`get`, &schema.CommercialPaperId{Issuer: cp.Issuer, PaperNumber: cp.PaperNumber}),
				&schema.CommercialPaper{}).(*schema.CommercialPaper)

			// state is updated
			Expect(cpaperFromCC.State).To(Equal(schema.CommercialPaper_TRADING))
			Expect(cpaperFromCC.Owner).To(Equal(`some-new-owner`))
		})

		It("Allow to delete entry", func() {

			cp := testdata.CPapers[0]
			toDelete := &schema.CommercialPaperId{Issuer: cp.Issuer, PaperNumber: cp.PaperNumber}

			expectcc.ResponseOk(cPaperCC.Invoke(`delete`, toDelete))
			cpapers := expectcc.PayloadIs(cPaperCC.Invoke(`list`), &[]schema.CommercialPaper{}).([]schema.CommercialPaper)

			Expect(len(cpapers)).To(Equal(2))
			expectcc.ResponseError(cPaperCC.Invoke(`get`, toDelete), state.ErrKeyNotFound)
		})
	})

})
