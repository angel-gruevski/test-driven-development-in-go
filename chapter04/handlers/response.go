package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter04/db"
)

type Response struct {
	Message string    `json:"message,omitempty"`
	Error   string    `json:"error,omitempty"`
	Books   []db.Book `json:"books,omitempty"`
	User    *db.User  `json:"user,omitempty"`
}

func writeResponse(w http.ResponseWriter, status int, resp *Response) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		fmt.Fprintf(w, "error encoding resp %v:%s", resp, err)
	}
}
