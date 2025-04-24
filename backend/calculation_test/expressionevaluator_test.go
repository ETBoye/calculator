package calculation_test

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"testing"

	"github.com/etboye/calculator/calculation"
)

type PanickyParser struct{}

func (p PanickyParser) Parse(input calculation.CalculationInput) (*calculation.Expression, error) {
	panic("This is my panic message")
}

var parsingErrorId string = "PARSING_ERROR"
var lexingErrorId string = "LEXING_ERROR"
var emptyInputErrorId string = "EMPTY_INPUT"

func TestParserPanicRecovery(t *testing.T) {
	log.SetOutput(io.Discard)

	calculatorWithPanickyParser := calculation.NewExpressionCalculatorWithParser(PanickyParser{})
	result := calculatorWithPanickyParser.Compute(calculation.CalculationInput{Input: "1+1"})

	if result.ErrorId == nil {
		t.Errorf("Expected errorId %s, got nil", parsingErrorId)
		return
	}

	if *result.ErrorId != parsingErrorId {
		t.Errorf("Expected errorId %s, got %s", parsingErrorId, *result.ErrorId)
		return
	}
}

func TestParsing(t *testing.T) {
	log.SetOutput(io.Discard)
	calculator := calculation.NewDefaultExpressionCalculator()

	assertExpectedError := func(inputString string, expectedErrorId string) {
		assertExpectedError(t, calculator, inputString, expectedErrorId)
	}

	assertExpectedError("(", parsingErrorId)
	assertExpectedError("+1", parsingErrorId)
	assertExpectedError("1-", parsingErrorId)
	assertExpectedError("(1+5/(2040*2*2+4)", parsingErrorId) // mismatched brackets
	assertExpectedError("//", parsingErrorId)
	assertExpectedError("+/+", parsingErrorId)
	assertExpectedError("", emptyInputErrorId)

	assertExpectedError("sdfadsf", lexingErrorId)
}

func TestCalculation(t *testing.T) {
	log.SetOutput(io.Discard)

	calculator := calculation.NewDefaultExpressionCalculator()

	// A bit of currying for readability
	assertExpectedOutput := func(inputString string, expectedOutput *big.Rat) {
		assertExpectedOutput(t, calculator, inputString, expectedOutput)
	}

	assertExpectedError := func(inputString string, expectedErrorId string) {
		assertExpectedError(t, calculator, inputString, expectedErrorId)
	}

	// We will use these to test that we can handle big constants
	bigNumberString1 := "8234098523049852035023940239450298435"
	bigNumberString2 := "234543298987120938120398213987348734959879345"

	bigNumber1 := intStringToRat(bigNumberString1)
	bigNumber2 := intStringToRat(bigNumberString2)

	tempRat := big.NewRat(1, 1) // used for in-place calculation, which big.Rat uses

	// Basic constants
	assertExpectedOutput("0", big.NewRat(0, 1))
	assertExpectedOutput("-0", big.NewRat(0, 1))

	assertExpectedOutput("-10", big.NewRat(-10, 1))

	assertExpectedOutput("123123", big.NewRat(123123, 1))
	assertExpectedOutput("0+ 123123", big.NewRat(123123, 1))

	// Basic operations
	assertExpectedOutput("0+0", big.NewRat(0, 1))
	assertExpectedOutput("0+-0", big.NewRat(0, 1))
	assertExpectedOutput("0+-1", big.NewRat(-1, 1))
	assertExpectedOutput("-1+0", big.NewRat(-1, 1))

	assertExpectedOutput("1+ 1", big.NewRat(2, 1)) // There was a bug with whitespace earlier..
	assertExpectedOutput("1 + 1", big.NewRat(2, 1))
	assertExpectedOutput(" 1 + 1 ", big.NewRat(2, 1))

	assertExpectedOutput(bigNumberString1+" + "+bigNumberString2, tempRat.Add(bigNumber1, bigNumber2))
	assertExpectedOutput(bigNumberString1+" - "+bigNumberString2, tempRat.Sub(bigNumber1, bigNumber2))
	assertExpectedOutput(bigNumberString1+" / "+bigNumberString2, tempRat.Quo(bigNumber1, bigNumber2))
	assertExpectedOutput(bigNumberString1+" * "+bigNumberString2, tempRat.Mul(bigNumber1, bigNumber2))

	// Order of operations
	assertExpectedOutput("2+3*4", big.NewRat(14, 1))  // (2+3)*4 is 20, which would be the wrong order
	assertExpectedOutput("2+3/4", big.NewRat(11, 4))  // (2+3)/4 is 5/4, which would be the wrong order
	assertExpectedOutput("2-3*4", big.NewRat(-10, 1)) // (2-3)*4 is -4, which would be the wrong order

	// Left to right behaviour
	// Not sure we want this behaviour, but now that it's there, we don't want it to change
	assertExpectedOutput("2/3/4/5", big.NewRat(1, 30))
	assertExpectedOutput("2/(3/(4/5))", big.NewRat(8, 15))

	// Brackets
	assertExpectedOutput("(2+3)*4", big.NewRat(20, 1))
	assertExpectedOutput("(2+3)/4", big.NewRat(5, 4))
	assertExpectedOutput("(2-3)*4", big.NewRat(-4, 1))

	assertExpectedOutput("(2-(3*65*2)*4+5)*300", big.NewRat(-465900, 1))
	assertExpectedOutput("(3/5+5432/83)*(-432)*(1+5/(2040*2*2+4))", big.NewRat(-24181645068, 847015))

	// Failures
	assertExpectedError("1/0", "DIVISION_BY_ZERO")
	assertExpectedError("(1/0)", "DIVISION_BY_ZERO")   // Triggers different error handling in coverage - at least as of now!
	assertExpectedError("1*(1/0)", "DIVISION_BY_ZERO") // Triggers different error handling in coverage - at least as of now!
	assertExpectedError("1+(1/0)", "DIVISION_BY_ZERO") // Triggers different error handling in coverage - at least as of now!

	// TODO: Test for big output
}

func intStringToRat(intString string) *big.Rat {
	newRat := big.NewRat(1, 1)
	newRat.SetString(fmt.Sprintf("%s/1", intString))
	return newRat
}

func assertExpectedError(t *testing.T, calculator calculation.Calculator, inputString string, expectedErrorId string) {
	t.Run(fmt.Sprintf("Expect input %s to give errorId %s", inputString, expectedErrorId), func(t *testing.T) {
		calculationResult := calculator.Compute(calculation.CalculationInput{Input: inputString})

		if calculationResult.ErrorId == nil {
			t.Errorf("Expected calculation error, but there was none")
			return
		}

		if *calculationResult.ErrorId != expectedErrorId {
			t.Errorf("Expected errorId %s, got errorId %s", expectedErrorId, *calculationResult.ErrorId)
		}
	})
}

func assertExpectedOutput(t *testing.T, calculator calculation.Calculator, inputString string, expectedOutput *big.Rat) {
	t.Run(fmt.Sprintf("Expect input: %s to output rational equal to %s", inputString, expectedOutput.String()), func(t *testing.T) {
		calculationResult := calculator.Compute(calculation.CalculationInput{Input: inputString})

		if calculationResult.ErrorId != nil {
			t.Errorf("calculationResult had a non-null error id - expected nil, got %s", *calculationResult.ErrorId)
		}

		resultAsRat := big.NewRat(0, 1)
		resultAsRat.SetString(fmt.Sprintf("%s/%s", calculationResult.Result.Num, calculationResult.Result.Denom))
		if !ratsAreEqual(expectedOutput, resultAsRat) {
			t.Errorf("calculationResult was %s, expected %s", resultAsRat.String(), expectedOutput.String())

		}
	})
}

func ratsAreEqual(a, b *big.Rat) bool {
	return a.Cmp(b) == 0
}
