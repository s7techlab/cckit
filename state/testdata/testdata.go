package testdata

import (
	"github.com/s7techlab/cckit/state/testdata/schema"
)

var Books = []schema.Book{{
	Id:    `ISBN-111`,
	Title: `first title`, Chapters: []schema.BookChapter{
		{Pos: 1, Title: `chapter 111.1`}, {Pos: 2, Title: `chapter 111.2`}}},

	{
		Id: `ISBN-222`, Title: `second title`, Chapters: []schema.BookChapter{
			{Pos: 1, Title: `chapter 222.1`}, {Pos: 2, Title: `chapter 222.2`}, {Pos: 3, Title: `chapter 222.3`}}},

	{
		Id: `ISBN-333`, Title: `third title`, Chapters: []schema.BookChapter{
			{Pos: 1, Title: `chapter 333.1`}, {Pos: 2, Title: `chapter 333.2`}, {Pos: 3, Title: `chapter 333.3`}, {Pos: 4, Title: `chapter 333.4`}}},
}
