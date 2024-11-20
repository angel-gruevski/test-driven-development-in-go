package calculator

import (
	"fmt"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter02/format"
)

// Operation is the wrapper object that contains
// the operator and operand of a mathematical expression.
type Operation struct {
	Expression string
	Operator   string
	Operands   []float64
}

type Engine struct {
	expectedLength  int
	validOperations map[string]func(x, y float64) float64
}

func NewEngine() *Engine {
	engine := Engine{
		expectedLength:  2,
		validOperations: make(map[string]func(x float64, y float64) float64),
	}
	engine.validOperations["+"] = engine.Add
	return &engine
}

// GetNumOperands returns the expected number of operands that the engine can process.
func (e *Engine) GetNumOperands() int {
	return e.expectedLength
}

func (e *Engine) GetValidOperators() []string {
	var operators []string
	for o := range e.validOperations {
		operators = append(operators, o)
	}

	return operators
}

// Add is the function that processes the addition operation
func (e *Engine) Add(x, y float64) float64 {
	return x + y
}

func (e *Engine) ProcessOperation(operation Operation) (*string, error) {
	// Comma ok idiom
	f, ok := e.validOperations[operation.Operator]
	if !ok {
		err := fmt.Errorf("no operation for operator %s found", operation.Operator)
		return nil, format.Error(operation.Expression, err)
	}
	res := f(operation.Operands[0], operation.Operands[1])
	fres := format.Result(operation.Expression, res)
	return &fres, nil
}
