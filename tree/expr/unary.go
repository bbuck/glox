package expr

import "github.com/bbuck/glox/token"

// Unary represents a unary expression where there is an operator and
// one other operand.
type Unary struct {
	Operator *token.T
	Right    Expr
}

// NewUnary constructs and returns a new Unary expression.
func NewUnary(op *token.T, right Expr) *Unary {
	return &Unary{
		Operator: op,
		Right:    right,
	}
}

// Accept for Unary calls the VisitUnary method on the visitor.
func (u *Unary) Accept(v Visitor) {
	v.VisitUnary(u)
}
