package mapping_test

import (
	"testing"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric/protos/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/opencontainers/runc/Godeps/_workspace/src/github.com/golang/protobuf/proto"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/examples/cpaper"
	cpaper_schema "github.com/s7techlab/cckit/examples/cpaper/schema"
	cpaper_testdata "github.com/s7techlab/cckit/examples/cpaper/testdata"
	"github.com/s7techlab/cckit/examples/cpaper_extended"
	cpaper_extended_schema "github.com/s7techlab/cckit/examples/cpaper_extended/schema"
	cpaper_extended_testdata "github.com/s7techlab/cckit/examples/cpaper_extended/testdata"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/state"
	"github.com/s7techlab/cckit/state/mapping"
	"github.com/s7techlab/cckit/state/mapping/testdata"
	"github.com/s7techlab/cckit/state/mapping/testdata/schema"
	state_schema "github.com/s7techlab/cckit/state/schema"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var (
	actors                                             identity.Actors
	cPaperCC, cPaperExtendedCC, complexIdCC, sliceIdCC *testcc.MockStub
	err                                                error
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

		cPaperExtendedCC = testcc.NewMockStub(`cpapers_extended`, cpaper_extended.NewCC())
		cPaperExtendedCC.From(actors[`owner`]).Init()

		complexIdCC = testcc.NewMockStub(`complexid`, testdata.NewComplexIdCC())
		complexIdCC.From(actors[`owner`]).Init()

		sliceIdCC = testcc.NewMockStub(`sliceid`, testdata.NewSliceIdCC())
		sliceIdCC.From(actors[`owner`]).Init()
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

	Describe(`Commercial paper extended, protobuf based schema with additional keys`, func() {

		var cpaper1 = &cpaper_extended_testdata.CPapers[0]

		It("Allow to add data to chaincode state", func(done Done) {
			events := cPaperExtendedCC.EventSubscription()
			expectcc.ResponseOk(cPaperExtendedCC.Invoke(`issue`, cpaper1))

			Expect(<-events).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: `IssueCommercialPaper`,
				Payload:   testcc.MustProtoMarshal(cpaper1),
			}))

			close(done)
		}, 0.2)

		It("Disallow to add data to chaincode state with same primary AND  uniq key fields", func() {
			expectcc.ResponseError(cPaperExtendedCC.Invoke(`issue`, cpaper1), mapping.ErrMappingUniqKeyExists)
		})

		It("Disallow to add data to chaincode state with same uniq key fields", func() {
			// change PK
			cpChanged1 := proto.Clone(cpaper1).(*cpaper_extended_schema.IssueCommercialPaper)
			cpChanged1.PaperNumber = `some-new-number`

			// errored on checking uniq key
			expectcc.ResponseError(cPaperExtendedCC.Invoke(`issue`, cpChanged1), mapping.ErrMappingUniqKeyExists)
		})

		It("Disallow to add data to chaincode state with same primary key fields", func() {
			// change Uniq Key
			cpChanged2 := proto.Clone(cpaper1).(*cpaper_extended_schema.IssueCommercialPaper)
			cpChanged2.ExternalId = `some-new-external-id`

			// errored obn checkong primary key
			expectcc.ResponseError(cPaperExtendedCC.Invoke(`issue`, cpChanged2), state.ErrKeyAlreadyExists)
		})

		It("Allow to find data by uniq key", func() {

			cpaperFromCCByExtId := expectcc.PayloadIs(
				cPaperExtendedCC.Query(`getByExternalId`, cpaper1.ExternalId),
				&cpaper_extended_schema.CommercialPaper{}).(*cpaper_extended_schema.CommercialPaper)

			cpaperFromCC := expectcc.PayloadIs(
				cPaperExtendedCC.Query(`get`,
					&cpaper_extended_schema.CommercialPaperId{Issuer: cpaper1.Issuer, PaperNumber: cpaper1.PaperNumber}),
				&cpaper_extended_schema.CommercialPaper{}).(*cpaper_extended_schema.CommercialPaper)

			Expect(cpaperFromCCByExtId).To(BeEquivalentTo(cpaperFromCC))
		})

		It("Disallow to find data by non existent uniq key", func() {
			expectcc.ResponseError(
				cPaperExtendedCC.Query(`getByExternalId`, `some-non-existent-id`), state.ErrKeyNotFound)
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

	Describe(`Entity with slice id`, func() {

		ent2 := &schema.EntityWithSliceId{Id: []string{`aa`, `bb`}, SomeDate: ptypes.TimestampNow()}

		It("Allow to add data to chaincode state", func() {
			expectcc.ResponseOk(sliceIdCC.Invoke(`entityInsert`, ent2))
			keys := expectcc.PayloadIs(sliceIdCC.From(actors[`owner`]).Invoke(`debugStateKeys`, []string{`EntityWithSliceId`}), &[]string{}).([]string)

			Expect(len(keys)).To(Equal(1))

			// from hyperledger/fabric/core/chaincode/shim/chaincode.go
			Expect(keys[0]).To(Equal("\x00" + `EntityWithSliceId` + string(0) + ent2.Id[0] + string(0) + ent2.Id[1] + string(0)))
		})

		It("Allow to get entity", func() {
			// use Id as key
			ent1FromCC := expectcc.ResponseOk(sliceIdCC.Query(`entityGet`, state.StringsIdToStr(ent2.Id))).Payload
			Expect(ent1FromCC).To(Equal(testcc.MustProtoMarshal(ent2)))
		})

		It("Allow to list entity", func() {
			// use Id as key
			listFromCC := expectcc.PayloadIs(sliceIdCC.Query(`entityList`), &state_schema.List{}).(*state_schema.List)
			Expect(listFromCC.Items).To(HaveLen(1))

			Expect(listFromCC.Items[0].Value).To(Equal(testcc.MustProtoMarshal(ent2)))
		})
	})
})
