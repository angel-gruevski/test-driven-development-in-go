package input_test

import (
	"fmt"
	"testing"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter02/calculator"
	"github.com/angel-gruevski/test-driven-development-in-go/chapter02/input"
	"github.com/angel-gruevski/test-driven-development-in-go/chapter02/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProcessExpression(t *testing.T) {

	t.Run("valid input", func(t *testing.T) {
		// Arrange
		expr := "2 + 3"
		operator := "+"
		operands := []float64{2.0, 3.0}
		expectedResult := "2 + 3 = 5.5"
		engine := mocks.NewOperationProcessor(t)
		validator := mocks.NewValidationHelper(t)
		parser := input.NewParser(engine, validator)

		validator.On("CheckInput", operator, operands).Return(nil).Once()
		engine.On("ProcessOperation", calculator.Operation{
			Expression: expr,
			Operator:   operator,
			Operands:   operands,
		}).Return(&expectedResult, nil).Once()

		// Act
		result, error := parser.ProcessExpression(expr)

		// Assert
		require.Nil(t, error)
		require.NotNil(t, result)
		assert.Contains(t, *result, expectedResult)
		assert.Contains(t, *result, expr)
		validator.AssertExpectations(t)
		engine.AssertExpectations(t)
	})

	t.Run("incorrect number of operands", func(t *testing.T) {
		// Arrangement
		expr := "2 + 4 + 6"
		engine := mocks.NewOperationProcessor(t)
		validator := mocks.NewValidationHelper(t)
		parser := input.NewParser(engine, validator)
		want := "CALCULATION ERROR: expression " + expr + " is invalid: incorrect expression length:got 5, want 3"

		// Act
		result, error := parser.ProcessExpression(expr)

		// Assert
		require.Nil(t, result)
		require.NotNil(t, error)
		assert.EqualError(t, error, want)
	})

	t.Run("incorrect left operand", func(t *testing.T) {
		// Arrangement
		expr := "! + 4"
		engine := mocks.NewOperationProcessor(t)
		validator := mocks.NewValidationHelper(t)
		parser := input.NewParser(engine, validator)
		want := "CALCULATION ERROR: expression " + expr + " is invalid: unable to process expression"

		// Act
		result, error := parser.ProcessExpression(expr)

		// Assert
		require.Nil(t, result)
		require.NotNil(t, error)
		assert.EqualError(t, error, want)

	})

	t.Run("incorrect right operand", func(t *testing.T) {
		// Arrangement
		expr := "4 + !"
		engine := mocks.NewOperationProcessor(t)
		validator := mocks.NewValidationHelper(t)
		parser := input.NewParser(engine, validator)
		want := "CALCULATION ERROR: expression " + expr + " is invalid: unable to process expression"

		// Act
		result, error := parser.ProcessExpression(expr)

		// Assert
		require.Nil(t, result)
		require.NotNil(t, error)
		assert.EqualError(t, error, want)

	})

	t.Run("validation fail", func(t *testing.T) {
		// Arrangement
		expr := "4 + 10"
		engine := mocks.NewOperationProcessor(t)
		validator := mocks.NewValidationHelper(t)
		parser := input.NewParser(engine, validator)
		validator.On("CheckInput", mock.Anything, mock.Anything).Return(fmt.Errorf("invalid operands"))
		want := "CALCULATION ERROR: expression 4 + 10 is invalid: invalid operands"

		// Act
		result, error := parser.ProcessExpression(expr)

		// Assert
		require.Nil(t, result)
		require.NotNil(t, error)
		assert.EqualError(t, error, want)
		validator.AssertExpectations(t)
	})
}
