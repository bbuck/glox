package printer_test

import (
	"testing"

	"github.com/bbuck/glox/token"
	"github.com/bbuck/glox/tree/expr"
	"github.com/bbuck/glox/tree/printer"
)

var ex = expr.NewBinary(
	expr.NewUnary(
		token.New(token.Minus, "-", nil, 1),
		expr.NewLiteral(expr.NumberLiteral, float64(123)),
	),
	token.New(token.Star, "*", nil, 1),
	expr.NewGrouping(
		expr.NewLiteral(expr.NumberLiteral, float64(45.67)),
	),
)

func Test_Print(t *testing.T) {
	result := printer.Print(ex)
	expected := "(* (- 123) (group 45.67))"
	if result != expected {
		t.Errorf("expected %q but got %q", expected, result)
	}
}

func Test_PrintRPN(t *testing.T) {
	expected := "123 - 45.67 *"
	result := printer.PrintRPN(ex)
	if result != expected {
		t.Errorf("expected %q but got %q", expected, result)
	}
}
