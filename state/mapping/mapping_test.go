package mapping_test

import (
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-protos-go/peer"

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
	compositeIDCC, complexIDCC, sliceIDCC, indexesCC, configCC *testcc.MockStub

	Owner = identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)
)
var _ = Describe(`State mapping in chaincode`, func() {

	BeforeSuite(func() {

		compositeIDCC = testcc.NewMockStub(`proto`, testdata.NewCompositeIdCC())
		compositeIDCC.From(Owner).Init()

		complexIDCC = testcc.NewMockStub(`complex_id`, testdata.NewComplexIdCC())
		complexIDCC.From(Owner).Init()

		sliceIDCC = testcc.NewMockStub(`slice_id`, testdata.NewSliceIdCC())
		sliceIDCC.From(Owner).Init()

		indexesCC = testcc.NewMockStub(`indexes`, testdata.NewIndexesCC())
		indexesCC.From(Owner).Init()

		configCC = testcc.NewMockStub(`config`, testdata.NewCCWithConfig())
		configCC.From(Owner).Init()
	})

	Describe(`Entity with composite id`, func() {
		create1 := testdata.CreateEntityWithCompositeId[0]
		create2 := testdata.CreateEntityWithCompositeId[1]
		create3 := testdata.CreateEntityWithCompositeId[2]

		It("Allow to get mapping data by namespace", func() {
			mapping, err := testdata.EntityWithCompositeIdStateMapping.GetByNamespace(testdata.EntityCompositeIdNamespace)
			Expect(err).NotTo(HaveOccurred())
			Expect(mapping.Schema()).To(BeEquivalentTo(&schema.EntityWithCompositeId{}))

			key, err := mapping.PrimaryKey(&schema.EntityWithCompositeId{
				IdFirstPart:  create1.IdFirstPart,
				IdSecondPart: create1.IdSecondPart,
				IdThirdPart:  create1.IdThirdPart,
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(key).To(Equal(
				testdata.EntityCompositeIdNamespace.Append(
					state.Key{create1.IdFirstPart, create1.IdSecondPart, testdata.Dates[0]})))
		})

		It("Allow to add data to chaincode state", func(done Done) {
			events := compositeIDCC.EventSubscription()
			expectcc.ResponseOk(compositeIDCC.Invoke(`create`, create1))

			Expect(<-events).To(BeEquivalentTo(&peer.ChaincodeEvent{
				EventName: `CreateEntityWithCompositeId`,
				Payload:   testcc.MustProtoMarshal(create1),
			}))

			expectcc.ResponseOk(compositeIDCC.Invoke(`create`, create2))
			expectcc.ResponseOk(compositeIDCC.Invoke(`create`, create3))

			close(done)
		})

		It("Disallow to insert entries with same primary key", func() {
			expectcc.ResponseError(compositeIDCC.Invoke(`create`, create1), state.ErrKeyAlreadyExists)
		})

		It("Allow to get entry list", func() {
			entities := expectcc.PayloadIs(compositeIDCC.Query(`list`),
				&schema.EntityWithCompositeIdList{}).(*schema.EntityWithCompositeIdList)
			Expect(len(entities.Items)).To(Equal(3))
			Expect(entities.Items[0].Name).To(Equal(create1.Name))
			Expect(entities.Items[0].Value).To(BeNumerically("==", create1.Value))
		})

		It("Allow to get entry raw protobuf", func() {
			dataFromCC := compositeIDCC.Query(`get`,
				&schema.EntityCompositeId{
					IdFirstPart:  create1.IdFirstPart,
					IdSecondPart: create1.IdSecondPart,
					IdThirdPart:  create1.IdThirdPart,
				},
			).Payload

			e := &schema.EntityWithCompositeId{
				IdFirstPart:  create1.IdFirstPart,
				IdSecondPart: create1.IdSecondPart,
				IdThirdPart:  create1.IdThirdPart,

				Name:  create1.Name,
				Value: create1.Value,
			}
			Expect(dataFromCC).To(Equal(testcc.MustProtoMarshal(e)))
		})

		It("Allow update data in chaincode state", func() {
			expectcc.ResponseOk(compositeIDCC.Invoke(`update`, &schema.UpdateEntityWithCompositeId{
				IdFirstPart:  create1.IdFirstPart,
				IdSecondPart: create1.IdSecondPart,
				IdThirdPart:  create1.IdThirdPart,
				Name:         `New name`,
				Value:        1000,
			}))

			entityFromCC := expectcc.PayloadIs(
				compositeIDCC.Query(`get`, &schema.EntityCompositeId{
					IdFirstPart:  create1.IdFirstPart,
					IdSecondPart: create1.IdSecondPart,
					IdThirdPart:  create1.IdThirdPart,
				}),
				&schema.EntityWithCompositeId{}).(*schema.EntityWithCompositeId)

			// state is updated
			Expect(entityFromCC.Name).To(Equal(`New name`))
			Expect(entityFromCC.Value).To(BeNumerically("==", 1000))
		})

		It("Allow to delete entry", func() {
			toDelete := &schema.EntityCompositeId{
				IdFirstPart:  create1.IdFirstPart,
				IdSecondPart: create1.IdSecondPart,
				IdThirdPart:  create1.IdThirdPart,
			}

			expectcc.ResponseOk(compositeIDCC.Invoke(`delete`, toDelete))
			ee := expectcc.PayloadIs(
				compositeIDCC.Invoke(`list`),
				&schema.EntityWithCompositeIdList{}).(*schema.EntityWithCompositeIdList)

			Expect(len(ee.Items)).To(Equal(2))
			expectcc.ResponseError(compositeIDCC.Invoke(`get`, toDelete), state.ErrKeyNotFound)
		})

		It("Allow to insert entry once more time", func() {
			expectcc.ResponseOk(compositeIDCC.Invoke(`create`, create1))
		})

		It("Check entry keying", func() {

		})

	})

	Describe(`Entity with complex id`, func() {

		ent1 := &schema.EntityWithComplexId{Id: &schema.EntityComplexId{
			IdPart1: []string{`aaa`, `bb`},
			IdPart2: `ccc`,
			IdPart3: testcc.MustTime(`2020-01-28T17:00:00Z`),
		}}

		It("Allow to add data to chaincode state", func() {
			expectcc.ResponseOk(complexIDCC.Invoke(`entityInsert`, ent1))
			keys := expectcc.PayloadIs(complexIDCC.From(Owner).Invoke(
				`debugStateKeys`, `EntityWithComplexId`), &[]string{}).([]string)
			Expect(len(keys)).To(Equal(1))

			timeStr := time.Unix(ent1.Id.IdPart3.GetSeconds(), int64(ent1.Id.IdPart3.GetNanos())).Format(`2006-01-02`)
			// from hyperledger/fabric/core/chaincode/shim/chaincode.go
			Expect(keys[0]).To(Equal(
				string(rune(0)) +
					`EntityWithComplexId` + string(rune(0)) +
					ent1.Id.IdPart1[0] + string(rune(0)) +
					ent1.Id.IdPart1[1] + string(rune(0)) +
					ent1.Id.IdPart2 + string(rune(0)) +
					timeStr + string(rune(0))))
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
				`debugStateKeys`, `EntityWithSliceId`), &[]string{}).([]string)

			Expect(len(keys)).To(Equal(1))

			// from hyperledger/fabric/core/chaincode/shim/chaincode.go
			Expect(keys[0]).To(Equal("\x00" + `EntityWithSliceId` + string(rune(0)) + ent2.Id[0] + string(rune(0)) + ent2.Id[1] + string(rune(0))))
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

	Describe(`Entity with indexes`, func() {

		create1 := testdata.CreateEntityWithIndexes[0]
		create2 := testdata.CreateEntityWithIndexes[1]

		It("Allow to add data with single external id", func() {
			expectcc.ResponseOk(indexesCC.Invoke(`create`, create1))
		})

		It("Disallow to add data to chaincode state with same uniq key fields", func() {
			createWithNewId := proto.Clone(create1).(*schema.CreateEntityWithIndexes)
			createWithNewId.Id = `abcdef` // id is really new

			// errored on checking uniq key
			expectcc.ResponseError(
				indexesCC.Invoke(`create`, create1),
				mapping.ErrMappingUniqKeyExists)
		})

		It("Allow finding data by uniq key", func() {
			fromCCByExtId := expectcc.PayloadIs(
				indexesCC.Query(`getByExternalId`, create1.ExternalId),
				&schema.EntityWithIndexes{}).(*schema.EntityWithIndexes)

			fromCCById := expectcc.PayloadIs(
				indexesCC.Query(`get`, create1.Id),
				&schema.EntityWithIndexes{}).(*schema.EntityWithIndexes)

			Expect(fromCCByExtId).To(BeEquivalentTo(fromCCById))
		})

		It("Allow to get idx state key by uniq key", func() {
			idxKey, err := testdata.EntityWithIndexesStateMapping.IdxKey(
				&schema.EntityWithIndexes{}, `ExternalId`, []string{create1.ExternalId})
			Expect(err).NotTo(HaveOccurred())

			Expect(idxKey).To(BeEquivalentTo([]string{
				mapping.KeyRefNamespace,
				strings.Join(mapping.SchemaNamespace(&schema.EntityWithIndexes{}), `-`),
				`ExternalId`,
				create1.ExternalId,
			}))
		})

		It("Disallow finding data by non existent uniq key", func() {
			expectcc.ResponseError(
				indexesCC.Query(`getByExternalId`, `some-non-existent-id`),
				mapping.ErrIndexReferenceNotFound)
		})

		It("Allow to add data with multiple external id", func() {
			expectcc.ResponseOk(indexesCC.Invoke(`create`, create2))
		})

		It("Allow to find data by multi key", func() {
			fromCCByExtId1 := expectcc.PayloadIs(
				indexesCC.Query(`getByOptMultiExternalId`, create2.OptionalExternalIds[0]),
				&schema.EntityWithIndexes{}).(*schema.EntityWithIndexes)

			fromCCByExtId2 := expectcc.PayloadIs(
				indexesCC.Query(`getByOptMultiExternalId`, create2.OptionalExternalIds[1]),
				&schema.EntityWithIndexes{}).(*schema.EntityWithIndexes)

			fromCCById := expectcc.PayloadIs(
				indexesCC.Query(`get`, create2.Id),
				&schema.EntityWithIndexes{}).(*schema.EntityWithIndexes)

			Expect(fromCCByExtId1).To(BeEquivalentTo(fromCCById))
			Expect(fromCCByExtId2).To(BeEquivalentTo(fromCCById))
		})

		It("Allow update indexes value", func() {
			update2 := &schema.UpdateEntityWithIndexes{
				Id:                  create2.Id,
				ExternalId:          `some_new_external_id`,
				OptionalExternalIds: []string{create2.OptionalExternalIds[0], `AND SOME NEW`},
			}
			expectcc.ResponseOk(indexesCC.Invoke(`update`, update2))
		})

		It("Allow to find data by updated multi key", func() {
			fromCCByExtId1 := expectcc.PayloadIs(
				indexesCC.Query(`getByOptMultiExternalId`, create2.OptionalExternalIds[0]),
				&schema.EntityWithIndexes{}).(*schema.EntityWithIndexes)

			fromCCByExtId2 := expectcc.PayloadIs(
				indexesCC.Query(`getByOptMultiExternalId`, `AND SOME NEW`),
				&schema.EntityWithIndexes{}).(*schema.EntityWithIndexes)

			Expect(fromCCByExtId1.Id).To(Equal(create2.Id))
			Expect(fromCCByExtId2.Id).To(Equal(create2.Id))

			Expect(fromCCByExtId2.OptionalExternalIds).To(
				BeEquivalentTo([]string{create2.OptionalExternalIds[0], `AND SOME NEW`}))
		})

		It("Disallow to find data by previous multi key", func() {
			expectcc.ResponseError(
				indexesCC.Query(`getByOptMultiExternalId`, create2.OptionalExternalIds[1]),
				mapping.ErrIndexReferenceNotFound)
		})

		It("Allow to find data by updated uniq key", func() {
			fromCCByExtId := expectcc.PayloadIs(
				indexesCC.Query(`getByExternalId`, `some_new_external_id`),
				&schema.EntityWithIndexes{}).(*schema.EntityWithIndexes)

			Expect(fromCCByExtId.Id).To(Equal(create2.Id))
			Expect(fromCCByExtId.ExternalId).To(Equal(`some_new_external_id`))
		})

		It("Disallow to find data by previous uniq key", func() {
			expectcc.ResponseError(
				indexesCC.Query(`getByExternalId`, create2.ExternalId),
				mapping.ErrIndexReferenceNotFound)
		})

		It("Allow to delete entry", func() {
			expectcc.ResponseOk(indexesCC.Invoke(`delete`, create2.Id))

			ee := expectcc.PayloadIs(
				indexesCC.Invoke(`list`),
				&schema.EntityWithIndexesList{}).(*schema.EntityWithIndexesList)

			Expect(len(ee.Items)).To(Equal(1))
			expectcc.ResponseError(indexesCC.Invoke(`get`, create2.Id), state.ErrKeyNotFound)
		})

		It("Allow to insert entry once more time", func() {
			expectcc.ResponseOk(indexesCC.Invoke(`create`, create2))
		})

	})

	Describe(`Entity with static key`, func() {
		configSample := &schema.Config{
			Field1: `aaa`,
			Field2: `bbb`,
		}

		It("Disallow to get config before set", func() {
			expectcc.ResponseError(configCC.Invoke(`configGet`), `state entry not found: Config`)
		})

		It("Allow to set config", func() {
			expectcc.ResponseOk(configCC.Invoke(`configSet`, configSample))
		})

		It("Allow to get config", func() {
			confFromCC := expectcc.PayloadIs(configCC.Invoke(`configGet`), &schema.Config{}).(*schema.Config)
			Expect(confFromCC.Field1).To(Equal(configSample.Field1))
			Expect(confFromCC.Field2).To(Equal(configSample.Field2))
		})

	})
})
