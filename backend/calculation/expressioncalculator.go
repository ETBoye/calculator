package calculation

import (
	"errors"
	"log"

	"github.com/etboye/calculator/errorid"
	"github.com/etboye/calculator/util"
)

type Parser interface {
	Parse(input string) (*Expression, error)
}

type ExpressionCalculator struct {
	parser Parser
}

func NewExpressionCalculatorWithParser(parser Parser) ExpressionCalculator {
	return ExpressionCalculator{parser: parser}
}

func NewDefaultExpressionCalculator() ExpressionCalculator {
	parser := newParticipleParser()
	return ExpressionCalculator{parser: &parser}
}

func (calculator ExpressionCalculator) Compute(input string) CalculationResult {
	// I think I've seen the parser panic sometimes on really bad inputs
	// I can't seem to reproduce this - but we try to protect against it anyways
	expression, err := util.RecoverFromPanicWithError(
		func() (*Expression, error) { return calculator.parser.Parse(input) },
		nil, errors.New(errorid.PARSING_OR_LEXING_PANIC_ERROR), "Recovered from parsing panic")

	if err != nil {
		errorId := err.Error()
		return CalculationResult{Input: &input, ErrorId: &errorId}
	}

	resultAsRat, err := expression.Eval()

	if err != nil {
		errorId := err.Error()
		log.Println("Calculation returned error with id", errorId)
		return CalculationResult{
			Input:   &input,
			ErrorId: &errorId,
		}
	}

	result := RationalNumber{
		Num:      resultAsRat.Num().String(),
		Denom:    resultAsRat.Denom().String(),
		Estimate: resultAsRat.FloatString(PRECISION_DIGITS),
	}

	return CalculationResult{Input: &input, Result: &result}
}
