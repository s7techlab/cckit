package testdata

import (
	"errors"
	"fmt"

	"github.com/s7techlab/cckit/extensions/debug"
	"github.com/s7techlab/cckit/extensions/owner"
	"github.com/s7techlab/cckit/router"
	p "github.com/s7techlab/cckit/router/param"
	"github.com/s7techlab/cckit/state"
	"github.com/s7techlab/cckit/state/testdata/schema"
)

const collection = "SampleCollection"

func NewBooksCC() *router.Chaincode {
	r := router.New(`books`)
	debug.AddHandlers(r, `debug`, owner.Only)

	r.Init(owner.InvokeSetFromCreator).
		Invoke(`bookList`, bookList).
		Invoke(`bookListPaginated`, bookListPaginated, p.Struct(`in`, &schema.BookListRequest{})).
		Invoke(`bookIds`, bookIds).
		Invoke(`bookGet`, bookGet, p.String(`id`)).
		Invoke(`bookInsert`, bookInsert, p.Struct(`book`, &schema.Book{})).
		Invoke(`bookUpsert`, bookUpsert, p.Struct(`book`, &schema.Book{})).
		Invoke(`bookUpsertWithCache`, bookUpsertWithCache, p.Struct(`book`, &schema.Book{})).
		Invoke(`bookDelete`, bookDelete, p.String(`id`)).
		Invoke(`privateBookList`, privateBookList).
		Invoke(`privateBookGet`, privateBookGet, p.String(`id`)).
		Invoke(`privateBookInsert`, privateBookInsert, p.Struct(`book`, &schema.PrivateBook{})).
		Invoke(`privateBookUpsert`, privateBookUpsert, p.Struct(`book`, &schema.PrivateBook{})).
		Invoke(`privateBookDelete`, privateBookDelete, p.String(`id`))

	return router.NewChaincode(r)
}

func bookList(c router.Context) (interface{}, error) {
	return c.State().List(schema.BookEntity, &schema.Book{})
}

func bookListPaginated(c router.Context) (interface{}, error) {
	in, ok := c.Param(`in`).(schema.BookListRequest)
	if !ok {
		return nil, fmt.Errorf("unexpected argument type")
	}

	list, md, err := c.State().ListPaginated(schema.BookEntity, in.PageSize, in.Bookmark, &schema.Book{})
	if err != nil {
		return nil, err
	}

	var books []*schema.Book
	for _, item := range list.([]interface{}) {
		var b = item.(schema.Book)
		books = append(books, &b)
	}

	return schema.BookList{
		Items: books,
		Next:  md.Bookmark,
	}, nil
}

func bookIds(c router.Context) (interface{}, error) {
	return c.State().Keys(schema.BookEntity)
}

func bookInsert(c router.Context) (interface{}, error) {
	book := c.Param(`book`)
	return book, c.State().Insert(book)
}

func bookUpsert(c router.Context) (interface{}, error) {
	book := c.Param(`book`).(schema.Book)

	// update data in state
	if err := c.State().Put(book); err != nil {
		return nil, err
	}

	//try to read new data in same transaction
	upsertedBook, err := c.State().Get(schema.Book{Id: book.Id}, &schema.Book{})
	if err != nil {
		return nil, err
	}

	// state read in same tx must return PREVIOUS value
	if book.Title == upsertedBook.(schema.Book).Title {
		return nil, errors.New(`read after write in same tx must return previous value`)
	}

	return book, err
}

func bookUpsertWithCache(c router.Context) (interface{}, error) {
	book := c.Param(`book`).(schema.Book)

	stateCached := state.WithCache(c.State())

	// update data in state
	if err := stateCached.Put(book); err != nil {
		return nil, err
	}

	//try to read new data in same transaction
	upsertedBook, err := stateCached.Get(schema.Book{Id: book.Id}, &schema.Book{})
	if err != nil {
		return nil, err
	}

	// state read in same tx with state caching must return NEW value
	if book.Title != upsertedBook.(schema.Book).Title {
		return nil, errors.New(`read after write with tx state caching must return same value`)
	}

	return book, err
}

func bookGet(c router.Context) (interface{}, error) {
	return c.State().Get(schema.Book{Id: c.ParamString(`id`)})
}

func bookDelete(c router.Context) (interface{}, error) {
	return nil, c.State().Delete(schema.Book{Id: c.ParamString(`id`)})
}

func privateBookList(c router.Context) (interface{}, error) {
	return c.State().ListPrivate(collection, false, schema.PrivateBookEntity, &schema.PrivateBook{})
}

func privateBookInsert(c router.Context) (interface{}, error) {
	book := c.Param(`book`)
	err := c.State().Insert(book, "{}")
	if err != nil {
		return book, err
	}
	return book, c.State().InsertPrivate(collection, book)
}

func privateBookUpsert(c router.Context) (interface{}, error) {
	book := c.Param(`book`)
	err := c.State().Put(book, "{}")
	if err != nil {
		return book, err
	}
	return book, c.State().PutPrivate(collection, book)
}

func privateBookGet(c router.Context) (interface{}, error) {
	return c.State().GetPrivate(collection, schema.PrivateBook{Id: c.ParamString(`id`)})
}

func privateBookDelete(c router.Context) (interface{}, error) {
	err := c.State().Delete(schema.PrivateBook{Id: c.ParamString(`id`)})
	if err != nil {
		return nil, err
	}
	return nil, c.State().DeletePrivate(collection, schema.PrivateBook{Id: c.ParamString(`id`)})
}
