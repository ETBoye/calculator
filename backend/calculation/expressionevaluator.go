package calculation

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const PRECISION_DIGITS = 5

var PARSING_ERROR_ID string = "PARSING_ERROR"
var LEXING_ERROR_ID string = "LEXING_ERROR"
var EMPTY_INPUT_ERROR_ID string = "EMPTY_INPUT"
var DIVISION_BY_ZERO_ERROR_ID string = "DIVISION_BY_ZERO"

// Borrowing heavily from https://github.com/alecthomas/participle/blob/master/_examples/expr/main.go
type Operator int

const (
	OpMul Operator = iota
	OpQuo
	OpAdd
	OpSub
)

var operatorMap = map[string]Operator{"+": OpAdd, "-": OpSub, "*": OpMul, "/": OpQuo}

func (o *Operator) Capture(s []string) error {
	*o = operatorMap[s[0]]
	return nil
}

type Value struct {
	IntegerWithoutSign *string     `@UnsignedInteger`
	IntegerWithSign    *string     `| OpSub@UnsignedInteger`
	Subexpression      *Expression `| StartParen@@EndParen`
}

type OpFactor struct {
	Operator Operator `@(OpMul|OpQuo)`
	Factor   *Value   `@@`
}

type Term struct {
	Left  *Value      `@@`
	Right []*OpFactor `@@*`
}

type OpTerm struct {
	Operator Operator `@(OpAdd|OpSub)`
	Term     *Term    `@@`
}

type Expression struct {
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
}

// Display

func (o Operator) String() string {
	switch o {
	case OpMul:
		return "*"
	case OpQuo:
		return "/"
	case OpSub:
		return "-"
	case OpAdd:
		return "+"
	}
	panic("unsupported operator")
}

func (v *Value) String() string {
	if v.IntegerWithoutSign != nil {
		return *v.IntegerWithoutSign
	} else if v.IntegerWithSign != nil {
		return "-" + *v.IntegerWithSign
	}
	return "(" + v.Subexpression.String() + ")"
}

func (o *OpFactor) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Factor)
}

func (t *Term) String() string {
	out := []string{t.Left.String()}
	for _, r := range t.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

func (o *OpTerm) String() string {
	return fmt.Sprintf("%s %s", o.Operator, o.Term)
}

func (e *Expression) String() string {
	out := []string{e.Left.String()}
	for _, r := range e.Right {
		out = append(out, r.String())
	}
	return strings.Join(out, " ")
}

// EVAL

// precondition: digitStringWithSign matches `-\d+`
func ratFromIntegerString(digitStringWithSign string) *big.Rat {
	rat := big.NewRat(1, 1)

	bigInt := big.NewInt(1)
	bigInt.SetString(digitStringWithSign, 10)
	rat.SetInt(bigInt)

	return rat
}

func (t *Term) Eval() (*big.Rat, error) {
	result, err := t.Left.Eval()

	if err != nil {
		return result, err
	}

	if t.Right == nil {
		return result, nil
	}

	for _, opTerm := range t.Right {
		factorEval, err := opTerm.Factor.Eval()

		if err != nil {
			return big.NewRat(1, 1), err
		}

		if opTerm.Operator == OpMul {
			result.Mul(result, factorEval)
		} else {
			if factorEval.Cmp(big.NewRat(0, 1)) == 0 {
				return big.NewRat(1, 1), errors.New(DIVISION_BY_ZERO_ERROR_ID)
			}

			result.Quo(result, factorEval) // TODO: Zero division and test for the same
		}
	}

	return result, nil
}

func (e *Expression) Eval() (*big.Rat, error) {
	result, err := e.Left.Eval()

	if err != nil {
		return result, err
	}

	if e.Right == nil {
		return result, nil
	}

	for _, opTerm := range e.Right {
		termEval, err := opTerm.Term.Eval()

		if err != nil {
			return big.NewRat(1, 1), err
		}

		if opTerm.Operator == OpAdd {
			result.Add(result, termEval)
		} else {
			result.Sub(result, termEval)
		}
	}

	return result, nil
}

func (v *Value) Eval() (*big.Rat, error) {
	if v.IntegerWithSign != nil {
		return ratFromIntegerString("-" + strings.TrimSpace(*v.IntegerWithSign)), nil
	} else if v.IntegerWithoutSign != nil {
		return ratFromIntegerString(strings.TrimSpace(*v.IntegerWithoutSign)), nil
	}

	return v.Subexpression.Eval()
}

type ExpressionCalculator struct {
	parser        *participle.Parser[Expression]
	lexingSymbols map[lexer.TokenType]string
}

func NewExpressionCalculator() ExpressionCalculator {

	var myLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "UnsignedInteger", Pattern: `\d+`},
		{Name: "OpAdd", Pattern: `\+`},
		{Name: "OpSub", Pattern: `-`},
		{Name: "OpMul", Pattern: `\*`},
		{Name: "OpQuo", Pattern: `/`},
		{Name: "StartParen", Pattern: `\(`},
		{Name: "EndParen", Pattern: `\)`},
		{Name: "WhiteSpace", Pattern: `[\s]*`},
	})

	parser := participle.MustBuild[Expression](
		participle.Lexer(myLexer),
		participle.Elide("WhiteSpace"), // The parser should ignore any whitespace

	)
	return ExpressionCalculator{parser: parser, lexingSymbols: lexer.SymbolsByRune(myLexer)}
}

func (calculator ExpressionCalculator) Compute(input CalculationInput) CalculationResult {
	if len(strings.TrimSpace(input.Input)) == 0 {
		log.Printf("Received empty input")
		return CalculationResult{
			ErrorId: &EMPTY_INPUT_ERROR_ID,
		}
	}

	tokens, lexerError := calculator.parser.Lex("", strings.NewReader(input.Input))

	if lexerError != nil {
		log.Println("Lexing error thrown:", lexerError.Error())

		return CalculationResult{
			ErrorId: &LEXING_ERROR_ID,
		}
	}

	logLexingResult(input, tokens, calculator.lexingSymbols)

	expression, err := calculator.parser.ParseString("", input.Input) // This actually lexes again. We accept this..

	// TODO: error handling with test
	if err != nil {
		log.Println("Parsing error thrown:", err.Error())

		return CalculationResult{
			ErrorId: &PARSING_ERROR_ID,
		}
	}

	log.Printf("Parsing succesful. Has parsed to %s", expression)

	resultAsRat, err := expression.Eval()

	if err != nil {
		errorId := err.Error()
		log.Println("Calculation returned error with id", errorId)
		return CalculationResult{
			ErrorId: &errorId,
		}
	}

	result := RationalNumber{
		Num:      resultAsRat.Num().String(),
		Denom:    resultAsRat.Denom().String(),
		Estimate: resultAsRat.FloatString(PRECISION_DIGITS),
	}

	return CalculationResult{Result: &result}
}

func logLexingResult(input CalculationInput, tokens []lexer.Token, lexingSymbols map[lexer.TokenType]string) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Lexed input \"%s\". Tokens: ", input.Input))

	for _, token := range tokens {
		sb.WriteString(fmt.Sprintf("%s(%s)", lexingSymbols[token.Type], token.Value))
	}

	log.Println(sb.String())
}
