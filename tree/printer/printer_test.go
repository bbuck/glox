package printer_test

import (
	"testing"

	"github.com/bbuck/glox/token"
	"github.com/bbuck/glox/tree/expr"
	"github.com/bbuck/glox/tree/printer"
)

var (
	minus = token.New(token.Minus, "-", nil, 1)
	star  = token.New(token.Star, "*", nil, 1)
	plus  = token.New(token.Plus, "+", nil, 1)

	ex = expr.NewBinary(
		expr.NewBinary(
			number(8),
			plus,
			number(10),
		),
		star,
		expr.NewUnary(
			minus,
			number(8),
		),
	)

	tex = expr.NewTernary(
		expr.NewBinary(
			number(88),
			minus,
			number(44),
		),
		expr.NewUnary(
			minus,
			number(1),
		),
		ex,
	)
)

func number(n float64) *expr.Literal {
	return expr.NewLiteral(expr.NumberLiteral, n)
}

func Test_Print(t *testing.T) {
	expect(
		t,
		"(* (+ 8 10) (- 8))",
		printer.Print(ex),
	)
}

func Test_Print_Ternary(t *testing.T) {
	expect(
		t,
		"(if (- 88 44) (- 1) (* (+ 8 10) (- 8)))",
		printer.Print(tex),
	)
}

func Test_PrintRPN(t *testing.T) {
	expect(
		t,
		"8 10 + 8 - *",
		printer.PrintRPN(ex),
	)
}

func Test_PrintRPN_Ternary(t *testing.T) {
	expect(
		t,
		"88 44 - ? 1 - : 8 10 + 8 - * ;",
		printer.PrintRPN(tex),
	)
}

func expect(t *testing.T, expected, result string) {
	if result != expected {
		t.Errorf("expected %q but got %q", expected, result)
	}
}
