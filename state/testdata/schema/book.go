package schema

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
