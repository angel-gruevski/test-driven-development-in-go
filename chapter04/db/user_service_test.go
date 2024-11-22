package db_test

import (
	"testing"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter04/db"
	"github.com/angel-gruevski/test-driven-development-in-go/chapter04/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	t.Run("existing-user", func(t *testing.T) {
		// Arrange
		userId := uuid.New().String()
		userOne := db.User{
			ID:       userId,
			Name:     "User One",
			Address:  "Miami Boulevard",
			PostCode: "1000",
			Country:  "Florida, US",
		}
		users := []db.User{userOne}
		book := db.Book{
			ID:      uuid.New().String(),
			Name:    "Book One",
			Author:  "Author One",
			OwnerID: userId,
			Status:  db.Available.String(),
		}
		bookOperationService := mocks.NewBookOperationsService(t)
		bookOperationService.On("ListByUser", userId).Return([]db.Book{book})
		userService := db.NewUserService(users, bookOperationService)

		// Act
		user, books, error := userService.Get(userId)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, user)
		require.NotNil(t, books)
		assert.Equal(t, book.Name, books[0].Name)
		bookOperationService.AssertExpectations(t)
	})

	t.Run("non-existing-user", func(t *testing.T) {
		// Arrange
		userId := uuid.New().String()
		userOne := db.User{
			ID:       uuid.New().String(),
			Name:     "User One",
			Address:  "Miami Boulevard",
			PostCode: "1000",
			Country:  "Florida, US",
		}
		initial := []db.User{userOne}
		bookOperationService := mocks.NewBookOperationsService(t)
		var books []db.Book = make([]db.Book, 0)
		userService := db.NewUserService(initial, bookOperationService)

		// Act
		user, books, error := userService.Get(userId)

		// Assert
		require.Nil(t, user)
		require.Nil(t, books)
		require.NotNil(t, error)
		assert.EqualError(t, error, "user does not exist")
		bookOperationService.AssertExpectations(t)
	})
}

func TestExists(t *testing.T) {
	eu := db.User{
		ID:   uuid.New().String(),
		Name: "Existing user",
	}
	bs := mocks.NewBookOperationsService(t)
	t.Run("user-exists", func(t *testing.T) {
		// Arrange
		us := db.NewUserService([]db.User{eu}, bs)

		// Act
		error := us.Exists(eu.ID)

		// Assert
		require.Nil(t, error)
	})

	t.Run("user-not-exists", func(t *testing.T) {
		// Arrange
		us := db.NewUserService([]db.User{eu}, bs)

		// Act
		error := us.Exists(uuid.New().String())

		// Assert
		require.NotNil(t, error)
		assert.EqualError(t, error, "no user found")
	})

	t.Run("empty users", func(t *testing.T) {
		// Arrange
		us := db.NewUserService([]db.User{}, bs)

		// Act
		error := us.Exists(eu.ID)

		// Assert
		require.NotNil(t, error)
		assert.EqualError(t, error, "no user found")
	})
}

func TestUpsertUser(t *testing.T) {
	eu := db.User{
		ID:   uuid.New().String(),
		Name: "Existing user",
	}
	bs := mocks.NewBookOperationsService(t)

	t.Run("user-exists", func(t *testing.T) {
		// Arrange
		us := db.NewUserService([]db.User{eu}, bs)
		updatedUser := db.User{
			ID:   eu.ID,
			Name: "Updated user",
		}
		// Act
		user, error := us.Upsert(updatedUser)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, user)
		assert.Equal(t, updatedUser.Name, user.Name)
	})
}
