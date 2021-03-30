package state_test

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	identitytestdata "github.com/s7techlab/cckit/identity/testdata"
	"github.com/s7techlab/cckit/state"
	"github.com/s7techlab/cckit/state/testdata"
	"github.com/s7techlab/cckit/state/testdata/schema"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

var (
	booksCC *testcc.MockStub

	Owner = identitytestdata.Certificates[0].MustIdentity(`SOME_MSP`)
)
var _ = Describe(`State`, func() {

	BeforeSuite(func() {

		//Create books chaincode mock - struct based schema
		booksCC = testcc.NewMockStub(`books`, testdata.NewBooksCC())
		booksCC.From(Owner).Init()
	})

	Describe(`Struct based schema`, func() {

		It("Allow to insert entries", func() {
			expectcc.ResponseOk(booksCC.From(Owner).Invoke(`bookInsert`, &testdata.Books[0]))
			expectcc.ResponseOk(booksCC.From(Owner).Invoke(`bookInsert`, &testdata.Books[1]))
			expectcc.ResponseOk(booksCC.From(Owner).Invoke(`bookInsert`, &testdata.Books[2]))
		})

		It("Disallow to insert entries with same keys", func() {
			expectcc.ResponseError(booksCC.From(Owner).Invoke(`bookInsert`, &testdata.Books[2]))
		})

		It("Allow to get entry list", func() {
			books := expectcc.PayloadIs(booksCC.Invoke(`bookList`), &[]schema.Book{}).([]schema.Book)
			Expect(len(books)).To(Equal(3))
			Expect(books[0]).To(Equal(testdata.Books[0]))
			Expect(books[1]).To(Equal(testdata.Books[1]))
			Expect(books[2]).To(Equal(testdata.Books[2]))
		})

		It("Allow to get entry converted to target type", func() {
			book1FromCC := expectcc.PayloadIs(booksCC.Invoke(`bookGet`, testdata.Books[0].Id), &schema.Book{}).(schema.Book)
			Expect(book1FromCC).To(Equal(testdata.Books[0]))
		})

		It("Allow to get entry json", func() {
			book2JsonFromCC := booksCC.Invoke(`bookGet`, testdata.Books[2].Id).Payload
			book2Json, _ := json.Marshal(testdata.Books[2])
			Expect(book2JsonFromCC).To(Equal(book2Json))
		})

		It("Allow to upsert entry", func() {
			bookToUpdate := testdata.Books[2]
			bookToUpdate.Title = `thirdiest title`

			bookUpdated := expectcc.PayloadIs(booksCC.Invoke(`bookUpsert`, &bookToUpdate), &schema.Book{}).(schema.Book)
			Expect(bookUpdated.Title).To(Equal(bookToUpdate.Title))

			bookFromCC := expectcc.PayloadIs(booksCC.Invoke(`bookGet`, bookToUpdate.Id), &schema.Book{}).(schema.Book)
			Expect(bookFromCC).To(Equal(bookUpdated))
		})

		It("Allow to upsert entry with tx state caching", func() {
			bookToUpdate := testdata.Books[1]
			bookToUpdate.Title = `once more strange uniq title`

			bookUpdated := expectcc.PayloadIs(booksCC.Invoke(`bookUpsertWithCache`, &bookToUpdate), &schema.Book{}).(schema.Book)
			Expect(bookUpdated.Title).To(Equal(bookToUpdate.Title))

			bookFromCC := expectcc.PayloadIs(booksCC.Invoke(`bookGet`, bookToUpdate.Id), &schema.Book{}).(schema.Book)
			Expect(bookFromCC).To(Equal(bookToUpdate))
		})

		It("Allow to delete entry", func() {
			expectcc.ResponseOk(booksCC.From(Owner).Invoke(`bookDelete`, testdata.Books[0].Id))
			books := expectcc.PayloadIs(booksCC.Invoke(`bookList`), &[]schema.Book{}).([]schema.Book)
			Expect(len(books)).To(Equal(2))

			expectcc.ResponseError(booksCC.Invoke(`bookGet`, testdata.Books[0].Id), state.ErrKeyNotFound)
		})
	})

})
