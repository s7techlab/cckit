package state_test

import (
	"encoding/json"
	"testing"

	"github.com/s7techlab/cckit/state"
	"github.com/s7techlab/cckit/state/testdata/schema"

	"github.com/s7techlab/cckit/state/testdata"

	expectcc "github.com/s7techlab/cckit/testing/expect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPrivateState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Private state suite")
}

// var (
// 	actors  testcc.Identities
// 	booksCC *testcc.MockStub
// 	err     error
// )
var _ = Describe(`PrivateState`, func() {

	// BeforeSuite(func() {

	// 	actors, err = testcc.IdentitiesFromFiles(`SOME_MSP`, map[string]string{
	// 		`owner`: `s7techlab.pem`,
	// 	}, examplecert.Content)

	// 	Expect(err).To(BeNil())

	// 	//Create books chaincode mock - struct based schema
	// 	booksCC = testcc.NewMockStub(`books`, testdata.NewBooksCC())
	// 	booksCC.From(actors[`owner`]).Init()
	// })

	Describe(`Struct based schema private`, func() {

		It("Allow to insert entries", func() {
			expectcc.ResponseOk(booksCC.From(actors[`owner`]).Invoke(`privateBookInsert`, &testdata.PrivateBooks[0]))
			expectcc.ResponseOk(booksCC.From(actors[`owner`]).Invoke(`privateBookInsert`, &testdata.PrivateBooks[1]))
			expectcc.ResponseOk(booksCC.From(actors[`owner`]).Invoke(`privateBookInsert`, &testdata.PrivateBooks[2]))
		})

		It("Disallow to insert entries with same keys", func() {
			expectcc.ResponseError(booksCC.From(actors[`owner`]).Invoke(`privateBookInsert`, &testdata.PrivateBooks[2]))
		})

		It("Allow to get entry list", func() {
			books := expectcc.PayloadIs(booksCC.Invoke(`privateBookList`), &[]schema.PrivateBook{}).([]schema.PrivateBook)
			Expect(len(books)).To(Equal(3))
			Expect(books[0]).To(Equal(testdata.PrivateBooks[0]))
			Expect(books[1]).To(Equal(testdata.PrivateBooks[1]))
			Expect(books[2]).To(Equal(testdata.PrivateBooks[2]))
		})

		It("Allow to get entry converted to target type", func() {
			book1FromCC := expectcc.PayloadIs(booksCC.Invoke(`privateBookGet`, testdata.PrivateBooks[0].Id), &schema.PrivateBook{}).(schema.PrivateBook)
			Expect(book1FromCC).To(Equal(testdata.PrivateBooks[0]))
		})

		It("Allow to get entry json", func() {
			book2JsonFromCC := booksCC.Invoke(`privateBookGet`, testdata.PrivateBooks[2].Id).Payload
			book2Json, _ := json.Marshal(testdata.PrivateBooks[2])
			Expect(book2JsonFromCC).To(Equal(book2Json))
		})

		It("Allow to upsert entry", func() {
			book2Updated := testdata.PrivateBooks[2]
			book2Updated.Title = `thirdiest title`

			updateRes := expectcc.PayloadIs(booksCC.Invoke(`privateBookUpsert`, &book2Updated), &schema.PrivateBook{}).(schema.PrivateBook)
			Expect(updateRes.Title).To(Equal(book2Updated.Title))

			book3FromCC := expectcc.PayloadIs(booksCC.Invoke(`privateBookGet`, testdata.PrivateBooks[2].Id), &schema.PrivateBook{}).(schema.PrivateBook)
			Expect(book3FromCC).To(Equal(book2Updated))
		})

		It("Allow to delete entry", func() {
			expectcc.ResponseOk(booksCC.From(actors[`owner`]).Invoke(`privateBookDelete`, testdata.PrivateBooks[0].Id))
			books := expectcc.PayloadIs(booksCC.Invoke(`privateBookList`), &[]schema.PrivateBook{}).([]schema.PrivateBook)
			Expect(len(books)).To(Equal(2))

			expectcc.ResponseError(booksCC.Invoke(`privateBookGet`, testdata.PrivateBooks[0].Id), state.ErrKeyNotFound)
		})
	})

})
