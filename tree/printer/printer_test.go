package printer_test

import (
	"testing"

	"github.com/bbuck/glox/token"
	"github.com/bbuck/glox/tree/expr"
	"github.com/bbuck/glox/tree/printer"
)

func Test_Print(t *testing.T) {
	ex := expr.NewBinary(
		expr.NewUnary(
			token.New(token.Minus, "-", nil, 1),
			expr.NewLiteral(expr.NumberLiteral, float64(123)),
		),
		token.New(token.Star, "*", nil, 1),
		expr.NewGrouping(
			expr.NewLiteral(expr.NumberLiteral, float64(45.67)),
		),
	)
	if printer.Print(ex) != "(* (- 123) (group 45.67))" {
		t.Fail()
	}
}
