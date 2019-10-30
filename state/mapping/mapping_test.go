package mapping_test

import (
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric/protos/peer"
	identitytestdata "github.com/s7techlab/cckit/identity/testdata"
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
	protoCC, complexIDCC, sliceIDCC *testcc.MockStub
	err                             error

	Owner = identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)
)
var _ = Describe(`Mapping`, func() {

	BeforeSuite(func() {

		protoCC = testcc.NewMockStub(`cpapers`, testdata.NewProtoCC())
		protoCC.From(Owner).Init()

		complexIDCC = testcc.NewMockStub(`complexid`, testdata.NewComplexIdCC())
		complexIDCC.From(Owner).Init()

		sliceIDCC = testcc.NewMockStub(`sliceid`, testdata.NewSliceIdCC())
		sliceIDCC.From(Owner).Init()
	})

	Describe(`Commercial paper extended, protobuf based schema with additional keys`, func() {
		issueMock1 := testdata.ProtoIssueMocks[0]
		issueMock2 := testdata.ProtoIssueMocks[1]
		issueMock3 := testdata.ProtoIssueMocks[2]
		issueMockExistingExternal := testdata.ProtoIssueMockExistingExternal
		issueMockExistingPrimary := testdata.ProtoIssueMockExistingPrimary

		It("Allow to get mapping data by namespace", func() {
			mapping, err := testdata.ProtoStateMapping.GetByNamespace(state.Key{`ProtoEntity`})
			Expect(err).NotTo(HaveOccurred())
			Expect(mapping.Schema()).To(BeEquivalentTo(&schema.ProtoEntity{}))
		})

		It("Allow to add data to chaincode state", func(done Done) {
			events := protoCC.EventSubscription()
			expectcc.ResponseOk(protoCC.Invoke(`issue`, &issueMock1))

			Expect(<-events).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: `IssueProtoEntity`,
				Payload:   testcc.MustProtoMarshal(&issueMock1),
			}))

			expectcc.ResponseOk(protoCC.Invoke(`issue`, &issueMock2))
			expectcc.ResponseOk(protoCC.Invoke(`issue`, &issueMock3))

			close(done)
		}, 0.2)

		It("Disallow to insert entries with same uniq AND primary keys", func() {
			expectcc.ResponseError(protoCC.Invoke(`issue`, &issueMock1))
		})

		It("Disallow to add data to chaincode state with same uniq key fields", func() {
			// errored on checking uniq key
			expectcc.ResponseError(
				protoCC.Invoke(`issue`, &issueMockExistingExternal),
				mapping.ErrMappingUniqKeyExists)
		})

		It("Disallow adding data to chaincode state with same primary key fields", func() {
			// errored obn checkong primary key
			expectcc.ResponseError(
				protoCC.Invoke(`issue`, &issueMockExistingPrimary),
				state.ErrKeyAlreadyExists)
		})

		It("Allow to get entry list", func() {
			entities := expectcc.PayloadIs(protoCC.Query(`list`),
				&schema.ProtoEntityList{}).(*schema.ProtoEntityList)
			Expect(len(entities.Items)).To(Equal(3))
			Expect(entities.Items[0].Name).To(Equal(issueMock1.Name))
			Expect(entities.Items[0].Value).To(BeNumerically("==", 0))
			Expect(entities.Items[0].ExternalId).To(Equal(issueMock1.ExternalId))
		})

		It("Allow finding data by uniq key", func() {

			cpaperFromCCByExtID := expectcc.PayloadIs(
				protoCC.Query(`getByExternalId`, issueMock1.ExternalId),
				&schema.ProtoEntity{}).(*schema.ProtoEntity)

			cpaperFromCC := expectcc.PayloadIs(
				protoCC.Query(`get`, &schema.ProtoEntityId{
					IdFirstPart:  issueMock1.IdFirstPart,
					IdSecondPart: issueMock1.IdSecondPart},
				),
				&schema.ProtoEntity{}).(*schema.ProtoEntity)

			Expect(cpaperFromCCByExtID).To(BeEquivalentTo(cpaperFromCC))
		})

		It("Allow to get idx state key by uniq key", func() {
			idxKey, err := testdata.ProtoStateMapping.IdxKey(&schema.ProtoEntity{}, `ExternalId`, []string{issueMock1.ExternalId})
			Expect(err).NotTo(HaveOccurred())

			Expect(idxKey).To(BeEquivalentTo([]string{
				mapping.KeyRefNamespace,
				strings.Join(mapping.SchemaNamespace(&schema.ProtoEntity{}), `-`),
				`ExternalId`,
				issueMock1.ExternalId,
			}))
		})

		It("Disallow finding data by non existent uniq key", func() {
			expectcc.ResponseError(
				protoCC.Query(`getByExternalId`, `some-non-existent-id`), `uniq index`)
		})

		It("Allow to get entry raw protobuf", func() {
			cpaperProtoFromCC := protoCC.Query(`get`,
				&schema.ProtoEntityId{
					IdFirstPart:  issueMock1.IdFirstPart,
					IdSecondPart: issueMock1.IdSecondPart},
			).Payload

			stateProtoEntity := &schema.ProtoEntity{
				IdFirstPart:  issueMock1.IdFirstPart,
				IdSecondPart: issueMock1.IdSecondPart,
				Name:         issueMock1.Name,
				Value:        0,
				ExternalId:   issueMock1.ExternalId,
			}
			Expect(cpaperProtoFromCC).To(Equal(testcc.MustProtoMarshal(stateProtoEntity)))
		})

		It("Allow update data in chaincode state", func() {
			expectcc.ResponseOk(protoCC.Invoke(`increment`, &schema.IncrementProtoEntity{
				IdFirstPart:  issueMock1.IdFirstPart,
				IdSecondPart: issueMock1.IdSecondPart,
				Name:         issueMock1.Name,
			}))

			entityFromCC := expectcc.PayloadIs(
				protoCC.Query(`get`, &schema.ProtoEntityId{
					IdFirstPart:  issueMock1.IdFirstPart,
					IdSecondPart: issueMock1.IdSecondPart,
				}),
				&schema.ProtoEntity{}).(*schema.ProtoEntity)

			// state is updated
			Expect(entityFromCC.Value).To(BeNumerically("==", 1))
		})

		It("Allow to delete entry", func() {
			toDelete := &schema.ProtoEntityId{
				IdFirstPart:  issueMock1.IdFirstPart,
				IdSecondPart: issueMock1.IdSecondPart,
			}

			expectcc.ResponseOk(protoCC.Invoke(`delete`, toDelete))
			cpapers := expectcc.PayloadIs(
				protoCC.Invoke(`list`),
				&schema.ProtoEntityList{},
			).(*schema.ProtoEntityList)

			Expect(len(cpapers.Items)).To(Equal(2))
			expectcc.ResponseError(protoCC.Invoke(`get`, toDelete), state.ErrKeyNotFound)
		})

		It("Allow to insert entry once more time", func() {
			expectcc.ResponseOk(protoCC.Invoke(`issue`, &issueMock1))

			cpaperFromCCByExtID := expectcc.PayloadIs(
				protoCC.Query(`getByExternalId`, issueMock1.ExternalId),
				&schema.ProtoEntity{}).(*schema.ProtoEntity)

			Expect(cpaperFromCCByExtID.IdFirstPart).To(Equal(issueMock1.IdFirstPart))
		})

	})

	Describe(`Entity with complex id`, func() {

		ent1 := &schema.EntityWithComplexId{Id: &schema.EntityComplexId{IdPart1: `aaa`, IdPart2: `bbb`}}

		It("Allow to add data to chaincode state", func() {
			expectcc.ResponseOk(complexIDCC.Invoke(`entityInsert`, ent1))
			keys := expectcc.PayloadIs(complexIDCC.From(Owner).Invoke(
				`debugStateKeys`, []string{`EntityWithComplexId`}), &[]string{}).([]string)
			Expect(len(keys)).To(Equal(1))

			// from hyperledger/fabric/core/chaincode/shim/chaincode.go
			Expect(keys[0]).To(Equal(
				"\x00" + `EntityWithComplexId` + string(0) + ent1.Id.IdPart1 + string(0) + ent1.Id.IdPart2 + string(0)))
		})

		It("Allow to get entity", func() {
			// use Id as key
			ent1FromCC := expectcc.ResponseOk(complexIDCC.Query(`entityGet`, ent1.Id)).Payload
			Expect(ent1FromCC).To(Equal(testcc.MustProtoMarshal(ent1)))
		})

		It("Allow to list entity", func() {
			// use Id as key
			listFromCC := expectcc.PayloadIs(complexIDCC.Query(`entityList`), &state_schema.List{}).(*state_schema.List)
			Expect(listFromCC.Items).To(HaveLen(1))

			Expect(listFromCC.Items[0].Value).To(Equal(testcc.MustProtoMarshal(ent1)))
		})
	})

	Describe(`Entity with slice id`, func() {

		ent2 := &schema.EntityWithSliceId{Id: []string{`aa`, `bb`}, SomeDate: ptypes.TimestampNow()}

		It("Allow to add data to chaincode state", func() {
			expectcc.ResponseOk(sliceIDCC.Invoke(`entityInsert`, ent2))
			keys := expectcc.PayloadIs(sliceIDCC.From(Owner).Invoke(
				`debugStateKeys`, []string{`EntityWithSliceId`}), &[]string{}).([]string)

			Expect(len(keys)).To(Equal(1))

			// from hyperledger/fabric/core/chaincode/shim/chaincode.go
			Expect(keys[0]).To(Equal("\x00" + `EntityWithSliceId` + string(0) + ent2.Id[0] + string(0) + ent2.Id[1] + string(0)))
		})

		It("Allow to get entity", func() {
			// use Id as key
			ent1FromCC := expectcc.ResponseOk(sliceIDCC.Query(`entityGet`, state.StringsIdToStr(ent2.Id))).Payload
			Expect(ent1FromCC).To(Equal(testcc.MustProtoMarshal(ent2)))
		})

		It("Allow to list entity", func() {
			// use Id as key
			listFromCC := expectcc.PayloadIs(sliceIDCC.Query(`entityList`), &state_schema.List{}).(*state_schema.List)
			Expect(listFromCC.Items).To(HaveLen(1))

			Expect(listFromCC.Items[0].Value).To(Equal(testcc.MustProtoMarshal(ent2)))
		})
	})
})
