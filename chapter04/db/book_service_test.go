package db_test

import (
	"errors"
	"testing"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter04/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Run("initial books", func(t *testing.T) {
		// Arrange
		book := db.Book{
			ID:     uuid.New().String(),
			Name:   "Existing book",
			Status: db.Available.String(),
		}

		// The difference between null and undefined in JS is that when you assign null to a variable, it points to void, non-existent location,
		// while undefined is the zero-value for the data types and it means that memory address is allocated but the value that is stored is undefined !!!
		bookService := db.NewBookService([]db.Book{book}, nil)

		tests := map[string]struct {
			id      string
			want    db.Book
			wantErr error
		}{
			"existing books": {
				id:   book.ID,
				want: book,
			},
			"no book found": {
				id:      "not-found",
				wantErr: errors.New("no book found"),
			},
			"empty-id": {
				id:      "",
				wantErr: errors.New("no book found"),
			},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				// Act
				book, err := bookService.Get(tc.id)

				// Assert
				if tc.wantErr != nil {
					require.Nil(t, book)
					require.NotNil(t, err)
					assert.EqualError(t, err, tc.wantErr.Error())
					return
				}
				require.Nil(t, err)
				require.NotNil(t, book)
				assert.Equal(t, tc.want, *book)
			})
		}
	})

	t.Run("empty books", func(t *testing.T) {
		// Arrange
		bs := db.NewBookService([]db.Book{}, nil)
		bookId := "1122-3322"

		// Act
		b, err := bs.Get(bookId)

		// Assert
		assert.Equal(t, errors.New("no book found"), err)
		assert.Nil(t, b)
	})
}

func TestUpsert(t *testing.T) {
	t.Run("existing-book", func(t *testing.T) {
		// Arrange
		book := db.Book{
			ID:      uuid.New().String(),
			Name:    "Existing book",
			Status:  db.Available.String(),
			OwnerID: uuid.New().String(),
		}

		bookService := db.NewBookService([]db.Book{book}, nil)
		updatedBook := db.Book{
			ID:     book.ID,
			Name:   "Updated",
			Status: book.Status,
		}

		// Act
		returnedBook := bookService.Upsert(updatedBook)

		// Assert
		require.NotNil(t, returnedBook)
		assert.Equal(t, updatedBook, returnedBook)
	})

	t.Run("new-book", func(t *testing.T) {
		// Arrange
		book := db.Book{
			Name:    "New",
			OwnerID: uuid.New().String(),
		}
		bookService := db.NewBookService([]db.Book{}, nil)

		// Act
		returnedBook := bookService.Upsert(book)

		// Assert
		require.NotNil(t, returnedBook)
		assert.NotEmpty(t, returnedBook.ID)
		assert.Equal(t, book.Name, returnedBook.Name)
		assert.Equal(t, book.OwnerID, returnedBook.OwnerID)
		assert.Equal(t, db.Available.String(), returnedBook.Status)
	})
}

func TestList(t *testing.T) {
	t.Run("existing-available-books", func(t *testing.T) {
		// Arrange
		bookOne := db.Book{
			ID:      uuid.New().String(),
			Name:    "Book 1",
			Author:  "Author 1",
			OwnerID: uuid.New().String(),
			Status:  db.Available.String(),
		}

		bookTwo := db.Book{
			ID:      uuid.New().String(),
			Name:    "Book 2",
			Author:  "Author 2",
			OwnerID: uuid.New().String(),
			Status:  db.Available.String(),
		}

		bookThree := db.Book{
			ID:      uuid.New().String(),
			Name:    "Book 3",
			Author:  "Author 3",
			OwnerID: uuid.New().String(),
			Status:  db.Swapped.String(),
		}

		bookService := db.NewBookService([]db.Book{bookOne, bookTwo, bookThree}, nil)

		// Act
		returnedAvailableBooks := bookService.List()

		// Assert
		assert.NotEmpty(t, returnedAvailableBooks)
		assert.Len(t, returnedAvailableBooks, 2)
		assert.Contains(t, returnedAvailableBooks, bookOne)
		assert.Contains(t, returnedAvailableBooks, bookTwo)
		assert.NotContains(t, returnedAvailableBooks, bookThree)
	})

	t.Run("no-avaiable-books", func(t *testing.T) {
		bookService := db.NewBookService([]db.Book{}, nil)

		// Act
		returnedAvailableBooks := bookService.List()

		assert.NotNil(t, returnedAvailableBooks)
		assert.Len(t, returnedAvailableBooks, 0)
	})
}

func TestListByUser(t *testing.T) {
	t.Run("existing-book", func(t *testing.T) {
		// Arrange
		ownerId := uuid.New().String()
		bookOne := db.Book{
			ID:      uuid.New().String(),
			Name:    "Book One",
			Author:  "Author One",
			OwnerID: ownerId,
			Status:  db.Available.String(),
		}
		bookTwo := db.Book{
			ID:      uuid.New().String(),
			Name:    "Book Two",
			Author:  "Author Two",
			OwnerID: ownerId,
			Status:  db.Available.String(),
		}
		bookThree := db.Book{
			ID:      uuid.New().String(),
			Name:    "Book Three",
			Author:  "Author Three",
			OwnerID: uuid.New().String(),
			Status:  db.Available.String(),
		}

		bookService := db.NewBookService([]db.Book{bookOne, bookTwo, bookThree}, nil)

		// Act
		books := bookService.ListByUser(ownerId)

		// Assert
		assert.Len(t, books, 2)
		assert.Contains(t, books, bookOne)
		assert.Contains(t, books, bookTwo)
		assert.NotContains(t, books, bookThree)
	})
}

func TestSwapBook(t *testing.T) {
	t.Run("available-book-for-swap", func(t *testing.T) {
		// Arrange
		ownerId := uuid.New().String()
		bookId := uuid.New().String()
		bookOne := db.Book{
			ID:      bookId,
			Name:    "Book One",
			Author:  "Author One",
			OwnerID: uuid.New().String(),
			Status:  db.Available.String(),
		}
		bookService := db.NewBookService([]db.Book{bookOne}, nil)

		// Act
		book, error := bookService.SwapBook(bookId, ownerId)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, book)
		assert.Equal(t, bookOne.Author, book.Author)
		assert.Equal(t, bookOne.ID, book.ID)
		assert.Equal(t, bookOne.Name, book.Name)
		assert.Equal(t, ownerId, book.OwnerID)
		assert.Equal(t, db.Swapped.String(), book.Status)
	})

	t.Run("book-does-not-exist", func(t *testing.T) {
		// Arrange
		ownerId := uuid.New().String()
		bookId := uuid.New().String()
		bookService := db.NewBookService([]db.Book{}, nil)

		// Act
		book, error := bookService.SwapBook(bookId, ownerId)

		// Assert
		require.Nil(t, book)
		require.NotNil(t, error)
		assert.EqualError(t, error, "book doesn't exist")
	})

	t.Run("book-not-swappable", func(t *testing.T) {
		// Arrange
		ownerId := uuid.New().String()
		bookId := uuid.New().String()
		bookOne := db.Book{
			ID:      bookId,
			Name:    "Book One",
			Author:  "Author One",
			OwnerID: uuid.New().String(),
			Status:  db.Swapped.String(),
		}
		bookService := db.NewBookService([]db.Book{bookOne}, nil)

		// Act
		book, error := bookService.SwapBook(bookId, ownerId)

		// Assert
		require.Nil(t, book)
		require.NotNil(t, error)
		assert.EqualError(t, error, "book is not available")
	})
}
