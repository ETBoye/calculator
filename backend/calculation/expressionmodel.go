package calculation

import (
	"fmt"
	"strings"
)

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
