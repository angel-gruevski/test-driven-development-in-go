package main

import (
	"flag"
	"log"

	"github.com/angel-gruevski/test-driven-development-in-go/chapter02/calculator"
	"github.com/angel-gruevski/test-driven-development-in-go/chapter02/input"
)

func main() {
	expr := flag.String("expression", "", "mathematical expression to parse")
	flag.Parse()

	engine := calculator.NewEngine()
	validator := input.NewValidator(engine.GetNumOperands(), engine.GetValidOperators())
	parser := input.NewParser(engine, validator)
	result, err := parser.ProcessExpression(*expr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*result)
}
