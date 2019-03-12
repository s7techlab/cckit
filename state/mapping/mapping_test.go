package mapping_test

import (
	"testing"

	state_schema "github.com/s7techlab/cckit/state/schema"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric/protos/peer"

	"github.com/s7techlab/cckit/state"
	"github.com/s7techlab/cckit/state/mapping/testdata/schema"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/examples/cpaper"
	cpaper_schema "github.com/s7techlab/cckit/examples/cpaper/schema"
	cpaper_testdata "github.com/s7techlab/cckit/examples/cpaper/testdata"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/state/mapping/testdata"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var (
	actors                identity.Actors
	cPaperCC, complexIdCC *testcc.MockStub
	err                   error
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

		complexIdCC = testcc.NewMockStub(`complexid`, testdata.NewComplexIdCC())
		complexIdCC.From(actors[`owner`]).Init()
	})

	Describe(`Commercial paper, protobuf based schema`, func() {

		var cpaper1 = &cpaper_testdata.CPapers[0]
		var cpaper2 = &cpaper_testdata.CPapers[1]
		var cpaper3 = &cpaper_testdata.CPapers[2]

		It("Allow to add data to chaincode state", func(done Done) {

			events := cPaperCC.EventSubscription()
			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, cpaper1))

			Expect(<-events).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: `IssueCommercialPaper`,
				Payload:   testcc.MustProtoMarshal(cpaper1),
			}))

			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, cpaper2))
			expectcc.ResponseOk(cPaperCC.Invoke(`issue`, cpaper3))

			close(done)
		}, 0.2)

		It("Disallow to insert entries with same keys", func() {
			expectcc.ResponseError(cPaperCC.Invoke(`issue`, cpaper1))
		})

		It("Allow to get entry list", func() {
			cpapers := expectcc.PayloadIs(cPaperCC.Query(`list`), &cpaper_schema.CommercialPaperList{}).(*cpaper_schema.CommercialPaperList)
			Expect(len(cpapers.Items)).To(Equal(3))
			Expect(cpapers.Items[0].Issuer).To(Equal(cpaper1.Issuer))
			Expect(cpapers.Items[0].PaperNumber).To(Equal(cpaper1.PaperNumber))
		})

		It("Allow to get entry raw protobuf", func() {
			cpaperProtoFromCC := cPaperCC.Query(`get`, &cpaper_schema.CommercialPaperId{Issuer: cpaper1.Issuer, PaperNumber: cpaper1.PaperNumber}).Payload

			stateCpaper := &cpaper_schema.CommercialPaper{
				Issuer:       cpaper1.Issuer,
				PaperNumber:  cpaper1.PaperNumber,
				Owner:        cpaper1.Issuer,
				IssueDate:    cpaper1.IssueDate,
				MaturityDate: cpaper1.MaturityDate,
				FaceValue:    cpaper1.FaceValue,
				State:        cpaper_schema.CommercialPaper_ISSUED, // initial state
			}
			Expect(cpaperProtoFromCC).To(Equal(testcc.MustProtoMarshal(stateCpaper)))
		})

		It("Allow update data in chaincode state", func() {

			expectcc.ResponseOk(cPaperCC.Invoke(`buy`, &cpaper_schema.BuyCommercialPaper{
				Issuer:       cpaper1.Issuer,
				PaperNumber:  cpaper1.PaperNumber,
				CurrentOwner: cpaper1.Issuer,
				NewOwner:     `some-new-owner`,
				Price:        cpaper1.FaceValue - 10,
				PurchaseDate: ptypes.TimestampNow(),
			}))

			cpaperFromCC := expectcc.PayloadIs(
				cPaperCC.Query(`get`, &cpaper_schema.CommercialPaperId{Issuer: cpaper1.Issuer, PaperNumber: cpaper1.PaperNumber}),
				&cpaper_schema.CommercialPaper{}).(*cpaper_schema.CommercialPaper)

			// state is updated
			Expect(cpaperFromCC.State).To(Equal(cpaper_schema.CommercialPaper_TRADING))
			Expect(cpaperFromCC.Owner).To(Equal(`some-new-owner`))
		})

		It("Allow to delete entry", func() {
			toDelete := &cpaper_schema.CommercialPaperId{Issuer: cpaper1.Issuer, PaperNumber: cpaper1.PaperNumber}

			expectcc.ResponseOk(cPaperCC.Invoke(`delete`, toDelete))
			cpapers := expectcc.PayloadIs(cPaperCC.Invoke(`list`), &state_schema.List{}).(*state_schema.List)

			Expect(len(cpapers.Items)).To(Equal(2))
			expectcc.ResponseError(cPaperCC.Invoke(`get`, toDelete), state.ErrKeyNotFound)
		})
	})

	Describe(`Entity with complex id`, func() {

		ent1 := &schema.EntityWithComplexId{Id: &schema.EntityComplexId{IdPart1: `aaa`, IdPart2: `bbb`}}

		It("Allow to add data to chaincode state", func() {

			expectcc.ResponseOk(complexIdCC.Invoke(`entityInsert`, ent1))

			keys := expectcc.PayloadIs(complexIdCC.From(actors[`owner`]).Invoke(`debugStateKeys`, []string{`EntityWithComplexId`}), &[]string{}).([]string)
			Expect(len(keys)).To(Equal(1))

			// from hyperledger/fabric/core/chaincode/shim/chaincode.go
			Expect(keys[0]).To(Equal("\x00" + `EntityWithComplexId` + string(0) + ent1.Id.IdPart1 + string(0) + ent1.Id.IdPart2 + string(0)))
		})

		It("Allow to get entity", func() {
			// use Id as key
			ent1FromCC := expectcc.ResponseOk(complexIdCC.Query(`entityGet`, ent1.Id)).Payload
			Expect(ent1FromCC).To(Equal(testcc.MustProtoMarshal(ent1)))
		})

		It("Allow to list entity", func() {
			// use Id as key
			listFromCC := expectcc.PayloadIs(complexIdCC.Query(`entityList`), &state_schema.List{}).(*state_schema.List)
			Expect(listFromCC.Items).To(HaveLen(1))

			Expect(listFromCC.Items[0].Value).To(Equal(testcc.MustProtoMarshal(ent1)))
		})
	})

})
