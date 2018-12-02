package expr

// Grouping is an expression wrapped in parenthesis, so it represents
// a grouped single or set of expressions.
type Grouping struct {
	Expression Expr
}

// NewGrouping returns a new Grouping expression containing the wrapped
// expression.
func NewGrouping(expr Expr) *Grouping {
	return &Grouping{
		Expression: expr,
	}
}

// Accept for Grouping calls thes VisitGrouping method on the visitor.
func (g *Grouping) Accept(v Visitor) {
	v.VisitGrouping(g)
}
