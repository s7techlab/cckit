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

const PrivateBookEntity = `PRIVATE_BOOK`

type PrivateBook struct {
	Id       string
	Title    string
	Chapters []BookChapter
}

func (pb PrivateBook) Key() ([]string, error) {
	return []string{PrivateBookEntity, pb.Id}, nil
}
