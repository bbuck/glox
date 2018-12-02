package expr

import "github.com/bbuck/glox/token"

// Binary represents an expression that contains two sub-expressions
// combined with an operator.
type Binary struct {
	Left     Expr
	Operator *token.T
	Right    Expr
}

// NewBinary constructs and returns a new binary expression.
func NewBinary(left Expr, op *token.T, right Expr) *Binary {
	return &Binary{
		Left:     left,
		Operator: op,
		Right:    right,
	}
}

// Accept for Binary calls the VisitBinary function on the visitor.
func (b *Binary) Accept(v Visitor) {
	v.VisitBinary(b)
}
