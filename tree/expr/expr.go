package expr

// Expr is an expression interface that defines an Accept method
type Expr interface {
	Accept(Visitor)
}
