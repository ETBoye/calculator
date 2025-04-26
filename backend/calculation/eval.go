package calculation

import (
	"errors"
	"math/big"
	"strings"

	"github.com/etboye/calculator/errorid"
)

const PRECISION_DIGITS = 5

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
				return big.NewRat(1, 1), errors.New(errorid.DIVISION_BY_ZERO_ERROR)
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
