package calculator_test

import (
	"testing"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter02/calculator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNumOperands(t *testing.T) {
	// Arrange
	engine := calculator.NewEngine()

	// Act
	numOperands := engine.GetNumOperands()

	// Assert
	assert.Equal(t, 2, numOperands)
}

func TestGetValidOperators(t *testing.T) {
	// Arrange
	engine := calculator.NewEngine()

	// Act
	validOperators := engine.GetValidOperators()

	// Assert
	assert.Equal(t, 4, len(validOperators))
	assert.Contains(t, validOperators, "+")
	assert.Contains(t, validOperators, "-")
	assert.Contains(t, validOperators, "*")
	assert.Contains(t, validOperators, "/")
}

func TestAdd(t *testing.T) {
	// Arrange
	engine := calculator.NewEngine()

	// Act
	t.Run("positive input", func(t *testing.T) {
		// Arrange
		x, y := 2.0, 3.5
		want := 5.5

		// Act
		result, err := engine.Add(x, y)

		// Assert
		require.Nil(t, err)
		require.NotNil(t, result)
		assert.Equal(t, want, *result)
	})

	t.Run("negative input", func(t *testing.T) {
		// Arrange
		x, y := -3.5, -2.0
		want := -5.5

		// Act
		result, err := engine.Add(x, y)

		// Assert
		require.NotNil(t, result)
		require.Nil(t, err)
		assert.Equal(t, want, *result)
	})
}

func TestSub(t *testing.T) {
	// Arrange
	engine := calculator.NewEngine()

	t.Run("positive input", func(t *testing.T) {
		// Arrange
		x, y := 2.0, 3.5
		want := -1.5

		// Act
		result, error := engine.Sub(x, y)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, result)
		assert.Equal(t, want, *result)
	})

	t.Run("negative input", func(t *testing.T) {
		// Arrange
		x, y := -2.0, -3.5
		want := 1.5

		// Act
		result, error := engine.Sub(x, y)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, result)
		assert.Equal(t, want, *result)
	})
}

func TestMul(t *testing.T) {
	// Arrange
	engine := calculator.NewEngine()

	t.Run("positive input", func(t *testing.T) {
		// Arrange
		x, y := 2., 3.0
		want := 6.0

		// Act
		result, error := engine.Mult(x, y)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, result)
		assert.Equal(t, want, *result)
	})

	t.Run("negative input", func(t *testing.T) {
		// Arrange
		x, y := -2.0, -5.0
		want := 10.0

		// Act
		result, error := engine.Mult(x, y)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, result)
		assert.Equal(t, want, *result)
	})
}

func TestDiv(t *testing.T) {
	// Arrange
	engine := calculator.NewEngine()

	t.Run("positive input", func(t *testing.T) {
		// Arrange
		x, y := 10.0, 2.0
		want := 5.0

		// Act
		result, error := engine.Div(x, y)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, result)
		assert.Equal(t, want, *result)
	})

	t.Run("negative integer", func(t *testing.T) {
		// Arrange
		x, y := -10.0, -2.0
		want := 5.0

		// Act
		result, error := engine.Div(x, y)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, result)
		assert.Equal(t, want, *result)
	})

	t.Run("division by zero", func(t *testing.T) {
		// Arrange
		x, y := 2.0, 0.0

		// Act
		result, error := engine.Div(x, y)

		// Assert
		require.Nil(t, result)
		require.NotNil(t, error)
		assert.EqualError(t, error, "cannot divide by zero")
	})
}

func TestProcessOperation(t *testing.T) {
	// Arrange
	engine := calculator.NewEngine()

	t.Run("incorrect operator", func(t *testing.T) {
		// Arrange
		expression := "2.0 % 3.0"
		operator := "%"
		operands := []float64{2.0, 3.0}
		operation := calculator.Operation{
			Expression: expression,
			Operator:   operator,
			Operands:   operands,
		}

		// Act
		result, error := engine.ProcessOperation(operation)

		// Assert
		require.Nil(t, result)
		require.NotNil(t, error)
		assert.EqualError(t, error, "CALCULATION ERROR: expression "+expression+" is invalid: no operation for operator "+operator+" found")
	})

	t.Run("wrong number of operands", func(t *testing.T) {
		// Arrange
		expression := "2.0 - 1.5 - 0.5"
		operator := "-"
		operands := []float64{2.0, 1.5, 0.5}
		operation := calculator.Operation{
			Expression: expression,
			Operator:   operator,
			Operands:   operands,
		}

		// Act
		result, error := engine.ProcessOperation(operation)

		// Assert
		require.Nil(t, result)
		require.NotNil(t, error)
		assert.EqualError(t, error, "CALCULATION ERROR: expression "+expression+" is invalid: incorrect number of operands")
	})

	t.Run("correct input", func(t *testing.T) {
		// Arrange
		expression := "2.0 - 1.5"
		operator := "-"
		operands := []float64{2.0, 1.5}
		operation := calculator.Operation{
			Expression: expression,
			Operator:   operator,
			Operands:   operands,
		}
		want := "CALCULATION SUCCESS: " + expression + " = 0.50"

		// Act
		result, error := engine.ProcessOperation(operation)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, result)
		assert.Equal(t, want, *result)
	})

	t.Run("division by zero", func(t *testing.T) {
		// Arrange
		expression := "2.0 / 0.0"
		operator := "/"
		operands := []float64{2.0, 0.0}
		operation := calculator.Operation{
			Expression: expression,
			Operator:   operator,
			Operands:   operands,
		}
		want := "CALCULATION ERROR: expression 2.0 / 0.0 is invalid: cannot divide by zero"

		// Act
		result, error := engine.ProcessOperation(operation)

		// Assert
		require.Nil(t, result)
		require.NotNil(t, error)
		assert.EqualError(t, error, want)
	})
}
