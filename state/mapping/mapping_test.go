package mapping_test

import (
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/s7techlab/cckit/examples/cpaper/schema"
	"github.com/s7techlab/cckit/examples/cpaper/testdata"

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
		It("Allow to insert entries", func() {
			expectcc.ResponseOk(cPaperCC.From(actors[`owner`]).Invoke(`cpaperInsert`, &testdata.CPapers[0]))
		})

		It("Disallow to insert entries with same keys", func() {
			expectcc.ResponseError(cPaperCC.From(actors[`owner`]).Invoke(`cpaperInsert`, &testdata.CPapers[0]))
		})

		It("Allow to get entry list", func() {
			cpapers := expectcc.PayloadIs(cPaperCC.Invoke(`cpaperList`), &[]schema.CommercialPaper{}).([]schema.CommercialPaper)
			Expect(len(cpapers)).To(Equal(1))
			Expect(cpapers[0].Paper).To(Equal(testdata.CPapers[0].Paper))
		})

		It("Allow to get entry raw protobuf", func() {
			cpaperProtoFromCC := cPaperCC.Invoke(`cpaperGet`, testdata.CPapers[0].Paper).Payload
			book2proto, _ := proto.Marshal(&testdata.CPapers[0])
			Expect(cpaperProtoFromCC).To(Equal(book2proto))
		})
	})

})
