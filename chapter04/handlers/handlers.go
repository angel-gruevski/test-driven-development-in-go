package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter04/db"
	"github.com/gorilla/mux"
)

type Handler struct {
	bs *db.BookService
	us *db.UserService
}

func NewHandler(bs *db.BookService, us *db.UserService) *Handler {
	return &Handler{
		bs: bs,
		us: us,
	}
}

func (handler *Handler) Index(w http.ResponseWriter, r *http.Request) {
	// Send an HTTP status & a hardcoded message
	resp := &Response{
		Message: "Welcome to the BookSwap service!",
		Books:   handler.bs.List(),
	}
	writeResponse(w, http.StatusOK, resp)
}

func (handler *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	resp := &Response{
		Books: handler.bs.List(),
	}
	writeResponse(w, http.StatusOK, resp)
}

// UserUpsert is invoked by HTTP POST /users
func (handler *Handler) UserUpsert(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := readRequestBody(r)
	// Handle any errors & write an error HTTP status & response
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: fmt.Errorf("invalid user body: %v", err).Error(),
		})
		return
	}
	// Initialize a user to unmarshal request body into
	var user db.User
	if err := json.Unmarshal(body, &user); err != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response{
			Error: fmt.Errorf("invalid user body:%v", err).Error(),
		})
		return
	}
	// Call the repository method corresponding to the operation
	u, err := handler.us.Upsert(user)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(),
		})
		return
	}

	writeResponse(w, http.StatusOK, &Response{
		User: &u,
	})
}

// ListUserByID is invoked by HTTP GET /users/{id}.
func (handler *Handler) ListUserByID(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]
	user, book, error := handler.us.Get(userID)
	if error != nil {
		writeResponse(w, http.StatusNotFound, &Response{
			Error: error.Error(),
		})
	}
	writeResponse(w, http.StatusOK, &Response{
		Books: book,
		User:  user,
	})
}

// SwapBook is invoked by POST /books/{id}
func (h *Handler) SwapBook(w http.ResponseWriter, r *http.Request) {
	bookID := mux.Vars(r)["id"]
	userID := r.URL.Query().Get("user")
	if err := h.us.Exists(userID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(),
		})
		return
	}
	_, err := h.bs.SwapBook(bookID, userID)
	if err != nil {
		writeResponse(w, http.StatusNotFound, &Response{
			Error: err.Error(),
		})
		return
	}

	user, books, err := h.us.Get(userID)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: err.Error(),
		})
		return
	}

	writeResponse(w, http.StatusOK, &Response{
		User:  user,
		Books: books,
	})
}

// BookUpsert is invoked by HTTP POST /books.
func (h *Handler) BookUpsert(w http.ResponseWriter, r *http.Request) {
	body, err := readRequestBody(r)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, &Response{
			Error: fmt.Errorf("invalid book body:%v", err).Error(),
		})
		return
	}

	var book db.Book
	if error := json.Unmarshal(body, &book); error != nil {
		writeResponse(w, http.StatusUnprocessableEntity, &Response{
			Error: fmt.Errorf("invalid book body:%v", err).Error(),
		})
	}

	if err := h.us.Exists(book.OwnerID); err != nil {
		writeResponse(w, http.StatusBadRequest, &Response{
			Error: err.Error(),
		})
	}

	// Call the repository method corresponding to the operation
	book = h.bs.Upsert(book)
	// Send an HTTP success status & the return value from the repo
	writeResponse(w, http.StatusOK, &Response{
		Books: []db.Book{book},
	})
}

// readRequestBody is a helper method that
// allows to read a request body and return any errors.
func readRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return []byte{}, err
	}
	if err := r.Body.Close(); err != nil {
		return []byte{}, err
	}
	return body, err
}
