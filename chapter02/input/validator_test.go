package input_test

import (
	"testing"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter02/input"
)

func TestCheckInput(t *testing.T) {
	// Arrange
	expectedLength := 2
	validOperators := []string{"+"}

	// Act & Assert
	t.Run("valid case", func(t *testing.T) {
		v := setup(t, validOperators, expectedLength)
		err := v.CheckInput(validOperators[0], []float64{2.5, 3.5})
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("invalid operands", func(t *testing.T) {
		v := setup(t, validOperators, expectedLength)
		err := v.CheckInput(validOperators[0], []float64{2.5, 2.2, 3.3})
		if err == nil {
			t.Fatal(err)
		}
	})

	t.Run("invalid operator", func(t *testing.T) {
		v := setup(t, validOperators, expectedLength)
		err := v.CheckInput("-", []float64{2.5, 55})
		if err == nil {
			t.Fatal(err)
		}
	})
}

func setup(t *testing.T, validOps []string, expectedLength int) *input.Validator {
	t.Helper()
	return input.NewValidator(expectedLength, validOps)
}
