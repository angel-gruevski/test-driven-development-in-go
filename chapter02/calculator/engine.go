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
	validOperations map[string]func(x, y float64) (*float64, error)
}

func NewEngine() *Engine {
	engine := Engine{
		expectedLength:  2,
		validOperations: make(map[string]func(x, y float64) (*float64, error)),
	}
	engine.validOperations["+"] = engine.Add
	engine.validOperations["-"] = engine.Sub
	engine.validOperations["*"] = engine.Mult
	engine.validOperations["/"] = engine.Div
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
func (e *Engine) Add(x, y float64) (*float64, error) {
	result := x + y
	return &result, nil
}

// Sub is the function that processes the subtract operation
func (e *Engine) Sub(x, y float64) (*float64, error) {
	result := x - y
	return &result, nil
}

// Sub is the function that processes the subtract operation
func (e *Engine) Mult(x, y float64) (*float64, error) {
	result := x * y
	return &result, nil
}

// Div is the function that processes the divide operation
func (e *Engine) Div(x, y float64) (*float64, error) {
	if y == 0 {
		return nil, fmt.Errorf("cannot divide by zero")
	}
	result := x / y
	return &result, nil
}

// Interface to allow mocking the Engine struct
type OperationProcessor interface {
	ProcessOperation(operation Operation) (*string, error)
}

func (e *Engine) ProcessOperation(operation Operation) (*string, error) {
	// Comma ok idiom
	f, ok := e.validOperations[operation.Operator]
	if !ok {
		err := fmt.Errorf("no operation for operator %s found", operation.Operator)
		return nil, format.Error(operation.Expression, err)
	}

	operandsLength := len(operation.Operands)
	if operandsLength != e.GetNumOperands() {
		err := fmt.Errorf("incorrect number of operands")
		return nil, format.Error(operation.Expression, err)
	}

	res, err := f(operation.Operands[0], operation.Operands[1])
	if err != nil {
		return nil, format.Error(operation.Expression, err)
	}
	fres := format.Result(operation.Expression, *res)
	return &fres, nil
}
