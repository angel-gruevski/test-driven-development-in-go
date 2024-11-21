package input

import "fmt"

type Validator struct {
	expectedLength int
	validOperators []string
}

func NewValidator(expectedLength int, validOps []string) *Validator {
	return &Validator{
		expectedLength: expectedLength,
		validOperators: validOps,
	}
}

// ValidationHelper to help us mock Validator
type ValidationHelper interface {
	CheckInput(operator string, operands []float64) error
}

func (v *Validator) CheckInput(operator string, operands []float64) error {
	operandsLength := len(operands)
	if operandsLength != v.expectedLength {
		return fmt.Errorf("unexpected operands length: got %d, want %d", operandsLength, v.expectedLength)
	}

	return v.checkOperator(operator)
}

// checkOperator validates the operator is supported
func (v *Validator) checkOperator(operator string) error {
	for _, o := range v.validOperators {
		if o == operator {
			return nil
		}
	}

	return fmt.Errorf("invalid operator: %s", operator)
}
