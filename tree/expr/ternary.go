package expr

// Ternary represents an inline 3-expression operation with a
// condition and two other expressions. If the result of the
// condition is true-thy then the result of the Positive
// expression will be returned, if not then the Negative
// expression's result will be returned.
type Ternary struct {
	Condition Expr
	Positive  Expr
	Negative  Expr
}

// NewTernary constructs a new Ternary expression from the
// given expressions.
func NewTernary(cond, pos, neg Expr) *Ternary {
	return &Ternary{
		Condition: cond,
		Positive:  pos,
		Negative:  neg,
	}
}

// Accept for a Ternary calls VisitTernary on the Visitor.
func (t *Ternary) Accept(v Visitor) {
	v.VisitTernary(t)
}
