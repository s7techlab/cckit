package state_test

import (
	"testing"

	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	examplecert "github.com/s7techlab/cckit/examples/cert"
	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/identity"
	"github.com/s7techlab/cckit/router"
	"github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/state"
	testcc "github.com/s7techlab/cckit/testing"
	expectcc "github.com/s7techlab/cckit/testing/expect"
)

func TestState(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "State suite")
}

const BookEntity = `BOOK`

type Book struct {
	Id       string
	Title    string
	Chapters []BookChapter
}

func (b Book) Key() ([]string, error) {
	return []string{BookEntity, b.Id}, nil
}

type BookChapter struct {
	Pos   int
	Title string
}

type BooksChaincode struct {
	router *router.Group
}

func (cc *BooksChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return owner.SetFromCreator(cc.router.Context(`init`, stub))
}

func (cc *BooksChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// delegate handling to router
	return cc.router.Handle(stub)
}

func New() *BooksChaincode {
	r := router.New(`books`)
	debug.AddHandlers(r, `debug`, owner.Only)

	r.Invoke(`bookList`, bookList)
	r.Invoke(`bookGet`, bookGet, param.String(`id`))
	r.Invoke(`bookInsert`, bookInsert, param.Struct(`book`, &Book{}))
	r.Invoke(`bookUpsert`, bookUpsert, param.Struct(`book`, &Book{}))
	r.Invoke(`bookDelete`, bookDelete, param.String(`id`))

	return &BooksChaincode{r}
}

func bookList(c router.Context) (interface{}, error) {
	return c.State().List(BookEntity, &Book{})
}

func bookInsert(c router.Context) (interface{}, error) {
	book := c.Arg(`book`).(Book)
	return book, c.State().Insert(book)
}

func bookUpsert(c router.Context) (interface{}, error) {
	book := c.Arg(`book`).(Book)
	return book, c.State().Put(book)
}

func bookGet(c router.Context) (interface{}, error) {
	return c.State().Get(Book{Id: c.ArgString(`id`)})
}

func bookDelete(c router.Context) (interface{}, error) {
	return nil, c.State().Delete(Book{Id: c.ArgString(`id`)})
}

var _ = Describe(`CRUD`, func() {

	book1 := Book{
		Id: `ISBN-111`, Title: `first title`, Chapters: []BookChapter{
			{Pos: 1, Title: `chapter 111.1`}, {Pos: 2, Title: `chapter 111.2`}}}

	book2 := Book{
		Id: `ISBN-222`, Title: `second title`, Chapters: []BookChapter{
			{Pos: 1, Title: `chapter 222.1`}, {Pos: 2, Title: `chapter 222.2`}, {Pos: 3, Title: `chapter 222.3`}}}

	book3 := Book{
		Id: `ISBN-333`, Title: `third title`, Chapters: []BookChapter{
			{Pos: 1, Title: `chapter 333.1`}, {Pos: 2, Title: `chapter 333.2`}, {Pos: 3, Title: `chapter 333.3`}, {Pos: 4, Title: `chapter 333.4`}}}

	//Create chaincode mock
	cc := testcc.NewMockStub(`debuggable`, New())
	actors, err := identity.ActorsFromPemFile(`SOME_MSP`, map[string]string{
		`owner`: `s7techlab.pem`,
	}, examplecert.Content)
	if err != nil {
		panic(err)
	}

	owner := actors[`owner`]
	cc.From(owner).Init()

	It("Allow to insert entries", func() {
		expectcc.ResponseOk(cc.From(owner).Invoke(`bookInsert`, &book1))
		expectcc.ResponseOk(cc.From(owner).Invoke(`bookInsert`, &book2))
		expectcc.ResponseOk(cc.From(owner).Invoke(`bookInsert`, &book3))
	})

	It("Disallow to insert entries with same keys", func() {
		expectcc.ResponseError(cc.From(owner).Invoke(`bookInsert`, &book3))
	})

	It("Allow to get entry list", func() {
		books := expectcc.PayloadIs(cc.Invoke(`bookList`), &[]Book{}).([]Book)
		Expect(len(books)).To(Equal(3))
		Expect(books[0]).To(Equal(book1))
		Expect(books[1]).To(Equal(book2))
		Expect(books[2]).To(Equal(book3))
	})

	It("Allow to get entry converted to target type", func() {
		book1FromCC := expectcc.PayloadIs(cc.Invoke(`bookGet`, book1.Id), &Book{}).(Book)
		Expect(book1FromCC).To(Equal(book1))
	})

	It("Allow to get entry json", func() {
		book2JsonFromCC := cc.Invoke(`bookGet`, book2.Id).Payload
		book2Json, _ := json.Marshal(book2)
		Expect(book2JsonFromCC).To(Equal(book2Json))
	})

	It("Allow to upsert entry", func() {
		book3Updated := book3
		book3Updated.Title = `thirdiest title`

		updateRes := expectcc.PayloadIs(cc.Invoke(`bookUpsert`, &book3Updated), &Book{}).(Book)
		Expect(updateRes.Title).To(Equal(book3Updated.Title))

		book3FromCC := expectcc.PayloadIs(cc.Invoke(`bookGet`, book3.Id), &Book{}).(Book)
		Expect(book3FromCC).To(Equal(book3Updated))
	})

	It("Allow to delete entry", func() {
		expectcc.ResponseOk(cc.From(owner).Invoke(`bookDelete`, book1.Id))
		books := expectcc.PayloadIs(cc.Invoke(`bookList`), &[]Book{}).([]Book)
		Expect(len(books)).To(Equal(2))

		expectcc.ResponseError(cc.Invoke(`bookGet`, book1.Id), state.ErrKeyNotFound)
	})

})
