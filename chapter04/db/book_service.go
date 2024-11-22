package db

import (
	"errors"

	"github.com/google/uuid"
)

type Book struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Author  string `json:"author"`
	OwnerID string `json:"owner_id"`
	Status  string `json:"status"`
}

type BookService struct {
	books map[string]Book
	ps    PostingService
}

func NewBookService(initial []Book, ps PostingService) *BookService {
	books := make(map[string]Book)
	for _, b := range initial {
		books[b.ID] = b
	}
	return &BookService{
		books: books,
		ps:    ps,
	}
}

// Get returns a given book or error if none exists.
func (bs *BookService) Get(id string) (*Book, error) {
	// To Do: Implement
	book, ok := bs.books[id]
	if !ok {
		return nil, errors.New("no book found")
	}
	return &book, nil
}

// Upsert creates or updates a book.
func (bs *BookService) Upsert(b Book) Book {
	_, ok := bs.books[b.ID]
	if !ok {
		b.ID = uuid.NewString()
		b.Status = Available.String()
	}
	bs.books[b.ID] = b
	return b
}

// List returns the list of available books.
func (bs *BookService) List() []Book {
	var items []Book = make([]Book, 0)
	for _, b := range bs.books {
		if b.Status == Available.String() {
			items = append(items, b)
		}
	}
	return items
}

// ListByUser returns the list of books for a given user.
func (bs *BookService) ListByUser(userID string) []Book {
	var items = make([]Book, 0)
	for _, b := range bs.books {
		if b.OwnerID == userID {
			items = append(items, b)
		}
	}
	return items
}

// SwapBook checks whether a book is available and, if possible, marks it as swapped.
func (bs *BookService) SwapBook(bookID, userID string) (*Book, error) {
	book, ok := bs.books[bookID]
	if !ok {
		return nil, errors.New("book doesn't exist")
	}
	if book.Status == Swapped.String() {
		return nil, errors.New("book is not available")
	}
	book.Status = Swapped.String()
	book.OwnerID = userID
	return &book, nil
}
