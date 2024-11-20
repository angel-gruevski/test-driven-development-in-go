package format_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter02/format"
)

func TestError(t *testing.T) {
	// Arrange
	initialErr := errors.New("error msg")
	expr := "2%3"

	// Act
	wrappedError := format.Error(expr, initialErr)

	// Assert
	if !strings.Contains(wrappedError.Error(), initialErr.Error()) {
		t.Errorf("error does not contain: got %s, want %s", wrappedError.Error(), initialErr.Error())
	}
	if !strings.Contains(wrappedError.Error(), expr) {
		t.Errorf("error does not contain: got %s, want %s", wrappedError.Error(), expr)
	}
}
